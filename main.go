package main

import (
	"github.com/gen1us2k/log"
	"github.com/maddevsio/ariadna/common"
	"github.com/maddevsio/ariadna/config"
	"github.com/maddevsio/ariadna/models"
	"github.com/maddevsio/ariadna/osm"
	"gopkg.in/urfave/cli.v1"
)

var (
	CitiesAndTowns, Roads []models.JsonWay
	configPath            string
	indexSettingsPath     string
	customDataPath        string
)

func main() {
	app := config.New()
	app.App().Action = func(ctx *cli.Context) error {
		o, err := osm.New(app.Get())
		if err != nil {
			return err
		}
		err = o.DownloadFile()
		log.Info("Starting to download file")
		if err != nil {
			return err
		}
		tags := common.BuildTags("place~city,place~village,place~suburb,place~town,place~neighbourhood")
		err = o.Prepare(tags)

		//	err := updater.DownloadOSMFile(common.AC.DownloadUrl, common.AC.FileName)
		//	if err != nil {
		//		importer.Logger.Fatal(err.Error())
		//	}
		//	actionImport(ctx)
		//	return nil

		return err
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

	//
	//app.App().Commands = []cli.Command{
	//	{
	//		Name:      "import",
	//		Aliases:   []string{"i"},
	//		Usage:     "Import OSM file to ElasticSearch",
	//		Action:    actionImport,
	//		ArgsUsage: "<filename>",
	//	},
	//	{
	//		Name:    "update",
	//		Aliases: []string{"u"},
	//		Usage:   "Download OSM file and update index",
	//		Action:  actionUpdate,
	//	},
	//	{
	//		Name:    "http",
	//		Aliases: []string{"h"},
	//		Usage:   "Run http server",
	//		Action:  actionHttp,
	//	},
	//	{
	//		Name:    "custom",
	//		Aliases: []string{"c"},
	//		Usage:   "Import custom data",
	//		Action:  actionCustom,
	//	},
	//	{
	//		Name:    "intersections",
	//		Aliases: []string{"i"},
	//		Usage:   "Process intersections only",
	//		Action:  actionIntersection,
	//	},
	//	{
	//		Name:    "test",
	//		Aliases: []string{"t"},
	//		Usage:   "For testing",
	//		Action:  actionTest,
	//	},
	//}

	//app.App().Before = func(context *cli.Context) error {
	//	if configPath == "" {
	//		configPath = "config.json"
	//	}
	//	if indexSettingsPath == "" {
	//		indexSettingsPath = "index.json"
	//	}
	//	if customDataPath == "" {
	//		customDataPath = "custom.json"
	//	}
	//	return nil
	//}
	//
	//if err := app.Run(os.Args); err != nil {
	//	importer.Logger.Fatal("error on run app, %v", err)
	//}

}

//func actionImport(ctx *cli.Context) error {
//
//	indexSettings, err := ioutil.ReadFile(indexSettingsPath)
//	if err != nil {
//		importer.Logger.Fatal(err.Error())
//	}
//	db := importer.OpenLevelDB("db")
//	defer db.Close()
//
//	file := importer.OpenFile(common.AC.FileName)
//	defer file.Close()
//	decoder := getDecoder(file)
//
//	client, err :=
//
//	if err != nil {
//		importer.Logger.Fatal("Failed to create Elastic search client with error %s", err)
//	}
//
//	indexVersion := fmt.Sprintf("%s_%d", common.AC.IndexName, time.Now().Unix())
//	importer.Logger.Info("Creating index with name %s", indexVersion)
//	_, err = client.CreateIndex(indexVersion).BodyString(string(indexSettings)).Do()
//
//	common.AC.ElasticSearchIndexUrl = indexVersion
//
//	if err != nil {
//		importer.Logger.Error(err.Error())
//	}
//
//	importer.Logger.Info("Searching cities, villages, towns and districts")
//	tags := importer.BuildTags("place~city,place~village,place~suburb,place~town,place~neighbourhood")
//	CitiesAndTowns, _ = importer.Run(decoder, db, tags)
//
//	importer.Logger.Info("Cities, villages, towns and districts found")
//
//	file = importer.OpenFile(common.AC.FileName)
//	defer file.Close()
//	decoder = getDecoder(file)
//
//	importer.Logger.Info("Searching addresses")
//	tags = importer.BuildTags("addr:street+addr:housenumber,amenity+name,building+name,addr:housenumber,shop+name,office+name,public_transport+name,cuisine+name,railway+name,sport+name,natural+name,tourism+name,leisure+name,historic+name,man_made+name,landuse+name,waterway+name,aerialway+name,aeroway+name,craft+name,military+name")
//	AddressWays, AddressNodes := importer.Run(decoder, db, tags)
//	importer.Logger.Info("Addresses found")
//	importer.JsonWaysToES(AddressWays, CitiesAndTowns, client)
//	importer.JsonNodesToEs(AddressNodes, CitiesAndTowns, client)
//	file = importer.OpenFile(common.AC.FileName)
//	defer file.Close()
//	decoder = getDecoder(file)
//
//	if !common.AC.DontImportIntersections {
//		tags = importer.BuildTags("highway+name")
//		Roads, _ = importer.Run(decoder, db, tags)
//		importer.RoadsToPg(Roads)
//		importer.Logger.Info("Searching all roads intersecitons")
//		Intersections := importer.GetRoadIntersectionsFromPG()
//		importer.JsonNodesToEs(Intersections, CitiesAndTowns, client)
//	}
//
//	importer.Logger.Info("Removing indices from alias")
//	_, err = client.Alias().Add(indexVersion, common.AC.IndexName).Do()
//	if err != nil {
//		return err
//	}
//	res, err := client.Aliases().Index("_all").Do()
//	if err != nil {
//		return err
//	}
//	for _, index := range res.IndicesByAlias(common.AC.IndexName) {
//		if strings.HasPrefix(index, common.AC.IndexName) && index != indexVersion {
//			_, err = client.Alias().Remove(index, common.AC.IndexName).Do()
//			if err != nil {
//				importer.Logger.Error("Failed to delete index alias: %s", err.Error())
//			}
//			_, err = client.DeleteIndex(index).Do()
//			if err != nil {
//				importer.Logger.Error("Failed to delete index: %s", err.Error())
//			}
//		}
//	}
//	return nil
//}
//
//func actionUpdate(ctx *cli.Context) error {
//	err := updater.DownloadOSMFile(common.AC.DownloadUrl, common.AC.FileName)
//	if err != nil {
//		importer.Logger.Fatal(err.Error())
//	}
//	actionImport(ctx)
//	return nil
//}
//
//func actionTest(ctx *cli.Context) error {
//	client, err := elastic.NewClient(
//		elastic.SetURL(common.AC.ElasticSearchHost),
//	)
//	if err != nil {
//		return err
//	}
//	currentIndex, err := getCurrentIndexName(client)
//
//	if err != nil {
//		return err
//	}
//	fmt.Println(currentIndex)
//	return nil
//}
//func actionHttp(ctx *cli.Context) error {
//	err := web.StartServer()
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//type CustomData struct {
//	ID   int64   `json:"id"`
//	Name string  `json:"name"`
//	Lat  float64 `json:"lat"`
//	Lon  float64 `json:"lon"`
//}
//type Custom []CustomData
//
//func actionCustom(ctx *cli.Context) error {
//	var custom Custom
//	data, err := ioutil.ReadFile(customDataPath)
//	if err != nil {
//		importer.Logger.Fatal(err.Error())
//	}
//	err = json.Unmarshal(data, &custom)
//	if err != nil {
//		importer.Logger.Fatal(err.Error())
//	}
//
//	client, err := elastic.NewClient(
//		elastic.SetURL(common.AC.ElasticSearchHost),
//	)
//	bulkClient := client.Bulk()
//	indexVersion, err := getCurrentIndexName(client)
//	if err != nil {
//		return err
//	}
//	for _, item := range custom {
//		centroid := make(map[string]float64)
//		centroid["lat"] = item.Lat
//		centroid["lon"] = item.Lon
//		marshall := importer.JsonEsIndex{
//			Name:     item.Name,
//			Centroid: centroid,
//			Custom:   true,
//		}
//		index := elastic.NewBulkIndexRequest().
//			Index(indexVersion).
//			Type(common.AC.IndexType).
//			Id(strconv.FormatInt(item.ID, 10)).
//			Doc(marshall)
//		bulkClient = bulkClient.Add(index)
//	}
//	_, err = bulkClient.Do()
//	if err != nil {
//		importer.Logger.Error(err.Error())
//	}
//	return nil
//}
//
//func actionIntersection(ctx *cli.Context) error {
//	db := importer.OpenLevelDB("db")
//	defer db.Close()
//
//	file := importer.OpenFile(common.AC.FileName)
//	defer file.Close()
//	decoder := getDecoder(file)
//
//	client, err := elastic.NewClient(
//		elastic.SetURL(common.AC.ElasticSearchHost),
//	)
//
//	if err != nil {
//		importer.Logger.Fatal("Failed to create Elastic search client with error %s", err)
//	}
//
//	indexVersion, err := getCurrentIndexName(client)
//
//	common.AC.ElasticSearchIndexUrl = indexVersion
//	importer.Logger.Info("Creating index with name %s", common.AC.ElasticSearchIndexUrl)
//
//	importer.Logger.Info("Searching cities, villages, towns and districts")
//	tags := importer.BuildTags("place~city,place~village,place~suburb,place~town,place~neighbourhood")
//	CitiesAndTowns, _ = importer.Run(decoder, db, tags)
//
//	importer.Logger.Info("Cities, villages, towns and districts found")
//
//	file = importer.OpenFile(common.AC.FileName)
//	defer file.Close()
//	decoder = getDecoder(file)
//
//	tags = importer.BuildTags("highway+name")
//	Roads, _ = importer.Run(decoder, db, tags)
//	importer.BuildIndex(Roads)
//	Intersections := importer.SearchIntersections(Roads)
//	importer.JsonNodesToEs(Intersections, CitiesAndTowns, client)
//	return nil
//}
