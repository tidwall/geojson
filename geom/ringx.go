package geom

import (
	"math"
	"strconv"
)

const complexRingMinPoints = 16

// RingX ...
type RingX = Series

func newRingX(points []Point) RingX {
	series := makeSeries(points, true, true)
	return &series
}

func ringxString(ring RingX) string {
	var buf []byte
	buf = append(buf, '[')
	n := ring.NumPoints()
	for i := 0; i < n; i++ {
		if i == 0 {
			buf = append(buf, '\n')
		}
		pt := ring.PointAt(i)
		buf = append(buf, '\t', '[')
		buf = strconv.AppendFloat(buf, pt.X, 'f', -1, 64)
		buf = append(buf, ',')
		buf = strconv.AppendFloat(buf, pt.Y, 'f', -1, 64)
		buf = append(buf, ']')
		if i < n-1 {
			buf = append(buf, ',')
		}
		buf = append(buf, '\n')
	}

	buf = append(buf, ']')
	return string(buf)
}

type ringxResult struct {
	hit bool // contains/intersects
	idx int  // edge index
}

func ringxContainsPoint(ring RingX, point Point, allowOnEdge bool) ringxResult {
	// println("A")
	var idx = -1
	// find all intersecting segments on the y-axis
	var in bool
	ring.Search(
		Rect{Point{math.Inf(-1), point.Y}, Point{math.Inf(+1), point.Y}},
		func(seg Segment, index int) bool {
			// fmt.Printf("%v %v\n", point, seg)
			// perform a raycast operation on the segments
			res := seg.Raycast(point)
			if res.On {
				// println(1)
				in = allowOnEdge
				idx = index
				return false
			}
			if res.In {
				// println(2)
				in = !in
			}
			return true
		},
	)
	return ringxResult{hit: in, idx: idx}
}

func ringxIntersectsPoint(ring RingX, point Point, allowOnEdge bool) ringxResult {
	return ringxContainsPoint(ring, point, allowOnEdge)
}

// func segmentsIntersects(seg, other Segment, allowOnEdge bool) bool {
// 	if seg.IntersectsSegment(other) {

// 	}
// 	return false
// }

func ringxContainsSegment(ring RingX, seg Segment, allowOnEdge bool) bool {
	// Test that segment points are contained in the ring.
	resA := ringxContainsPoint(ring, seg.A, allowOnEdge)
	if !resA.hit {
		// seg A is not inside ring
		return false
	}
	resB := ringxContainsPoint(ring, seg.B, allowOnEdge)
	if !resB.hit {
		// seg B is not inside ring
		return false
	}
	if ring.Convex() {
		// ring is convex so the segment must be contained
		return true
	}
	// The ring is concave so it's possible that the segment crosses over the
	// edge of the ring.
	if allowOnEdge {
		// do some logic around seg points that are on the edge of the ring.
		if resA.idx != -1 {
			// seg A is on a ring segment
			if resB.idx != -1 {
				// seg B is on a ring segment
				if resB.idx == resA.idx {
					// case (3)
					// seg A and B share the same ring segment, so it must be
					// on the inside.
					return true
				}
				// case (1)
				// seg A and seg B are on different segments.
				// determine if the space that the seg passes over is inside or
				// outside of the ring. To do so we create a ring from the two
				// ring segments and check if that ring winding order matches
				// the winding order of the ring.
				// -- create a ring
				rSegA := ring.SegmentAt(resA.idx)
				rSegB := ring.SegmentAt(resB.idx)
				if resB.idx < resA.idx {
					rSegA, rSegB = rSegB, rSegA
				}
				pts := [5]Point{rSegA.A, rSegA.B, rSegB.A, rSegB.B, rSegA.A}
				// -- calc winding order
				var cwc float64
				for i := 0; i < len(pts)-1; i++ {
					a, b := pts[i], pts[i+1]
					cwc += (b.X - a.X) * (b.Y + a.Y)
				}
				clockwise := cwc > 0
				if clockwise != ring.Clockwise() {
					// -- on the outside
					return false
				}
				// the passover space is on the inside of the ring.
				// check if seg intersects any ring segments where A and B are
				// not on.
				var intersects bool
				ring.Search(seg.Rect(), func(seg2 Segment, index int) bool {
					if seg.IntersectsSegment(seg2) {
						if !seg2.Raycast(seg.A).On && !seg2.Raycast(seg.B).On {
							intersects = true
							return false
						}
					}
					return true
				})
				return !intersects
			}
			// case (4)
			// seg A is on a ring segment, but seg B is not.
			// check if seg intersects any ring segments where A is not on.
			var intersects bool
			ring.Search(seg.Rect(), func(seg2 Segment, index int) bool {
				if seg.IntersectsSegment(seg2) {
					if !seg2.Raycast(seg.A).On {
						intersects = true
						return false
					}
				}
				return true
			})
			return !intersects
		} else if resB.idx != -1 {
			// case (2)
			// seg B is on a ring segment, but seg A is not.
			// check if seg intersects any ring segments where B is not on.
			var intersects bool
			ring.Search(seg.Rect(), func(seg2 Segment, index int) bool {
				if seg.IntersectsSegment(seg2) {
					if !seg2.Raycast(seg.B).On {
						intersects = true
						return false
					}
				}
				return true
			})
			return !intersects
		}
		// case (5) (15)
		var intersects bool
		ring.Search(seg.Rect(), func(seg2 Segment, index int) bool {
			if seg.IntersectsSegment(seg2) {
				if !seg.Raycast(seg2.A).On && !seg.Raycast(seg2.B).On {
					intersects = true
					return false
				}
			}
			return true
		})
		return !intersects
	}

	// allowOnEdge is false. (not allow on edge)
	var intersects bool
	ring.Search(seg.Rect(), func(seg2 Segment, index int) bool {
		if seg.IntersectsSegment(seg2) {
			// if seg.Raycast(seg2.A).On || seg.Raycast(seg2.B).On {
			intersects = true
			// 	return false
			// }
			return false
		}
		return true
	})
	return !intersects
}

// ringxIntersectsSegment detect if the segment intersects the ring
func ringxIntersectsSegment(ring RingX, seg Segment, allowOnEdge bool) bool {
	// Quick check that either point is inside of the ring
	if ringxContainsPoint(ring, seg.A, allowOnEdge).hit {
		return true
	}
	if ringxContainsPoint(ring, seg.B, allowOnEdge).hit {
		return true
	}
	// Neither point A or B is inside the the ring. It's possible that both
	// are on the outside and are passing over segments. If the segment passes
	// over at least two ring segments then it's intersecting.
	var count int
	var segAOn bool
	var segBOn bool
	ring.Search(seg.Rect(), func(seg2 Segment, index int) bool {
		if seg.IntersectsSegment(seg2) {
			if !allowOnEdge {
				// for segments that are not allowed on the edge, extra care
				// must be taken.
				if !(seg.CollinearPoint(seg2.A) && seg.CollinearPoint(seg2.B)) {
					if !segAOn {
						if seg.A == seg2.A || seg.A == seg2.B {
							segAOn = true
							return true
						}
					}
					if !segBOn {
						if seg.B == seg2.A || seg.B == seg2.B {
							segBOn = true
							return true
						}
					}
					count++
				}
			} else {
				count++
			}
		}
		return count < 2
	})
	return count >= 2
}

func ringxContainsRing(ring, other RingX, allowOnEdge bool) bool {
	if ring.Empty() || other.Empty() {
		return false
	}
	if other.NumPoints() >= complexRingMinPoints {
		// inner ring has a lot of points, and is convex, so let just check if
		// the rect ring is fully contained before we do the complicated stuff.
		if ringxContainsRing(ring, other.Rect(), allowOnEdge) {
			return true
		}
	}
	// test if the inner rect does not contain the outer rect
	if !ring.Rect().ContainsRect(other.Rect()) {
		// not contained so it's not possible for the outer ring to contain
		// the inner ring
		return false
	}
	if ring.Convex() {
		// outer ring is convex so test that all inner points are inside of
		// the outer ring
		otherNumPoints := other.NumPoints()
		for i := 0; i < otherNumPoints; i++ {
			if !ringxContainsPoint(ring, other.PointAt(i), allowOnEdge).hit {
				// point is on the outside the outer ring
				return false
			}
		}
	} else {
		// outer ring is concave so let's make sure that all inner segments are
		// fully contained inside of the outer ring.
		otherNumSegments := other.NumSegments()
		for i := 0; i < otherNumSegments; i++ {
			if !ringxContainsSegment(ring, other.SegmentAt(i), allowOnEdge) {
				return false
			}
		}
	}
	return true
}

func ringxIntersectsRing(ring, other RingX, allowOnEdge bool) bool {
	if ring.Empty() || other.Empty() {
		return false
	}
	// check outer and innter rects intersection first
	if !ring.Rect().IntersectsRect(other.Rect()) {
		return false
	}
	if other.Rect().Area() > ring.Rect().Area() {
		// swap the rings so that the inner ring is smaller than the outer ring
		ring, other = other, ring
	}
	if ring.Convex() && allowOnEdge {
		// outer ring is convex so test that any inner points are inside of
		// the outer ring
		otherNumPoints := other.NumPoints()
		for i := 0; i < otherNumPoints; i++ {
			if ringxContainsPoint(ring, other.PointAt(i), allowOnEdge).hit {
				return true
			}
		}
	} else {
		// outer ring is concave so let's make sure that all inner segments are
		// fully contained inside of the outer ring.
		otherNumSegments := other.NumSegments()
		for i := 0; i < otherNumSegments; i++ {
			if ringxIntersectsSegment(ring, other.SegmentAt(i), allowOnEdge) {
				return true
			}
		}
	}
	return false
}

func ringxContainsLine(ring RingX, line *Line, allowOnEdge bool) bool {
	// shares the same logic
	return ringxContainsRing(ring, RingX(&line.baseSeries), allowOnEdge)
}

func ringxIntersectsLine(ring RingX, line *Line, allowOnEdge bool) bool {
	if ring.Empty() || line.Empty() {
		return false
	}
	// check outer and innter rects intersection first
	if !ring.Rect().IntersectsRect(line.Rect()) {
		return false
	}
	// check if any points are inside ring
	lineNumPoints := line.NumPoints()
	for i := 0; i < lineNumPoints; i++ {
		if ringxContainsPoint(ring, line.PointAt(i), allowOnEdge).hit {
			return true
		}
	}
	lineNumSegments := line.NumSegments()
	for i := 0; i < lineNumSegments; i++ {
		if ringxIntersectsSegment(ring, line.SegmentAt(i), allowOnEdge) {
			return true
		}
	}
	return false
}
