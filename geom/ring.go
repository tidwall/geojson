package geom

import "math"

// Ring ...
type Ring interface {
	Rect() Rect
	Empty() bool
	Convex() bool
	ForEachPoint(iter func(point Point) bool)
	ForEachSegment(iter func(seg Segment, idx int) bool)
	Search(rect Rect, iter func(seg Segment, idx int) bool)
}

func newRing(points []Point) Ring {
	series := MakeSeries(points, true, true)
	return &series
}

func ringContainsPoint(ring Ring, point Point, allowOnEdge bool) bool {
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

func ringIntersectsPoint(ring Ring, point Point, allowOnEdge bool) bool {
	return ringContainsPoint(ring, point, allowOnEdge)
}

func ringContainsSegment(ring Ring, seg Segment, allowOnEdge bool) bool {
	if !ringContainsPoint(ring, seg.A, allowOnEdge) {
		return false
	}
	if !ringContainsPoint(ring, seg.B, allowOnEdge) {
		return false
	}
	if !ring.Convex() {
		if ringIntersectsSegment(ring, seg, false) {
			return false
		}
	}
	return true
}

func ringIntersectsSegment(ring Ring, seg Segment, allowOnEdge bool) bool {
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

func ringIntersectsPoly(ring Ring, poly *Poly, allowOnEdge bool) bool {
	// 1) ring must intersect poly exterior
	if !ringIntersectsRing(poly.Exterior, ring, allowOnEdge) {
		return false
	}
	// 2) ring cannot be contained by a poly hole
	for _, polyHole := range poly.Holes {
		if ringContainsRing(polyHole, ring, false) {
			return false
		}
	}
	return true
}

func ringContainsPoly(ring Ring, poly *Poly, allowOnEdge bool) bool {
	return ringContainsRing(ring, poly.Exterior, allowOnEdge)
}

func ringContainsRing(ring, other Ring, allowOnEdge bool) bool {
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
		if !ringContainsPoint(outer, point, allowOnEdge) {
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
			if ringIntersectsSegment(outer, seg, false) {
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

func ringIntersectsRing(ring, other Ring, allowOnEdge bool) bool {
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
		if ringContainsPoint(outer, seg.A, allowOnEdge) {
			// point from inner is inside outer. they intersect, stop now
			intersects = true
			return false
		}
		if ringIntersectsSegment(outer, seg, allowOnEdge) {
			// segment from inner intersects outer. they intersect, stop now
			intersects = true
			return false
		}
		return true
	})
	return intersects
}

func ringContainsRect(ring Ring, rect Rect, allowOnEdge bool) bool {
	return ringContainsRing(ring, rect, allowOnEdge)
}

func ringIntersectsRect(ring Ring, rect Rect, allowOnEdge bool) bool {
	return ringIntersectsRing(ring, rect, allowOnEdge)
}

func ringContainsLine(ring Ring, line *Line, allowOnEdge bool) bool {
	panic("not ready")
}

func ringIntersectsLine(ring Ring, line *Line, allowOnEdge bool) bool {
	panic("not ready")
}
