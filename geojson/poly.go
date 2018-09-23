package geojson

import (
	"unsafe"

	"github.com/tidwall/geojson/poly"
)

func polyPoint(posn Position) poly.Point {
	return *(*poly.Point)(unsafe.Pointer(&posn))
}
func polyRect(rect Rect) poly.Rect {
	return *(*poly.Rect)(unsafe.Pointer(&rect))
}
func polyLine(line []Position) poly.Line {
	return *(*poly.Line)(unsafe.Pointer(&line))
}
func polyPolygon(polygon [][]Position) poly.Polygon {
	var newPoly poly.Polygon
	if len(polygon) > 0 {
		newPoly.Exterior = *(*poly.Ring)(unsafe.Pointer(&polygon[0]))
		if len(polygon) > 1 {
			newPoly.Holes = (*(*[]poly.Ring)(unsafe.Pointer(&polygon)))[1:]
		}
	}
	return newPoly
}
