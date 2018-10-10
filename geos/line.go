// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geos

// Line is a open series of points
type Line struct {
	baseSeries
}

// NewLine creates a new Line
func NewLine(points []Point, index int) *Line {
	line := new(Line)
	line.baseSeries = makeSeries(points, true, false, index)
	return line
}

// Move ...
func (line *Line) Move(deltaX, deltaY float64) *Line {
	nline := new(Line)
	nline.baseSeries = *line.baseSeries.Move(deltaX, deltaY).(*baseSeries)
	return nline
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
	lineNumSegments := line.NumSegments()
	for i := 0; i < lineNumSegments; i++ {
		segA := line.SegmentAt(i)
		var intersects bool
		other.Search(segA.Rect(), func(segB Segment, _ int) bool {
			if segA.IntersectsSegment(segB) {
				intersects = true
				return false
			}
			return true
		})
		if intersects {
			return true
		}
	}
	return false
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
	// polygon can fit in a straight (vertial or horizontal) line
	points := [2]Point{rect.Min, rect.Max}
	var other Line
	other.baseSeries.points = points[:]
	other.baseSeries.rect = rect
	return line.ContainsLine(&other)
}

// IntersectsPoly ...
func (line *Line) IntersectsPoly(poly *Poly) bool {
	return poly.IntersectsLine(line)
}