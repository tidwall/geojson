package geojson

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/pretty"
)

func init() {
	seed := time.Now().UnixNano()
	println(seed)
	rand.Seed(seed)
}

func R(minX, minY, maxX, maxY float64) geometry.Rect {
	return geometry.Rect{
		Min: geometry.Point{X: minX, Y: minY},
		Max: geometry.Point{X: maxX, Y: maxY},
	}
}
func P(x, y float64) geometry.Point {
	return geometry.Point{X: x, Y: y}
}

func PO(x, y float64) *Point {
	return NewPoint(P(x, y))
}

func MPO(points []geometry.Point) *MultiPoint {
	return NewMultiPoint(points)
}

func RO(minX, minY, maxX, maxY float64) *Rect {
	return NewRect(R(minX, minY, maxX, maxY))
}

func LO(points []geometry.Point) *LineString {
	return NewLineString(geometry.NewLine(points, nil))
}

func L(points []geometry.Point) *geometry.Line {
	return geometry.NewLine(points, geometry.DefaultIndexOptions)
}

func MLO(lines []*geometry.Line) *MultiLineString {
	return NewMultiLineString(lines)
}

func PPO(exterior []geometry.Point, holes [][]geometry.Point) *Polygon {
	return NewPolygon(geometry.NewPoly(exterior, holes, nil))
}

func expectJSON(t testing.TB, data string, expect interface{}) Object {
	if t != nil {
		t.Helper()
	}
	return expectJSONOpts(t, data, expect, nil)
}
func expectJSONOpts(t testing.TB, data string, expect interface{}, opts *ParseOptions) Object {
	if t != nil {
		t.Helper()
	}
	var exerr error
	var exstr string
	switch expect := expect.(type) {
	case string:
		exstr = expect
	case error:
		exerr = expect
	case nil:
		exstr = data
	}
	obj, err := Parse(data, opts)
	if err != exerr {
		if t == nil {
			panic(fmt.Sprintf("expected '%v', got '%v'", exerr, err))
		} else {
			t.Fatalf("expected '%v', got '%v'", exerr, err)
		}
	}
	if exstr != "" {
		if cleanJSON(exstr) != cleanJSON(string(obj.AppendJSON(nil))) {
			if t == nil {
				panic("json mismatch")
			} else {
				t.Fatal("json mismatch")
			}
		}
	}
	return obj
}

func expect(t testing.TB, what bool) {
	if t != nil {
		t.Helper()
	}
	if !what {
		if t == nil {
			panic("expectation failure")
		} else {
			t.Fatal("expectation failure")
		}
	}
}

func cleanJSON(data string) string {
	var v interface{}
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		println(string(data))
		panic(err)
	}
	dst, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	opts := *pretty.DefaultOptions
	opts.Width = 99999999
	return string(pretty.PrettyOptions(dst, &opts))
}

func TestBinary(t *testing.T) {
	jsons := []string{
		`G:{"type":"Point","coordinates":[10,20]}`,
		`G:{"type":"Point","coordinates":[10,20,30]}`,
		`G:{"type":"Point","coordinates":[10,20],"a":[1,2,3],"b":{"c":"d"}}`,
		`G:{"type":"Point","coordinates":[10,20,30],"a":[1,2,3],"b":{"c":"d"}}`,
		`G:{"type":"LineString","coordinates":[[10,20],[30,40]]}`,
		`G:{"type":"LineString","coordinates":[[10,20],[30,40]],"a":[1,2,3],"b":{"c":"d"}}`,
		`G:{"type":"LineString","coordinates":[[10,20,33],[30,40,44]],"a":[1,2,3],"b":{"c":"d"}}`,
		`G:{"type":"Polygon","coordinates":[[[0,0],[190,0],[10,10],[0,10],[0,0]],[[2,2],[8,2],[8,8],[2,8],[2,2]]]}`,
		`G:{"type":"Polygon","coordinates":[[[0,0],[190,0],[10,10],[0,10],[0,0]],[[2,2],[8,2],[8,8],[2,8],[2,2]]],"a":[1,2,3],"b":{"c":"d"}}`,
		`G:{"type":"Polygon","coordinates":[[[0,0,0],[190,0,1],[10,10,2],[0,10,3],[0,0,4]],[[2,2,5],[8,2,6],[8,8,7],[2,8,8],[2,2,9]],[[3,3,5],[9,3,6],[9,9,7],[3,9,8],[3,3,9]]],"a":[1,2,3],"b":{"c":"d"}}`,
		`G:{"type":"MultiPoint","coordinates":[]}`,
		`G:{"type":"MultiPoint","coordinates":[[1,2],[3,4]]}`,
		`G:{"type":"MultiPoint","coordinates":[[1,2,3]]}`,
		`G:{"type":"MultiPoint","coordinates":[[1,2,9],[3,4,10]]}`,
		`G:{"type":"MultiPoint","coordinates":[[1,2,9],[3,4,10]],"a":[1,2,3],"b":{"c":"d"}}`,
		`G:{"type":"MultiLineString","coordinates":[[[10,10],[20,20]],[[50,50],[100,100]]]}`,
		`G:{"type":"MultiLineString","coordinates":[[[10,10,5],[20,20,30]],[[50,50,6],[100,100,7]]]}`,
		`G:{"type":"MultiLineString","coordinates":[[[10,10,5],[20,20,30]],[[50,50,6],[100,100,7]]],"a":[1,2,3],"b":{"c":"d"}}`,
		`G:{"type":"MultiLineString","coordinates":[[[10,10],[20,20]],[[50,50],[100,100]]]}`,
		`G:{"type":"MultiPolygon","coordinates":[[[[0,0],[10,0],[10,10],[0,10],[0,0]],[[2,2],[8,2],[8,8],[2,8],[2,2]]],[[[0,0],[10,0],[10,10],[0,10],[0,0]],[[2,2],[8,2],[8,8],[2,8],[2,2]]]]}`,
		`G:{"type":"MultiPolygon","coordinates":[[[[0,0,1],[10,0,2],[10,10,3],[0,10,4],[0,0,5]],[[2,2,6],[8,2,7],[8,8,8],[2,8,9],[2,2,10]]],[[[0,0],[10,0],[10,10],[0,10],[0,0]],[[2,2],[8,2],[8,8],[2,8],[2,2]]]]}`,
		`G:{"type":"MultiPolygon","coordinates":[[[[0,0,1],[10,0,2],[10,10,3],[0,10,4],[0,0,5]],[[2,2,6],[8,2,7],[8,8,8],[2,8,9],[2,2,10]]],[[[0,0],[10,0],[10,10],[0,10],[0,0]],[[2,2],[8,2],[8,8],[2,8],[2,2]]]],"a":[1,2,3],"b":{"c":"d"}}`,
		`G:{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,2]},{"type":"LineString","coordinates":[[10,20],[30,40]]}]}`,
		`G:{"type":"GeometryCollection","geometries":[{"type":"Point","coordinates":[1,2]},{"type":"LineString","coordinates":[[10,20],[30,40]]}],"a":[1,2,3],"b":{"c":"d"}}`,
		`G:{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]},"bbox":[1,2,3,4],"properties":{}}`,
		`G:{"type":"FeatureCollection","features":[
			{"type":"Feature","id":"A","geometry":{"type":"Point","coordinates":[1,2]},"properties":{}},
			{"type":"Feature","id":"B","geometry":{"type":"Point","coordinates":[3,4]},"properties":{}},
			{"type":"Feature","id":"C","geometry":{"type":"Point","coordinates":[5,6]},"properties":{}},
			{"type":"Feature","id":"D","geometry":{"type":"Point","coordinates":[7,8]},"properties":{}}
		],"hello":"jello"}`,

		// Rect
		`X:{"type":"Polygon","coordinates":[[[10,20],[30,20],[30,40],[10,40],[10,20]]]}`,

		// SimplePoint
		`X:{"type":"Point","coordinates":[10,20]}`,

		// Circle
		`X:{"type":"Feature","geometry":{"type":"Point","coordinates":[-112,33]},"properties":{"type":"Circle","radius":5000,"radius_units":"m"}}`,
	}
	for _, json := range jsons {
		opts := *DefaultParseOptions
		if json[:2] == "X:" {
			opts.AllowRects = true
			opts.AllowSimplePoints = true
		}
		json = json[2:]
		obj, err := Parse(json, &opts)
		if err != nil {
			t.Fatal(err)
		}
		bin := obj.Binary()
		if bin == nil {
			t.Fatal("expected not nil")
		}
		obj2, n := ParseBinary(bin, nil)
		if n != len(bin) {
			t.Fatalf("expected %d, got %d", len(bin), n)
		}
		if obj2 == nil {
			t.Fatal("expected not nil")
		}
		json1 := string(pretty.Ugly([]byte(cleanJSON(json))))
		json2 := string(pretty.Ugly([]byte(cleanJSON(obj2.JSON()))))
		if json2 != json1 {
			println("GOT", json2)
			println("EXPECTED", json1)
			t.Fatal("json mismatch")
		}

	}
}
