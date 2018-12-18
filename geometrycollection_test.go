package geojson

import "testing"

func TestGeometryCollection(t *testing.T) {
	p := expectJSON(t, `{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,2,3]}]}`, nil)
	expect(t, p.Center() == P(1, 2))
	expectJSON(t, `{"type":"GeometryCollection"}`, errGeometriesMissing)
	expectJSON(t, `{"type":"GeometryCollection","geometries":null}`, errGeometriesInvalid)
	expectJSON(t, `{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,2,3]}],"bbox":null}`, nil)
	expectJSON(t, `{"type":"GeometryCollection","geometries":[{"type":"Point"}]}`, errCoordinatesMissing)
}

func TestGeometryCollectionPoly(t *testing.T) {
	p := expectJSON(t, `{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,2]}]}`, nil)
	expect(t, p.Intersects(PO(1, 2)))
	expect(t, p.Contains(PO(1, 2)))
}

func TestGeometryCollectionValid(t *testing.T) {
	json := `{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,200]}]}`
	expectJSON(t, json, nil)
	expectJSONOpts(t, json, errCoordinatesInvalid, &ParseOptions{RequireValid: true})
}
