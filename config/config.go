package config

import (
	"github.com/gen1us2k/log"
	"gopkg.in/urfave/cli.v1"
	"os"
)

type AriadnaConfig struct {
	PGConnString            string
	ElasticSearchHost       string
	IndexName               string
	IndexType               string
	FileName                string
	DownloadUrl             string
	ElasticSearchIndexUrl   string
	DontImportIntersections bool
	LevelDBPath             string
}

var (
	Version                 string = "dev"
	LevelDBPath             string
	FileName                string
	configPath              string
	indexSettingsPath       string
	customDataPath          string
	ElasticSearchIndexName  string
	PGConnString            string
	ElasticSearchHost       string
	IndexType               string
	DownloadUrl             string
	LogLevel                string
	DontImportIntersections bool
)

type Configuration struct {
	data *AriadnaConfig
	app  *cli.App
}

// NewConfigurator is constructor and creates a new copy of Configuration
func New() *Configuration {
	Version = "0.1dev"
	app := cli.NewApp()
	app.Name = "Ariadna"
	app.Usage = "OSM Geocoder"
	return &Configuration{
		data: &AriadnaConfig{},
		app:  app,
	}
}

func (c *Configuration) fillConfig() *AriadnaConfig {
	return &AriadnaConfig{
		IndexType:               IndexType,
		PGConnString:            PGConnString,
		ElasticSearchHost:       ElasticSearchHost,
		IndexName:               ElasticSearchIndexName,
		FileName:                FileName,
		DownloadUrl:             DownloadUrl,
		DontImportIntersections: DontImportIntersections,
		LevelDBPath:             LevelDBPath,
	}
}

// Run is wrapper around cli.App
func (c *Configuration) Run() error {
	c.app.Before = func(ctx *cli.Context) error {
		log.SetLevel(log.MustParseLevel(LogLevel))
		return nil
	}
	c.app.Flags = c.setupFlags()
	return c.app.Run(os.Args)
}

// App is public method for Configuration.app
func (c *Configuration) App() *cli.App {
	return c.app
}

func (c *Configuration) setupFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "config",
			Usage:       "Config file path",
			Destination: &configPath,
		},
		cli.StringFlag{
			Name:        "index_settings",
			Usage:       "ElasticSearch Index settings",
			Destination: &indexSettingsPath,
		},
		cli.StringFlag{
			Name:        "custom_data",
			Usage:       "Custom data file path",
			Destination: &customDataPath,
		},
		cli.StringFlag{
			Name:        "leveldb",
			Usage:       "Leveldb database path",
			Value:       "db",
			EnvVar:      "LEVEL_DB_PATH",
			Destination: &LevelDBPath,
		},
		cli.StringFlag{
			Name:        "es_index_name",
			Usage:       "Specify custom elasticsearch index name",
			Value:       "addresses",
			EnvVar:      "ES_INDEX_NAME",
			Destination: &ElasticSearchIndexName,
		},
		cli.StringFlag{
			Name:        "es_pg_conn_url",
			Usage:       "Specify custom PG connection URL",
			Destination: &PGConnString,
			Value:       "host=localhost user=geo password=geo dbname=geo sslmode=disable",
			EnvVar:      "PG_CONN_URL",
		},
		cli.StringFlag{
			Name:        "es_url",
			Usage:       "Custom url for elasticsearch e.g http://192.168.0.1:9200",
			Destination: &ElasticSearchHost,
			Value:       "http://localhost:9200/",
			EnvVar:      "ELASTICSEARCH_HOST",
		},
		cli.StringFlag{
			Name:        "es_index_type",
			Usage:       "ElasticSearch index type",
			Destination: &IndexType,
			Value:       "address",
			EnvVar:      "INDEX_TYPE",
		},
		cli.StringFlag{
			Name:        "filename",
			Usage:       "filename for storing osm.pbf file",
			Destination: &FileName,
			Value:       "xxx",
			EnvVar:      "FILENAME",
		},
		cli.StringFlag{
			Name:        "download_url",
			Usage:       "Geofabrik url to download file",
			Destination: &DownloadUrl,
			Value:       "xxx",
			EnvVar:      "DOWNLOAD_URL",
		},
		cli.StringFlag{
			Name:        "log_level",
			Usage:       "Set log level",
			Destination: &LogLevel,
			Value:       "debug",
			EnvVar:      "LOG_LEVEL",
		},
		cli.BoolFlag{
			Name:        "dont_import_intersections",
			Usage:       "if checked, then ariadna won't import intersections",
			Destination: &DontImportIntersections,
			EnvVar:      "DONT_IMPORT_INTERSECTIONS",
		},
	}

}

// Get returns filled AriadnaConfig
func (c *Configuration) Get() *AriadnaConfig {
	c.data = c.fillConfig()
	return c.data
}
