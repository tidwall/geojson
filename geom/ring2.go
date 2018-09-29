package geom

import "math"

// Ring2 ...
type Ring2 struct {
	Series
}

// NewRing2 ...
func NewRing2(points []Point) *Ring2 {
	ring := new(Ring2)
	ring.Series = MakeSeries(points, true, true)
	return ring
}

// ContainsPoint ...
func (ring *Ring2) ContainsPoint(point Point, allowOnEdge bool) bool {
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
func (ring *Ring2) IntersectsPoint(point Point, allowOnEdge bool) bool {
	return ring.ContainsPoint(point, allowOnEdge)
}

// ContainsSegment ...
func (ring *Ring2) ContainsSegment(seg Segment, allowOnEdge bool) bool {
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
func (ring *Ring2) IntersectsSegment(seg Segment, allowOnEdge bool) bool {
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
func (ring *Ring2) ContainsRect(rect Rect, allowOnEdge bool) bool {
	points := rect.ringPoints()
	rectRing := &Ring2{MakeSeries(points[:], false, true)}
	return ring.ContainsRing(rectRing, allowOnEdge)
}

// IntersectsRect ...
func (ring *Ring2) IntersectsRect(rect Rect, allowOnEdge bool) bool {
	points := rect.ringPoints()
	rectRing := &Ring2{MakeSeries(points[:], false, true)}
	return ring.IntersectsRing(rectRing, allowOnEdge)
}

// ContainsRing ...
func (ring *Ring2) ContainsRing(other *Ring2, allowOnEdge bool) bool {
	outer, inner := ring, other
	outerRect := outer.Rect()
	innerRect := inner.Rect()
	// 1) check if the rects intersect each other
	if !outerRect.ContainsRect(innerRect) {
		// not contained, stop now
		return false
	}
	// 2) test if points are inside
	points := inner.Points()
	for _, point := range points {
		if !outer.ContainsPoint(point, allowOnEdge) {
			// not contained, stop now
			return false
		}
	}
	// 3) check intersecting segments if outer is convex
	if !outer.Convex() {
		var intersects bool
		inner.Scan(func(seg Segment, idx int) bool {
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
func (ring *Ring2) IntersectsRing(other *Ring2, allowOnEdge bool) bool {
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
	inner.Scan(func(seg Segment, idx int) bool {
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
func (ring *Ring2) ContainsPoly(poly *Poly2, allowOnEdge bool) bool {
	return ring.ContainsRing(poly.Exterior, allowOnEdge)
}

// IntersectsPoly ...
func (ring *Ring2) IntersectsPoly(poly *Poly2, allowOnEdge bool) bool {
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
