package storage

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/qedus/osmpbf"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDBStorage struct {
	conn *leveldb.DB
}

func NewLevelDBStorage(dbPath string) (*LevelDBStorage, error) {
	conn, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBStorage{
		conn: conn,
	}, nil
}

func (db *LevelDBStorage) formatLevelDB(node *osmpbf.Node) (id string, val []byte) {
	stringid := strconv.FormatInt(node.ID, 10)
	var bufval bytes.Buffer
	bufval.WriteString(strconv.FormatFloat(node.Lat, 'f', 16, 64))
	bufval.WriteString(":")
	bufval.WriteString(strconv.FormatFloat(node.Lon, 'f', 16, 64))
	byteval := []byte(bufval.String())
	return stringid, byteval
}

// queue a leveldb write in a batch
func (db *LevelDBStorage) CacheQueue(batch *leveldb.Batch, node *osmpbf.Node) {
	id, val := db.formatLevelDB(node)
	batch.Put([]byte(id), []byte(val))
}

// flush a leveldb batch to database and reset batch to 0
func (db *LevelDBStorage) CacheFlush(batch *leveldb.Batch) error {
	err := db.conn.Write(batch, nil)
	if err != nil {
		return err
	}
	batch.Reset()
	return nil
}

func (db *LevelDBStorage) CacheLookup(way *osmpbf.Way) ([]map[string]string, error) {

	var container []map[string]string

	for _, each := range way.NodeIDs {
		stringid := strconv.FormatInt(each, 10)

		data, err := db.conn.Get([]byte(stringid), nil)
		if err != nil {
			return container, fmt.Errorf("denormalize failed for way: %d node not found: %s", way.ID, stringid)
		}

		s := string(data)
		spl := strings.Split(s, ":")

		latlon := make(map[string]string)
		lat, lon := spl[0], spl[1]
		latlon["lat"] = lat
		latlon["lon"] = lon

		container = append(container, latlon)

	}

	return container, nil
}
