package geojson

import "testing"

func TestMultiPolygon(t *testing.T) {
	json := `{"type":"MultiPolygon","coordinates":[
		[
			[[0,0],[10,0],[10,10],[0,10],[0,0]],
			[[2,2],[8,2],[8,8],[2,8],[2,2]]
		],[
			[[0,0],[10,0],[10,10],[0,10],[0,0]],
			[[2,2],[8,2],[8,8],[2,8],[2,2]]
		]
	]}`
	p := expectJSON(t, json, nil)
	if p.Center() != P(5, 5) {
		t.Fatalf("expected '%v', got '%v'", P(5, 5), p.Center())
	}
	if cleanJSON(string(p.AppendJSON(nil))) != cleanJSON(json) {
		t.Fatalf("expectect '%v', got '%v'", cleanJSON(json), cleanJSON(string(p.AppendJSON(nil))))
	}
	expectJSON(t, `{"type":"MultiPolygon","coordinates":[[[[0,0],[10,0],[5,10],[0,0]],[[1,1]]]],"bbox":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiPolygon"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"MultiPolygon","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiPolygon","coordinates":[1,null]}`, errCoordinatesInvalid)
}

func TestMultiPolygonPoly(t *testing.T) {
	p := expectJSON(t, `{"type":"MultiPolygon","coordinates":[
		[
			[[10,10],[20,20],[30,10],[10,10]]
		],[
			[[100,100],[200,200],[300,100],[100,100]]
		]
	]}`, nil)
	expect(t, p.Intersects(PO(15, 15)))
	expect(t, p.Contains(PO(15, 15)))
	expect(t, p.Contains(PO(150, 150)))
	expect(t, !p.Contains(PO(40, 40)))
	expect(t, p.Within(RO(-100, -100, 1000, 1000)))
}
