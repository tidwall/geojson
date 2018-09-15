package geojson

import (
	"testing"

	"github.com/tidwall/gjson"
)

// https://gist.github.com/tidwall/5524c468fa4212b89e9c3532a5b1f355
var detectJSON = `{"type":"FeatureCollection","features":[
	{"type":"Feature","properties":{"id":"point1"},"geometry":{"type":"Point","coordinates":[-73.98549556732178,40.72198994979771]}},
	{"type":"Feature","properties":{"id":"polygon1"},"geometry":{"type":"Polygon","coordinates":[[[-74.0035629272461,40.71994085251552],[-73.98914337158203,40.71994085251552],[-73.98914337158203,40.72755146730012],[-74.0035629272461,40.72755146730012],[-74.0035629272461,40.71994085251552]]]}},
	{"type":"Feature","properties":{"id":"linestring1"},"geometry":{"type":"LineString","coordinates":[[-73.98382186889648,40.73652697126574],[-73.98821868896482,40.73652697126574],[-73.97420883178711,40.72943772441242]]}},
	{"type":"Feature","properties":{"id":"linestring2"},"geometry":{"type":"LineString","coordinates":[[-73.98146152496338,40.72716120053256],[-73.99098873138428,40.724754504892424]]}},
	{"type":"Feature","properties":{"id":"linestring3"},"geometry":{"type":"LineString","coordinates":[[-73.98386478424072,40.72696606629052],[-73.98090362548828,40.72501469240076],[-73.97837162017821,40.72621804639551]]}},
	{"type":"Feature","properties":{"id":"polygon2","fill":"#0433ff"},"geometry":{"type":"Polygon","coordinates":[[[-73.98661136627197,40.72540497175607],[-73.99064540863037,40.71938791069558],[-73.98807048797607,40.71779411151556],[-73.97571086883545,40.72338850378556],[-73.98017406463623,40.72960033028089],[-73.98661136627197,40.72540497175607]]]}},
	{"type":"Feature","properties":{"id":"polygon3"},"geometry":{"type":"Polygon","coordinates":[[[-73.98352146148682,40.72550254123727],[-73.98579597473145,40.72088409560772],[-73.97914409637451,40.72251034541217],[-73.98017406463623,40.72599038649773],[-73.98352146148682,40.72550254123727]],[[-73.98300647735596,40.72439674540761],[-73.98111820220947,40.72446179272971],[-73.98141860961913,40.7221525738643],[-73.98300647735596,40.72439674540761]]]}},
	{"type":"Feature","properties":{"id":"multipoint1","marker-color":"#941751"},"geometry":{"type":"MultiPoint","coordinates":[[-73.98957252502441,40.72049378974239],[-73.9897871017456,40.720233584560724],[-73.9897871017456,40.721664700472566],[-73.99085998535155,40.720916620993194],[-73.9912462234497,40.720331161623065]]}},
	{"type":"Feature","properties":{"id":"multilinestring1","stroke":"#941751"},"geometry":{"type":"MultiLineString","coordinates":[[[-73.98442268371582,40.72459188718318],[-73.98463726043701,40.72384384060296],[-73.98382186889648,40.72355112443509]],[[-73.9850664138794,40.72358364851732],[-73.98476600646973,40.72485207532725],[-73.9854097366333,40.72491712220435]]]}},
	{"type":"Feature","properties":{"id":"multipolygon1","fill":"#941751"},"geometry":{"type":"MultiPolygon","coordinates":[[[[-73.98021697998047,40.72429917430525],[-73.97892951965332,40.7250472157678],[-73.98000240325926,40.72524235563617],[-73.98021697998047,40.72429917430525]]],[[[-73.97901535034178,40.72452683998823],[-73.9788007736206,40.72345355209305],[-73.97764205932617,40.72410403167144],[-73.97901535034178,40.72452683998823]]]]}},
	{"type":"Feature","properties":{"id":"point2"},"geometry":{"type":"Point","coordinates":[-73.98326396942139,40.723681220668624]}},
	{"type":"Feature","properties":{"id":"point3"},"geometry":{"type":"Point","coordinates":[-73.98196396942139,40.723681220668624]}},
	{"type":"Feature","properties":{"id":"point4"},"geometry":{"type":"Point","coordinates":[-73.9785396942139,40.7238220668624]}}
]}`

func getByID(id string) Object {
	var r gjson.Result
	gjson.Get(detectJSON, "features").ForEach(func(_, v gjson.Result) bool {
		if v.Get("properties.id").String() == id {
			r = v.Get("geometry")
			return false
		}
		return true
	})
	if !r.Exists() {
		panic("not found '" + id + "'")
	}
	o, err := ObjectJSON(r.String())
	if err != nil {
		panic(err)
	}
	if p, ok := o.(SimplePoint); ok {
		o = Point{Coordinates: Position{X: p.X, Y: p.Y}}
	}
	return o
}

func toSimplePoint(p Object) Object {
	return SimplePoint{X: p.(Point).Coordinates.X, Y: p.(Point).Coordinates.Y}
}

// Basic geometry detections
// Point -> Point
// Point -> MultiPoint
// Point -> LineString
// Point -> MultiLineString
// Point -> Polygon
// Point -> MultiPolygon
// MultiPoint -> Point
// MultiPoint -> MultiPoint
// MultiPoint -> LineString
// MultiPoint -> MultiLineString
// MultiPoint -> Polygon
// MultiPoint -> MultiPolygon
// LineString -> Point
// LineString -> MultiPoint
// LineString -> LineString
// LineString -> MultiLineString
// LineString -> Polygon
// LineString -> MultiPolygon
// MultiLineString -> Point
// MultiLineString -> MultiPoint
// MultiLineString -> LineString
// MultiLineString -> MultiLineString
// MultiLineString -> Polygon
// MultiLineString -> MultiPolygon
// Polygon -> Point
// Polygon -> MultiPoint
// Polygon -> LineString
// Polygon -> MultiLineString
// Polygon -> Polygon
// Polygon -> MultiPolygon
// MultiPolygon -> Point
// MultiPolygon -> MultiPoint
// MultiPolygon -> LineString
// MultiPolygon -> MultiLineString
// MultiPolygon -> Polygon
// MultiPolygon -> MultiPolygon

func TestDetectSimplePointSimplePoint(t *testing.T) {
	p1 := toSimplePoint(getByID("point1"))
	p2 := toSimplePoint(getByID("point2"))
	if p1.Intersects(p2) {
		t.Fatal("expected false")
	}
	if !p1.Intersects(p1) {
		t.Fatal("expected true")
	}
	if p1.Within(p2) {
		t.Fatal("expected false")
	}
	if !p1.Within(p1) {
		t.Fatal("expected true")
	}
}
func TestDetectPointPoint(t *testing.T) {
	p1 := getByID("point1")
	p2 := getByID("point2")
	if p1.Intersects(p2) {
		t.Fatal("expected false")
	}
	if !p1.Intersects(p1) {
		t.Fatal("expected true")
	}
	if p1.Within(p2) {
		t.Fatal("expected false")
	}
	if !p1.Within(p1) {
		t.Fatal("expected true")
	}
}
func TestDetectPointMultPoint(t *testing.T) {
	p1 := getByID("point1")
	mp1 := getByID("multipoint1")
	if p1.Intersects(mp1) {
		t.Fatal("expected false")
	}
	if p1.Within(mp1) {
		t.Fatal("expected false")
	}
	pp := Point{Coordinates: mp1.(MultiPoint).Coordinates[0]}
	if !pp.Intersects(mp1) {
		t.Fatal("expected true")
	}
	if !pp.Within(mp1) {
		t.Fatal("expected true")
	}
}
func TestDetectPointLineString(t *testing.T) {
	p1 := getByID("point1")
	ls1 := getByID("linestring1")
	if p1.Intersects(ls1) {
		t.Fatal("expected false")
	}
	if p1.Within(ls1) {
		t.Fatal("expected false")
	}
	pp := Point{Coordinates: ls1.(LineString).Coordinates[0]}
	if !pp.Intersects(ls1) {
		t.Fatal("expected true")
	}
	if !pp.Within(ls1) {
		t.Fatal("expected true")
	}
}
func TestDetectPointMultiLineString(t *testing.T) {
	p1 := getByID("point1")
	mls1 := getByID("multilinestring1")
	if p1.Intersects(mls1) {
		t.Fatal("expected false")
	}
	if p1.Within(mls1) {
		t.Fatal("expected false")
	}
	pp := Point{Coordinates: mls1.(MultiLineString).Coordinates[0][1]}
	if !pp.Intersects(mls1) {
		t.Fatal("expected true")
	}
	if !pp.Within(mls1) {
		t.Fatal("expected true")
	}
}
func TestDetectPointPolygon(t *testing.T) {
	p1 := getByID("point1")
	pl3 := getByID("polygon3")
	if p1.Intersects(pl3) {
		t.Fatal("expected false")
	}
	if p1.Within(pl3) {
		t.Fatal("expected false")
	}
	p2 := getByID("point2")
	if !p2.Intersects(pl3) {
		t.Fatal("expected true")
	}
	if !p2.Within(pl3) {
		t.Fatal("expected true")
	}
	p3 := getByID("point3")
	if p3.Intersects(pl3) {
		t.Fatal("expected false")
	}
	if p3.Within(pl3) {
		t.Fatal("expected false")
	}
}

func TestDetectPointMultiPolygon(t *testing.T) {
	p3 := getByID("point3")
	p4 := getByID("point4")
	mp1 := getByID("multipolygon1")
	if p3.Intersects(mp1) {
		t.Fatal("expected false")
	}
	if p3.Within(mp1) {
		t.Fatal("expected false")
	}
	if !p4.Intersects(mp1) {
		t.Fatal("expected true")
	}
	if !p4.Within(mp1) {
		t.Fatal("expected true")
	}
}

func TestDetectMultiPointPoint(t *testing.T) {
	p1 := getByID("point1")
	mp1 := getByID("multipoint1")
	if mp1.Intersects(p1) {
		t.Fatal("expected false")
	}
	if mp1.Within(p1) {
		t.Fatal("expected false")
	}
	pp := Point{Coordinates: mp1.(MultiPoint).Coordinates[0]}
	if !mp1.Intersects(pp) {
		t.Fatal("expected true")
	}
	if mp1.Within(pp) {
		t.Fatal("expected false")
	}
}

func TestDetectMultiPointPolygon(t *testing.T) {
	mp1 := getByID("multipoint1")
	pl1 := getByID("polygon1")
	pl2 := getByID("polygon2")
	pl3 := getByID("polygon3")
	if !mp1.Intersects(pl1) {
		t.Fatal("expected true")
	}
	if !mp1.Within(pl1) {
		t.Fatal("expected true")
	}
	if !mp1.Intersects(pl2) {
		t.Fatal("expected true")
	}
	if mp1.Within(pl2) {
		t.Fatal("expected false")
	}
	if mp1.Intersects(pl3) {
		t.Fatal("expected false")
	}
	if mp1.Within(pl3) {
		t.Fatal("expected false")
	}
}

func TestDetectLineStringLineString(t *testing.T) {
	ls1 := getByID("linestring1")
	ls2 := getByID("linestring2")
	ls3 := getByID("linestring3")
	if ls1.Intersects(ls2) {
		t.Fatal("expected false")
	}
	if ls2.Intersects(ls1) {
		t.Fatal("expected false")
	}
	if !ls2.Intersects(ls3) {
		t.Fatal("expected true")
	}
	if !ls2.Intersects(ls3) {
		t.Fatal("expected true")
	}
}
func TestDetectLineStringPolygon(t *testing.T) {
	ls1 := getByID("linestring1")
	ls2 := getByID("linestring2")
	ls3 := getByID("linestring3")
	pl1 := getByID("polygon1")
	pl2 := getByID("polygon2")
	pl3 := getByID("polygon3")

	if ls1.Intersects(pl1) {
		t.Fatal("expected false")
	}
	if pl1.Intersects(ls1) {
		t.Fatal("expected false")
	}
	if ls1.Intersects(pl3) {
		t.Fatal("expected false")
	}
	if pl3.Intersects(ls1) {
		t.Fatal("expected false")
	}
	if !ls2.Intersects(pl1) {
		t.Fatal("expected true")
	}
	if !pl2.Intersects(ls2) {
		t.Fatal("expected true")
	}
	if ls2.Within(pl1) {
		t.Fatal("expected false")
	}
	if ls2.Intersects(pl3) {
		t.Fatal("expected false")
	}
	if pl3.Intersects(ls2) {
		t.Fatal("expected false")
	}
	if !ls3.Intersects(pl2) {
		t.Fatal("expected true")
	}
	if !pl2.Intersects(ls3) {
		t.Fatal("expected true")
	}
	if !ls3.Within(pl2) {
		t.Fatal("expected true")
	}
	if ls2.Intersects(pl3) {
		t.Fatal("expected false")
	}
	if pl3.Intersects(ls2) {
		t.Fatal("expected false")
	}
	if ls2.Intersects(pl3) {
		t.Fatal("expected false")
	}
	if pl3.Intersects(ls2) {
		t.Fatal("expected false")
	}
	if !ls3.Intersects(pl3) {
		t.Fatal("expected true")
	}
	if !pl3.Intersects(ls3) {
		t.Fatal("expected true")
	}
}
func TestDetectMultiLineStringPolygon(t *testing.T) {
	mls1 := getByID("multilinestring1")
	pl1 := getByID("polygon1")
	pl2 := getByID("polygon2")
	pl3 := getByID("polygon3")
	if mls1.Intersects(pl1) {
		t.Fatal("expected false")
	}
	if mls1.Within(pl1) {
		t.Fatal("expected false")
	}
	if !mls1.Intersects(pl2) {
		t.Fatal("expected true")
	}
	if !mls1.Within(pl2) {
		t.Fatal("expected true")
	}
	if !mls1.Intersects(pl3) {
		t.Fatal("expected true")
	}
	if mls1.Within(pl3) {
		t.Fatal("expected false")
	}
}
