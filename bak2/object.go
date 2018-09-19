package geojson

import (
	"unsafe"

	"github.com/tidwall/geojson/poly"
)

// Object ...
type Object interface {
	// Rect is the outer bounding rectangle of the object
	Rect() Rect
	// Center of the object
	Center() Position
	// AppendJSON appends a json representation to destination
	AppendJSON(dst []byte) []byte
	// JSON representation of the object
	JSON() string
	// String representation of the object
	String() string
	// Stats of the object
	Stats() Stats
	// Contains test if object contains another object
	Contains(other Object) bool
	// Intersects test if object intersects another object
	Intersects(other Object) bool
	// // WithinPolyLine test if object is within a polyline
	// WithinPolyLine(line []Position) bool
	// IntersectsPolyLine test if object intersect a polyline
	IntersectsPolyLine(line []Position) bool
}

func mustConformObject() {
	var obj Object
	// primatives
	//obj = SimplePoint{}
	obj = Point{}
	obj = LineString{}
	obj = Polygon{}

	// obj = Feature{}
	// obj = FeatureCollection{}
	// obj = GeomCollection{}
	obj = String("")
	_ = obj
}

// objectContainsRect tests if object contains another object, but only checks
// the outer bounding rectangles
// returns false for 'certain' when it's not known for sure.
func objectContainsRect(gBBox BBox, g, o Object) (contains, certain bool) {
	if !g.Rect().ContainsRect(o.Rect()) {
		return false, true
	}
	if gBBox != nil && gBBox.Defined() {
		return true, true
	}
	return false, false
}

// objectIntersectsRect tests if object intersects another object, but only
// checks the outer bounding rectangles
// returns false for 'certain' when it's not known for sure.
func objectIntersectsRect(gBBox BBox, g, o Object) (intersects, certain bool) {
	if !g.Rect().IntersectsRect(o.Rect()) {
		return false, true
	}
	if gBBox != nil && gBBox.Defined() {
		return gBBox.Rect().Intersects(o), true
	}
	return false, false
}

func polyLine(coords []Position) poly.Polygon {
	return poly.Polygon(*(*[]poly.Point)(unsafe.Pointer(&coords)))
}

func polyRect(rect Rect) poly.Rect {
	return *(*poly.Rect)(unsafe.Pointer(&rect))
}

func polyPoint(posn Position) poly.Point {
	return poly.Point(posn)
}
