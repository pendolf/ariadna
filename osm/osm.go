package osm

import (
	"fmt"
	"github.com/gen1us2k/log"
	"github.com/maddevsio/ariadna/config"
	"github.com/maddevsio/ariadna/storage"
	"github.com/pkg/errors"
	"github.com/qedus/osmpbf"
	"io"
	"net/http"
	"os"
	"runtime"
)

type (
	OSM interface {
		DownloadFile(string, string) error
	}
	OSMWorker struct {
		OSM
		decoder *osmpbf.Decoder
		levelDB *storage.LevelDBStorage
		logger  log.Logger
	}
)

func New(conf *config.AriadnaConfig) (*OSMWorker, error) {
	if len(conf.FileName) < 1 {
		return nil, errors.New("Invalid file: you must specify a pbf path as arg[1]")
	}
	// try to open the file
	file, err := os.Open(conf.FileName)
	if err != nil {
		return nil, err
	}
	decoder := osmpbf.NewDecoder(file)
	err = decoder.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		return nil, err
	}
	db, err := storage.NewLevelDBStorage(config.LevelDBPath)
	if err != nil {
		return nil, err
	}

	return &OSMWorker{
		decoder: decoder,
		levelDB: db,
		logger:  log.NewLogger("osm"),
	}, nil
}
func (osm *OSMWorker) Run() error {
	//batch := &leveldb.Batch{}
	for {
		v, err := osm.decoder.Decode()
		if err == io.EOF {
			break
		}
		if err != nil {
			osm.logger.Error(err)
			continue
		}
		switch v := v.(type) {
		case *osmpbf.Node:
			osm.logger.Info("Node")
		case *osmpbf.Way:
			osm.logger.Info("Way")
		case *osmpbf.Relation:
			osm.logger.Info("Relation")
		default:
			osm.logger.Error("Unknown")

		}
	}
	return nil
}

func (osm *OSMWorker) DownloadFile(source, destination string) error {
	err := osm.downloadFile(source, destination)
	return err
}

func (osm *OSMWorker) downloadFile(source, destination string) error {

	output, err := os.Create(destination)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error while creating file %s", destination))
	}

	defer output.Close()

	client := &http.Client{}

	request, err := http.NewRequest("GET", source, nil)

	if err != nil {
		return errors.Wrap(err, "Caught error on create request")
	}

	request.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.3",
	)

	response, err := client.Do(request)

	if err != nil {
		return errors.Wrap(err, "Error while downloading")
	}
	defer response.Body.Close()

	_, err := io.Copy(output, response.Body)

	if err != nil {
		return errors.Wrap(err, "Error while copying data")
	}
	return nil

}
