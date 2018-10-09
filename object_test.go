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

func RO(minX, minY, maxX, maxY float64) *Rect {
	return NewRect(minX, minY, maxX, maxY)
}

func expectJSON(t testing.TB, data string, expect interface{}) Object {
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
	obj, err := Parse(data, nil)
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
			panic("exception failure")
		} else {
			t.Fatal("expection failure")
		}
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
