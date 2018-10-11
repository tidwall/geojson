package geojson

import (
	"github.com/tidwall/boxtree/d2"
	"github.com/tidwall/geojson/geometry"
)

// Collection is a searchable with children
type Collection interface {
	Children() []Object
	Indexed() bool
	Search(rect geometry.Rect, iter func(child Object) bool)
}

type collection struct {
	children []Object
	extra    *extra
	tree     *d2.BoxTree
	prect    geometry.Rect
	pempty   bool
}

func (g *collection) Indexed() bool {
	return g.tree != nil
}

func (g *collection) Children() []Object {
	return g.children
}

// forEach ...
func (g *collection) forEach(iter func(geom Object) bool) bool {
	for _, child := range g.children {
		if !child.forEach(iter) {
			return false
		}
	}
	return true
}

func (g *collection) Search(rect geometry.Rect, iter func(child Object) bool) {
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

// Empty ...
func (g *collection) Empty() bool {
	return g.pempty
}

// Rect ...
func (g *collection) Rect() geometry.Rect {
	return g.prect
}

// Center ...
func (g *collection) Center() geometry.Point {
	return g.Rect().Center()
}

// AppendJSON ...
func (g *collection) AppendJSON(dst []byte) []byte {
	// this should never be called
	return append(dst, "null"...)
}

// String ...
func (g *collection) String() string {
	return string(g.AppendJSON(nil))
}

// Within ...
func (g *collection) Within(obj Object) bool {
	return obj.Contains(g)
}

// Contains ...
func (g *collection) Contains(obj Object) bool {
	if g.Empty() {
		return false
	}
	// all of obj must be contained by any number of the collection children
	var objContained bool
	obj.forEach(func(geom Object) bool {
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

func (g *collection) withinRect(rect geometry.Rect) bool {
	if g.Empty() {
		return false
	}
	var withinCount int
	g.Search(rect, func(child Object) bool {
		if child.withinRect(rect) {
			withinCount++
			return true
		}
		return false
	})
	return withinCount == len(g.children)
}

func (g *collection) withinPoint(point geometry.Point) bool {
	if g.Empty() {
		return false
	}
	var withinCount int
	g.Search(point.Rect(), func(child Object) bool {
		if child.withinPoint(point) {
			withinCount++
			return true
		}
		return false
	})
	return withinCount == len(g.children)
}

func (g *collection) withinLine(line *geometry.Line) bool {
	if g.Empty() {
		return false
	}
	var withinCount int
	g.Search(line.Rect(), func(child Object) bool {
		if child.withinLine(line) {
			withinCount++
			return true
		}
		return false
	})
	return withinCount == len(g.children)
}

func (g *collection) withinPoly(poly *geometry.Poly) bool {
	if g.Empty() {
		return false
	}
	var withinCount int
	g.Search(poly.Rect(), func(child Object) bool {
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
	// check if any of obj intersects with any of collection
	var intersects bool
	obj.forEach(func(geom Object) bool {
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
		if intersects {
			return false
		}
		return true
	})
	return intersects
}

func (g *collection) intersectsPoint(point geometry.Point) bool {
	var intersects bool
	g.Search(point.Rect(), func(child Object) bool {
		if child.intersectsPoint(point) {
			intersects = true
			return false
		}
		return true
	})
	return intersects
}

func (g *collection) intersectsRect(rect geometry.Rect) bool {
	var intersects bool
	g.Search(rect, func(child Object) bool {
		if child.intersectsRect(rect) {
			intersects = true
			return false
		}
		return true
	})
	return intersects
}

func (g *collection) intersectsLine(line *geometry.Line) bool {
	var intersects bool
	g.Search(line.Rect(), func(child Object) bool {
		if child.intersectsLine(line) {
			intersects = true
			return false
		}
		return true
	})
	return intersects
}

func (g *collection) intersectsPoly(poly *geometry.Poly) bool {
	var intersects bool
	g.Search(poly.Rect(), func(child Object) bool {
		if child.intersectsPoly(poly) {
			intersects = true
			return false
		}
		return true
	})
	return intersects
}

// NumPoints ...
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
	if opts.IndexChildren != 0 && count >= opts.IndexChildren {
		g.tree = new(d2.BoxTree)
		for _, child := range g.children {
			if child.Empty() {
				continue
			}
			rect := child.Rect()
			g.tree.Insert(
				[]float64{rect.Min.X, rect.Min.Y},
				[]float64{rect.Max.X, rect.Max.Y},
				child,
			)
		}
	}
}

// Clipped ...
func (g *collection) Clipped(obj Object) Object {
	var newChildren []Object
	for _, child := range g.children {
		newChild := child.Clipped(obj)
		if _, ok := newChild.(*Feature); !ok {
			newChild = &Feature{base: newChild}
		}
		newChildren = append(newChildren, newChild)
	}
	multi := new(FeatureCollection)
	multi.children = newChildren
	multi.parseInitRectIndex(DefaultParseOptions)
	return multi
}

// Distance ...
func (g *collection) Distance(obj Object) float64 {
	return obj.distancePoint(g.Center())
}
func (g *collection) distancePoint(point geometry.Point) float64 {
	return geoDistancePoints(g.Center(), point)
}
func (g *collection) distanceRect(rect geometry.Rect) float64 {
	return geoDistancePoints(g.Center(), rect.Center())
}
func (g *collection) distanceLine(line *geometry.Line) float64 {
	return geoDistancePoints(g.Center(), line.Rect().Center())
}
func (g *collection) distancePoly(poly *geometry.Poly) float64 {
	return geoDistancePoints(g.Center(), poly.Rect().Center())
}
