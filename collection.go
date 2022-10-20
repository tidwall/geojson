package geojson

import (
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/rtree"
)

type collection struct {
	children []Object
	extra    *extra
	tree     *rtree.RTree
	prect    geometry.Rect
	pempty   bool
}

func (g *collection) Indexed() bool {
	return g.tree != nil
}

func (g *collection) Children() []Object {
	return g.children
}

func (g *collection) ForEach(iter func(geom Object) bool) bool {
	for _, child := range g.children {
		if !child.ForEach(iter) {
			return false
		}
	}
	return true
}

func (g *collection) Base() []Object {
	return g.children
}

func (g *collection) Search(rect geometry.Rect, iter func(child Object) bool) {
	if g.tree != nil {
		g.tree.Search(
			[2]float64{rect.Min.X, rect.Min.Y},
			[2]float64{rect.Max.X, rect.Max.Y},
			func(_, _ [2]float64, value interface{}) bool {
				return iter(value.(Object))
			},
		)
	} else {
		for _, child := range g.children {
			if child.Empty() {
				continue
			}
			if child.Rect().IntersectsRect(rect) {
				if !iter(child) {
					break
				}
			}
		}
	}
}

func (g *collection) Empty() bool {
	return g.pempty
}

func (g *collection) Valid() bool {
	return g.Rect().Valid()
}

func (g *collection) Rect() geometry.Rect {
	return g.prect
}

func (g *collection) Center() geometry.Point {
	return g.Rect().Center()
}

func (g *collection) AppendJSON(dst []byte) []byte {
	// this should never be called
	return append(dst, "null"...)
}

func (g *collection) JSON() string {
	return string(g.AppendJSON(nil))
}

func (g *collection) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

func (g *collection) String() string {
	return string(g.AppendJSON(nil))
}

func (g *collection) Within(obj Object) bool {
	return obj.Contains(g)
}

func (g *collection) Contains(obj Object) bool {
	if g.Empty() {
		return false
	}
	// all of obj must be contained by any number of the collection children
	var objContained bool
	obj.ForEach(func(geom Object) bool {
		if geom.Empty() {
			// ignore empties
			return true
		}
		var geomContained bool
		g.Search(geom.Rect(), func(child Object) bool {
			if child.Contains(geom) {
				// found a child object that contains geom, end inner loop
				geomContained = true
				return false
			}
			return true
		})
		if !geomContained {
			// unmark and quit the loop
			objContained = false
			return false
		}
		// mark that at least one geom is contained
		objContained = true
		return true
	})
	return objContained
}

func (g *collection) Spatial() Spatial { return g }

func (g *collection) WithinRect(rect geometry.Rect) bool {
	if g.Empty() {
		return false
	}
	var withinCount int
	g.Search(rect, func(child Object) bool {
		if child.Spatial().WithinRect(rect) {
			withinCount++
			return true
		}
		return false
	})
	return withinCount == len(g.children)
}

func (g *collection) WithinPoint(point geometry.Point) bool {
	if g.Empty() {
		return false
	}
	var withinCount int
	g.Search(point.Rect(), func(child Object) bool {
		if child.Spatial().WithinPoint(point) {
			withinCount++
			return true
		}
		return false
	})
	return withinCount == len(g.children)
}

func (g *collection) WithinLine(line *geometry.Line) bool {
	if g.Empty() {
		return false
	}
	var withinCount int
	g.Search(line.Rect(), func(child Object) bool {
		if child.Spatial().WithinLine(line) {
			withinCount++
			return true
		}
		return false
	})
	return withinCount == len(g.children)
}

func (g *collection) WithinPoly(poly *geometry.Poly) bool {
	if g.Empty() {
		return false
	}
	var withinCount int
	g.Search(poly.Rect(), func(child Object) bool {
		if child.Spatial().WithinPoly(poly) {
			withinCount++
			return true
		}
		return false
	})
	return withinCount == len(g.children)
}

func (g *collection) Intersects(obj Object) bool {
	// check if any of obj intersects with any of collection
	var intersects bool
	obj.ForEach(func(geom Object) bool {
		if geom.Empty() {
			// ignore the empties
			return true
		}
		g.Search(geom.Rect(), func(child Object) bool {
			if child.Intersects(geom) {
				intersects = true
				return false
			}
			return true
		})
		return !intersects
	})
	return intersects
}

func (g *collection) IntersectsPoint(point geometry.Point) bool {
	var intersects bool
	g.Search(point.Rect(), func(child Object) bool {
		if child.Spatial().IntersectsPoint(point) {
			intersects = true
			return false
		}
		return true
	})
	return intersects
}

func (g *collection) IntersectsRect(rect geometry.Rect) bool {
	var intersects bool
	g.Search(rect, func(child Object) bool {
		if child.Spatial().IntersectsRect(rect) {
			intersects = true
			return false
		}
		return true
	})
	return intersects
}

func (g *collection) IntersectsLine(line *geometry.Line) bool {
	var intersects bool
	g.Search(line.Rect(), func(child Object) bool {
		if child.Spatial().IntersectsLine(line) {
			intersects = true
			return false
		}
		return true
	})
	return intersects
}

func (g *collection) IntersectsPoly(poly *geometry.Poly) bool {
	var intersects bool
	g.Search(poly.Rect(), func(child Object) bool {
		if child.Spatial().IntersectsPoly(poly) {
			intersects = true
			return false
		}
		return true
	})
	return intersects
}

func (g *collection) NumPoints() int {
	var n int
	for _, child := range g.children {
		n += child.NumPoints()
	}
	return n
}

func (g *collection) parseInitRectIndex(opts *ParseOptions) {
	g.pempty = true
	var count int
	for _, child := range g.children {
		if child.Empty() {
			continue
		}
		if g.pempty && !child.Empty() {
			g.pempty = false
		}
		if count == 0 {
			g.prect = child.Rect()
		} else {
			if len(g.children) == 1 {
				g.prect = child.Rect()
			} else {
				g.prect = unionRects(g.prect, child.Rect())
			}
		}
		count++
	}
	if count > 0 && opts.IndexChildren != 0 && count >= opts.IndexChildren {
		g.tree = new(rtree.RTree)
		for _, child := range g.children {
			if child.Empty() {
				continue
			}
			rect := child.Rect()
			g.tree.Insert(
				[2]float64{rect.Min.X, rect.Min.Y},
				[2]float64{rect.Max.X, rect.Max.Y},
				child,
			)
		}
	}
}

func (g *collection) Distance(obj Object) float64 {
	return obj.Spatial().DistancePoint(g.Center())
}
func (g *collection) DistancePoint(point geometry.Point) float64 {
	return geoDistancePoints(g.Center(), point)
}
func (g *collection) DistanceRect(rect geometry.Rect) float64 {
	return geoDistancePoints(g.Center(), rect.Center())
}
func (g *collection) DistanceLine(line *geometry.Line) float64 {
	return geoDistancePoints(g.Center(), line.Rect().Center())
}
func (g *collection) DistancePoly(poly *geometry.Poly) float64 {
	return geoDistancePoints(g.Center(), poly.Rect().Center())
}

func (g *collection) Members() string {
	if g.extra != nil {
		return g.extra.members
	}
	return ""
}
