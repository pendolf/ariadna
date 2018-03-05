package common

import (
	log "github.com/pendolf/ariadna/logger"
)

var logger log.Logger

func init() {
	logger = log.L("common")
}
