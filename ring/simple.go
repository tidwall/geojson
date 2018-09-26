package ring

type simpleRing struct {
	points []Point
}

func newSimpleRing(points []Point) Ring {
	return &simpleRing{points: points}
}

func (ring *simpleRing) Points() []Point {
	return ring.points
}

func (ring *simpleRing) Search(
	rect Rect, iter func(seg Segment, index int) bool,
) {
	var index int
	ring.Scan(func(seg Segment) bool {
		if seg.Rect().IntersectsRect(rect) {
			if !iter(seg, index) {
				return false
			}
		}
		index++
		return true
	})
}
func (ring *simpleRing) Scan(iter func(seg Segment) bool) {
	for i := 0; i < len(ring.points); i++ {
		var seg Segment
		seg.A = ring.points[i]
		if i == len(ring.points)-1 {
			if seg.A == ring.points[0] {
				break
			}
			seg.B = ring.points[0]
		} else {
			seg.B = ring.points[i+1]
		}
		if !iter(seg) {
			return
		}
	}
}

func (ring *simpleRing) Rect() Rect {
	return pointsRect(ring.points)
}

func (ring *simpleRing) Convex() bool {
	return pointsConvex(ring.points)
}

func (ring *simpleRing) IntersectsSegment(seg Segment, allowOnEdge bool) bool {
	for i := 0; i < len(ring.points); i++ {
		var other Segment
		other.A = ring.points[i]
		if i == len(ring.points)-1 {
			if other.A == ring.points[0] {
				break
			}
			other.B = ring.points[0]
		} else {
			other.B = ring.points[i+1]
		}
		if segmentsIntersect(seg.A, seg.B, other.A, other.B) {
			return true
		}
	}
	return false
}

func (ring *simpleRing) ContainsPoint(point Point, allowOnEdge bool) bool {
	in := false
	var a, b Point
	for i := 0; i < len(ring.points); i++ {
		a = ring.points[i]
		if i == len(ring.points)-1 {
			b = ring.points[0]
		} else {
			b = ring.points[i+1]
		}
		res := raycast(point, a, b)
		if res.on {
			return allowOnEdge
		}
		if res.in {
			in = !in
		}
	}
	return in
}

func (ring *simpleRing) ContainsRing(other Ring) bool {

	panic("not ready")
}

func (ring *simpleRing) IntersectsRing(other Ring) bool {
	return ringsIntersect(ring, other)
}

func ringsIntersect(inner, outer Ring) bool {
	outerRect := outer.Rect()
	innerRect := inner.Rect()
	// 1) make sure the outer rect area is greater to inner rect area
	if outerRect.Area() < innerRect.Area() {
		outer, inner = inner, outer
		outerRect, innerRect = innerRect, outerRect
	}
	// 2) check if the rects intersect each other
	if !outerRect.IntersectsRect(innerRect) {
		// they do not intersect so stop now
		return false
	}
	// 3) test is points and segment intersection
	var intersects bool
	inner.Scan(func(seg Segment) bool {
		if outer.ContainsPoint(seg.A, true) {
			// point from inner is inside outer. they intersect so stop now
			intersects = true
			return false
		}
		return true
	})
	return intersects
}
