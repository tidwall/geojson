package geom

import (
	"github.com/tidwall/geojson/geo"
)

// CircleRing ...
func CircleRing(center Point, meters float64, segments int) *Ring {
	if segments < 3 {
		segments = 3
	}
	points := make([]Point, segments+1)
	step := 360.0 / float64(segments)
	i := 0
	for deg := 360.0; deg > 0; deg -= step {
		lat, lon := geo.DestinationPoint(center.Y, center.X, meters, deg)
		points[i] = Point{X: lon, Y: lat}
		i++
	}
	points[i] = points[0]
	return NewRing(points)
}
