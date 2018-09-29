package geom

import (
	"math"

	"github.com/tidwall/boxtree/d2"
)

type treeLine struct {
	points []Point
	rect   Rect
	tree   d2.BoxTree
}

func newTreeLine(points []Point) *treeLine {
	var ring treeLine
	ring.points = make([]Point, len(points))
	copy(ring.points, points)
	ring.init()
	return &ring
}

func (ring *treeLine) move(deltaX, deltaY float64) *treeLine {
	points := make([]Point, len(ring.points))
	for i := 0; i < len(ring.points); i++ {
		points[i].X = ring.points[i].X + deltaX
		points[i].Y = ring.points[i].Y + deltaY
	}
	return newTreeLine(points)
}

func (ring *treeLine) init() {
	ring.rect = pointsRect(ring.points)
	for i := 0; i < len(ring.points)-1; i++ {
		rect := (Segment{A: ring.points[i], B: ring.points[i+1]}).Rect()
		ring.tree.Insert(
			[]float64{rect.Min.X, rect.Min.Y},
			[]float64{rect.Max.X, rect.Max.Y},
			i,
		)
	}
}

func (ring *treeLine) Points() []Point {
	return ring.points
}

func (ring *treeLine) IsClosed() bool {
	return true
}

func (ring *treeLine) Rect() Rect {
	return ring.rect
}

func (ring *treeLine) Search(
	rect Rect, iter func(seg Segment, index int) bool,
) {
	ring.tree.Search(
		[]float64{rect.Min.X, rect.Min.Y},
		[]float64{rect.Max.X, rect.Max.Y},
		func(_, _ []float64, value interface{}) bool {
			index := value.(int)
			seg := Segment{A: ring.points[index], B: ring.points[index+1]}
			if !iter(seg, index) {
				return false
			}
			return true
		},
	)
}

func (ring *treeLine) Scan(iter func(seg Segment) bool) {
	for i := 0; i < len(ring.points); i++ {
		seg := Segment{A: ring.points[i], B: ring.points[i+1]}
		if !iter(seg) {
			return
		}
	}
}

func (ring *treeLine) IntersectsSegment(seg Segment, allowOnEdge bool) bool {
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

func (ring *treeLine) ContainsPoint(point Point, allowOnEdge bool) bool {
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

// func (ring *treeLine) IntersectsPoint(point Point, allowOnEdge bool) bool {
// 	return ring.ContainsPoint(point, allowOnEdge)
// }

// func (ring *treeLine) ContainsSegment(seg Segment, allowOnEdge bool) bool {
// 	return ringContainsSegment(ring, seg, allowOnEdge)
// }

// func (ring *treeLine) ContainsRing(other Ring, allowOnEdge bool) bool {
// 	return ringContainsRing(ring, other, allowOnEdge)
// }

// func (ring *treeLine) IntersectsRing(other Ring, allowOnEdge bool) bool {
// 	return ringIntersectsRing(ring, other, allowOnEdge)
// }

// func (ring *treeLine) ContainsRect(rect Rect, allowOnEdge bool) bool {
// 	return ringContainsRect(ring, rect, allowOnEdge)
// }

// func (ring *treeLine) IntersectsRect(rect Rect, allowOnEdge bool) bool {
// 	return ringIntersectsRect(ring, rect, allowOnEdge)
// }

// func (ring *treeLine) ContainsPoly(poly Poly, allowOnEdge bool) bool {
// 	return ring.ContainsRing(poly.Exterior(), allowOnEdge)
// }

// func (ring *treeLine) IntersectsPoly(poly Poly, allowOnEdge bool) bool {
// 	return ringIntersectsPoly(ring, poly, allowOnEdge)
// }
