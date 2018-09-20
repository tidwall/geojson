package poly

import "strconv"

// Ring is series of points that make up a closed shape
type Ring []Point

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

// IntersectsRect test if ring intersects rect
func (ring Ring) IntersectsRect(rect Rect) bool {
	if len(ring) == 0 {
		return false
	}
	return ring.Intersects(rect.Polygon(), nil)
}

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

// Inside returns true if shape is inside of exterior and not in a hole.
func (ring Ring) Inside(exterior Ring, holes []Ring) bool {
	var ok bool
	for _, p := range ring {
		ok = p.Inside(exterior, holes)
		if !ok {
			return false
		}
	}
	ok = true
	for _, hole := range holes {
		if hole.Inside(ring, nil) {
			return false
		}
	}
	return ok
}

// Intersects detects if a polygon intersects another polygon
func (ring Ring) Intersects(exterior Ring, holes []Ring) bool {
	return ring.doesIntersects(false, exterior, holes)
}

// LineStringIntersectsLineString detects if a linestring intersects a
// linestring assume shape and exterior are actually linestrings
func (ring Ring) LineStringIntersectsLineString(exterior Ring) bool {
	for i := 0; i < len(ring)-1; i++ {
		for j := 0; j < len(exterior)-1; j++ {
			if lineintersects(
				ring[i], ring[i+1],
				exterior[i], exterior[i+1],
			) {
				return true
			}
		}
	}
	return false
}

// LineStringIntersectsPoint detects if a linestring intersects a point. The
// point will need to be exactly on a segment of the linestring
func (ring Ring) LineStringIntersectsPoint(point Point) bool {
	return point.IntersectsLineString(ring)
}

// LineStringIntersectsRect detects if a linestring intersects a rect
func (ring Ring) LineStringIntersectsRect(rect Rect) bool {
	return ring.LineStringIntersects(rect.Polygon(), nil)
}

// LineStringIntersects detects if a polygon intersects a linestring
// assume shape is a linestring
func (ring Ring) LineStringIntersects(
	exterior Ring, holes []Ring,
) bool {
	return ring.doesIntersects(true, exterior, holes)
}
func (ring Ring) doesIntersects(
	isLineString bool, exterior Ring, holes []Ring,
) bool {
	switch len(ring) {
	case 0:
		return false
	case 1:
		switch len(exterior) {
		case 0:
			return false
		case 1:
			return ring[0].X == exterior[0].X && ring[0].Y == ring[0].Y
		default:
			return ring[0].Inside(exterior, holes)
		}
	default:
		switch len(exterior) {
		case 0:
			return false
		case 1:
			return exterior[0].Inside(ring, holes)
		}
	}
	if !ring.Rect().IntersectsRect(exterior.Rect()) {
		return false
	}
	for i := 0; i < len(ring); i++ {
		for j := 0; j < len(exterior); j++ {
			if lineintersects(
				ring[i], ring[(i+1)%len(ring)],
				exterior[j], exterior[(j+1)%len(exterior)],
			) {
				return true
			}
		}
	}
	for _, hole := range holes {
		if ring.Inside(hole, nil) {
			return false
		}
	}
	if ring.Inside(exterior, nil) {
		return true
	}
	if !isLineString {
		if exterior.Inside(ring, nil) {
			return true
		}
	}
	return false
}
