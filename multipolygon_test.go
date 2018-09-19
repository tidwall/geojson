package geojson

import "testing"

func TestMultiPolygon(t *testing.T) {
	json := `{"type":"MultiPolygon","coordinates":[
		[
			[[0,0],[10,0],[10,10],[0,10],[0,0]],
			[[2,2],[8,2],[8,8],[2,8],[2,2]]
		],[
			[[0,0],[10,0],[10,10],[0,10],[0,0]],
			[[2,2],[8,2],[8,8],[2,8],[2,2]]
		]
	]}`
	p := expectJSON(t, json, nil)
	if p.Center() != P(5, 5) {
		t.Fatalf("expected '%v', got '%v'", P(5, 5), p.Center())
	}
	if cleanJSON(string(p.AppendJSON(nil))) != cleanJSON(json) {
		t.Fatalf("expectect '%v', got '%v'", cleanJSON(json), cleanJSON(string(p.AppendJSON(nil))))
	}
	expectJSON(t, `{"type":"MultiPolygon","coordinates":[[[[0,0],[10,0],[5,10],[0,0]],[[1,1]]]],"bbox":null}`, errBBoxInvalid)
	expectJSON(t, `{"type":"MultiPolygon"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"MultiPolygon","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiPolygon","coordinates":[1,null]}`, errCoordinatesInvalid)
}
