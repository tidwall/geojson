package geojson

import (
	"testing"
)

func TestPolygonParse(t *testing.T) {
	json := `{"type":"Polygon","coordinates":[[[0,0],[10,0],[10,10],[0,10],[0,0]]]}`
	expectJSON(t, json, nil)
	json = `{"type":"Polygon","coordinates":[
		[[0,0],[10,0],[10,10],[0,10],[0,0]],
		[[2,2],[8,2],[8,8],[2,8],[2,2]]
	]}`
	g := expectJSON(t, json, nil)
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

func TestEmptyPolygon(t *testing.T) {
	p := NewPolygon(nil)
	expect(t, p.JSON() == `{"type":"Polygon","coordinates":[]}`)
}

// https://github.com/tidwall/tile38/issues/664
func TestIssue664(t *testing.T) {
	// original geojson from issue
	p1, err := Parse(`{"type":"Polygon","coordinates":[[
			[-0.104364362074306,51.515197601239528],
			[-0.100183436878063,51.511898267797733],
			[-0.095787000073766,51.509618100991439],
			[-0.097554195259807,51.505459855954911],
			[-0.106390171190011,51.504842793711688],
			[-0.115829579622766,51.507418302507524],
			[-0.115527863371491,51.511737318590235],
			[-0.108157366376052,51.513722319074922],
			[-0.104364362074306,51.515197601239528]
		],[
			[-0.108114264054441,51.51141541846934],
			[-0.10289888313954,51.511066690771599],
			[-0.102036836707325,51.507954848516555],
			[-0.108674594235381,51.507874367018005],
			[-0.108114264054441,51.51141541846934]
		]]}`, nil)
	if err != nil {
		t.Fatal(err)
	}

	// converted to right-hand rule
	p2, err := Parse(`{"type":"Polygon","coordinates":[[
			[-0.104364362074306,51.515197601239528],
			[-0.108157366376052,51.513722319074922],
			[-0.115527863371491,51.511737318590235],
			[-0.115829579622766,51.507418302507524],
			[-0.106390171190011,51.504842793711688],
			[-0.097554195259807,51.505459855954911],
			[-0.095787000073766,51.509618100991439],
			[-0.100183436878063,51.511898267797733],
			[-0.104364362074306,51.515197601239528]
		],[
			[-0.108114264054441,51.51141541846934],
			[-0.10289888313954,51.511066690771599],
			[-0.102036836707325,51.507954848516555],
			[-0.108674594235381,51.507874367018005],
			[-0.108114264054441,51.51141541846934]
		]]}`, nil)
	if err != nil {
		t.Fatal(err)
	}

	// input lines
	lines := []string{
		`{"type": "LineString","coordinates": [[-0.11203657532102,51.509805883746516],[-0.098071423119136,51.509618100991439]]}`,
		`{"type": "LineString","coordinates": [[-0.106778092084508,51.510235098565985],[-0.104148850466252,51.510208272758227]]}`,
		`{"type": "LineString","coordinates": [[-0.099235185802626,51.512327462904388],[-0.099235185802626,51.510798436879945]]}`,
		`{"type": "LineString","coordinates": [[-0.111741852049867,51.510945273705836],[-0.10999491495498,51.512082479942833]]}`,
	}

	// expected results
	expect := []bool{
		true,
		false,
		true,
		true,
	}

	for i, line := range lines {
		l, err := Parse(line, nil)
		if err != nil {
			t.Fatal(err)
		}
		t1 := p1.Intersects(l)
		t2 := l.Intersects(p1)
		if t1 != expect[i] || t2 != expect[i] {
			t.Fatalf("line %d: expected %v/%v, got %v/%v",
				i, expect[i], expect[i], t1, t2)
		}
	}

	for i, line := range lines {
		l, err := Parse(line, nil)
		if err != nil {
			t.Fatal(err)
		}
		t1 := p2.Intersects(l)
		t2 := l.Intersects(p2)
		if t1 != expect[i] || t2 != expect[i] {
			t.Fatalf("line %d: expected %v/%v, got %v/%v",
				i, expect[i], expect[i], t1, t2)
		}
	}
}

func TestIssue714(t *testing.T) {
	_, err := Parse(`{"type":"Polygon","coordinates":[[[0,0],[10,0],[0,10],[0,0]],[[0,0,0],[0,10,0],[10,0,0],[0,0,0]]]}`, nil)
	if err.Error() != "invalid coordinates" {
		t.Fatalf("expected '%v', got '%v'", "invalid coordinates", err)
	}
	_, err = Parse(`{"type":"Polygon","coordinates":[[[0,0],[10,0,1],[0,10],[0,0]]]}`, nil)
	if err.Error() != "invalid coordinates" {
		t.Fatalf("expected '%v', got '%v'", "invalid coordinates", err)
	}
	_, err = Parse(`{"type":"LineString","coordinates":[[0,0],[10,0,1],[0,10],[0,0]]}`, nil)
	if err.Error() != "invalid coordinates" {
		t.Fatalf("expected '%v', got '%v'", "invalid coordinates", err)
	}
}
