package geojson

import "testing"

func TestSimplePointNotSimple(t *testing.T) {
	p := expectJSONOpts(t, `{"type":"Point","coordinates":[1,2,3]}`, nil, &ParseOptions{AllowSimplePoints: true})
	expect(t, p.Center() == P(1, 2))
	expectJSONOpts(t, `{"type":"Point","coordinates":[1,null]}`, errCoordinatesInvalid, &ParseOptions{AllowSimplePoints: true})
	expectJSONOpts(t, `{"type":"Point","coordinates":[1,2],"bbox":null}`, nil, &ParseOptions{AllowSimplePoints: true})
	expectJSONOpts(t, `{"type":"Point"}`, errCoordinatesMissing, &ParseOptions{AllowSimplePoints: true})
	expectJSONOpts(t, `{"type":"Point","coordinates":null}`, errCoordinatesInvalid, &ParseOptions{AllowSimplePoints: true})
	expectJSONOpts(t, `{"type":"Point","coordinates":[1,2,3,4,5]}`, `{"type":"Point","coordinates":[1,2,3,4]}`, &ParseOptions{AllowSimplePoints: true})
	expectJSONOpts(t, `{"type":"Point","coordinates":[1]}`, errCoordinatesInvalid, &ParseOptions{AllowSimplePoints: true})
	expectJSONOpts(t, `{"type":"Point","coordinates":[1,2,3],"bbox":[1,2,3,4]}`, `{"type":"Point","coordinates":[1,2,3],"bbox":[1,2,3,4]}`, &ParseOptions{AllowSimplePoints: true})
}

func TestSimplePointParseValid(t *testing.T) {
	json := `{"type":"Point","coordinates":[190,90]}`
	p := expectJSONOpts(t, json, nil, &ParseOptions{AllowSimplePoints: true})
	expect(t, !p.(*SimplePoint).Empty())
	p = expectJSONOpts(t, json, nil, &ParseOptions{AllowSimplePoints: false})
	expect(t, !p.(*Point).Empty())
	p = expectJSONOpts(t, json, errCoordinatesInvalid, &ParseOptions{RequireValid: true, AllowSimplePoints: true})
	expect(t, p == nil)
}

func TestSimplePointVarious(t *testing.T) {
	var g Object = PO(10, 20)
	expect(t, string(g.AppendJSON(nil)) == `{"type":"Point","coordinates":[10,20]}`)
	expect(t, g.Rect() == R(10, 20, 10, 20))
	expect(t, g.Center() == P(10, 20))
	expect(t, !g.Empty())
}

func TestSimplePointValid(t *testing.T) {
	var g Object = PO(0, 20)
	expect(t, g.Valid())

	var g1 Object = PO(10, 20)
	expect(t, g1.Valid())
}

func TestSimplePointInvalidLargeX(t *testing.T) {
	var g Object = PO(10, 91)
	expect(t, !g.Valid())
}

func TestSimplePointInvalidLargeY(t *testing.T) {
	var g Object = PO(181, 20)
	expect(t, !g.Valid())
}

func TestSimplePointValidLargeX(t *testing.T) {
	var g Object = PO(180, 20)
	expect(t, g.Valid())
}

func TestSimplePointValidLargeY(t *testing.T) {
	var g Object = PO(180, 90)
	expect(t, g.Valid())
}
