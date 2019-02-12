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
