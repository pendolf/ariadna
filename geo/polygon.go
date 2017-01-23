package geo

import "math"

type Polygon struct {
	points []*Point
}

func NewPolygon(points []*Point) *Polygon {
	return &Polygon{
		points: points,
	}
}

// Returns whether or not the current Polygon contains the passed in Point.
func (p *Polygon) Contains(point *Point) bool {
	if !p.IsClosed() {
		return false
	}

	start := len(p.points) - 1
	end := 0

	contains := p.intersectsWithRaycast(point, p.points[start], p.points[end])

	for i := 1; i < len(p.points); i++ {
		if p.intersectsWithRaycast(point, p.points[i-1], p.points[i]) {
			contains = !contains
		}
	}

	return contains
}

// Returns whether or not the polygon is closed.
// TODO:  This can obviously be improved, but for now,
//        this should be sufficient for detecting if points
//        are contained using the raycast algorithm.
func (p *Polygon) IsClosed() bool {
	if len(p.points) < 3 {
		return false
	}

	return true
}

// Using the raycast algorithm, this returns whether or not the passed in point
// Intersects with the edge drawn by the passed in start and end points.
// Original implementation: http://rosettacode.org/wiki/Ray-casting_algorithm#Go
func (p *Polygon) intersectsWithRaycast(point *Point, start *Point, end *Point) bool {
	// Always ensure that the the first point
	// has a y coordinate that is less than the second point
	if start.Lat() > end.Lat() {

		// Switch the points if otherwise.
		start, end = end, start

	}

	// Move the point's y coordinate
	// outside of the bounds of the testing region
	// so we can start drawing a ray
	for point.Lat() == start.Lat() || point.Lat() == end.Lat() {
		newLng := math.Nextafter(point.Lat(), math.Inf(1))
		point = NewPoint(point.Lon(), newLng)
	}

	// If we are outside of the polygon, indicate so.
	if point.Lat() < start.Lat() || point.Lat() > end.Lat() {
		return false
	}

	if start.Lon() > end.Lon() {
		if point.Lon() > start.Lon() {
			return false
		}
		if point.Lon() < end.Lon() {
			return true
		}

	} else {
		if point.Lon() > end.Lon() {
			return false
		}
		if point.Lon() < start.Lon() {
			return true
		}
	}

	raySlope := (point.Lat() - start.Lat()) / (point.Lon() - start.Lon())
	diagSlope := (end.Lat() - start.Lat()) / (end.Lon() - start.Lon())

	return raySlope >= diagSlope
}
