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

// ContainsRing ...
func (rect Rect) ContainsRing(ring Ring) bool {
	return rect.ContainsRect(ring.Rect())
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

// IntersectsRing ...
func (rect Rect) IntersectsRing(ring Ring) bool {
	rectPoints := rect.ringPoints()
	rectRing := &ringSimple{points: rectPoints[:]}
	return rectRing.IntersectsRing(ring, true)
}

// // ContainsPoly ...
// func (rect Rect) ContainsPoly(poly Poly) bool {
// 	rectPoints := rect.ringPoints()
// 	rectRing := &simpleRing{points: rectPoints[:]}
// 	return rectRing.ContainsPoly(ring, true)
// }
