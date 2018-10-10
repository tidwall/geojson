package geojson

import (
	"github.com/tidwall/geojson/geo"
	"github.com/tidwall/geojson/geometry"
)

// CircleFromCenter returns an circle object
func CircleFromCenter(x, y, meters float64, steps int) Object {
	// TODO: account for the pole and antimerdian
	if steps < 3 {
		steps = 3
	}
	var points []geometry.Point
	step := 360.0 / float64(steps)
	i := 0
	for deg := 360.0; deg > 0; deg -= step {
		lat, lon := geo.DestinationPoint(y, x, meters, deg)
		points = append(points, geometry.Point{X: lon, Y: lat})
		i++
	}
	points = append(points, points[0])
	poly := geometry.NewPoly(points, nil, 0)
	g := new(Polygon)
	g.base = *poly
	return g
}
