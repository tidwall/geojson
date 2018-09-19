package geojson

import "testing"

func TestMultiLineString(t *testing.T) {
	p := expectJSON(t, `{"type":"MultiLineString","coordinates":[[[1,2,3]]]}`, nil)
	if p.Center() != P(1, 2) {
		t.Fatalf("expected '%v', got '%v'", P(1, 2), p.Center())
	}
	expectJSON(t, `{"type":"MultiLineString","coordinates":[[[1,2]]],"bbox":null}`, errBBoxInvalid)
	expectJSON(t, `{"type":"MultiLineString"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"MultiLineString","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiLineString","coordinates":[1,null]}`, errCoordinatesInvalid)
}
