package poly

// Ring ...
type Ring []Point

// ContainsPoint ...
func (ring Ring) ContainsPoint(point Point) bool {
	if len(ring) == 0 {
		return false
	}
	if len(ring) == 1 {
		return ring[0].ContainsPoint(point)
	}
	if len(ring) == 2 {
		return Line(ring).ContainsPoint(point)
	}
	return algoPointInRing(point, ring, true)
}

// ContainsLine ...
func (ring Ring) ContainsLine(line Line) bool {
	if len(ring) == 0 {
		return false
	}
	if len(ring) == 1 {
		return ring[0].ContainsLine(line)
	}
	if len(ring) == 2 {
		return Line(ring).ContainsLine(line)
	}
	// make sure all the other points are inside the ring
	for _, point := range line {
		if !algoPointInRing(point, ring, true) {
			return false
		}
	}
	if !ring.IsConvex() {
		// make sure that no lines intersect
		if algoAnySegmentIntersects(ring, line, true, false) {
			return false
		}
	}
	return true
}

// ContainsRing ...
func (ring Ring) ContainsRing(other Ring) bool {
	if len(ring) == 0 {
		return false
	}
	if len(ring) == 1 {
		return ring[0].ContainsRing(other)
	}
	if len(ring) == 2 {
		return Line(ring).ContainsRing(other)
	}
	// make sure all the other points are inside the ring
	for _, point := range other {
		if !algoPointInRing(point, ring, true) {
			return false
		}
	}
	if !ring.IsConvex() {
		// make sure that no lines intersect
		if algoAnySegmentIntersects(ring, other, true, true) {
			return false
		}
	}
	return true
}

// IsConvex ...
func (ring Ring) IsConvex() bool {
	if len(ring) < 3 {
		return false
	}
	var dir int
	var a, b, c Point
	for i := 0; i < len(ring); i++ {
		a = ring[i]
		if i == len(ring)-1 {
			b = ring[0]
			c = ring[1]
		} else if i == len(ring)-2 {
			b = ring[i+1]
			c = ring[0]
		} else {
			b = ring[i+1]
			c = ring[i+2]
		}
		dx1 := b.X - a.X
		dy1 := b.Y - a.Y
		dx2 := c.X - b.X
		dy2 := c.Y - b.Y
		zCrossProduct := dx1*dy2 - dy1*dx2
		if dir == 0 {
			if zCrossProduct < 0 {
				dir = -1
			} else if zCrossProduct > 0 {
				dir = 1
			}
		} else if zCrossProduct < 0 {
			if dir == 1 {
				return false
			}
		} else if zCrossProduct > 0 {
			if dir == -1 {
				return false
			}
		}
	}
	return true
}
