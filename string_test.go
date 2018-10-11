package geojson

import (
	"testing"

	"github.com/tidwall/geojson/geometry"
)

func TestString(t *testing.T) {
	s := NewString("hello")
	expect(t, s.Empty())
	expect(t, s.String() == "hello")
	expect(t, !s.Contains(NewString("")))
	expect(t, !s.Within(NewString("")))
	expect(t, !s.Intersects(NewString("")))
	expect(t, s.Distance(NewString("")) == 0)
	expect(t, s.Rect() == geometry.Rect{})
	expect(t, s.Center() == geometry.Point{})
	var g Object
	s.forEach(func(o Object) bool {
		expect(t, g == nil)
		g = o
		return true
	})
	expect(t, g == s)
	expect(t, string(s.AppendJSON(nil)) == `"hello"`)

	expect(t, s.NumPoints() == 0)
	expect(t, s.Clipped(nil) == s)

	expect(t, !(&Point{}).Contains(s))
	expect(t, !(&Rect{}).Contains(s))
	expect(t, !(&LineString{}).Contains(s))
	expect(t, !(&Polygon{}).Contains(s))

	expect(t, !(&Point{}).Intersects(s))
	expect(t, !(&Rect{}).Intersects(s))
	expect(t, !(&LineString{}).Intersects(s))
	expect(t, !(&Polygon{}).Intersects(s))

	expect(t, (&Point{}).Distance(s) == 0)
	expect(t, (&Rect{}).Distance(s) == 0)
	expect(t, (&LineString{}).Distance(s) == 0)
	expect(t, (&Polygon{}).Distance(s) == 0)

}
