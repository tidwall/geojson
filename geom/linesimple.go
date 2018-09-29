package geom

type lineSimple struct {
	points []Point
}

func (line *lineSimple) move(deltaX, deltaY float64) *lineSimple {
	points := make([]Point, len(line.points))
	for i := 0; i < len(line.points); i++ {
		points[i].X = line.points[i].X + deltaX
		points[i].Y = line.points[i].Y + deltaY
	}
	return newLineSimple(points)
}

func newLineSimple(points []Point) *lineSimple {
	ring := new(lineSimple)
	ring.points = make([]Point, len(points))
	copy(ring.points, points)
	return ring
}

func (line *lineSimple) Points() []Point {
	return line.points
}

func (line *lineSimple) Search(
	rect Rect, iter func(seg Segment, index int) bool,
) {
	var index int
	line.Scan(func(seg Segment) bool {
		if seg.Rect().IntersectsRect(rect) {
			if !iter(seg, index) {
				return false
			}
		}
		index++
		return true
	})
}
func (line *lineSimple) Scan(iter func(seg Segment) bool) {
	for i := 0; i < len(line.points)-1; i++ {
		if !iter(Segment{A: line.points[i], B: line.points[i+1]}) {
			return
		}
	}
}

func (line *lineSimple) Rect() Rect {
	return pointsRect(line.points)
}

func (line *lineSimple) IsClosed() bool {
	return false
}

func (line *lineSimple) ContainsPoint(point Point) bool {
	for i := 0; i < len(line.points)-1; i++ {
		if raycast(point, line.points[i], line.points[i+1]).on {
			return true
		}
	}
	return false
}

func (line *lineSimple) IntersectsPoint(point Point) bool {
	return line.ContainsPoint(point)
}

// func (line *simpleLine) IntersectsPoint(point Point, allowOnEdge bool) bool {
// 	return ring.ContainsPoint(point, allowOnEdge)
// }

// func (line *simpleLine) ContainsSegment(seg Segment, allowOnEdge bool) bool {
// 	return ringContainsSegment(ring, seg, allowOnEdge)
// }

// func (line *simpleLine) ContainsRing(other Ring, allowOnEdge bool) bool {
// 	return ringContainsRing(ring, other, allowOnEdge)
// }

// func (line *simpleLine) IntersectsRing(other Ring, allowOnEdge bool) bool {
// 	return ringIntersectsRing(ring, other, allowOnEdge)
// }

// func (line *simpleLine) ContainsRect(rect Rect, allowOnEdge bool) bool {
// 	return ringContainsRect(ring, rect, allowOnEdge)
// }

// func (line *simpleLine) IntersectsRect(rect Rect, allowOnEdge bool) bool {
// 	return ringIntersectsRect(ring, rect, allowOnEdge)
// }

// func (line *simpleLine) ContainsPoly(poly Poly, allowOnEdge bool) bool {
// 	return ring.ContainsRing(poly.Exterior(), allowOnEdge)
// }

// func (line *simpleLine) IntersectsPoly(poly Poly, allowOnEdge bool) bool {
// 	return ringIntersectsPoly(ring, poly, allowOnEdge)
// }
