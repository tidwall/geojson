package geojson

import (
	"github.com/tidwall/geojson/geometry"
)

// Rect ...
type Rect struct {
	base geometry.Rect
}

// NewRect ...
func NewRect(rect geometry.Rect) *Rect {
	return &Rect{base: rect}
}

// ForEach ...
func (g *Rect) ForEach(iter func(geom Object) bool) bool {
	return iter(g)
}

// Empty ...
func (g *Rect) Empty() bool {
	return g.base.Empty()
}

// Valid ...
func (g *Rect) Valid() bool {
	return g.base.Valid()
}

// Rect ...
func (g *Rect) Rect() geometry.Rect {
	return g.base
}

// Base ...
func (g *Rect) Base() geometry.Rect {
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

// MarshalJSON ...
func (g *Rect) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

// String ...
func (g *Rect) String() string {
	return string(g.AppendJSON(nil))
}

// Contains ...
func (g *Rect) Contains(obj Object) bool {
	return obj.Spatial().WithinRect(g.base)
}

// Within ...
func (g *Rect) Within(obj Object) bool {
	return obj.Contains(g)
}

// WithinRect ...
func (g *Rect) WithinRect(rect geometry.Rect) bool {
	return rect.ContainsRect(g.base)
}

// WithinPoint ...
func (g *Rect) WithinPoint(point geometry.Point) bool {
	return point.ContainsRect(g.base)
}

// WithinLine ...
func (g *Rect) WithinLine(line *geometry.Line) bool {
	return line.ContainsRect(g.base)
}

// WithinPoly ...
func (g *Rect) WithinPoly(poly *geometry.Poly) bool {
	return poly.ContainsRect(g.base)
}

// Intersects ...
func (g *Rect) Intersects(obj Object) bool {
	return obj.Spatial().IntersectsRect(g.base)
}

// IntersectsPoint ...
func (g *Rect) IntersectsPoint(point geometry.Point) bool {
	return g.base.IntersectsPoint(point)
}

// IntersectsRect ...
func (g *Rect) IntersectsRect(rect geometry.Rect) bool {
	return g.base.IntersectsRect(rect)
}

// IntersectsLine ...
func (g *Rect) IntersectsLine(line *geometry.Line) bool {
	return g.base.IntersectsLine(line)
}

// IntersectsPoly ...
func (g *Rect) IntersectsPoly(poly *geometry.Poly) bool {
	return g.base.IntersectsPoly(poly)
}

// NumPoints ...
func (g *Rect) NumPoints() int {
	return 2
}

// Spatial ...
func (g *Rect) Spatial() Spatial {
	return g
}

// Distance ...
func (g *Rect) Distance(obj Object) float64 {
	return obj.Spatial().DistanceRect(g.base)
}

// DistancePoint ...
func (g *Rect) DistancePoint(point geometry.Point) float64 {
	return geoDistancePoints(g.Center(), point)
}

// DistanceRect ...
func (g *Rect) DistanceRect(rect geometry.Rect) float64 {
	return geoDistancePoints(g.Center(), rect.Center())
}

// DistanceLine ...
func (g *Rect) DistanceLine(line *geometry.Line) float64 {
	return geoDistancePoints(g.Center(), line.Rect().Center())
}

// DistancePoly ...
func (g *Rect) DistancePoly(poly *geometry.Poly) float64 {
	return geoDistancePoints(g.Center(), poly.Rect().Center())
}
