package web

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/olivere/elastic"
	"github.com/pendolf/ariadna/common"
)

var es *elastic.Client

func StartServer() error {
	var err error
	es, err = elastic.NewClient(
		elastic.SetURL(common.AC.ElasticSearchHost),
		elastic.SetBasicAuth(common.AC.ESUsername, common.AC.ESPassword),
	)
	if err != nil {
		return err
	}
	router := httprouter.New()
	router.GET("/api/search/:query", geoCoder)
	router.GET("/api/reverse/:lat/:lon", reverseGeoCode)
	// router.NotFound = http.FileServer(http.Dir("public"))
	http.ListenAndServe(":8080", router)
	return nil
}
