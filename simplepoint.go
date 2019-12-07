package geojson

import "github.com/tidwall/geojson/geometry"

// SimplePoint ...
type SimplePoint struct {
	geometry.Point
}

// NewSimplePoint returns a new SimplePoint object.
func NewSimplePoint(point geometry.Point) *SimplePoint {
	return &SimplePoint{Point: point}
}

// ForEach ...
func (g *SimplePoint) ForEach(iter func(geom Object) bool) bool {
	return iter(g)
}

// Empty ...
func (g *SimplePoint) Empty() bool {
	return g.Point.Empty()
}

// Valid ...
func (g *SimplePoint) Valid() bool {
	return g.Point.Valid()
}

// Rect ...
func (g *SimplePoint) Rect() geometry.Rect {
	return g.Point.Rect()
}

// Spatial ...
func (g *SimplePoint) Spatial() Spatial {
	return g
}

// Center ...
func (g *SimplePoint) Center() geometry.Point {
	return g.Point
}

// Base ...
func (g *SimplePoint) Base() geometry.Point {
	return g.Point
}

// AppendJSON ...
func (g *SimplePoint) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Point","coordinates":`...)
	dst = appendJSONPoint(dst, g.Point, nil, 0)
	dst = append(dst, '}')
	return dst
}

// JSON ...
func (g *SimplePoint) JSON() string {
	return string(g.AppendJSON(nil))
}

// MarshalJSON ...
func (g *SimplePoint) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

// String ...
func (g *SimplePoint) String() string {
	return string(g.AppendJSON(nil))
}

// Within ...
func (g *SimplePoint) Within(obj Object) bool {
	return obj.Contains(g)
}

// Contains ...
func (g *SimplePoint) Contains(obj Object) bool {
	return obj.Spatial().WithinPoint(g.Point)
}

// Intersects ...
func (g *SimplePoint) Intersects(obj Object) bool {
	if obj, ok := obj.(*Circle); ok {
		return obj.Contains(g)
	}
	return obj.Spatial().IntersectsPoint(g.Point)
}

// WithinRect ...
func (g *SimplePoint) WithinRect(rect geometry.Rect) bool {
	return rect.ContainsPoint(g.Point)
}

// WithinPoint ...
func (g *SimplePoint) WithinPoint(point geometry.Point) bool {
	return point.ContainsPoint(g.Point)
}

// WithinLine ...
func (g *SimplePoint) WithinLine(line *geometry.Line) bool {
	return line.ContainsPoint(g.Point)
}

// WithinPoly ...
func (g *SimplePoint) WithinPoly(poly *geometry.Poly) bool {
	return poly.ContainsPoint(g.Point)
}

// IntersectsPoint ...
func (g *SimplePoint) IntersectsPoint(point geometry.Point) bool {
	return g.Point.IntersectsPoint(point)
}

// IntersectsRect ...
func (g *SimplePoint) IntersectsRect(rect geometry.Rect) bool {
	return g.Point.IntersectsRect(rect)
}

// IntersectsLine ...
func (g *SimplePoint) IntersectsLine(line *geometry.Line) bool {
	return g.Point.IntersectsLine(line)
}

// IntersectsPoly ...
func (g *SimplePoint) IntersectsPoly(poly *geometry.Poly) bool {
	return g.Point.IntersectsPoly(poly)
}

// NumPoints ...
func (g *SimplePoint) NumPoints() int {
	return 1
}

// Distance ...
func (g *SimplePoint) Distance(obj Object) float64 {
	return obj.Spatial().DistancePoint(g.Point)
}

// DistancePoint ...
func (g *SimplePoint) DistancePoint(point geometry.Point) float64 {
	return geoDistancePoints(g.Center(), point)
}

// DistanceRect ...
func (g *SimplePoint) DistanceRect(rect geometry.Rect) float64 {
	return geoDistancePoints(g.Center(), rect.Center())
}

// DistanceLine ...
func (g *SimplePoint) DistanceLine(line *geometry.Line) float64 {
	return geoDistancePoints(g.Center(), line.Rect().Center())
}

// DistancePoly ...
func (g *SimplePoint) DistancePoly(poly *geometry.Poly) float64 {
	return geoDistancePoints(g.Center(), poly.Rect().Center())
}
