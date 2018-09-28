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

	IntersectsSegment(seg Segment, allowOnEdge bool) bool
	ContainsPoint(point Point, allowOnEdge bool) bool
	IntersectsRing(ring Ring, allowOnEdge bool) bool
	ContainsRing(ring Ring, allowOnEdge bool) bool
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
