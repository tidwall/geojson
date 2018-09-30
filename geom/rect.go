package geom

// Rect ...
type Rect struct {
	Min, Max Point
}

// Center ...
func (rect Rect) Center() Point {
	return Point{(rect.Max.X + rect.Min.X) / 2, (rect.Max.Y + rect.Min.Y) / 2}
}

// Area ...
func (rect Rect) Area() float64 {
	return (rect.Max.X - rect.Min.X) * (rect.Max.Y - rect.Min.Y)
}

func (rect Rect) ringPoints() [5]Point {
	return [5]Point{
		{rect.Min.X, rect.Min.Y},
		{rect.Max.X, rect.Min.Y},
		{rect.Max.X, rect.Max.Y},
		{rect.Min.X, rect.Max.Y},
		{rect.Min.X, rect.Min.Y},
	}
}

func (rect Rect) ring() *Ring {
	points := rect.ringPoints()
	series := Series{closed: true, convex: true, rect: rect, points: points[:]}
	return &Ring{series}
}

// Empty ...
func (rect Rect) Empty() bool {
	return false
}

// Rect ...
func (rect Rect) Rect() Rect {
	return rect
}

// ContainsPoint ...
func (rect Rect) ContainsPoint(point Point) bool {
	return point.X >= rect.Min.X && point.X <= rect.Max.X &&
		point.Y >= rect.Min.Y && point.Y <= rect.Max.Y
}

// IntersectsPoint ...
func (rect Rect) IntersectsPoint(point Point) bool {
	return rect.ContainsPoint(point)
}

// ContainsRect ...
func (rect Rect) ContainsRect(other Rect) bool {
	if other.Min.X < rect.Min.X || other.Max.X > rect.Max.X {
		return false
	}
	if other.Min.Y < rect.Min.Y || other.Max.Y > rect.Max.Y {
		return false
	}
	return true
}

// IntersectsRect ...
func (rect Rect) IntersectsRect(other Rect) bool {
	if rect.Min.Y > other.Max.Y || rect.Max.Y < other.Min.Y {
		return false
	}
	if rect.Min.X > other.Max.X || rect.Max.X < other.Min.X {
		return false
	}
	return true
}

// ContainsLine ...
func (rect Rect) ContainsLine(line *Line) bool {
	return !line.Empty() && rect.ContainsRect(line.Rect())
}

// IntersectsLine ...
func (rect Rect) IntersectsLine(line *Line) bool {
	return rect.ring().IntersectsLine(line, true)
}

// ContainsPoly ...
func (rect Rect) ContainsPoly(poly *Poly) bool {
	return !poly.Empty() && rect.ContainsRect(poly.Rect())
}

// IntersectsPoly ...
func (rect Rect) IntersectsPoly(poly *Poly) bool {
	// TODO: optimize
	return rect.ring().IntersectsPoly(poly, true)
}
