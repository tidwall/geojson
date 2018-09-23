package poly

import (
	"strconv"
)

// Ring is series of points that make up a closed shape
type Ring []Point

// String returns a string representation of the polygon.
func (ring Ring) String() string {
	var b []byte
	b = append(b, '[')
	for i, p := range ring {
		if i > 0 {
			b = append(b, ',')
		}
		b = append(b, '[')
		b = strconv.AppendFloat(b, p.X, 'f', -1, 64)
		b = append(b, ',')
		b = strconv.AppendFloat(b, p.Y, 'f', -1, 64)
		b = append(b, ']')
	}
	b = append(b, ']')
	return string(b)
}

// Rect returns the bounding box rectangle for the polygon
func (ring Ring) Rect() Rect {
	var bbox Rect
	for i, p := range ring {
		if i == 0 {
			bbox.Min = p
			bbox.Max = p
		} else {
			if p.X < bbox.Min.X {
				bbox.Min.X = p.X
			} else if p.X > bbox.Max.X {
				bbox.Max.X = p.X
			}
			if p.Y < bbox.Min.Y {
				bbox.Min.Y = p.Y
			} else if p.Y > bbox.Max.Y {
				bbox.Max.Y = p.Y
			}
		}
	}
	return bbox
}

// InsidePoint tests if ring is inside a point
func (ring Ring) InsidePoint(point Point) bool {
	if len(ring) == 0 {
		return false
	}
	for _, p := range ring {
		if !p.InsidePoint(point) {
			return false
		}
	}
	return true
}

// InsideRect tests if ring is inside of a rect
func (ring Ring) InsideRect(rect Rect) bool {
	if len(ring) == 0 {
		return false
	}
	for _, p := range ring {
		if !p.InsideRect(rect) {
			return false
		}
	}
	return true
}

// InsideLine tests if ring is inside of a line
func (ring Ring) InsideLine(line Line) bool {
	if len(ring) == 0 {
		return false
	}
	return ring.Rect().InsideLine(line)
}

// InsideRing detects if a rect is inside of another ring
func (ring Ring) InsideRing(other Ring) bool {
	return ringInPolygon(ring, Polygon{other, nil})
}

// InsidePolygon returns true if shape is inside of exterior and not in a hole.
func (ring Ring) InsidePolygon(polygon Polygon) bool {
	return ringInPolygon(ring, polygon)
}

// IntersectsPoint test if ring intersects rect
func (ring Ring) IntersectsPoint(point Point) bool {
	return point.IntersectsRing(ring)
}

// IntersectsRect test if ring intersects rect
func (ring Ring) IntersectsRect(rect Rect) bool {
	return rect.IntersectsRing(ring)
}

// IntersectsLine test if ring intersects rect
func (ring Ring) IntersectsLine(line Line) bool {
	return doesIntersect(line, true, Polygon{ring, nil})
}

// IntersectsRing test if ring intersects other ring
func (ring Ring) IntersectsRing(other Ring) bool {
	return doesIntersect(ring, true, Polygon{other, nil})
}

// IntersectsPolygon detects if a polygon intersects another polygon
func (ring Ring) IntersectsPolygon(polygon Polygon) bool {
	return doesIntersect(ring, false, polygon)
}
