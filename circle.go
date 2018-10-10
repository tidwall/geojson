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
func NewCircle(center geometry.Point, meters float64, steps int) *Circle {
	g := new(Circle)
	g.center = center
	g.meters = meters
	g.steps = steps
	if steps < 3 {
		steps = 3
	}
	var points []geometry.Point
	step := 360.0 / float64(steps)
	i := 0
	for deg := 360.0; deg > 0; deg -= step {
		lat, lon := geo.DestinationPoint(center.Y, center.X, meters, deg)
		points = append(points, geometry.Point{X: lon, Y: lat})
		i++
	}
	// TODO: account for the pole and antimerdian. In most cases only a polygon
	// is needed, but when the circle bounds passes the 90/180 lines, we need
	// to create a multipolygon
	points = append(points, points[0])
	poly := geometry.NewPoly(points, nil, 0)
	gPoly := new(Polygon)
	gPoly.base = *poly
	g.Object = gPoly
	return g
}

// var circleCache = func() *lru.Cache {
// 	l, _ := lru.New(512)
// 	return l
// }()

// type circleCacheKey struct {
// 	center geometry.Point
// 	meters float64
// 	steps  int
// }

// // loadCircle will do an lru look up on cached circles.
// func loadCircle(center geometry.Point, meters float64, steps int) *Circle {
// 	key := circleCacheKey{center, meters, steps}
// 	value, ok := circleCache.Get(key)
// 	if ok {
// 		return value.(*Circle)
// 	}
// 	circle := NewCircle(center, meters, steps)
// 	circleCache.Add(key, circle)
// 	return circle
// }
