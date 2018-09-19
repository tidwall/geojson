package geojson

import "testing"

func TestFeatureCollection(t *testing.T) {
	var g FeatureCollection
	g = mustParseGeoJSON(
		`{"type":"FeatureCollection","features":[
			{"type":"Point","coordinates":[10,11]},
			{"type":"Point","coordinates":[12,13]}
		]}`,
	).(FeatureCollection)
	expect(g.Rect(), R(10, 11, 12, 13))
	g = mustParseGeoJSON(
		`{"type":"FeatureCollection","features":[]}`,
	).(FeatureCollection)
	expect(g.Rect(), R(0, 0, 0, 0))
	g = mustParseGeoJSON(
		`{"type":"FeatureCollection","features":[],"bbox":[10,11,12,13]}`,
	).(FeatureCollection)
	expect(g.Rect(), R(10, 11, 12, 13))
}
