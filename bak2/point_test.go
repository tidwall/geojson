package geojson

import "testing"

func TestPoint(t *testing.T) {
	mustParseGeoJSON(
		`{"type":"Point","coordinates":[10,11]}`,
	)
	mustParseGeoJSON(
		`{"type":"Point","coordinates":[10,11,12]}`,
	)
	mustParseGeoJSON(
		`{"type":"Point","coordinates":[10,11,13]}`,
	)
	mustParseGeoJSON(
		`{"type":"Point","coordinates":[10,11],"bbox":[1,2,3,4]}`,
	)
	mustParseGeoJSON(
		`{"type":"Point","coordinates":[10,11,12],"bbox":[1,2,3,4]}`,
	)
	mustParseGeoJSON(
		`{"type":"Point","coordinates":[10,11,12,13],"bbox":[1,2,3,4]}`,
	)
	mustParseGeoJSON(
		`{"type":"Point","coordinates":[10,11,12,13],"bbox":[1,2,3,4,5,6]}`,
	)
	mustParseGeoJSON(
		`{"type":"Point","coordinates":[10,11,12,13],"bbox":[1,2,3,4,5,6,7,8]}`,
	)
}
