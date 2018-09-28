package geom

type simpleRing struct {
	points []Point
}

func (ring *simpleRing) move(deltaX, deltaY float64) *simpleRing {
	points := make([]Point, len(ring.points))
	for i := 0; i < len(ring.points); i++ {
		points[i].X = ring.points[i].X + deltaX
		points[i].Y = ring.points[i].Y + deltaY
	}
	return newSimpleRing(points)
}

func newSimpleRing(points []Point) *simpleRing {
	ring := new(simpleRing)
	ring.points = make([]Point, len(points))
	copy(ring.points, points)
	return ring
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

func (ring *simpleRing) IsClosed() bool {
	return true
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
			if !allowOnEdge {
				if raycast(seg.A, other.A, other.B).on ||
					raycast(seg.B, other.A, other.B).on ||
					raycast(other.A, seg.A, seg.B).on ||
					raycast(other.B, seg.A, seg.B).on {
					return false
				}
			}
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

func (ring *simpleRing) ContainsRing(other Ring, allowOnEdge bool) bool {
	return ringContainsRing(ring, other, allowOnEdge)
}

func (ring *simpleRing) IntersectsRing(other Ring, allowOnEdge bool) bool {
	return ringIntersectsRing(ring, other, allowOnEdge)
}
