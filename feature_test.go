package geojson

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestFeatureParse(t *testing.T) {
	p := expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]},"properties":{}}`, nil)
	expect(t, p.Center() == P(1, 2))
	expectJSON(t, `{"type":"Feature"}`, errGeometryMissing)
	expectJSON(t, `{"type":"Feature","geometry":null}`, errDataInvalid)
	expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]},"bbox":null,"properties":{}}`, nil)
	expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"id":[4,true],"properties":{}}`, nil)
	expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"id":"15","properties":{"a":"b"}}`, nil)
	expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]},"bbox":[1,2,3,4],"properties":{}}`, nil)
	expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2],"bbox":[1,2,3,4]},"id":[4,true],"properties":{}}`, nil)
}

func TestFeatureVarious(t *testing.T) {
	var g = expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]},"properties":{}}`, nil)
	expect(t, string(g.AppendJSON(nil)) == `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]},"properties":{}}`)
	expect(t, g.Rect() == R(1, 2, 1, 2))
	expect(t, g.Center() == P(1, 2))
	expect(t, !g.Empty())

	g = expectJSONOpts(t,
		`{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]},"bbox":[1,2,3,4],"properties":{}}`,
		nil, nil)
	expect(t, !g.Empty())
	expect(t, g.Rect() == R(1, 2, 1, 2))
	expect(t, g.Center() == P(1, 2))

}

func TestFeatureProperties(t *testing.T) {
	obj, err := Parse(`{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]}}`, nil)
	if err != nil {
		t.Fatal(err)
	}
	json := obj.JSON()
	if !gjson.Valid(json) {
		t.Fatal("invalid json")
	}
	if !gjson.Get(json, "properties").Exists() {
		t.Fatal("expected 'properties' member")
	}

	obj, err = Parse(`{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"properties":true}`, nil)
	if err != nil {
		t.Fatal(err)
	}
	json = obj.JSON()
	if !gjson.Valid(json) {
		t.Fatal("invalid json")
	}
	if gjson.Get(json, "properties").Type != gjson.True {
		t.Fatal("expected 'properties' member to be 'true'")
	}

	obj, err = Parse(`{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"id":{}}`, nil)
	if err != nil {
		t.Fatal(err)
	}
	json = obj.JSON()
	if !gjson.Valid(json) {
		t.Fatal("invalid json")
	}
	if !gjson.Get(json, "properties").Exists() {
		t.Fatal("expected 'properties' member")
	}
	if gjson.Get(json, "id").String() != "{}" {
		t.Fatal("expected 'id' member")
	}

}
