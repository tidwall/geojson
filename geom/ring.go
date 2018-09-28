package geom

// DefaultIndex ...
const (
	Index     = 0
	NoIndex   = -1
	AutoIndex = 48
)

// Ring ...
type Ring interface {
	Scan(iter func(seg Segment) bool)
	Search(rect Rect, iter func(seg Segment, index int) bool)
	Points() []Point
	Rect() Rect
	IsClosed() bool
	Convex() bool

	ContainsPoint(point Point, allowOnEdge bool) bool

	ContainsSegment(seg Segment, allowOnEdge bool) bool
	IntersectsSegment(seg Segment, allowOnEdge bool) bool

	ContainsRect(rect Rect, allowOnEdge bool) bool
	IntersectsRect(rect Rect, allowOnEdge bool) bool

	ContainsRing(ring Ring, allowOnEdge bool) bool
	IntersectsRing(ring Ring, allowOnEdge bool) bool

	ContainsPoly(poly Poly, allowOnEdge bool) bool
	IntersectsPoly(poly Poly, allowOnEdge bool) bool
}

// NewRing returns a new ring. index of zero reutrns simple ring
func NewRing(points []Point, index int) Ring {
	if index >= 0 && len(points) > index {
		return newTreeRing(points)
	}
	return newSimpleRing(points)
}

func ringIntersectsRing(outer, inner Ring, allowOnEdge bool) bool {
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
	inner.Scan(func(seg Segment) bool {
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

func ringContainsRing(outer, inner Ring, allowOnEdge bool) bool {
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
		inner.Scan(func(seg Segment) bool {
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

func ringContainsRect(ring Ring, rect Rect, allowOnEdge bool) bool {
	points := rect.ringPoints()
	rectRing := &simpleRing{points: points[:]}
	return ringContainsRing(ring, rectRing, allowOnEdge)
}

func ringIntersectsRect(ring Ring, rect Rect, allowOnEdge bool) bool {
	points := rect.ringPoints()
	rectRing := &simpleRing{points: points[:]}
	return ringIntersectsRing(ring, rectRing, allowOnEdge)
}

func ringContainsSegment(ring Ring, seg Segment, allowOnEdge bool) bool {
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

func ringIntersectsPoly(ring Ring, poly Poly, allowOnEdge bool) bool {
	// 1) ring must intersect poly exterior
	if !poly.Exterior().IntersectsRing(ring, allowOnEdge) {
		return false
	}
	// 2) ring cannot be contained by a poly hole
	for _, polyHole := range poly.Holes() {
		if polyHole.ContainsRing(ring, false) {
			return false
		}
	}
	return true
}
