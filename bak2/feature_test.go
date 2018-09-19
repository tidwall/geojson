package geojson

import "testing"

func TestFeature(t *testing.T) {
	var g Feature
	g = mustParseGeoJSON(
		`{"type":"Feature","geometry":{"type":"Point","coordinates":[10,11]}}`,
	).(Feature)
	expect(g.Rect(), R(10, 11, 10, 11))
	expect(g.Center(), P(10, 11))
	expect(g.Rect().Center(), P(10, 11))
	expect(g.ID.Exists(), false)
	expect(g.Properties.Exists(), false)
	g = mustParseGeoJSON(
		`{"type":"Feature",
			"geometry":{"type":"Point","coordinates":[10,11]},
			"id":"123"}`,
	).(Feature)
	expect(g.ID.String(), "123")
	expect(g.Properties.Exists(), false)
	g = mustParseGeoJSON(
		`{"type":"Feature",
			"geometry":{"type":"Point","coordinates":[10,11]},
			"properties":{"hello":"world"},
			"id":"123"}`,
	).(Feature)
	expect(g.ID.String(), "123")
	expect(g.Properties.Get("hello"), "world")
}
