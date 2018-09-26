package poly

// Ring ...
type Ring []Point

// Rect ...
func (ring Ring) Rect() Rect {
	return pointsRect(ring)
}

func (ring Ring) move(deltaX, deltaY float64) Ring {
	newRing := make(Ring, len(ring))
	for i := 0; i < len(ring); i++ {
		newRing[i].X = ring[i].X + deltaX
		newRing[i].Y = ring[i].Y + deltaY
	}
	return newRing
}

// ContainsRing ...
func (ring Ring) ContainsRing(other Ring) bool {
	if len(ring) == 0 {
		return false
	}
	if len(ring) == 1 {
		panic("Point(ring[0]).ContainsRing(other)")
	}
	if len(ring) == 2 {
		panic("Line(ring).ContainsRing(other)")
	}
	if !ring.Rect().ContainsRect(other.Rect()) {
		return false
	}
	// make sure all the other points are inside the ring
	for _, point := range other {
		if !pointInRing(point, ring, true) {
			return false
		}
	}
	if !ringConvex(ring) {
		// make sure that no lines intersect
		if anySegmentsIntersect(ring, other, true, true, false, true) {
			return false
		}
	}
	return true
}

// IntersectsRing ...
func (ring Ring) IntersectsRing(other Ring) bool {
	if len(ring) == 0 {
		return false
	}
	if len(ring) == 1 {
		panic("Point(ring[0]).IntersectsRing(other)")
	}
	if len(ring) == 2 {
		panic("Line(ring).IntersectsRing(other)")
	}
	return ring.intersectsRing(other, true)
}

func (ring Ring) intersectsRing(other Ring, exterior bool) bool {
	if !ring.Rect().IntersectsRect(other.Rect()) {
		return false
	}
	// check if any points from a ring are within the other ring
	for _, point := range other {
		if pointInRing(point, ring, exterior) {
			return true
		}
	}
	for _, point := range ring {
		if pointInRing(point, other, exterior) {
			return true
		}
	}
	return anySegmentsIntersect(ring, other, true, true, exterior, exterior)
}
