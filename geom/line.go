package geom

// Line is a open series of points
type Line struct {
	baseSeries
}

// NewLine creates a new Line
func NewLine(points []Point) *Line {
	line := new(Line)
	line.baseSeries = makeSeries(points, true, false)
	return line
}

// Move ...
func (line *Line) Move(deltaX, deltaY float64) *Line {
	nline := new(Line)
	nline.baseSeries = *line.baseSeries.Move(deltaX, deltaY).(*baseSeries)
	return nline
}

// Clockwise ...
func (line *Line) Clockwise() bool {
	return line.Clockwise()
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
	// Convert rect into a poly
	return line.ContainsPoly(&Poly{Exterior: rect})
}

// IntersectsRect ...
func (line *Line) IntersectsRect(rect Rect) bool {
	return rect.IntersectsLine(line)
}

// ContainsLine ...
func (line *Line) ContainsLine(other *Line) bool {
	if line.Empty() || other.Empty() {
		return false
	}
	// locate the first "other" segment that contains the first "line" segment.
	lineNumSegments := line.NumSegments()
	segIdx := -1
	for j := 0; j < lineNumSegments; j++ {
		if line.SegmentAt(j).ContainsSegment(other.SegmentAt(0)) {
			segIdx = j
			break
		}
	}
	if segIdx == -1 {
		return false
	}
	otherNumSegments := other.NumSegments()
	for i := 1; i < otherNumSegments; i++ {
		lineSeg := line.SegmentAt(segIdx)
		otherSeg := other.SegmentAt(i)
		if lineSeg.ContainsSegment(otherSeg) {
			continue
		}
		if otherSeg.A == lineSeg.A {
			// reverse it
			if segIdx == 0 {
				return false
			}
			segIdx--
			i--
		} else if otherSeg.A == lineSeg.B {
			// forward it
			if segIdx == lineNumSegments-1 {
				return false
			}
			segIdx++
			i--
		}
	}
	return true
}

// ContainsSegment ...
func (line *Line) ContainsSegment(seg Segment) bool {
	var contains bool
	line.Search(seg.Rect(), func(other Segment, index int) bool {
		if other.Raycast(seg.A).On && other.Raycast(seg.B).On {
			contains = true
			return false
		}
		return true
	})
	return contains
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
	var intersects bool
	seriesForEachSegment(line, func(segA Segment) bool {
		other.Search(segA.Rect(), func(segB Segment, _ int) bool {
			if segA.IntersectsSegment(segB) {
				intersects = true
				return false
			}
			return true
		})
		return !intersects
	})
	return intersects
}

// ContainsPoly ...
func (line *Line) ContainsPoly(poly *Poly) bool {
	if line.Empty() || poly.Empty() {
		return false
	}
	rect := poly.Rect()
	if rect.Min.X != rect.Max.X && rect.Min.Y != rect.Max.Y {
		return false
	}

	return line.ContainsSegment(Segment{A: rect.Min, B: rect.Max})
}

// IntersectsPoly ...
func (line *Line) IntersectsPoly(poly *Poly) bool {
	return poly.IntersectsLine(line)
}
