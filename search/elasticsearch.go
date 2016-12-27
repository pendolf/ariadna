package search

import (
	"fmt"
	"github.com/dhconnelly/rtreego"
	"github.com/gen1us2k/go-translit"
	"gopkg.in/olivere/elastic.v3"
	"strings"
)

type Search interface {
	GetCurrentIndexName() (string, error)
}

type ElasticSearch struct {
	Search
	client *elastic.Client
}

func NewElasticSearch() (*ElasticSearch, error) {
	// TODO: Remove hardcoded value
	elasticHost := ""
	e := &ElasticSearch{}
	client, err := elastic.NewClient(
		elastic.SetURL(elasticHost),
	)
	if err != nil {
		return nil, err
	}
	e.client = client
	return e, nil
}

func (es *ElasticSearch) GetCurrentIndexName() (string, error) {
	// TODO: remove hardcoded value
	indexName := ""
	res, err := es.client.Aliases().Index("_all").Do()
	if err != nil {
		return "", err
	}
	for _, index := range res.IndicesByAlias(indexName) {
		if strings.HasPrefix(index, indexName) {
			return index, nil
		}
	}
	return "", nil
}

func JsonWaysToES(Addresses []JsonWay, CitiesAndTowns []JsonWay, client *elastic.Client) {
	Logger.Info("Populating elastic search index")
	bulkClient := client.Bulk()
	Logger.Info("Creating bulk client")
	for _, address := range Addresses {
		cityName, villageName, suburbName, townName := "", "", "", ""
		var lat, _ = strconv.ParseFloat(address.Centroid["lat"], 64)
		var lng, _ = strconv.ParseFloat(address.Centroid["lon"], 64)
		for _, city := range CitiesAndTowns {
			polygon := geo.NewPolygon(city.Nodes)

			if polygon.Contains(geo.NewPoint(lat, lng)) {
				switch city.Tags["place"] {
				case "city":
					cityName = city.Tags["name"]
				case "village":
					villageName = city.Tags["name"]
				case "suburb":
					suburbName = city.Tags["name"]
				case "town":
					townName = city.Tags["name"]
				case "neighbourhood":
					suburbName = city.Tags["name"]
				}
			}
		}
		var points [][][]float64
		for _, point := range address.Nodes {
			points = append(points, [][]float64{[]float64{point.Lat(), point.Lng()}})
		}

		pg := gj.NewPolygonFeature(points)
		centroid := make(map[string]float64)
		centroid["lat"] = lat
		centroid["lon"] = lng
		name := cleanAddress(address.Tags["name"])
		translated := ""

		if latinre.Match([]byte(name)) {
			word := make(map[string]string)
			word["original"] = name

			trans := strings.Split(name, " ")
			for _, k := range trans {
				s := synonims[k]
				if s == "" {
					s = translit.Translit(k)
				}
				translated += fmt.Sprintf("%s ", s)
			}

			word["trans"] = translated
		}
		housenumber := translit.Translit(address.Tags["addr:housenumber"])
		marshall := JsonEsIndex{
			Country:           "KG",
			City:              cityName,
			Village:           villageName,
			Town:              townName,
			District:          suburbName,
			Street:            cleanAddress(address.Tags["addr:street"]),
			HouseNumber:       housenumber,
			Name:              name,
			OldName:           address.Tags["old_name"],
			HouseName:         address.Tags["housename"],
			PostCode:          address.Tags["postcode"],
			LocalName:         address.Tags["loc_name"],
			AlternativeName:   address.Tags["alt_name"],
			InternationalName: address.Tags["int_name"],
			NationalName:      address.Tags["nat_name"],
			OfficialName:      address.Tags["official_name"],
			RegionalName:      address.Tags["reg_name"],
			ShortName:         address.Tags["short_name"],
			SortingName:       address.Tags["sorting_name"],
			TranslatedName:    translated,
			Centroid:          centroid,
			Geom:              pg,
			Custom:            false,
		}
		index := elastic.NewBulkIndexRequest().
			Index(common.AC.ElasticSearchIndexUrl).
			Type(common.AC.IndexType).
			Id(strconv.FormatInt(address.ID, 10)).
			Doc(marshall)
		bulkClient = bulkClient.Add(index)
	}
	Logger.Info("Starting to insert many data to elasticsearch")
	_, err := bulkClient.Do()
	Logger.Info("Data insert")
	if err != nil {
		Logger.Error(err.Error())
	}
}

func JsonNodesToEs(Addresses []JsonNode, CitiesAndTowns []JsonWay, client *elastic.Client) {
	Logger.Info("Populating elastic search index with Nodes")
	bulkClient := client.Bulk()
	Logger.Info("Created bulk request to elasticsearch")
	for _, address := range Addresses {
		cityName, villageName, suburbName, townName := "", "", "", ""
		for _, city := range CitiesAndTowns {
			polygon := geo.NewPolygon(city.Nodes)

			if polygon.Contains(geo.NewPoint(address.Lat, address.Lon)) {
				switch city.Tags["place"] {
				case "city":
					cityName = city.Tags["name"]
				case "village":
					villageName = city.Tags["name"]
				case "suburb":
					suburbName = city.Tags["name"]
				case "town":
					townName = city.Tags["name"]
				case "neighbourhood":
					suburbName = city.Tags["name"]
				}
			}
		}

		centroid := make(map[string]float64)
		centroid["lat"] = address.Lat
		centroid["lon"] = address.Lon
		name := cleanAddress(address.Tags["name"])
		translated := ""
		if latinre.Match([]byte(name)) {
			word := make(map[string]string)
			word["original"] = name

			trans := strings.Split(name, " ")
			for _, k := range trans {
				s := synonims[k]
				if s == "" {
					s = translit.Translit(k)
				}
				translated += fmt.Sprintf("%s ", s)
			}

			word["trans"] = translated
		}
		housenumber := translit.Translit(address.Tags["addr:housenumber"])

		marshall := JsonEsIndex{
			Country:           "KG",
			City:              cityName,
			Village:           villageName,
			Town:              townName,
			District:          suburbName,
			Street:            cleanAddress(address.Tags["addr:street"]),
			HouseNumber:       housenumber,
			Name:              name,
			TranslatedName:    translated,
			OldName:           address.Tags["old_name"],
			HouseName:         address.Tags["housename"],
			PostCode:          address.Tags["postcode"],
			LocalName:         address.Tags["loc_name"],
			AlternativeName:   address.Tags["alt_name"],
			InternationalName: address.Tags["int_name"],
			NationalName:      address.Tags["nat_name"],
			OfficialName:      address.Tags["official_name"],
			RegionalName:      address.Tags["reg_name"],
			ShortName:         address.Tags["short_name"],
			SortingName:       address.Tags["sorting_name"],
			Centroid:          centroid,
			Geom:              nil,
			Custom:            false,
			Intersection:      address.Intersection,
		}

		index := elastic.NewBulkIndexRequest().
			Index(common.AC.ElasticSearchIndexUrl).
			Type(common.AC.IndexType).
			Id(strconv.FormatInt(address.ID, 10)).
			Doc(marshall)
		bulkClient = bulkClient.Add(index)
	}
	Logger.Info("Started to bulk insert to elasticsearch")
	_, err := bulkClient.Do()
	Logger.Info("Data inserted")
	if err != nil {
		Logger.Error(err.Error())
	}

}

type JsonEsIndex struct {
	Country           string             `json:"country"`
	City              string             `json:"city"`
	Village           string             `json:"village"`
	Town              string             `json:"town"`
	District          string             `json:"district"`
	Street            string             `json:"street"`
	HouseNumber       string             `json:"housenumber"`
	Name              string             `json:"name"`
	OldName           string             `json:"old_name"`
	HouseName         string             `json:"housename"`
	PostCode          string             `json:"postcode"`
	LocalName         string             `json:"local_name"`
	AlternativeName   string             `json:"alternative_name"`
	InternationalName string             `json:"international"`
	NationalName      string             `json:"national"`
	OfficialName      string             `json:"official"`
	RegionalName      string             `json:"regional"`
	ShortName         string             `json:"short_name"`
	SortingName       string             `json:"sorting"`
	TranslatedName    string             `json:"translated"`
	Custom            bool               `json:"custom"`
	Intersection      bool               `json:"intersection"`
	Centroid          map[string]float64 `json:"centroid"`
	Geom              interface{}        `json:"geom"`
}

type JsonNode struct {
	ID           int64             `json:"id"`
	Type         string            `json:"type"`
	Lat          float64           `json:"lat"`
	Lon          float64           `json:"lon"`
	Tags         map[string]string `json:"tags"`
	Intersection bool              `json:"-"`
}

type JsonRelation struct {
	ID       int64               `json:"id"`
	Type     string              `json:"type"`
	Tags     map[string]string   `json:"tags"`
	Centroid map[string]string   `json:"centroid"`
	Nodes    []map[string]string `json:"nodes"`
}

type JsonWay struct {
	ID       int64             `json:"id"`
	Type     string            `json:"type"`
	Tags     map[string]string `json:"tags"`
	Centroid map[string]string `json:"centroid"`
	Nodes    []*geo.Point      `json:"nodes"`
	Rect     *rtreego.Rect     `json:"-"`
}
