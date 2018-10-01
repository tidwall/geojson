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
	// Convert rect into a poly
	return line.ContainsPoly(&Poly{Exterior: rect})
}

// IntersectsRect ...
func (line *Line) IntersectsRect(rect Rect) bool {
	return rect.IntersectsLine(line)
}

func lineInsideLine(line, other *Line) bool {
	if line.Empty() || other.Empty() {
		return false
	}
	panic("123")
	// // locate the first "other" segment that contains the first "line" segment.

	// segIdx := -1
	// for i
	// for j := 0; j < len(other)-1; j++ {
	// 	if segmentOnSegment(line[0], line[1], other[j], other[j+1]) {
	// 		segIdx = j
	// 		break
	// 	}
	// }
	// if segIdx == -1 {
	// 	return false
	// }
	// for i := 1; i < len(line)-1; i++ {
	// 	if segmentOnSegment(line[i], line[i+1], other[segIdx], other[segIdx+1]) {
	// 		continue
	// 	}
	// 	if line[i] == other[segIdx] {
	// 		// reverse it
	// 		if segIdx == 0 {
	// 			return false
	// 		}
	// 		segIdx--
	// 		i--
	// 	} else if line[i] == other[segIdx+1] {
	// 		// forward it
	// 		if segIdx == len(other)-2 {
	// 			return false
	// 		}
	// 		segIdx++
	// 		i--
	// 	}
	// }
	// return true
}

// ContainsLine ...
func (line *Line) ContainsLine(other *Line) bool {
	// if line.Empty() || other.Empty() {
	// 	return false
	// }
	// if !line.Rect().ContainsRect(other.Rect()) {
	// 	return false
	// }
	// other.ForEachSegment(func(segA Segment, idx int) bool {
	// 	contains
	// 	line.Search(segA.Rect(), func(segB Segment, idx int) bool {
	// 		return true
	// 	})
	// 	return true
	// })
	panic("not ready")
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
	line.ForEachSegment(func(segA Segment, _ int) bool {
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
