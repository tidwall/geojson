package geojson

import (
	"github.com/tidwall/geojson/geometry"
)

type Rect struct {
	base geometry.Rect
}

func NewRect(rect geometry.Rect) *Rect {
	return &Rect{base: rect}
}

func (g *Rect) ForEach(iter func(geom Object) bool) bool {
	return iter(g)
}

func (g *Rect) Empty() bool {
	return g.base.Empty()
}

func (g *Rect) Valid() bool {
	return g.base.Valid()
}

func (g *Rect) Rect() geometry.Rect {
	return g.base
}

func (g *Rect) Base() geometry.Rect {
	return g.base
}

func (g *Rect) Center() geometry.Point {
	return g.base.Center()
}

// Polygon returns the Rect as a GeoJSON Polygon.
func (g *Rect) Polygon() Object {
	gPoly := new(Polygon)
	gPoly.base.Exterior = g.base
	return gPoly
}

func (g *Rect) AppendJSON(dst []byte) []byte {
	return g.Polygon().AppendJSON(dst)
}

func (g *Rect) JSON() string {
	return string(g.AppendJSON(nil))
}

func (g *Rect) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

func (g *Rect) String() string {
	return string(g.AppendJSON(nil))
}

func (g *Rect) Contains(obj Object) bool {
	return obj.Spatial().WithinRect(g.base)
}

func (g *Rect) Within(obj Object) bool {
	return obj.Contains(g)
}

func (g *Rect) WithinRect(rect geometry.Rect) bool {
	return rect.ContainsRect(g.base)
}

func (g *Rect) WithinPoint(point geometry.Point) bool {
	return point.ContainsRect(g.base)
}

func (g *Rect) WithinLine(line *geometry.Line) bool {
	return line.ContainsRect(g.base)
}

func (g *Rect) WithinPoly(poly *geometry.Poly) bool {
	return poly.ContainsRect(g.base)
}

func (g *Rect) Intersects(obj Object) bool {
	return obj.Spatial().IntersectsRect(g.base)
}

func (g *Rect) IntersectsPoint(point geometry.Point) bool {
	return g.base.IntersectsPoint(point)
}

func (g *Rect) IntersectsRect(rect geometry.Rect) bool {
	return g.base.IntersectsRect(rect)
}

func (g *Rect) IntersectsLine(line *geometry.Line) bool {
	return g.base.IntersectsLine(line)
}

func (g *Rect) IntersectsPoly(poly *geometry.Poly) bool {
	return g.base.IntersectsPoly(poly)
}

func (g *Rect) NumPoints() int {
	return 2
}

func (g *Rect) Spatial() Spatial {
	return g
}

func (g *Rect) Distance(obj Object) float64 {
	return obj.Spatial().DistanceRect(g.base)
}

func (g *Rect) DistancePoint(point geometry.Point) float64 {
	return geoDistancePoints(g.Center(), point)
}

func (g *Rect) DistanceRect(rect geometry.Rect) float64 {
	return geoDistancePoints(g.Center(), rect.Center())
}

func (g *Rect) DistanceLine(line *geometry.Line) float64 {
	return geoDistancePoints(g.Center(), line.Rect().Center())
}

func (g *Rect) DistancePoly(poly *geometry.Poly) float64 {
	return geoDistancePoints(g.Center(), poly.Rect().Center())
}

func (g *Rect) Members() string {
	return ""
}
