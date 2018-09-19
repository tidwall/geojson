package geojson

import "testing"

func TestFeatureCollection(t *testing.T) {
	p := expectJSON(t, `{"type":"FeatureCollection","features":[{"type":"Point","coordinates":[1,2,3]}]}`, nil)
	if p.Center() != P(1, 2) {
		t.Fatalf("expected '%v', got '%v'", P(1, 2), p.Center())
	}
	expectJSON(t, `{"type":"FeatureCollection"}`, errFeaturesMissing)
	expectJSON(t, `{"type":"FeatureCollection","features":null}`, errFeaturesInvalid)
	expectJSON(t, `{"type":"FeatureCollection","features":[{"type":"Point","coordinates":[1,2,3]}],"bbox":null}`, errBBoxInvalid)
	expectJSON(t, `{"type":"FeatureCollection","features":[{"type":"Point"}]}`, errCoordinatesMissing)
}
