package geojson

import (
	"github.com/tidwall/geojson/geo"
	"github.com/tidwall/geojson/geometry"
)

// Circle ...
type Circle struct {
	Object
	center geometry.Point
	meters float64
	steps  int
	extra  *extra
}

// NewCircle returns an circle object
func NewCircle(x, y, meters float64, steps int) *Circle {
	if steps < 3 {
		steps = 3
	}
	g := new(Circle)
	g.center.X = x
	g.center.Y = y
	g.meters = meters
	g.steps = steps
	if meters == 0 {
		g.Object = NewPoint(geometry.Point{X: x, Y: y})
	} else {
		var points []geometry.Point
		step := 360.0 / float64(steps)
		i := 0
		for deg := 360.0; deg > 0; deg -= step {
			lat, lon := geo.DestinationPoint(y, x, meters, deg)
			points = append(points, geometry.Point{X: lon, Y: lat})
			i++
		}
		// TODO: account for the pole and antimerdian. In most cases only a polygon
		// is needed, but when the circle bounds passes the 90/180 lines, we need
		// to create a multipolygon
		points = append(points, points[0])
		g.Object = NewPolygon(
			geometry.NewPoly(points, nil, geometry.DefaultIndex),
		)
	}
	return g
}
