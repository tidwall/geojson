package geojson

import "testing"

func TestLineString(t *testing.T) {
	g := expectJSON(t, `{"type":"LineString","coordinates":[[1,2,3]]}`, nil)
	if g.Center() != P(1, 2) {
		t.Fatalf("expected '%v', got '%v'", P(1, 2), g.Center())
	}
	expectJSON(t, `{"type":"LineString","coordinates":[[1,null]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[[1,2]],"bbox":null}`, errBBoxInvalid)
	expectJSON(t, `{"type":"LineString"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"LineString","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[[null]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[null]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[[1,2,3,4,5]]}`, nil)
	expectJSON(t, `{"type":"LineString","coordinates":[[1]]}`, errCoordinatesInvalid)
}

func TestLineStringPoly(t *testing.T) {
	ls := expectJSON(t, `{"type":"LineString","coordinates":[
		[10,10],[20,20],[20,10]
	]}`, nil)
	expect(t, !ls.BBoxDefined())
	ls.ForEachChild(func(Object) bool { panic("should not be reached") })
	expect(t, ls.(LineString).Contains(ls))
	expect(t, ls.Contains(P(10, 10)))
	expect(t, ls.Contains(P(15, 15)))
	expect(t, ls.Contains(P(20, 20)))
	expect(t, ls.Contains(P(20, 15)))
	expect(t, !ls.Contains(P(12, 13)))
	expect(t, !ls.Contains(R(10, 10, 20, 20)))
	expect(t, ls.Intersects(P(10, 10)))
	expect(t, ls.Intersects(P(15, 15)))
	expect(t, ls.Intersects(P(20, 20)))
	expect(t, !ls.Intersects(P(12, 13)))
	expect(t, ls.Intersects(R(10, 10, 20, 20)))
	expect(t, ls.Intersects(
		expectJSON(t, `{"type":"Point","coordinates":[15,15,0]}`, nil),
	))
	expect(t, ls.Intersects(ls))

	lsb := expectJSON(t, `{"type":"LineString","coordinates":[
		[10,10],[20,20],[20,10]
	],"bbox":[10,10,20,20]}`, nil)
	expect(t, lsb.Contains(P(12, 13)))

	expect(t, ls.Contains(Point{Coordinates: P(20, 20)}))
	expect(t, ls.Contains(Polygon{Coordinates: [][]Position{{P(20, 20)}}}))
	expect(t, ls.Intersects(Polygon{Coordinates: [][]Position{{P(20, 20)}}}))
	expect(t, !ls.(LineString).primativeContains(nil))
	expect(t, !ls.(LineString).primativeIntersects(nil))

}
