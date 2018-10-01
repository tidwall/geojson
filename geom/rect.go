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

// NumPoints ...
func (rect Rect) NumPoints() int {
	return 5
}

// ForEachPoint ...
func (rect Rect) ForEachPoint(iter func(point Point) bool) {
	points := [5]Point{
		{rect.Min.X, rect.Min.Y},
		{rect.Max.X, rect.Min.Y},
		{rect.Max.X, rect.Max.Y},
		{rect.Min.X, rect.Max.Y},
		{rect.Min.X, rect.Min.Y},
	}
	for _, point := range points {
		if !iter(point) {
			return
		}
	}
}

// ForEachSegment ...
func (rect Rect) ForEachSegment(iter func(seg Segment, idx int) bool) {
	segs := [4]Segment{
		{Point{rect.Min.X, rect.Min.Y}, Point{rect.Max.X, rect.Min.Y}},
		{Point{rect.Max.X, rect.Min.Y}, Point{rect.Max.X, rect.Max.Y}},
		{Point{rect.Max.X, rect.Max.Y}, Point{rect.Min.X, rect.Max.Y}},
		{Point{rect.Min.X, rect.Max.Y}, Point{rect.Min.X, rect.Min.Y}},
	}
	for i, seg := range segs {
		if !iter(seg, i) {
			return
		}
	}
}

// Search ...
func (rect Rect) Search(target Rect, iter func(seg Segment, idx int) bool) {
	rect.ForEachSegment(func(seg Segment, idx int) bool {
		if seg.Rect().IntersectsRect(rect) {
			if !iter(seg, idx) {
				return false
			}
		}
		return true
	})
}

// Empty ...
func (rect Rect) Empty() bool {
	return false
}

// Rect ...
func (rect Rect) Rect() Rect {
	return rect
}

// Convex ...
func (rect Rect) Convex() bool {
	return true
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
	return ringIntersectsLine(rect, line, true)
}

// ContainsPoly ...
func (rect Rect) ContainsPoly(poly *Poly) bool {
	return !poly.Empty() && rect.ContainsRect(poly.Rect())
}

// IntersectsPoly ...
func (rect Rect) IntersectsPoly(poly *Poly) bool {
	return ringIntersectsPoly(rect, poly, true)
}
