package geojson

import "testing"

func TestMultiPoint(t *testing.T) {
	p := expectJSON(t, `{"type":"MultiPoint","coordinates":[[1,2,3]]}`, nil)
	expect(t, p.Center() == P(1, 2))
	expectJSON(t, `{"type":"MultiPoint","coordinates":[1,null]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiPoint","coordinates":[[1,2]],"bbox":null}`, nil)
	expectJSON(t, `{"type":"MultiPoint"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"MultiPoint","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiPoint","coordinates":[[1,2,3],[4,5,6]],"bbox":[1,2,3,4]}`, nil)
}

// func TestMultiPointPoly(t *testing.T) {
// 	p := expectJSON(t, `{"type":"MultiPoint","coordinates":[[1,2],[2,2]]}`, nil)
// 	expect(t, p.Intersects(PO(1, 2)))
// 	expect(t, p.Contains(PO(1, 2)))
// 	expect(t, p.Contains(PO(2, 2)))
// 	expect(t, !p.Contains(PO(3, 2)))
// }
