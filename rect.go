package geojson

import (
	"github.com/tidwall/geojson/geos"
)

// Rect ...
type Rect struct {
	base geos.Rect
}

// NewRect ...
func NewRect(minX, minY, maxX, maxY float64) *Rect {
	return &Rect{base: geos.Rect{
		Min: geos.Point{X: minX, Y: minY},
		Max: geos.Point{X: maxX, Y: maxY},
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
func (g *Rect) Rect() geos.Rect {
	return g.base
}

// Center ...
func (g *Rect) Center() geos.Point {
	return g.base.Center()
}

// AppendJSON ...
func (g *Rect) AppendJSON(dst []byte) []byte {
	panic("not ready")
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

func (g *Rect) withinRect(rect geos.Rect) bool {
	return rect.ContainsRect(g.base)
}

func (g *Rect) withinPoint(point geos.Point) bool {
	return point.ContainsRect(g.base)
}

func (g *Rect) withinLine(line *geos.Line) bool {
	return line.ContainsRect(g.base)
}

func (g *Rect) withinPoly(poly *geos.Poly) bool {
	return poly.ContainsRect(g.base)
}

func (g *Rect) intersectsPoint(point geos.Point) bool {
	return g.base.IntersectsPoint(point)
}

func (g *Rect) intersectsRect(rect geos.Rect) bool {
	return g.base.IntersectsRect(rect)
}

func (g *Rect) intersectsLine(line *geos.Line) bool {
	return g.base.IntersectsLine(line)
}

func (g *Rect) intersectsPoly(poly *geos.Poly) bool {
	return g.base.IntersectsPoly(poly)
}
