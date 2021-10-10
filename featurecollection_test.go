package geojson

import (
	"testing"
)

func TestFeatureCollection(t *testing.T) {
	p := expectJSON(t, `{"type":"FeatureCollection","features":[{"type":"Point","coordinates":[1,2,3]}]}`, nil)
	if p.Center() != P(1, 2) {
		t.Fatalf("expected '%v', got '%v'", P(1, 2), p.Center())
	}
	expectJSON(t, `{"type":"FeatureCollection"}`, errFeaturesMissing)
	expectJSON(t, `{"type":"FeatureCollection","features":null}`, errFeaturesInvalid)
	expectJSON(t, `{"type":"FeatureCollection","features":[{"type":"Point","coordinates":[1,2,3]}],"bbox":null}`, nil)
	expectJSON(t, `{"type":"FeatureCollection","features":[{"type":"Point"}]}`, errCoordinatesMissing)
}

func TestFeatureCollectionPoly(t *testing.T) {
	p := expectJSON(t, `{"type":"FeatureCollection","features":[{"type":"Point","coordinates":[1,2]}]}`, nil)
	expect(t, p.Intersects(PO(1, 2)))
	expect(t, p.Contains(PO(1, 2)))
}

func TestFeatureCollectionValid(t *testing.T) {
	json := `{"type":"FeatureCollection","features":[{"type":"Point","coordinates":[1,200]}]}`
	expectJSON(t, json, nil)
	expectJSONOpts(t, json, errCoordinatesInvalid, &ParseOptions{RequireValid: true})
}

func TestForEach(t *testing.T) {
	json := `{"type":"FeatureCollection","features":[
		{"type":"Feature","id":"A","geometry":{"type":"Point","coordinates":[1,2]},"properties":{}},
		{"type":"Feature","id":"B","geometry":{"type":"Point","coordinates":[3,4]},"properties":{}},
		{"type":"Feature","id":"C","geometry":{"type":"Point","coordinates":[5,6]},"properties":{}},
		{"type":"Feature","id":"D","geometry":{"type":"Point","coordinates":[7,8]},"properties":{}}
	]}`

	g, _ := Parse(json, nil)
	objsA := g.(*FeatureCollection).Children()
	var objsB []Object
	g.ForEach(func(geom Object) bool {
		objsB = append(objsB, geom)
		return true
	})
	for i := 0; i < len(objsA) && i < len(objsB); i++ {
		expect(t, objsA[i].String() == objsB[i].String())
	}
}
