package geom

// Point ...
type Point struct {
	X, Y float64
}

// Move ...
func (point Point) Move(deltaX, deltaY float64) Point {
	return Point{X: point.X + deltaX, Y: point.Y + deltaY}
}

// Empty ...
func (point Point) Empty() bool {
	return false
}

// Rect ...
func (point Point) Rect() Rect {
	return Rect{point, point}
}

// ContainsPoint ...
func (point Point) ContainsPoint(other Point) bool {
	return point == other
}

// IntersectsPoint ...
func (point Point) IntersectsPoint(other Point) bool {
	return point == other
}

// ContainsRect ...
func (point Point) ContainsRect(rect Rect) bool {
	return point.Rect() == rect
}

// IntersectsRect ...
func (point Point) IntersectsRect(rect Rect) bool {
	return rect.ContainsPoint(point)
}

// ContainsLine ...
func (point Point) ContainsLine(line *Line) bool {
	return !line.Empty() && line.Rect() == point.Rect()
}

// IntersectsLine ...
func (point Point) IntersectsLine(line *Line) bool {
	return line.IntersectsPoint(point)
}

// ContainsPoly ...
func (point Point) ContainsPoly(poly *Poly) bool {
	return !poly.Empty() && poly.Rect() == point.Rect()
}

// IntersectsPoly ...
func (point Point) IntersectsPoly(poly *Poly) bool {
	return poly.IntersectsPoint(point)
}
