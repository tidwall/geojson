package geojson

import (
	"math"
	"strconv"

	"github.com/tidwall/geojson/geo"
	"github.com/tidwall/geojson/geometry"
)

// Circle ...
type Circle struct {
	object    Object
	center    geometry.Point
	meters    float64
	haversine float64
	steps     int
	km        bool
	extra     *extra
}

// NewCircle returns an circle object
func NewCircle(center geometry.Point, meters float64, steps int) *Circle {
	if steps < 3 {
		steps = 3
	}
	g := new(Circle)
	g.center = center
	g.meters = meters
	g.steps = steps
	if meters > 0 {
		meters = geo.NormalizeDistance(meters)
		g.haversine = geo.DistanceToHaversine(meters)
	}
	return g
}

// AppendJSON ...
func (g *Circle) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Feature","geometry":`...)
	dst = append(dst, `{"type":"Point","coordinates":[`...)
	dst = strconv.AppendFloat(dst, g.center.X, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, g.center.Y, 'f', -1, 64)
	dst = append(dst, `]},"properties":{"type":"Circle","radius":`...)
	dst = strconv.AppendFloat(dst, g.meters, 'f', -1, 64)
	dst = append(dst, `,"radius_units":"m"}}`...)
	return dst
}

// JSON ...
func (g *Circle) JSON() string {
	return string(g.AppendJSON(nil))
}

// MarshalJSON ...
func (g *Circle) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

// String ...
func (g *Circle) String() string {
	return string(g.AppendJSON(nil))
}

// Meters returns the circle's radius
func (g *Circle) Meters() float64 {
	return g.meters
}

// Center returns the circle's center point
func (g *Circle) Center() geometry.Point {
	return g.center
}

// Haversine returns the haversine corresponding to circle's radius
func (g *Circle) Haversine() float64 {
	return g.haversine
}

// HaversineTo returns the haversine from a given point to circle's center
func (g *Circle) HaversineTo(p geometry.Point) float64 {
	return geo.Haversine(p.Y, p.X, g.center.Y, g.center.X)
}

// Within returns true if circle is contained inside object
func (g *Circle) Within(obj Object) bool {
	return obj.Contains(g)
}

// containsPoint returns true if circle contains a given point
func (g *Circle) containsPoint(p geometry.Point) bool {
	h := geo.Haversine(p.Y, p.X, g.center.Y, g.center.X)
	return h <= g.haversine
}

// Contains returns true if the circle contains other object
func (g *Circle) Contains(obj Object) bool {
	switch other := obj.(type) {
	case *Point:
		return g.containsPoint(other.Center())
	case *SimplePoint:
		return g.containsPoint(other.Center())
	case *Circle:
		return other.Distance(g) < (other.meters + g.meters)
	case Collection:
		for _, p := range other.Children() {
			if !g.Contains(p) {
				return false
			}
		}
		return true
	default:
		// No simple cases, so using polygon approximation.
		return g.getObject().Contains(other)
	}
}

// Intersects returns true the circle intersects other object
func (g *Circle) Intersects(obj Object) bool {
	switch other := obj.(type) {
	case *Point:
		return g.containsPoint(other.Center())
	case *Circle:
		return other.Distance(g) <= (other.meters + g.meters)
	case Collection:
		for _, p := range other.Children() {
			if g.Intersects(p) {
				return true
			}
		}
		return false
	case *Feature:
		return g.Intersects(other.base)
	default:
		// No simple cases, so using polygon approximation.
		return g.getObject().Intersects(obj)
	}
}

// Empty ...
func (g *Circle) Empty() bool {
	return false
}

// Valid ...
func (g *Circle) Valid() bool {
	return g.getObject().Valid()
}

// ForEach ...
func (g *Circle) ForEach(iter func(geom Object) bool) bool {
	return iter(g)
}

// NumPoints ...
func (g *Circle) NumPoints() int {
	// should this be g.steps?
	return 1
}

// Distance ...
func (g *Circle) Distance(other Object) float64 {
	return g.getObject().Distance(other)
}

// Rect ...
func (g *Circle) Rect() geometry.Rect {
	return g.getObject().Rect()
}

// Spatial ...
func (g *Circle) Spatial() Spatial {
	return g.getObject().Spatial()
}

// Primative returns a primative GeoJSON object. Either a Polygon or Point.
func (g *Circle) Primative() Object {
	return g.getObject()
}

func (g *Circle) getObject() Object {
	if g.object != nil {
		return g.object
	}
	return makeCircleObject(g.center, g.meters, g.steps)
}

func makeCircleObject(center geometry.Point, meters float64, steps int) Object {
	if meters <= 0 {
		return NewPoint(center)
	}
	meters = geo.NormalizeDistance(meters)
	points := make([]geometry.Point, 0, steps+1)

	// calc the four corners
	maxY, _ := geo.DestinationPoint(center.Y, center.X, meters, 0)
	_, maxX := geo.DestinationPoint(center.Y, center.X, meters, 90)
	minY, _ := geo.DestinationPoint(center.Y, center.X, meters, 180)
	_, minX := geo.DestinationPoint(center.Y, center.X, meters, 270)

	// TODO: detect of pole and antimeridian crossing and generate a
	// valid multigeometry

	// use the half width of the lat and lon
	lons := (maxX - minX) / 2
	lats := (maxY - minY) / 2

	// generate the
	for th := 0.0; th <= 360.0; th += 360.0 / float64(steps) {
		radians := (math.Pi / 180) * th
		x := center.X + lons*math.Cos(radians)
		y := center.Y + lats*math.Sin(radians)
		points = append(points, geometry.Point{X: x, Y: y})
	}
	// add last connecting point, make a total of steps+1
	points = append(points, points[0])

	return NewPolygon(
		geometry.NewPoly(points, nil, &geometry.IndexOptions{
			Kind: geometry.None,
		}),
	)
}
