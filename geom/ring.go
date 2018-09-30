package geom

import "math"

// Ring ...
type Ring struct {
	Series
}

// NewRing ...
func NewRing(points []Point) *Ring {
	ring := new(Ring)
	ring.Series = MakeSeries(points, true, true)
	return ring
}

func (ring *Ring) move(deltaX, deltaY float64) *Ring {
	points := make([]Point, len(ring.points))
	for i := 0; i < len(ring.points); i++ {
		points[i].X = ring.points[i].X + deltaX
		points[i].Y = ring.points[i].Y + deltaY
	}
	return NewRing(points)
}

// ContainsPoint ...
func (ring *Ring) ContainsPoint(point Point, allowOnEdge bool) bool {
	in := false
	ring.Search(
		Rect{Point{math.Inf(-1), point.Y}, Point{math.Inf(+1), point.Y}},
		func(seg Segment, index int) bool {
			res := raycast(point, seg.A, seg.B)
			if res.on {
				in = allowOnEdge
				return false
			}
			if res.in {
				in = !in
			}
			return true
		},
	)
	return in
}

// IntersectsPoint ...
func (ring *Ring) IntersectsPoint(point Point, allowOnEdge bool) bool {
	return ring.ContainsPoint(point, allowOnEdge)
}

// ContainsSegment ...
func (ring *Ring) ContainsSegment(seg Segment, allowOnEdge bool) bool {
	if !ring.ContainsPoint(seg.A, allowOnEdge) {
		return false
	}
	if !ring.ContainsPoint(seg.B, allowOnEdge) {
		return false
	}
	if !ring.Convex() {
		if ring.IntersectsSegment(seg, false) {
			return false
		}
	}
	return true

}

// IntersectsSegment ...
func (ring *Ring) IntersectsSegment(seg Segment, allowOnEdge bool) bool {
	var intersects bool
	ring.Search(seg.Rect(), func(other Segment, index int) bool {
		if segmentsIntersect(seg.A, seg.B, other.A, other.B) {
			if !allowOnEdge {
				if raycast(seg.A, other.A, other.B).on ||
					raycast(seg.B, other.A, other.B).on ||
					raycast(other.A, seg.A, seg.B).on ||
					raycast(other.B, seg.A, seg.B).on {
					intersects = false
					return false
				}
			}
			intersects = true
			return false
		}
		return true
	})
	return intersects
}

// ContainsRect ...
func (ring *Ring) ContainsRect(rect Rect, allowOnEdge bool) bool {
	points := rect.ringPoints()
	rectRing := &Ring{MakeSeries(points[:], false, true)}
	return ring.ContainsRing(rectRing, allowOnEdge)
}

// IntersectsRect ...
func (ring *Ring) IntersectsRect(rect Rect, allowOnEdge bool) bool {
	points := rect.ringPoints()
	rectRing := &Ring{MakeSeries(points[:], false, true)}
	return ring.IntersectsRing(rectRing, allowOnEdge)
}

// ContainsLine ...
func (ring *Ring) ContainsLine(line *Line, allowOnEdge bool) bool {
	panic("not ready")
}

// IntersectsLine ...
func (ring *Ring) IntersectsLine(line *Line, allowOnEdge bool) bool {
	panic("not ready")
}

// ContainsRing ...
func (ring *Ring) ContainsRing(other *Ring, allowOnEdge bool) bool {
	if ring.Empty() || other.Empty() {
		return false
	}
	outer, inner := ring, other
	outerRect := outer.Rect()
	innerRect := inner.Rect()
	// 1) check if the rects intersect each other
	if !outerRect.ContainsRect(innerRect) {
		// not contained, stop now
		return false
	}
	// 2) test if points are inside
	inside := true
	inner.ForEachPoint(func(point Point) bool {
		if !outer.ContainsPoint(point, allowOnEdge) {
			// not contained, stop now
			inside = false
			return false
		}
		return true
	})
	if !inside {
		return false
	}
	// 3) check intersecting segments if outer is convex
	if !outer.Convex() {
		var intersects bool
		inner.ForEachSegment(func(seg Segment, idx int) bool {
			if outer.IntersectsSegment(seg, false) {
				intersects = true
				return false
			}
			return true
		})
		if intersects {
			return false
		}
	}
	return true
}

// IntersectsRing ...
func (ring *Ring) IntersectsRing(other *Ring, allowOnEdge bool) bool {
	if ring.Empty() || other.Empty() {
		return false
	}
	outer, inner := ring, other
	outerRect := outer.Rect()
	innerRect := inner.Rect()
	// 1) check if the rects intersect each other
	if !outerRect.IntersectsRect(innerRect) {
		// they do not intersect so stop now
		return false
	}
	// 2) make sure the outer rect area is greater or equal to inner rect area
	if outerRect.Area() < innerRect.Area() {
		outer, inner = inner, outer
		outerRect, innerRect = innerRect, outerRect
	}
	// 3) test if points or segment intersection
	var intersects bool
	inner.ForEachSegment(func(seg Segment, idx int) bool {
		if outer.ContainsPoint(seg.A, allowOnEdge) {
			// point from inner is inside outer. they intersect, stop now
			intersects = true
			return false
		}
		if outer.IntersectsSegment(seg, allowOnEdge) {
			// segment from inner intersects outer. they intersect, stop now
			intersects = true
			return false
		}
		return true
	})
	return intersects
}

// ContainsPoly ...
func (ring *Ring) ContainsPoly(poly *Poly, allowOnEdge bool) bool {
	return ring.ContainsRing(poly.Exterior, allowOnEdge)
}

// IntersectsPoly ...
func (ring *Ring) IntersectsPoly(poly *Poly, allowOnEdge bool) bool {
	// 1) ring must intersect poly exterior
	if !poly.Exterior.IntersectsRing(ring, allowOnEdge) {
		return false
	}
	// 2) ring cannot be contained by a poly hole
	for _, polyHole := range poly.Holes {
		if polyHole.ContainsRing(ring, false) {
			return false
		}
	}
	return true
}
