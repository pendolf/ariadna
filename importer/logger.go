package importer

import (
	log "github.com/pendolf/ariadna/logger"
)

var Logger log.Logger

func init() {
	Logger = log.L("importer")
}
