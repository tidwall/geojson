package geojson

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/tidwall/geojson/geos"
	"github.com/tidwall/pretty"
)

func init() {
	seed := time.Now().UnixNano()
	println(seed)
	rand.Seed(seed)
}

func R(minX, minY, maxX, maxY float64) geos.Rect {
	return geos.Rect{
		Min: geos.Point{X: minX, Y: minY},
		Max: geos.Point{X: maxX, Y: maxY},
	}
}
func P(x, y float64) geos.Point {
	return geos.Point{X: x, Y: y}
}
func PO(x, y float64) *Point {
	return NewPoint(x, y)
}

func expectJSON(t testing.TB, data string, exp error) Object {
	if t != nil {
		t.Helper()
	}
	obj, err := Parse(data)
	if err != exp {
		if t == nil {
			panic(fmt.Sprintf("expected '%v', got '%v'", exp, err))
		} else {
			t.Fatalf("expected '%v', got '%v'", exp, err)
		}
	}
	return obj
}

func expect(t testing.TB, what bool) {
	t.Helper()
	if !what {
		t.Fatal("expection failure")
	}
}

func cleanJSON(data string) string {
	var v interface{}
	if err := json.Unmarshal([]byte(data), &v); err != nil {
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
