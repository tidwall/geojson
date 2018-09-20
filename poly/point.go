package poly

// Point is simple 2D point
// For geo locations: X is lat, Y is lon, and Z is elev or time measure.
type Point struct {
	X, Y float64
}

// InsideRect detects point is inside of another rect
func (p Point) InsideRect(rect Rect) bool {
	if p.X < rect.Min.X || p.X > rect.Max.X {
		return false
	}
	if p.Y < rect.Min.Y || p.Y > rect.Max.Y {
		return false
	}
	return true
}

// Inside returns true if point is inside of exterior and not in a hole.
// The validity of the exterior and holes must be done elsewhere and are
// assumed valid.
//   A valid exterior is a near-linear ring.
//   A valid hole is one that is full contained inside the exterior.
//   A valid hole may not share the same segment line as the exterior.
func (p Point) Inside(exterior Ring, holes []Ring) bool {
	if !insideshpext(p, exterior, true) {
		return false
	}
	for i := 0; i < len(holes); i++ {
		if insideshpext(p, holes[i], false) {
			return false
		}
	}
	return true
}

// IntersectsLineString detect if a point intersects a linestring
func (p Point) IntersectsLineString(exterior Ring) bool {
	for j := 0; j < len(exterior); j++ {
		if raycast(p, exterior[j], exterior[(j+1)%len(exterior)]).on {
			return true
		}
	}
	return false
}

// Intersects detects if a point intersects another polygon
func (p Point) Intersects(exterior Ring, holes []Ring) bool {
	return p.Inside(exterior, holes)
}
