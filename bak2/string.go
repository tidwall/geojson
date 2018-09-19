package geojson

import (
	"encoding/json"
)

// String is a plain string
type String string

// Rect containing the object
func (g String) Rect() Rect {
	return Rect{}
}

// Center of the object
func (g String) Center() Position {
	return Position{}
}

// AppendJSON appends a json representation to destination
func (g String) AppendJSON(dst []byte) []byte {
	b, _ := json.Marshal(string(g))
	dst = append(dst, b...)
	return dst
}

// JSON representation of the object
func (g String) JSON() string {
	return string(g.AppendJSON(nil))
}

// String representation of the object
func (g String) String() string {
	return string(g)
}

// Stats of the object
func (g String) Stats() Stats {
	return Stats{Weight: len(g)}
}

// Contains another object
func (g String) Contains(o Object) bool {
	return false
}

// Intersects another object
func (g String) Intersects(o Object) bool {
	return false
}

// IntersectsPolyLine test if object intersect a polyline
func (g String) IntersectsPolyLine(line []Position) bool {
	return false
}
