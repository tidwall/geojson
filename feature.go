package geojson

import "github.com/tidwall/geojson/geos"

// Feature ...
type Feature struct {
	base       Object
	extra      *extra
	id         string
	properties string
}

func (g *Feature) bbox() *geos.Rect {
	if g.extra != nil {
		return g.extra.bbox
	}
	return nil
}

// Empty ...
func (g *Feature) Empty() bool {
	if g.extra != nil && g.extra.bbox != nil {
		return false
	}
	return g.base.Empty()
}

// Rect ...
func (g *Feature) Rect() geos.Rect {
	if g.extra != nil && g.extra.bbox != nil {
		return *g.extra.bbox
	}
	return g.base.Rect()
}

// Center ...
func (g *Feature) Center() geos.Point {
	return g.Rect().Center()
}

// AppendJSON ...
func (g *Feature) AppendJSON(dst []byte) []byte {
	panic("not ready")
}

// Within ...
func (g *Feature) Within(obj Object) bool {
	return obj.Contains(g)
}

// Contains ...
func (g *Feature) Contains(obj Object) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return obj.withinRect(*g.extra.bbox)
	}
	return obj.Within(g.base)
}

func (g *Feature) withinRect(rect geos.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return rect.ContainsRect(*g.extra.bbox)
	}
	return g.base.withinRect(rect)
}

func (g *Feature) withinPoint(point geos.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return point.ContainsRect(*g.extra.bbox)
	}
	return g.base.withinPoint(point)
}

func (g *Feature) withinLine(line *geos.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return line.ContainsRect(*g.extra.bbox)
	}
	return g.base.withinLine(line)
}

func (g *Feature) withinPoly(poly *geos.Poly) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return poly.ContainsRect(*g.extra.bbox)
	}
	return g.base.withinPoly(poly)
}

// Intersects ...
func (g *Feature) Intersects(obj Object) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return obj.intersectsRect(*g.extra.bbox)
	}
	return obj.Intersects(g.base)
}

func (g *Feature) intersectsPoint(point geos.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoint(point)
	}
	return g.base.intersectsPoint(point)
}

func (g *Feature) intersectsRect(rect geos.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsRect(rect)
	}
	return g.base.intersectsRect(rect)
}

func (g *Feature) intersectsLine(line *geos.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsLine(line)
	}
	return g.base.intersectsLine(line)
}

func (g *Feature) intersectsPoly(poly *geos.Poly) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoly(poly)
	}
	return g.base.intersectsPoly(poly)
}
