package geom

// Line ...
type Line struct {
	Series
}

// NewLine ...
func NewLine(points []Point) *Line {
	line := new(Line)
	line.Series = MakeSeries(points, true, false)
	return line
}

func (line *Line) move(deltaX, deltaY float64) *Line {
	points := make([]Point, len(line.points))
	for i := 0; i < len(line.points); i++ {
		points[i].X = line.points[i].X + deltaX
		points[i].Y = line.points[i].Y + deltaY
	}
	return NewLine(points)
}

// ContainsPoint ...
func (line *Line) ContainsPoint(point Point) bool {
	contains := false
	line.Search(Rect{point, point}, func(seg Segment, index int) bool {
		if seg.Raycast(point).On {
			contains = true
			return false
		}
		return true
	})
	return contains
}

// IntersectsPoint ...
func (line *Line) IntersectsPoint(point Point) bool {
	return line.ContainsPoint(point)
}

// ContainsRect ...
func (line *Line) ContainsRect(rect Rect) bool {
	return rect.Min == rect.Max && line.ContainsPoint(rect.Min)
}

// IntersectsRect ...
func (line *Line) IntersectsRect(rect Rect) bool {
	return rect.IntersectsLine(line)
}

// ContainsLine ...
func (line *Line) ContainsLine(other *Line) bool {

	panic("not ready")
}

// IntersectsLine ...
func (line *Line) IntersectsLine(other *Line) bool {
	if line.Empty() || other.Empty() {
		return false
	}
	if !line.Rect().IntersectsRect(other.Rect()) {
		return false
	}
	if line.NumPoints() > other.NumPoints() {
		line, other = other, line
	}
	line.ForEachSegment(func(segA Segment, _ int) bool {
		other.Search(segA.Rect(), func(segB Segment, _ int) bool {

			return true
		})
		return true
	})
	// for i := 0; i < len(line)-1; i++ {
	// 	for j := 0; j < len(other)-1; j++ {
	// 		if segmentsIntersect(
	// 			line[i], line[i+1],
	// 			other[i], other[i+1],
	// 		) {
	// 			return true
	// 		}
	// 	}
	// }
	// return false
	panic("not ready")
}

// ContainsPoly ...
func (line *Line) ContainsPoly(poly *Poly) bool {
	panic("not ready")
}

// IntersectsPoly ...
func (line *Line) IntersectsPoly(poly *Poly) bool {
	panic("not ready")
}
