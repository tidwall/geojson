package geojson

import "testing"

func TestGeometryCollection(t *testing.T) {
	p := expectJSON(t, `{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,2,3]}]}`, nil)
	if p.Center() != P(1, 2) {
		t.Fatalf("expected '%v', got '%v'", P(1, 2), p.Center())
	}
	expectJSON(t, `{"type":"GeometryCollection"}`, errGeometriesMissing)
	expectJSON(t, `{"type":"GeometryCollection","geometries":null}`, errGeometriesInvalid)
	expectJSON(t, `{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,2,3]}],"bbox":null}`, errBBoxInvalid)
	expectJSON(t, `{"type":"GeometryCollection","geometries":[{"type":"Point"}]}`, errCoordinatesMissing)
}

func TestGeometryCollectionPoly(t *testing.T) {
	p := expectJSON(t, `{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,2]}]}`, nil)
	expect(t, p.Intersects(P(1, 2)))
	expect(t, p.Contains(P(1, 2)))
}
