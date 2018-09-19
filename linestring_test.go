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
