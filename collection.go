package geojson

import (
	"github.com/tidwall/boxtree/d2"
	"github.com/tidwall/geojson/geos"
)

// minTreeChildren are the minumum number of child objects required before it
// makes sense to index the children
const minTreeChildren = 32

type collection struct {
	children []Object
	extra    *extra
	tree     *d2.BoxTree
	prect    geos.Rect
	pempty   bool
}

func (g *collection) index() {
	g.tree = new(d2.BoxTree)
	for _, child := range g.children {
		rect := child.Rect()
		g.tree.Insert(
			[]float64{rect.Min.X, rect.Min.Y},
			[]float64{rect.Max.X, rect.Max.Y},
			child,
		)
	}
}

func (g *collection) search(rect geos.Rect, iter func(child Object) bool) {
	if g.tree != nil {
		g.tree.Search(
			[]float64{rect.Min.X, rect.Min.Y},
			[]float64{rect.Max.X, rect.Max.Y},
			func(_, _ []float64, value interface{}) bool {
				return iter(value.(Object))
			},
		)
	} else {
		for _, child := range g.children {
			if child.Rect().IntersectsRect(rect) {
				if !iter(child) {
					break
				}
			}
		}
	}
}

// Empty ...
func (g *collection) Empty() bool {
	return g.pempty
}

// Rect ...
func (g *collection) Rect() geos.Rect {
	if g.extra != nil && g.extra.bbox != nil {
		return *g.extra.bbox
	}
	return g.prect
}

// Center ...
func (g *collection) Center() geos.Point {
	return g.Rect().Center()
}

// AppendJSON ...
func (g *collection) AppendJSON(dst []byte) []byte {
	panic("not ready")
}

// Within ...
func (g *collection) Within(obj Object) bool {
	return obj.Contains(g)
}

// Contains ...
func (g *collection) Contains(obj Object) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return obj.withinRect(*g.extra.bbox)
	}
	objRect := obj.Rect()
	if !g.prect.ContainsRect(objRect) {
		return false
	}
	var contains bool
	g.search(objRect, func(child Object) bool {
		if child.Contains(obj) {
			contains = true
			return false
		}
		return true
	})
	return contains
}

func (g *collection) withinRect(rect geos.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return rect.ContainsRect(*g.extra.bbox)
	}
	if g.Empty() {
		return false
	}
	var withinCount int
	g.search(rect, func(child Object) bool {
		if child.withinRect(rect) {
			withinCount++
			return true
		}
		return false
	})
	return withinCount == len(g.children)
}

func (g *collection) withinPoint(point geos.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return point.ContainsRect(*g.extra.bbox)
	}
	if g.Empty() {
		return false
	}
	var withinCount int
	g.search(point.Rect(), func(child Object) bool {
		if child.withinPoint(point) {
			withinCount++
			return true
		}
		return false
	})
	return withinCount == len(g.children)
}

func (g *collection) withinLine(line *geos.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return line.ContainsRect(*g.extra.bbox)
	}
	if g.Empty() {
		return false
	}
	var withinCount int
	g.search(line.Rect(), func(child Object) bool {
		if child.withinLine(line) {
			withinCount++
			return true
		}
		return false
	})
	return withinCount == len(g.children)
}

func (g *collection) withinPoly(poly *geos.Poly) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return poly.ContainsRect(*g.extra.bbox)
	}
	if g.Empty() {
		return false
	}
	var withinCount int
	g.search(poly.Rect(), func(child Object) bool {
		if child.withinPoly(poly) {
			withinCount++
			return true
		}
		return false
	})
	return withinCount == len(g.children)
}

// Intersects ...
func (g *collection) Intersects(obj Object) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return obj.intersectsRect(*g.extra.bbox)
	}
	var intersects bool
	g.search(obj.Rect(), func(child Object) bool {
		if child.Intersects(obj) {
			intersects = true
			return false
		}
		return true
	})
	return intersects
}

func (g *collection) intersectsPoint(point geos.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoint(point)
	}
	var intersects bool
	g.search(point.Rect(), func(child Object) bool {
		if child.intersectsPoint(point) {
			intersects = true
			return false
		}
		return true
	})
	return intersects
}

func (g *collection) intersectsRect(rect geos.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsRect(rect)
	}
	var intersects bool
	g.search(rect, func(child Object) bool {
		if child.intersectsRect(rect) {
			intersects = true
			return false
		}
		return true
	})
	return intersects
}

func (g *collection) intersectsLine(line *geos.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsLine(line)
	}
	var intersects bool
	g.search(line.Rect(), func(child Object) bool {
		if child.intersectsLine(line) {
			intersects = true
			return false
		}
		return true
	})
	return intersects
}

func (g *collection) intersectsPoly(poly *geos.Poly) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoly(poly)
	}
	var intersects bool
	g.search(poly.Rect(), func(child Object) bool {
		if child.intersectsPoly(poly) {
			intersects = true
			return false
		}
		return true
	})
	return intersects
}
