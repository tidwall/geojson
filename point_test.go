package geojson

import (
	"math/rand"
	"testing"

	"github.com/tidwall/geojson/geometry"
)

func TestPointParse(t *testing.T) {
	p := expectJSON(t, `{"type":"Point","coordinates":[1,2,3]}`, nil)
	expect(t, p.Center() == P(1, 2))
	expectJSON(t, `{"type":"Point","coordinates":[1,null]}`, `{"type":"Point","coordinates":[1,null]}`)
	expectJSON(t, `{"type":"Point","coordinates":[1,"hello"]}`, errCoordinatesInvalid)
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

func BenchmarkPointValid(b *testing.B) {
	// debault b.N is too large
	points := make([]*Point, 1000000)
	for i := 0; i < 1000000; i++ {
		points[i] = NewPoint(geometry.Point{
			X: rand.Float64()*400 - 200, // some are out of bounds
			Y: rand.Float64()*200 - 100, // some are out of bounds
		})
	}
	b.ResetTimer()
	for i := 0; i < 1000000; i++ {
		points[i].Valid()
	}
}
