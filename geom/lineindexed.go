package geom

import (
	"github.com/tidwall/boxtree/d2"
)

type lineIndexed struct {
	points []Point
	rect   Rect
	tree   d2.BoxTree
}

func newLineIndexed(points []Point) *lineIndexed {
	var line lineIndexed
	line.points = make([]Point, len(points))
	copy(line.points, points)
	line.init()
	return &line
}

func (line *lineIndexed) move(deltaX, deltaY float64) *lineIndexed {
	points := make([]Point, len(line.points))
	for i := 0; i < len(line.points); i++ {
		points[i].X = line.points[i].X + deltaX
		points[i].Y = line.points[i].Y + deltaY
	}
	return newLineIndexed(points)
}

func (line *lineIndexed) init() {
	line.rect = pointsRect(line.points)
	for i := 0; i < len(line.points)-1; i++ {
		seg := Segment{line.points[i], line.points[i+1]}
		rect := seg.Rect()
		line.tree.Insert(
			[]float64{rect.Min.X, rect.Min.Y},
			[]float64{rect.Max.X, rect.Max.Y},
			i,
		)
	}
}

func (line *lineIndexed) Points() []Point {
	return line.points
}

func (line *lineIndexed) IsClosed() bool {
	return false
}

func (line *lineIndexed) Rect() Rect {
	return line.rect
}

func (line *lineIndexed) Search(
	rect Rect, iter func(seg Segment, index int) bool,
) {
	line.tree.Search(
		[]float64{rect.Min.X, rect.Min.Y},
		[]float64{rect.Max.X, rect.Max.Y},
		func(_, _ []float64, value interface{}) bool {
			index := value.(int)
			seg := Segment{line.points[index], line.points[index+1]}
			if !iter(seg, index) {
				return false
			}
			return true
		},
	)
}

func (line *lineIndexed) Scan(iter func(seg Segment) bool) {
	for i := 0; i < len(line.points)-1; i++ {
		seg := Segment{line.points[i], line.points[i+1]}
		if !iter(seg) {
			return
		}
	}
}

func (line *lineIndexed) ContainsPoint(point Point) bool {
	contains := false
	line.Search(Rect{point, point}, func(seg Segment, index int) bool {
		if raycast(point, seg.A, seg.B).on {
			contains = true
			return false
		}
		return true
	})
	return contains
}

func (line *lineIndexed) IntersectsPoint(point Point) bool {
	return line.ContainsPoint(point)
}

// func (ring *lineIndexed) IntersectsPoint(point Point, allowOnEdge bool) bool {
// 	return ring.ContainsPoint(point, allowOnEdge)
// }

// func (ring *lineIndexed) ContainsSegment(seg Segment, allowOnEdge bool) bool {
// 	return ringContainsSegment(ring, seg, allowOnEdge)
// }

// func (ring *lineIndexed) ContainsRing(other Ring, allowOnEdge bool) bool {
// 	return ringContainsRing(ring, other, allowOnEdge)
// }

// func (ring *lineIndexed) IntersectsRing(other Ring, allowOnEdge bool) bool {
// 	return ringIntersectsRing(ring, other, allowOnEdge)
// }

// func (ring *lineIndexed) ContainsRect(rect Rect, allowOnEdge bool) bool {
// 	return ringContainsRect(ring, rect, allowOnEdge)
// }

// func (ring *lineIndexed) IntersectsRect(rect Rect, allowOnEdge bool) bool {
// 	return ringIntersectsRect(ring, rect, allowOnEdge)
// }

// func (ring *lineIndexed) ContainsPoly(poly Poly, allowOnEdge bool) bool {
// 	return ring.ContainsRing(poly.Exterior(), allowOnEdge)
// }

// func (ring *lineIndexed) IntersectsPoly(poly Poly, allowOnEdge bool) bool {
// 	return ringIntersectsPoly(ring, poly, allowOnEdge)
// }
