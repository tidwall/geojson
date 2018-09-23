package poly

// Polygon is a closed shape that is of an exterior ring and interior holes.
type Polygon struct {
	Exterior Ring
	Holes    []Ring
}

// InsidePoint tests if polygon is inside a point
func (polygon Polygon) InsidePoint(point Point) bool {
	return polygon.Exterior.InsidePoint(point)
}

// InsideRect tests if polygon is inside a rect
func (polygon Polygon) InsideRect(rect Rect) bool {
	return polygon.Exterior.InsideRect(rect)
}

// InsideLine tests if polygon is inside a line
func (polygon Polygon) InsideLine(line Line) bool {
	return polygon.Exterior.InsideLine(line)
}

// InsideRing tests if polygon is inside a ring
func (polygon Polygon) InsideRing(ring Ring) bool {
	return polygon.Exterior.InsideRing(ring)
}

// InsidePolygon tests if polygon is inside another polygon
func (polygon Polygon) InsidePolygon(other Polygon) bool {
	// TODO: better hole detection
	return polygon.Exterior.InsidePolygon(other)
}

// IntersectsPoint tests if polygon intersects a point
func (polygon Polygon) IntersectsPoint(point Point) bool {
	return point.IntersectsPolygon(polygon)
}

// IntersectsRect tests if polygon intersects a rect
func (polygon Polygon) IntersectsRect(rect Rect) bool {
	return rect.IntersectsPolygon(polygon)
}

// IntersectsLine tests if polygon intersects a linestring
func (polygon Polygon) IntersectsLine(line Line) bool {
	return line.IntersectsPolygon(polygon)
}

// IntersectsRing tests if polygon intersects a ring
func (polygon Polygon) IntersectsRing(ring Ring) bool {
	return ring.IntersectsPolygon(polygon)
}

// IntersectsPolygon tests if polygon intersects another polygon
func (polygon Polygon) IntersectsPolygon(other Polygon) bool {
	// TODO: better hole detection
	return polygon.Exterior.IntersectsPolygon(polygon)
}
