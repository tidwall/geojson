package geom

import (
	"math"

	"github.com/tidwall/boxtree/d2"
)

type ringIndexed struct {
	points []Point
	rect   Rect
	convex bool
	tree   d2.BoxTree
}

func newRingIndexed(points []Point) *ringIndexed {
	var ring ringIndexed
	ring.points = make([]Point, len(points))
	copy(ring.points, points)
	ring.init()
	return &ring
}

func (ring *ringIndexed) move(deltaX, deltaY float64) *ringIndexed {
	points := make([]Point, len(ring.points))
	for i := 0; i < len(ring.points); i++ {
		points[i].X = ring.points[i].X + deltaX
		points[i].Y = ring.points[i].Y + deltaY
	}
	return newRingIndexed(points)
}

func (ring *ringIndexed) init() {
	ring.convex, ring.rect = pointsConvexRect(ring.points)
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
		rect := seg.Rect()
		ring.tree.Insert(
			[]float64{rect.Min.X, rect.Min.Y},
			[]float64{rect.Max.X, rect.Max.Y},
			i,
		)
	}
}

func (ring *ringIndexed) Points() []Point {
	return ring.points
}

func (ring *ringIndexed) IsClosed() bool {
	return true
}

func (ring *ringIndexed) Rect() Rect {
	return ring.rect
}

func (ring *ringIndexed) Convex() bool {
	return ring.convex
}

func (ring *ringIndexed) Search(
	rect Rect, iter func(seg Segment, index int) bool,
) {
	ring.tree.Search(
		[]float64{rect.Min.X, rect.Min.Y},
		[]float64{rect.Max.X, rect.Max.Y},
		func(_, _ []float64, value interface{}) bool {
			index := value.(int)
			var seg Segment
			seg.A = ring.points[index]
			if index == len(ring.points)-1 {
				seg.B = ring.points[0]
			} else {
				seg.B = ring.points[index+1]
			}
			if !iter(seg, index) {
				return false
			}
			return true
		},
	)
}

func (ring *ringIndexed) Scan(iter func(seg Segment) bool) {
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

func (ring *ringIndexed) IntersectsSegment(seg Segment, allowOnEdge bool) bool {
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

func (ring *ringIndexed) ContainsPoint(point Point, allowOnEdge bool) bool {
	rect := Rect{
		Min: Point{math.Inf(-1), point.Y},
		Max: Point{math.Inf(+1), point.Y},
	}
	in := false
	ring.Search(rect, func(seg Segment, index int) bool {
		res := raycast(point, seg.A, seg.B)
		if res.on {
			in = allowOnEdge
			return false
		}
		if res.in {
			in = !in
		}
		return true
	})
	return in
}

func (ring *ringIndexed) IntersectsPoint(point Point, allowOnEdge bool) bool {
	return ring.ContainsPoint(point, allowOnEdge)
}

func (ring *ringIndexed) ContainsSegment(seg Segment, allowOnEdge bool) bool {
	return ringContainsSegment(ring, seg, allowOnEdge)
}

func (ring *ringIndexed) ContainsRing(other Ring, allowOnEdge bool) bool {
	return ringContainsRing(ring, other, allowOnEdge)
}

func (ring *ringIndexed) IntersectsRing(other Ring, allowOnEdge bool) bool {
	return ringIntersectsRing(ring, other, allowOnEdge)
}

func (ring *ringIndexed) ContainsRect(rect Rect, allowOnEdge bool) bool {
	return ringContainsRect(ring, rect, allowOnEdge)
}

func (ring *ringIndexed) IntersectsRect(rect Rect, allowOnEdge bool) bool {
	return ringIntersectsRect(ring, rect, allowOnEdge)
}

func (ring *ringIndexed) ContainsPoly(poly Poly, allowOnEdge bool) bool {
	return ring.ContainsRing(poly.Exterior(), allowOnEdge)
}

func (ring *ringIndexed) IntersectsPoly(poly Poly, allowOnEdge bool) bool {
	return ringIntersectsPoly(ring, poly, allowOnEdge)
}
