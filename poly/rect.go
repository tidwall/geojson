package poly

// Rect is rectangle
type Rect struct {
	Min, Max Point
}

// Ring returns an exterior ring
func (rect Rect) Ring() Ring {
	return Ring{
		{rect.Min.X, rect.Min.Y},
		{rect.Max.X, rect.Min.Y},
		{rect.Max.X, rect.Max.Y},
		{rect.Min.X, rect.Max.Y},
		{rect.Min.X, rect.Min.Y},
	}
}

// Polygon returns a polygon for the rect
func (rect Rect) Polygon() Polygon {
	return Polygon{rect.Ring(), nil}
}

// InsidePoint tests if rect is inside of a point
func (rect Rect) InsidePoint(point Point) bool {
	return rect.Min == point && rect.Max == point
}

// InsideRect tests if rect is inside of another rect
func (rect Rect) InsideRect(other Rect) bool {
	return rectInRect(rect, other)
}

// InsideLine tests if a rect is inside of a line
func (rect Rect) InsideLine(line Line) bool {
	if rect.Min == rect.Max {
		return rect.Min.InsideLine(line)
	}
	if rect.Min.X == rect.Max.X || rect.Min.Y == rect.Max.Y {
		for i := 0; i < len(line)-1; i++ {
			if segmentOnSegment(rect.Min, rect.Max, line[i], line[i+1]) {
				return true
			}
		}
		return false
	}
	return false
}

// InsideRing tests if a rect is inside of a ring
func (rect Rect) InsideRing(ring Ring) bool {
	// all four points should be inside the ring
	return (Point{rect.Min.X, rect.Min.Y}).InsideRing(ring) &&
		(Point{rect.Max.X, rect.Min.Y}).InsideRing(ring) &&
		(Point{rect.Max.X, rect.Max.Y}).InsideRing(ring) &&
		(Point{rect.Min.X, rect.Max.Y}).InsideRing(ring)
}

// InsidePolygon tests if a rect is inside a polygon
func (rect Rect) InsidePolygon(polygon Polygon) bool {
	// all four points should be inside the polygon
	return (Point{rect.Min.X, rect.Min.Y}).InsidePolygon(polygon) &&
		(Point{rect.Max.X, rect.Min.Y}).InsidePolygon(polygon) &&
		(Point{rect.Max.X, rect.Max.Y}).InsidePolygon(polygon) &&
		(Point{rect.Min.X, rect.Max.Y}).InsidePolygon(polygon)
}

// IntersectsPoint tests if a rects intersects a point
func (rect Rect) IntersectsPoint(point Point) bool {
	return point.IntersectsRect(rect)
}

// IntersectsRect tests if a rect intersects another rect
func (rect Rect) IntersectsRect(other Rect) bool {
	return rectIntersectsRect(rect, other)
}

// IntersectsLine tests if a rect intersects a line
func (rect Rect) IntersectsLine(line Line) bool {
	ring := Ring{
		{rect.Min.X, rect.Min.Y},
		{rect.Max.X, rect.Min.Y},
		{rect.Max.X, rect.Max.Y},
		{rect.Min.X, rect.Max.Y},
		{rect.Min.X, rect.Min.Y},
	}
	return doesIntersect(line, true, Polygon{ring, nil})
}

// IntersectsRing tests if a rect intersects a ring
func (rect Rect) IntersectsRing(ring Ring) bool {
	other := Ring{
		{rect.Min.X, rect.Min.Y},
		{rect.Max.X, rect.Min.Y},
		{rect.Max.X, rect.Max.Y},
		{rect.Min.X, rect.Max.Y},
		{rect.Min.X, rect.Min.Y},
	}
	return doesIntersect(ring, false, Polygon{other, nil})
}

// IntersectsPolygon tests if a rect intersects another polygon
func (rect Rect) IntersectsPolygon(polygon Polygon) bool {
	ring := Ring{
		{rect.Min.X, rect.Min.Y},
		{rect.Max.X, rect.Min.Y},
		{rect.Max.X, rect.Max.Y},
		{rect.Min.X, rect.Max.Y},
		{rect.Min.X, rect.Min.Y},
	}
	return doesIntersect(ring, false, polygon)
}
