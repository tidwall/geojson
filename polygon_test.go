package geojson

import "testing"

func TestPolygon(t *testing.T) {
	g := expectJSON(t, `{"type":"Polygon","coordinates":[[[0,0],[10,0],[10,10],[0,10],[0,0]]]}`, nil)
	if g.Center() != P(5, 5) {
		t.Fatalf("expected '%v', got '%v'", P(5, 5), g.Center())
	}
	json := `{"type":"Polygon","coordinates":[
		[[0,0],[10,0],[10,10],[0,10],[0,0]],
		[[2,2],[8,2],[8,8],[2,8],[2,2]]
	]}`
	g = expectJSON(t, json, nil)
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
