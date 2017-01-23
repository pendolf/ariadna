package geo

type Point struct {
	X float64
	Y float64
}

func (p *Point) Lat() float64 {
	return p.Y
}
func (p *Point) Lon() float64 {
	return p.X
}

// Returns a new Point populated by the passed in latitude (lat) and longitude (lng) values.
func NewPoint(lat float64, lon float64) *Point {
	return &Point{X: lon, Y: lat}
}
