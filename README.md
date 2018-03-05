Ariadna consists of 3 parts:
* Importer: OSM data importer to elastic search
* Updater: Download and import data
* WebUI for searching data
* Custom importer

### Configure

```
COMMANDS:
    import              Import OSM file to ElasticSearch
    update              Download OSM file and update index
    http                Run http server
    custom              Import custom data
    intersections       Process intersections only

GLOBAL OPTIONS:
   --config                                                                             Config file path
   --index_settings                                                                     ElasticSearch Index settings
   --custom_data                                                                        Custom data file path
   --es_index_name "addresses"                                                          Specify custom elasticsearch index name [$ARIADNA_ES_INDEX_NAME]
   --es_pg_conn_url "host=localhost user=geo password=geo dbname=geo sslmode=disable"   Specify custom PG connection URL [$ARIADNA_PG_CONN_URL]
   --es_url "http://localhost:9200"                                                     Custom url for elasticsearch e.g http://192.168.0.1:9200 [$ARIADNA_ES_HOST]
   --es_index_type "address"                                                            ElasticSearch index type [$ARIADNA_INDEX_TYPE]
   --filename "xxx"                                                                     filename for storing osm.pbf file [$ARIADNA_FILE_NAME]
   --download_url "xxx"                                                                 Geofabrik url to download file [$ARIADNA_DOWNLOAD_URL]
   --dont_import_intersections                                                          if checked, then ariadna won't import intersections [$ARIADNA_DONT_IMPORT_INTERSECTIONS]
   --help, -h                                                                           show help
   --version, -v                                                                        print the version
```

### Usage
First import data. Download it from geofabrik.de and run
```
$ ./ariadna import
```
Or you can specify download_url and file_name into settings and run
```
$ ./ariadna update
```
This creates elasticsearch index

### WebUI
```
$ ./ariadna http
```
Open http://localhost:8080 in your browser and enjoy

### Http API
There is http api for geocode and reverse geocode

1. /api/search/:query
2. /api/reverse/:lat/:lon

### Docker
To start Postgres, Elasticsearch and Ariadna run
```
# edit index files if need
$ docker-compose up -d
$ docker-compose run --rm ariadna-export-suggest
$ docker-compose run --rm ariadna-export-address
```

