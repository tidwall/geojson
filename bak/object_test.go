package geojson

import "testing"

func TestObject(t *testing.T) {
	expectJSON(t, "", errDataInvalid)
	expectJSON(t, string([]byte{0, 1, 2, 3}), errDataInvalid)
	expectJSON(t, string([]byte{' ', 0}), errDataInvalid)
	expectJSON(t, `{}`, errTypeMissing)
	expectJSON(t, `{"}`, errDataInvalid)
	expectJSON(t, `{"type":null}`, errTypeInvalid)
	_, err := Parse(`{"type":"Square"}`)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestObjectVarious(t *testing.T) {
	c := expectJSON(t, `{"type":"GeometryCollection","geometries":[
		{"type":"Point","coordinates":[20,20]},
		{"type":"LineString","coordinates":[[10,10],[20,20],[30,10]]},
		{"type":"Point","coordinates":[30,30]}
	]}`, nil)
	p := expectJSON(t, `{"type":"LineString","coordinates":[[0,0],[20,30],[30,10]]}`, nil)
	expect(t, P(10, 10).Intersects(c))
	expect(t, P(20, 20).Intersects(c))
	expect(t, P(15, 10).Intersects(c))
	expect(t, !p.Contains(c))
	expect(t, !P(10, 10).Contains(c))
	expect(t, !P(20, 20).Contains(c))
	expect(t, !P(15, 15).Contains(c))
	expect(t, c.Intersects(P(15, 15)))

	c = expectJSON(t, `{"type":"GeometryCollection","geometries":[
		{"type":"Point","coordinates":[20,20]}
	],"bbox":[10,10,20,20]}`, nil)
	expect(t, c.Contains(P(15, 15)))
	expect(t, !c.Contains(P(50, 15)))
	expect(t, !c.Intersects(P(50, 15)))
	expect(t, c.Intersects(P(15, 15)))
}
