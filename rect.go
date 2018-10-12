package geojson

import (
	"github.com/tidwall/geojson/geometry"
)

// Rect ...
type Rect struct {
	base geometry.Rect
}

// NewRect ...
func NewRect(minX, minY, maxX, maxY float64) *Rect {
	return &Rect{base: geometry.Rect{
		Min: geometry.Point{X: minX, Y: minY},
		Max: geometry.Point{X: maxX, Y: maxY},
	}}
}

// forEach ...
func (g *Rect) forEach(iter func(geom Object) bool) bool {
	return iter(g)
}

// Empty ...
func (g *Rect) Empty() bool {
	return g.base.Empty()
}

// Rect ...
func (g *Rect) Rect() geometry.Rect {
	return g.base
}

// Center ...
func (g *Rect) Center() geometry.Point {
	return g.base.Center()
}

// AppendJSON ...
func (g *Rect) AppendJSON(dst []byte) []byte {
	var gPoly Polygon
	gPoly.base.Exterior = g.base
	return gPoly.AppendJSON(dst)
}

// JSON ...
func (g *Rect) JSON() string {
	return string(g.AppendJSON(nil))
}

// String ...
func (g *Rect) String() string {
	return string(g.AppendJSON(nil))
}

// IsSpatial ...
func (g *Rect) IsSpatial() bool {
	return true
}

// Within ...
func (g *Rect) Within(obj Object) bool {
	return obj.Contains(g)
}

// Contains ...
func (g *Rect) Contains(obj Object) bool {
	return obj.withinRect(g.base)
}

// Intersects ...
func (g *Rect) Intersects(obj Object) bool {
	return obj.intersectsRect(g.base)
}

func (g *Rect) withinRect(rect geometry.Rect) bool {
	return rect.ContainsRect(g.base)
}

func (g *Rect) withinPoint(point geometry.Point) bool {
	return point.ContainsRect(g.base)
}

func (g *Rect) withinLine(line *geometry.Line) bool {
	return line.ContainsRect(g.base)
}

func (g *Rect) withinPoly(poly *geometry.Poly) bool {
	return poly.ContainsRect(g.base)
}

func (g *Rect) intersectsPoint(point geometry.Point) bool {
	return g.base.IntersectsPoint(point)
}

func (g *Rect) intersectsRect(rect geometry.Rect) bool {
	return g.base.IntersectsRect(rect)
}

func (g *Rect) intersectsLine(line *geometry.Line) bool {
	return g.base.IntersectsLine(line)
}

func (g *Rect) intersectsPoly(poly *geometry.Poly) bool {
	return g.base.IntersectsPoly(poly)
}

// NumPoints ...
func (g *Rect) NumPoints() int {
	return 2
}

// // Clipped ...
// func (g *Rect) Clipped(obj Object) Object {
// 	if obj == nil {
// 		return g
// 	}
// 	// convert rect into a polygon
// 	points := make([]geometry.Point, g.base.NumPoints())
// 	for i := 0; i < len(points); i++ {
// 		points[i] = g.base.PointAt(i)
// 	}
// 	poly := geometry.NewPoly(points, nil, geometry.DefaultIndex)
// 	var polygon Polygon
// 	polygon.base = *poly
// 	return polygon.Clipped(obj)
// }

// Distance ...
func (g *Rect) Distance(obj Object) float64 {
	return obj.distanceRect(g.base)
}
func (g *Rect) distancePoint(point geometry.Point) float64 {
	return geoDistancePoints(g.Center(), point)
}
func (g *Rect) distanceRect(rect geometry.Rect) float64 {
	return geoDistancePoints(g.Center(), rect.Center())
}
func (g *Rect) distanceLine(line *geometry.Line) float64 {
	return geoDistancePoints(g.Center(), line.Rect().Center())
}
func (g *Rect) distancePoly(poly *geometry.Poly) float64 {
	return geoDistancePoints(g.Center(), poly.Rect().Center())
}
