package importer

import (
	"github.com/dhconnelly/rtreego"
	"github.com/kellydunn/golang-geo"
	"math"
)

type Tags struct {
	housenumber string
	street      string
}

type PGNode struct {
	ID      int64
	Name    string
	OldName string
	Lng     float64
	Lat     float64
}

type Translate struct {
	Original  string
	Translate string
}
