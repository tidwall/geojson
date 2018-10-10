package geojson

import (
	"encoding/json"

	"github.com/tidwall/geojson/geometry"
)

// String is a not a geojson object, but just a string
type String struct {
	s string
}

// NewString ...
func NewString(s string) *String {
	return &String{s: s}
}

// forEach ...
func (g *String) forEach(iter func(geom Object) bool) bool {
	return iter(g)
}

// Empty ...
func (g *String) Empty() bool {
	return true
}

// Rect ...
func (g *String) Rect() geometry.Rect {
	return geometry.Rect{}
}

// Center ...
func (g *String) Center() geometry.Point {
	return geometry.Point{}
}

// AppendJSON ...
func (g *String) AppendJSON(dst []byte) []byte {
	data, _ := json.Marshal(g.s)
	return append(dst, data...)
}

// String ...
func (g *String) String() string {
	return string(g.s)
}

// Within ...
func (g *String) Within(obj Object) bool {
	return obj.Contains(g)
}

// Contains ...
func (g *String) Contains(obj Object) bool {
	return false
}

// Intersects ...
func (g *String) Intersects(obj Object) bool {
	return false
}

func (g *String) withinRect(rect geometry.Rect) bool {
	return false
}

func (g *String) withinPoint(point geometry.Point) bool {
	return false
}

func (g *String) withinLine(line *geometry.Line) bool {
	return false
}

func (g *String) withinPoly(poly *geometry.Poly) bool {
	return false
}

func (g *String) intersectsPoint(point geometry.Point) bool {
	return false
}

func (g *String) intersectsRect(rect geometry.Rect) bool {
	return false
}

func (g *String) intersectsLine(line *geometry.Line) bool {
	return false
}

func (g *String) intersectsPoly(poly *geometry.Poly) bool {
	return false
}

// NumPoints ...
func (g *String) NumPoints() int {
	return 0
}
