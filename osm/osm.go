package osm

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
)

type (
	OSM interface {
		DownloadFile(string, string) error
	}
	OSMWorker struct {
		OSM
	}
)

func New() *OSMWorker {
	return &OSMWorker{}
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
