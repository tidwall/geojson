package geojson

import "testing"

func TestGeometryCollection(t *testing.T) {
	var g GeometryCollection
	g = mustParseGeoJSON(
		`{"type":"GeometryCollection","geometries":[
			{"type":"Point","coordinates":[10,11]},
			{"type":"Point","coordinates":[12,13]}
		]}`,
	).(GeometryCollection)
	expect(g.Rect(), R(10, 11, 12, 13))
	g = mustParseGeoJSON(
		`{"type":"GeometryCollection","geometries":[]}`,
	).(GeometryCollection)
	expect(g.Rect(), R(0, 0, 0, 0))
	g = mustParseGeoJSON(
		`{"type":"GeometryCollection","geometries":[],"bbox":[10,11,12,13]}`,
	).(GeometryCollection)
	expect(g.Rect(), R(10, 11, 12, 13))
}
