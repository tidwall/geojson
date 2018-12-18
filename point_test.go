package geojson

import "testing"

func TestPointParse(t *testing.T) {
	p := expectJSON(t, `{"type":"Point","coordinates":[1,2,3]}`, nil)
	expect(t, p.Center() == P(1, 2))
	expectJSON(t, `{"type":"Point","coordinates":[1,null]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Point","coordinates":[1,2],"bbox":null}`, nil)
	expectJSON(t, `{"type":"Point"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"Point","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Point","coordinates":[1,2,3,4,5]}`, `{"type":"Point","coordinates":[1,2,3,4]}`)
	expectJSON(t, `{"type":"Point","coordinates":[1]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Point","coordinates":[1,2,3],"bbox":[1,2,3,4]}`, `{"type":"Point","coordinates":[1,2,3],"bbox":[1,2,3,4]}`)
}
func TestPointParseValid(t *testing.T) {
	json := `{"type":"Point","coordinates":[190,90]}`
	expectJSON(t, json, nil)
	expectJSONOpts(t, json, errCoordinatesInvalid, &ParseOptions{RequireValid: true})
}

func TestPointVarious(t *testing.T) {
	var g Object = PO(10, 20)
	expect(t, string(g.AppendJSON(nil)) == `{"type":"Point","coordinates":[10,20]}`)
	expect(t, g.Rect() == R(10, 20, 10, 20))
	expect(t, g.Center() == P(10, 20))
	expect(t, !g.Empty())
}

func TestPointValid(t *testing.T) {
	var g Object = PO(0, 20)
	expect(t, g.Valid())

	var g1 Object = PO(10, 20)
	expect(t, g1.Valid())
}

func TestPointInvalidLargeX(t *testing.T) {
	var g Object = PO(10, 91)
	expect(t, !g.Valid())
}

func TestPointInvalidLargeY(t *testing.T) {
	var g Object = PO(181, 20)
	expect(t, !g.Valid())
}

func TestPointValidLargeX(t *testing.T) {
	var g Object = PO(180, 20)
	expect(t, g.Valid())
}

func TestPointValidLargeY(t *testing.T) {
	var g Object = PO(180, 90)
	expect(t, g.Valid())
}

// func TestPointPoly(t *testing.T) {
// 	p := expectJSON(t, `{"type":"Point","coordinates":[15,15,0]}`, nil)
// 	expect(t, p.Within(PO(15, 15)))
// 	expect(t, p.Contains(PO(15, 15)))
// 	expect(t, p.Contains(RO(15, 15, 15, 15)))
// 	expect(t, !p.Contains(RO(10, 10, 15, 15)))
// 	expect(t, !p.Contains(PO(10, 10)))
// 	expect(t, p.Intersects(PO(15, 15)))
// 	expect(t, p.Intersects(RO(10, 10, 20, 20)))
// 	expect(t, p.Intersects(
// 		expectJSON(t, `{"type":"Point","coordinates":[15,15,10]}`, nil),
// 	))
// 	expect(t, !p.Intersects(
// 		expectJSON(t, `{"type":"Point","coordinates":[9,15,10]}`, nil),
// 	))
// 	expect(t, p.Intersects(
// 		expectJSON(t, `{"type":"Point","coordinates":[9,15,10],"bbox":[10,10,20,20]}`, nil),
// 	))
// 	expect(t, p.Intersects(
// 		expectJSON(t, `{"type":"LineString","coordinates":[
// 			[10,10],[20,20]
// 		]}`, nil),
// 	))
// 	expect(t, !p.Intersects(
// 		expectJSON(t, `{"type":"LineString","coordinates":[
// 			[9,10],[20,20]
// 		]}`, nil),
// 	))
// 	expect(t, p.Intersects(
// 		expectJSON(t, `{"type":"Polygon","coordinates":[
// 			[[9,9],[9,21],[21,21],[21,9],[9,9]]
// 		]}`, nil),
// 	))

// 	expect(t, !p.Intersects(
// 		expectJSON(t, `{"type":"Polygon","coordinates":[
// 			[[9,9],[9,21],[21,21],[21,9],[9,9]],
// 			[[9.5,9.5],[9.5,20.5],[20.5,20.5],[20.5,9.5],[9.5,9.5]]
// 		]}`, nil),
// 	))
// 	expect(t, p.Intersects(
// 		expectJSON(t, `{"type":"Feature","geometry":
// 			{"type":"Point","coordinates":[15,15,10]}
// 		}`, nil),
// 	))
// 	expect(t, p.Intersects(
// 		expectJSON(t, `{"type":"Feature","geometry":
// 			{"type":"Polygon","coordinates":[
// 				[[9,9],[9,21],[21,21],[21,9],[9,9]]
// 			]}
// 		}`, nil),
// 	))
// 	expect(t, !p.Intersects(
// 		expectJSON(t, `{"type":"Feature","geometry":
// 			{"type":"Polygon","coordinates":[
// 				[[9,9],[9,21],[21,21],[21,9],[9,9]],
// 				[[9.5,9.5],[9.5,20.5],[20.5,20.5],[20.5,9.5],[9.5,9.5]]
// 			]}
// 		}`, nil),
// 	))
// 	expect(t, !expectJSON(t,
// 		`{"type":"Point","coordinates":[15,15],"bbox":[10,10,15,15]}`, nil,
// 	).Contains(PO(7, 7)))
// 	expect(t, expectJSON(t,
// 		`{"type":"Point","coordinates":[15,15],"bbox":[10,10,15,15]}`, nil,
// 	).Contains(PO(12, 12)))

// 	expect(t, !expectJSON(t,
// 		`{"type":"Point","coordinates":[15,15],"bbox":[10,10,15,15]}`, nil,
// 	).Intersects(PO(7, 7)))
// 	expect(t, expectJSON(t,
// 		`{"type":"Point","coordinates":[15,15],"bbox":[10,10,15,15]}`, nil,
// 	).Intersects(PO(12, 12)))
// }
