package geojson

import "testing"

func TestMultiLineString(t *testing.T) {
	expectJSON(t, `{"type":"MultiLineString","coordinates":[[[1,2,3]]]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiLineString","coordinates":[[[1,2]]],"bbox":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiLineString"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"MultiLineString","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"MultiLineString","coordinates":[1,null]}`, errCoordinatesInvalid)
}

func TestMultiLineStringValid(t *testing.T) {
	json := `{"type":"MultiLineString","coordinates":[
		[[10,10],[120,190]],
		[[50,50],[100,100]]
	]}`
	expectJSON(t, json, nil)
	expectJSONOpts(t, json, errCoordinatesInvalid, &ParseOptions{RequireValid: true})
}

func TestMultiLineStringPoly(t *testing.T) {
	p := expectJSON(t, `{"type":"MultiLineString","coordinates":[
		[[10,10],[20,20]],
		[[50,50],[100,100]]
	]}`, nil)
	expect(t, p.Intersects(PO(15, 15)))
	expect(t, p.Contains(PO(15, 15)))
	expect(t, p.Contains(PO(70, 70)))
	expect(t, !p.Contains(PO(40, 40)))
}
