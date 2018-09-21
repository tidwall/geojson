package poly

// Line is a series of points that make up a polyline
type Line []Point

// InsidePoint tests if line is inside point
func (line Line) InsidePoint(point Point) bool {
	if len(line) == 0 {
		return false
	}
	for _, p := range line {
		if p != point {
			return false
		}
	}
	return true
}

// InsideRect tests if line is inside rect
func (line Line) InsideRect(rect Rect) bool {
	if len(line) == 0 {
		return false
	}
	for _, p := range line {
		if !p.InsideRect(rect) {
			return false
		}
	}
	return true
}

// InsideLine tests if line is inside other line.
func (line Line) InsideLine(other Line) bool {
	if len(line) == 0 {
		return false
	}
	if len(other) == 0 {
		return false
	}
	if len(line) == 1 {
		if len(other) == 1 {
			return line[0] == other[0]
		}
		return line[0].InsideLine(other)
	}
	if len(other) == 1 {
		for _, p := range line {
			if p != other[0] {
				return false
			}
		}
		return true
	}
	// locate the first "other" segment that contains the first "line" segment.
	segIdx := -1
	for j := 0; j < len(other)-1; j++ {
		if segmentOnSegment(line[0], line[1], other[j], other[j+1]) {
			segIdx = j
			break
		}
	}
	if segIdx == -1 {
		return false
	}
	for i := 1; i < len(line)-1; i++ {
		if segmentOnSegment(line[i], line[i+1], other[segIdx], other[segIdx+1]) {
			continue
		}
		if line[i] == other[segIdx] {
			// reverse it
			if segIdx == 0 {
				return false
			}
			segIdx--
			i--
		} else if line[i] == other[segIdx+1] {
			// forward it
			if segIdx == len(other)-2 {
				return false
			}
			segIdx++
			i--
		}
	}
	return true
}

// InsideRing tests if line is inside a ring
func (line Line) InsideRing(ring Ring) bool {
	if len(line) == 0 {
		return false
	}
	if len(ring) == 0 {
		return false
	}
	for _, p := range line {
		if !p.InsideRing(ring) {
			return false
		}
	}
	return true
}

// InsidePolygon tests if line is inside a polygon
func (line Line) InsidePolygon(polygon Polygon) bool {
	if !line.InsideRing(polygon.Exterior) {
		return false
	}
	for _, hole := range polygon.Holes {
		if doesIntersect(line, true, Polygon{hole, nil}) {
			return false
		}
	}
	return true
}

// IntersectsPoint detects if a linestring intersects a point. The
// point will need to be exactly on a segment of the linestring
func (line Line) IntersectsPoint(point Point) bool {
	return point.IntersectsLine(line)
}

// IntersectsRect detects if a linestring intersects a rect
func (line Line) IntersectsRect(rect Rect) bool {
	return rect.IntersectsLine(line)
}

// IntersectsLine detects if a linestring intersects a
// polyline assume shape and exterior are actually polylines
func (line Line) IntersectsLine(other Line) bool {
	for i := 0; i < len(line)-1; i++ {
		for j := 0; j < len(other)-1; j++ {
			if segmentsIntersect(
				line[i], line[i+1],
				other[i], other[i+1],
			) {
				return true
			}
		}
	}
	return false
}

// IntersectsRing detects if a linestring intersects a ring.
func (line Line) IntersectsRing(ring Ring) bool {
	return doesIntersect(line, true, Polygon{ring, nil})
}

// IntersectsPolygon detects if a line intersects a polygon.
func (line Line) IntersectsPolygon(polygon Polygon) bool {
	return doesIntersect(line, true, polygon)
}
