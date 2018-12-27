package geojson

import "testing"

func TestPolygonParse(t *testing.T) {
	json := `{"type":"Polygon","coordinates":[[[0,0],[10,0],[10,10],[0,10],[0,0]]]}`
	g := expectJSON(t, json, nil)
	json = `{"type":"Polygon","coordinates":[
		[[0,0],[10,0],[10,10],[0,10],[0,0]],
		[[2,2],[8,2],[8,8],[2,8],[2,2]]
	]}`
	g = expectJSON(t, json, nil)
	if g.Center() != P(5, 5) {
		t.Fatalf("expected '%v', got '%v'", P(5, 5), g.Center())
	}
	expectJSON(t, `{"type":"Polygon","coordinates":[[1,null]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Polygon","coordinates":[]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Polygon","coordinates":[[[0,0],[10,0],[5,10],[0,0]],[[1,1]]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Polygon","coordinates":[[[0,0],[10,0],[5,10],[0,0]]],"bbox":null}`, nil)
	expectJSON(t, `{"type":"Polygon"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"Polygon","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Polygon","coordinates":[null]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Polygon","coordinates":[[null]]}`, errCoordinatesInvalid)
	expectJSON(t,
		`{"type":"Polygon","coordinates":[[[0,0,0,0,0],[10,0],[5,10],[0,0]]]}`,
		`{"type":"Polygon","coordinates":[[[0,0,0,0],[10,0,0,0],[5,10,0,0],[0,0,0,0]]]}`)
	expectJSON(t,
		`{"type":"Polygon","coordinates":[[[0,0,0],[10,0,4,5],[5,10],[0,0]]]}`,
		`{"type":"Polygon","coordinates":[[[0,0,0],[10,0,4],[5,10,0],[0,0,0]]]}`)
}
func TestPolygonParseValid(t *testing.T) {
	json := `{"type":"Polygon","coordinates":[
		[[0,0],[190,0],[10,10],[0,10],[0,0]],
		[[2,2],[8,2],[8,8],[2,8],[2,2]]
	]}`
	expectJSON(t, json, nil)
	expectJSONOpts(t, json, errCoordinatesInvalid, &ParseOptions{RequireValid: true})
}

func TestPolygonVarious(t *testing.T) {
	var g = expectJSON(t, `{"type":"Polygon","coordinates":[[[0,0],[10,0],[10,10],[0,10],[0,0]]]}`, nil)
	expect(t, string(g.AppendJSON(nil)) == `{"type":"Polygon","coordinates":[[[0,0],[10,0],[10,10],[0,10],[0,0]]]}`)
	expect(t, g.Rect() == R(0, 0, 10, 10))
	expect(t, g.Center() == P(5, 5))
	expect(t, !g.Empty())
}

// func TestPolygonPoly(t *testing.T) {
// 	json := `{"type":"Polygon","coordinates":[[[0,0],[10,0],[10,10],[0,10],[0,0]]]}`
// 	g := expectJSON(t, json, nil)
// 	expect(t, g.Contains(PO(5, 5)))
// 	expect(t, g.Contains(RO(5, 5, 6, 6)))
// 	expect(t, g.Contains(expectJSON(t, `{"type":"LineString","coordinates":[
// 		[5,5],[5,6],[6,5]
// 	]}`, nil)))
// 	expect(t, g.Intersects(PO(5, 5)))
// 	expect(t, g.Intersects(RO(5, 5, 6, 6)))
// 	expect(t, g.Intersects(expectJSON(t, `{"type":"LineString","coordinates":[
// 		[5,5],[5,6],[6,5],[50,50]
// 	]}`, nil)))
// 	expect(t, g.Intersects(expectJSON(t, `{"type":"Polygon","coordinates":[[
// 		[5,5],[5,6],[6,5],[50,50],[5,5]
// 	]]}`, nil)))
// 	expect(t, !g.Contains(expectJSON(t, `{"type":"Polygon","coordinates":[[
// 		[5,5],[5,6],[6,5],[50,50],[5,5]
// 	]]}`, nil)))
// 	expect(t, g.Contains(expectJSON(t, `{"type":"Polygon","coordinates":[[
// 		[5,5],[5,6],[6,5],[5,5]
// 	]]}`, nil)))
// }
