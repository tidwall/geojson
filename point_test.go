package geojson

import "testing"

func TestPoint(t *testing.T) {
	p := expectJSON(t, `{"type":"Point","coordinates":[1,2,3]}`, nil)
	if p.Center() != P(1, 2) {
		t.Fatalf("expected '%v', got '%v'", P(1, 2), p.Center())
	}
	expectJSON(t, `{"type":"Point","coordinates":[1,null]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Point","coordinates":[1,2],"bbox":null}`, errBBoxInvalid)
	expectJSON(t, `{"type":"Point"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"Point","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Point","coordinates":[1,2,3,4,5]}`, nil)
	expectJSON(t, `{"type":"Point","coordinates":[1]}`, errCoordinatesInvalid)
}
