package geojson

import "github.com/tidwall/geojson/geos"

// Polygon ...
type Polygon struct {
	base  geos.Poly
	extra *extra
}

// Empty ...
func (g *Polygon) Empty() bool {
	if g.extra != nil && g.extra.bbox != nil {
		return false
	}
	return g.base.Empty()
}

// Rect ...
func (g *Polygon) Rect() geos.Rect {
	if g.extra != nil && g.extra.bbox != nil {
		return *g.extra.bbox
	}
	return g.base.Rect()
}

// Center ...
func (g *Polygon) Center() geos.Point {
	return g.Rect().Center()
}

// AppendJSON ...
func (g *Polygon) AppendJSON(dst []byte) []byte {
	panic("not ready")
}

// Within ...
func (g *Polygon) Within(obj Object) bool {
	return obj.Contains(g)
}

// Contains ...
func (g *Polygon) Contains(obj Object) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return obj.withinRect(*g.extra.bbox)
	}
	return obj.withinPoly(&g.base)
}

func (g *Polygon) withinRect(rect geos.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return rect.ContainsRect(*g.extra.bbox)
	}
	return rect.ContainsPoly(&g.base)
}

func (g *Polygon) withinPoint(point geos.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return point.ContainsRect(*g.extra.bbox)
	}
	return point.ContainsPoly(&g.base)
}

func (g *Polygon) withinLine(line *geos.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return line.ContainsRect(*g.extra.bbox)
	}
	return line.ContainsPoly(&g.base)
}

func (g *Polygon) withinPoly(poly *geos.Poly) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return poly.ContainsRect(*g.extra.bbox)
	}
	return poly.ContainsPoly(&g.base)
}

// Intersects ...
func (g *Polygon) Intersects(obj Object) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return obj.intersectsRect(*g.extra.bbox)
	}
	return obj.intersectsPoly(&g.base)
}

func (g *Polygon) intersectsPoint(point geos.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoint(point)
	}
	return g.base.IntersectsPoint(point)
}

func (g *Polygon) intersectsRect(rect geos.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsRect(rect)
	}
	return g.base.IntersectsRect(rect)
}

func (g *Polygon) intersectsLine(line *geos.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsLine(line)
	}
	return g.base.IntersectsLine(line)
}

func (g *Polygon) intersectsPoly(poly *geos.Poly) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoly(poly)
	}
	return g.base.IntersectsPoly(poly)
}
