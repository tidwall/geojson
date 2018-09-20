package poly

// Line is a series of points that make up a polyline
type Line []Point

// IntersectsLine detects if a linestring intersects a
// polyline assume shape and exterior are actually polylines
func (line Line) IntersectsLine(other Line) bool {
	for i := 0; i < len(line)-1; i++ {
		for j := 0; j < len(other)-1; j++ {
			if lineintersects(
				line[i], line[i+1],
				other[i], other[i+1],
			) {
				return true
			}
		}
	}
	return false
}

// IntersectsPoint detects if a linestring intersects a point. The
// point will need to be exactly on a segment of the linestring
func (line Line) IntersectsPoint(point Point) bool {
	return point.IntersectsLine(line)
}

// IntersectsRect detects if a linestring intersects a rect
func (line Line) IntersectsRect(rect Rect) bool {
	return line.IntersectsPolygon(rect.Polygon())
}

// IntersectsPolygon detects if a polygon intersects a linestring
// assume shape is a linestring
func (line Line) IntersectsPolygon(polygon Polygon) bool {
	return doesIntersects(line, true, polygon.Exterior, polygon.Holes)
}
