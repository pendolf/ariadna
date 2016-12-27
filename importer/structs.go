package importer

import (
	"github.com/dhconnelly/rtreego"
	"github.com/kellydunn/golang-geo"
	"math"
)

func (way *JsonWay) Bounds() *rtreego.Rect {
	return way.Rect
}

func (way *JsonWay) GetXY() (x, y, z, j float64) {
	var maxlat, minlat, maxlon, minlon float64
	minlat = float64(99999999999)
	minlon = float64(99999999999)
	for _, point := range way.Nodes {
		x, y := getXY(point.Lat(), point.Lng())
		maxlat = math.Max(maxlat, x)
		minlat = math.Min(minlat, x)
		maxlon = math.Max(maxlon, y)
		minlon = math.Min(minlon, y)
	}
	return maxlat, minlat, maxlon, minlon
}
func getXY(lat, lng float64) (float64, float64) {
	LAT := (lat * math.Pi) / 180
	LON := (lng * math.Pi) / 180
	X := 6371 * math.Sin(LAT) * math.Sin(LON)
	Y := 6371 * math.Cos(LAT)
	return X, Y
}

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
