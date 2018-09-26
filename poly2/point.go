package poly

// Point ...
type Point struct {
	X, Y float64
}

// ContainsPoint ...
func (point Point) ContainsPoint(other Point) bool {
	return point == other
}

// ContainsRect ...
func (point Point) ContainsRect(rect Rect) bool {
	return point == rect.Min && point == rect.Max
}

// ContainsLine ...
func (point Point) ContainsLine(line Line) bool {
	if len(line) == 0 {
		return false
	}
	for _, ringLine := range line {
		if ringLine != point {
			return false
		}
	}
	return true
}

// ContainsRing ...
func (point Point) ContainsRing(ring Ring) bool {
	if len(ring) == 0 {
		return false
	}
	for _, ringPoint := range ring {
		if ringPoint != point {
			return false
		}
	}
	return true
}

// ContainsPolygon ...
func (point Point) ContainsPolygon(polygon Polygon) bool {
	return point.ContainsRing(polygon.Exterior)
}

// IntersectsPoint ...
func (point Point) IntersectsPoint(other Point) bool {
	return point == other
}

// IntersectsRect ...
func (point Point) IntersectsRect(rect Rect) bool {
	return point.X >= rect.Min.X && point.X <= rect.Max.X &&
		point.Y >= rect.Min.Y && point.Y <= rect.Max.Y
}

// IntersectsLine ...
func (point Point) IntersectsLine(line Line) bool {
	return algoPointOnLine(point, line)
}

// IntersectsRing ...
func (point Point) IntersectsRing(ring Ring) bool {
	return algoPointInRing(point, ring, true)
}

// IntersectsPolygon ...
func (point Point) IntersectsPolygon(polygon Polygon) bool {
	if !algoPointInRing(point, polygon.Exterior, true) {
		return false
	}
	for _, hole := range polygon.Holes {
		if algoPointInRing(point, hole, false) {
			return false
		}
	}
	return true
}
