package geojson

import "testing"

func TestPolygon(t *testing.T) {
	json := `{"type":"Polygon","coordinates":[[[0,0],[10,0],[10,10],[0,10],[0,0]]]}`
	g := expectJSON(t, json, nil)
	if cleanJSON(string(g.AppendJSON(nil))) != cleanJSON(json) {
		t.Fatalf("expected '%v', got '%v'", cleanJSON(json), cleanJSON(string(g.AppendJSON(nil))))
	}
	json = `{"type":"Polygon","coordinates":[
		[[0,0],[10,0],[10,10],[0,10],[0,0]],
		[[2,2],[8,2],[8,8],[2,8],[2,2]]
	]}`
	g = expectJSON(t, json, nil)
	if g.Center() != P(5, 5) {
		t.Fatalf("expected '%v', got '%v'", P(5, 5), g.Center())
	}
	if cleanJSON(string(g.AppendJSON(nil))) != cleanJSON(json) {
		t.Fatalf("expected '%v', got '%v'", cleanJSON(json), cleanJSON(string(g.AppendJSON(nil))))
	}
	expectJSON(t, `{"type":"Polygon","coordinates":[[1,null]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Polygon","coordinates":[]}`, errMustBeALinearRing)
	expectJSON(t, `{"type":"Polygon","coordinates":[[[0,0],[10,0],[5,10],[0,0]],[[1,1]]]}`, errMustBeALinearRing)
	expectJSON(t, `{"type":"Polygon","coordinates":[[[0,0],[10,0],[5,10],[0,0]]],"bbox":null}`, errBBoxInvalid)
	expectJSON(t, `{"type":"Polygon"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"Polygon","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Polygon","coordinates":[null]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Polygon","coordinates":[[null]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Polygon","coordinates":[[[0,0,0,0,0],[10,0],[5,10],[0,0]]]}`, nil)
}
func TestPolygonPoly(t *testing.T) {
	json := `{"type":"Polygon","coordinates":[[[0,0],[10,0],[10,10],[0,10],[0,0]]]}`
	g := expectJSON(t, json, nil)
	expect(t, g.Contains(P(5, 5)))
	expect(t, g.Contains(R(5, 5, 6, 6)))
	expect(t, g.Contains(Point{Coordinates: P(5, 5)}))
	expect(t, g.Contains(expectJSON(t, `{"type":"LineString","coordinates":[
		[5,5],[5,6],[6,5]
	]}`, nil)))
	expect(t, g.Intersects(P(5, 5)))
	expect(t, g.Intersects(R(5, 5, 6, 6)))
	expect(t, g.Intersects(Point{Coordinates: P(5, 5)}))
	expect(t, g.Intersects(expectJSON(t, `{"type":"LineString","coordinates":[
		[5,5],[5,6],[6,5],[50,50]
	]}`, nil)))
	expect(t, g.Intersects(expectJSON(t, `{"type":"Polygon","coordinates":[[
		[5,5],[5,6],[6,5],[50,50],[5,5]
	]]}`, nil)))
	expect(t, !g.Contains(expectJSON(t, `{"type":"Polygon","coordinates":[[
		[5,5],[5,6],[6,5],[50,50],[5,5]
	]]}`, nil)))
	expect(t, g.Contains(expectJSON(t, `{"type":"Polygon","coordinates":[[
		[5,5],[5,6],[6,5],[5,5]
	]]}`, nil)))
	expect(t, !g.(Polygon).primativeContains(nil))
	expect(t, !g.(Polygon).primativeIntersects(nil))

}
