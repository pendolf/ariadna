package storage

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/qedus/osmpbf"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"strconv"
	"strings"
)

func OpenFile(filename string) (*os.File, error) {
	// no file specified
	if len(filename) < 1 {
		return nil, errors.New("Invalid file: you must specify a pbf path as arg[1]")
	}
	// try to open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func OpenLevelDB(path string) (*leveldb.DB, error) {
	// try to open the db
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func formatLevelDB(node *osmpbf.Node) (id string, val []byte) {
	stringid := strconv.FormatInt(node.ID, 10)
	var bufval bytes.Buffer
	bufval.WriteString(strconv.FormatFloat(node.Lat, 'f', 16, 64))
	bufval.WriteString(":")
	bufval.WriteString(strconv.FormatFloat(node.Lon, 'f', 16, 64))
	byteval := []byte(bufval.String())
	return stringid, byteval
}

// queue a leveldb write in a batch
func cacheQueue(batch *leveldb.Batch, node *osmpbf.Node) {
	id, val := formatLevelDB(node)
	batch.Put([]byte(id), []byte(val))
}

// flush a leveldb batch to database and reset batch to 0
func cacheFlush(db *leveldb.DB, batch *leveldb.Batch) error {
	err := db.Write(batch, nil)
	if err != nil {
		return err
	}
	batch.Reset()
	return nil
}

func cacheLookup(db *leveldb.DB, way *osmpbf.Way) ([]map[string]string, error) {

	var container []map[string]string

	for _, each := range way.NodeIDs {
		stringid := strconv.FormatInt(each, 10)

		data, err := db.Get([]byte(stringid), nil)
		if err != nil {
			return fmt.Errorf("denormalize failed for way: %d node not found: %s", way.ID, stringid)
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
