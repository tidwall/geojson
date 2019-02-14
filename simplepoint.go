package geojson

import "github.com/tidwall/geojson/geometry"

// SimplePoint ...
type SimplePoint struct {
	base geometry.Point
}

// NewSimplePoint ...
func NewSimplePoint(point geometry.Point) *SimplePoint {
	return &SimplePoint{base: point}
}

// ForEach ...
func (g *SimplePoint) ForEach(iter func(geom Object) bool) bool {
	return iter(g)
}

// Empty ...
func (g *SimplePoint) Empty() bool {
	return g.base.Empty()
}

// Valid ...
func (g *SimplePoint) Valid() bool {
	return g.base.Valid()
}

// Rect ...
func (g *SimplePoint) Rect() geometry.Rect {
	return g.base.Rect()
}

// Spatial ...
func (g *SimplePoint) Spatial() Spatial {
	return g
}

// Center ...
func (g *SimplePoint) Center() geometry.Point {
	return g.base
}

// Base ...
func (g *SimplePoint) Base() geometry.Point {
	return g.base
}

// AppendJSON ...
func (g *SimplePoint) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Point","coordinates":`...)
	dst = appendJSONPoint(dst, g.base, nil, 0)
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
	return obj.Spatial().WithinPoint(g.base)
}

// Intersects ...
func (g *SimplePoint) Intersects(obj Object) bool {
	return obj.Spatial().IntersectsPoint(g.base)
}

// WithinRect ...
func (g *SimplePoint) WithinRect(rect geometry.Rect) bool {
	return rect.ContainsPoint(g.base)
}

// WithinPoint ...
func (g *SimplePoint) WithinPoint(point geometry.Point) bool {
	return point.ContainsPoint(g.base)
}

// WithinLine ...
func (g *SimplePoint) WithinLine(line *geometry.Line) bool {
	return line.ContainsPoint(g.base)
}

// WithinPoly ...
func (g *SimplePoint) WithinPoly(poly *geometry.Poly) bool {
	return poly.ContainsPoint(g.base)
}

// IntersectsPoint ...
func (g *SimplePoint) IntersectsPoint(point geometry.Point) bool {
	return g.base.IntersectsPoint(point)
}

// IntersectsRect ...
func (g *SimplePoint) IntersectsRect(rect geometry.Rect) bool {
	return g.base.IntersectsRect(rect)
}

// IntersectsLine ...
func (g *SimplePoint) IntersectsLine(line *geometry.Line) bool {
	return g.base.IntersectsLine(line)
}

// IntersectsPoly ...
func (g *SimplePoint) IntersectsPoly(poly *geometry.Poly) bool {
	return g.base.IntersectsPoly(poly)
}

// NumPoints ...
func (g *SimplePoint) NumPoints() int {
	return 1
}

// Distance ...
func (g *SimplePoint) Distance(obj Object) float64 {
	return obj.Spatial().DistancePoint(g.base)
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
