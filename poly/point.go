package poly

// Point is simple 2D point
// For geo locations: X is lat, Y is lon, and Z is elev or time measure.
type Point struct {
	X, Y float64
}

// InsidePoint tests if point is inside of another point
func (point Point) InsidePoint(other Point) bool {
	return point == other
}

// InsideRect tests if point is inside of a rect
func (point Point) InsideRect(rect Rect) bool {
	return pointInRect(point, rect)
}

// InsideLine tests if point is inside of a linestring
func (point Point) InsideLine(line Line) bool {
	return pointOnLine(point, line)
}

// InsideRing tests if point is inside of a ring
func (point Point) InsideRing(ring Ring) bool {
	return pointInPolygon(point, Polygon{ring, nil})
}

// InsidePolygon returns true if point is inside of exterior and not in a hole.
// The validity of the exterior and holes must be done elsewhere and are
// assumed valid.
//   A valid exterior is a near-linear ring.
//   A valid hole is one that is full contained inside the exterior.
//   A valid hole may not share the same segment line as the exterior.
func (point Point) InsidePolygon(polygon Polygon) bool {
	return pointInPolygon(point, polygon)
}

// IntersectsPoint tests if if a point intersects another point
func (point Point) IntersectsPoint(other Point) bool {
	return point.InsidePoint(other)
}

// IntersectsRect tests if if a point intersects a rect
func (point Point) IntersectsRect(rect Rect) bool {
	return point.InsideRect(rect)
}

// IntersectsLine tests if if a point intersects a linestring
func (point Point) IntersectsLine(line Line) bool {
	return point.InsideLine(line)
}

// IntersectsRing tests if if a point intersects a ring
func (point Point) IntersectsRing(ring Ring) bool {
	return point.InsideRing(ring)
}

// IntersectsPolygon tests if if a point intersects a polygon
func (point Point) IntersectsPolygon(polygon Polygon) bool {
	return point.InsidePolygon(polygon)
}
