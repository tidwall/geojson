package geojson

import (
	"math/rand"
	"testing"

	"github.com/tidwall/geojson/geometry"
)

func TestRect(t *testing.T) {
	rect := RO(10, 20, 30, 40)
	expect(t, !rect.Empty())
	expect(t, string(rect.AppendJSON(nil)) ==
		`{"type":"Polygon","coordinates":[[[10,20],[30,20],[30,40],[10,40],[10,20]]]}`)
	expect(t, rect.String() == string(rect.AppendJSON(nil)))
	// expect(t, !rect.Contains(NewString("")))
	// expect(t, !rect.Within(NewString("")))
	// expect(t, !rect.Intersects(NewString("")))
	// expect(t, rect.Distance(NewString("")) == 0)

	expect(t, rect.Rect() == R(10, 20, 30, 40))
	expect(t, rect.Center() == P(20, 30))
	var g Object
	rect.ForEach(func(o Object) bool {
		expect(t, g == nil)
		g = o
		return true
	})
	expect(t, g == rect)

	expect(t, rect.NumPoints() == 2)

	expect(t, !(&Point{}).Contains(rect))
	expect(t, !(&Rect{}).Contains(rect))
	expect(t, !(&LineString{}).Contains(rect))
	expect(t, !(&Polygon{}).Contains(rect))

	expect(t, !(&Point{}).Intersects(rect))
	expect(t, !(&Rect{}).Intersects(rect))
	expect(t, !(&LineString{}).Intersects(rect))
	expect(t, !(&Polygon{}).Intersects(rect))

	expect(t, (&Point{}).Distance(rect) != 0)
	expect(t, (&Rect{}).Distance(rect) != 0)
	expect(t, (&LineString{}).Distance(rect) != 0)
	expect(t, (&Polygon{}).Distance(rect) != 0)

}

func TestRectValid(t *testing.T) {
	json := `{"type":"Polygon","coordinates":[[[10,200],[30,200],[30,40],[10,40],[10,200]]]}`
	expectJSON(t, json, nil)
	expectJSONOpts(t, json, errCoordinatesInvalid, &ParseOptions{RequireValid: true})
}

func BenchmarkRectValid(b *testing.B) {
	rects := make([]*Rect, b.N)
	for i := 0; i < b.N; i++ {
		min := geometry.Point{
			X: rand.Float64()*400 - 200, // some are out of bounds
			Y: rand.Float64()*200 - 100, // some are out of bounds
		}
		max := geometry.Point{
			X: rand.Float64()*400 - 200, // some are out of bounds
			Y: rand.Float64()*200 - 100, // some are out of bounds
		}
		if min.X > max.X {
			min.X, max.X = max.X, min.X
		}
		if min.Y > max.Y {
			min.Y, max.Y = max.Y, min.Y
		}
		rects[i] = NewRect(geometry.Rect{Min: min, Max: max})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rects[i].Valid()
	}
}
