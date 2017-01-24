package osm

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gen1us2k/log"
	"github.com/maddevsio/ariadna/config"
	"github.com/maddevsio/ariadna/geo"
	"github.com/maddevsio/ariadna/models"
	"github.com/maddevsio/ariadna/storage"
	"github.com/pkg/errors"
	"github.com/qedus/osmpbf"
	"github.com/syndtr/goleveldb/leveldb"
	"strconv"
)

type (
	OSM interface {
		DownloadFile(string, string) error
	}
	OSMWorker struct {
		OSM
		decoder   *osmpbf.Decoder
		levelDB   *storage.LevelDBStorage
		logger    log.Logger
		batch     *leveldb.Batch
		tags      map[string][]string
		appConfig *config.AriadnaConfig
	}
)

func New(conf *config.AriadnaConfig) (*OSMWorker, error) {
	if len(conf.FileName) < 1 {
		return nil, errors.New("Invalid file: you must specify a pbf path as arg[1]")
	}
	//// try to open the file
	//file, err := os.Open(conf.FileName)
	//if err != nil {
	//	return nil, err
	//}
	//decoder := osmpbf.NewDecoder(file)
	//err = decoder.Start(runtime.GOMAXPROCS(-1))
	//if err != nil {
	//	return nil, err
	//}
	db, err := storage.NewLevelDBStorage(config.LevelDBPath)
	if err != nil {
		return nil, err
	}

	return &OSMWorker{
		//decoder:   decoder,
		levelDB:   db,
		batch:     &leveldb.Batch{},
		appConfig: conf,
		logger:    log.NewLogger("osm"),
	}, nil
}
func (osm *OSMWorker) SetTags(tags map[string][]string) {
	osm.tags = tags
}
func (osm *OSMWorker) Run() error {

	for {
		v, err := osm.decoder.Decode()
		if err == io.EOF {
			osm.logger.Info("got end of file. Breaking")
			break
		}
		if err != nil {
			osm.logger.Error(err)
			continue
		}
		switch v := v.(type) {
		case *osmpbf.Node:
			osm.logger.Info("Node")
			osm.onNode(v)
		case *osmpbf.Way:
			osm.logger.Info("Way")
			osm.onWay(v)
		case *osmpbf.Relation:
			osm.logger.Info("Relation")
		default:
			osm.logger.Error("Unknown")

		}
	}
	return nil
}
func (osm *OSMWorker) onNode(node *osmpbf.Node) {
	osm.levelDB.CacheQueue(osm.batch, node)
	// TODO: Remove hardcoded value
	if osm.batch.Len() > 50000 {
		osm.levelDB.CacheFlush(osm.batch)
	}
	if !osm.hasTags(node.Tags) {
		return
	}
	node.Tags = osm.trimTags(node.Tags)
	if osm.containsValidTags(node.Tags, osm.tags) {
		// TODO: Process it
		_ = osm.toJSONNode(node)
	}

}
func (osm *OSMWorker) onWay(way *osmpbf.Way) {
	// TODO: Remove hardcoded value
	if osm.batch.Len() > 1 {
		osm.levelDB.CacheFlush(osm.batch)
	}
	if !osm.hasTags(way.Tags) {
		return
	}
	way.Tags = osm.trimTags(way.Tags)
	if osm.containsValidTags(way.Tags, osm.tags) {
		latlons, err := osm.levelDB.CacheLookup(way)
		if err != nil {
			return
		}
		var centroid = geo.ComputeCentroid(latlons)
		// TODO: Handle ways
		_ = osm.toJsonWay(way, latlons, centroid)
	}
}
func (osm *OSMWorker) toJSONNode(v *osmpbf.Node) models.JsonNode {
	return models.JsonNode{}
}
func (osm *OSMWorker) toJsonWay(v *osmpbf.Way, latlons []map[string]string, centroid map[string]string) models.JsonWay {
	var points []*geo.Point
	for _, latlon := range latlons {
		var lat, _ = strconv.ParseFloat(latlon["lat"], 64)
		var lng, _ = strconv.ParseFloat(latlon["lon"], 64)
		points = append(points, geo.NewPoint(lat, lng))
	}
	return models.JsonWay{
		ID:       v.ID,
		Type:     "way",
		Tags:     v.Tags,
		Centroid: centroid,
		Nodes:    points,
	}
}
func (osm *OSMWorker) containsValidTags(tags map[string]string, group map[string][]string) bool {
	for _, list := range group {
		if osm.matchTagsAgainstCompulsoryTagList(tags, list) {
			return true
		}
	}
	return false
}
func (osm *OSMWorker) matchTagsAgainstCompulsoryTagList(tags map[string]string, tagList []string) bool {
	for _, name := range tagList {

		feature := strings.Split(name, "~")
		foundVal, foundKey := tags[feature[0]]

		// key check
		if !foundKey {
			return false
		}

		// value check
		if len(feature) > 1 {
			if foundVal != feature[1] {
				return false
			}
		}
	}

	return true
}

func (osm *OSMWorker) hasTags(tags map[string]string) bool {
	n := len(tags)
	if n == 0 {
		return false
	}
	return true
}
func (osm *OSMWorker) trimTags(tags map[string]string) map[string]string {
	trimmed := make(map[string]string)
	for k, v := range tags {
		trimmed[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return trimmed
}
func (osm *OSMWorker) DownloadFile() error {
	osm.logger.Info("Downloading file")
	fmt.Println("Downloading file")
	err := osm.downloadFile(osm.appConfig.DownloadUrl, osm.appConfig.FileName)
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

	_, err = io.Copy(output, response.Body)

	if err != nil {
		return errors.Wrap(err, "Error while copying data")
	}
	return nil

}
