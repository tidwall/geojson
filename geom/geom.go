package geom

type validGeometry interface {
	Rect() Rect
	Empty() bool
	ContainsPoint(point Point) bool
	IntersectsPoint(point Point) bool
	ContainsRect(rect Rect) bool
	IntersectsRect(rect Rect) bool
	ContainsLine(line *Line) bool
	IntersectsLine(line *Line) bool
	ContainsPoly(poly *Poly) bool
	IntersectsPoly(poly *Poly) bool
}

// require conformance
var _ = []validGeometry{Point{}, Rect{}, &Line{}, &Poly{}}
