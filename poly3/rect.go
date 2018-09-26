package poly

// Point ...
type Point struct{ X, Y float64 }

// ContainsRing ...
func (point Point) ContainsRing(ring Ring) bool {
	panic("asdf")
}

// Rect ...
type Rect struct {
	Min, Max Point
}

func (rect Rect) move(deltaX, deltaY float64) Rect {
	return Rect{
		Point{rect.Min.X + deltaX, rect.Min.Y + deltaY},
		Point{rect.Max.X + deltaX, rect.Max.Y + deltaY},
	}
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
