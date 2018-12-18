package geojson

import "testing"

func TestLineStringParse(t *testing.T) {
	expectJSON(t, `{"type":"LineString","coordinates":[[1,2,3]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[[1,null]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[[1,2]],"bbox":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"LineString","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[[null]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[null]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[[1,2,3,4,5]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"LineString","coordinates":[[1]]}`, errCoordinatesInvalid)
	g := expectJSON(t, `{"type":"LineString","coordinates":[[3,4],[1,2]]}`, nil)
	expect(t, g.Rect() == R(1, 2, 3, 4))
	expectJSON(t, `{"type":"LineString","coordinates":[[3,4],[1,2]],"bbox":null}`, nil)
	expectJSON(t, `{"type":"LineString","coordinates":[[3,4],[1,2]],"bbox":[1,2,3,4]}`, nil)
}

func TestLineStringParseValid(t *testing.T) {
	json := `{"type":"LineString","coordinates":[[1,2],[-12,-190]]}`
	expectJSON(t, json, nil)
	expectJSONOpts(t, json, errDataInvalid, &ParseOptions{RequireValid: true})
}

func TestLineStringVarious(t *testing.T) {
	var g = expectJSON(t, `{"type":"LineString","coordinates":[[3,4],[1,2]]}`, nil)
	expect(t, string(g.AppendJSON(nil)) == `{"type":"LineString","coordinates":[[3,4],[1,2]]}`)
	expect(t, g.Rect() == R(1, 2, 3, 4))
	expect(t, g.Center() == P(2, 3))
	expect(t, !g.Empty())
	g = expectJSON(t, `{"type":"LineString","coordinates":[[3,4],[1,2]],"bbox":[1,2,3,4]}`, nil)
	expect(t, !g.Empty())
	expect(t, g.Rect() == R(1, 2, 3, 4))
	expect(t, g.Center() == R(1, 2, 3, 4).Center())
}

func TestLineStringValid(t *testing.T) {
	var g = expectJSON(t, `{"type":"LineString","coordinates":[[3,4],[1,2]]}`, nil)
	expect(t, g.Valid())
}

func TestLineStringInvalid(t *testing.T) {
	var g = expectJSON(t, `{"type":"LineString","coordinates":[[3,4],[1,2],[0, 190]]}`, nil)
	expect(t, !g.Valid())
}

// func TestLineStringPoly(t *testing.T) {
// 	ls := expectJSON(t, `{"type":"LineString","coordinates":[
// 		[10,10],[20,20],[20,10]
// 	]}`, nil)
// 	expect(t, ls.(*LineString).Contains(ls))
// 	expect(t, ls.Contains(PO(10, 10)))
// 	expect(t, ls.Contains(PO(15, 15)))
// 	expect(t, ls.Contains(PO(20, 20)))
// 	expect(t, ls.Contains(PO(20, 15)))
// 	expect(t, !ls.Contains(PO(12, 13)))
// 	expect(t, !ls.Contains(RO(10, 10, 20, 20)))
// 	expect(t, ls.Intersects(PO(10, 10)))
// 	expect(t, ls.Intersects(PO(15, 15)))
// 	expect(t, ls.Intersects(PO(20, 20)))
// 	expect(t, !ls.Intersects(PO(12, 13)))
// 	expect(t, ls.Intersects(RO(10, 10, 20, 20)))
// 	expect(t, ls.Intersects(
// 		expectJSON(t, `{"type":"Point","coordinates":[15,15,0]}`, nil),
// 	))
// 	expect(t, ls.Intersects(ls))
// 	lsb := expectJSON(t, `{"type":"LineString","coordinates":[
// 		[10,10],[20,20],[20,10]
// 	],"bbox":[10,10,20,20]}`, nil)
// 	expect(t, lsb.Contains(PO(12, 13)))
// 	expect(t, ls.Contains(PO(20, 20)))
// }
