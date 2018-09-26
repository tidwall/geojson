package poly

// // forEachSegment will return each segment on a line, ring, or polygon. When
// // close is provided, the last segment will be guarenteed to be end where the
// // the first started. When there is only one point, a single segment will be
// // returned where point 'a' and 'b' are equal
// func forEachSegment(points []Point, close bool, iter func(a, b Point) bool) {
// 	if len(points) == 1 {
// 		if !iter(points[0], points[0]) {
// 			return
// 		}
// 	}
// 	//Math.atan2(p2.y - p1.y, p2.x - p1.x);

// }

// pointInRing return true when point is inside a ring. When then exterior
// param is provided, the function will return false if the point is on the
// edge of the ring.
func pointInRing(point Point, ring Ring, exterior bool) bool {
	// if len(shape) < 3 {
	// 	return false
	// }
	in := false
	for i := 0; i < len(ring); i++ {
		res := raycast(point, ring[i], ring[(i+1)%len(ring)])
		if res.on {
			return exterior
		}
		if res.in {
			in = !in
		}
	}
	return in
}

// segmentOnSegment returns true when a line segment is on another segment
func segmentOnSegment(a, b, c, d Point) bool {
	return raycast(a, c, d).on && raycast(b, c, d).on
}

// segmentsIntersect returns true when two line segments intersect with each other
func segmentsIntersect(
	a, b Point, // segment 1
	c, d Point, // segment 2
) bool {
	// do the bounding boxes intersect?
	// the following checks without swapping values.
	if a.Y > b.Y {
		if c.Y > d.Y {
			if b.Y > c.Y || a.Y < d.Y {
				return false
			}
		} else {
			if b.Y > d.Y || a.Y < c.Y {
				return false
			}
		}
	} else {
		if c.Y > d.Y {
			if a.Y > c.Y || b.Y < d.Y {
				return false
			}
		} else {
			if a.Y > d.Y || b.Y < c.Y {
				return false
			}
		}
	}
	if a.X > b.X {
		if c.X > d.X {
			if b.X > c.X || a.X < d.X {
				return false
			}
		} else {
			if b.X > d.X || a.X < c.X {
				return false
			}
		}
	} else {
		if c.X > d.X {
			if a.X > c.X || b.X < d.X {
				return false
			}
		} else {
			if a.X > d.X || b.X < c.X {
				return false
			}
		}
	}

	// the following code is from http://ideone.com/PnPJgb
	cmpx, cmpy := c.X-a.X, c.Y-a.Y
	rx, ry := b.X-a.X, b.Y-a.Y
	cmpxr := cmpx*ry - cmpy*rx
	if cmpxr == 0 {
		// Lines are collinear, and so intersect if they have any overlap
		if !(((c.X-a.X <= 0) != (c.X-b.X <= 0)) ||
			((c.Y-a.Y <= 0) != (c.Y-b.Y <= 0))) {
			return false
		}
		return true
	}
	sx, sy := d.X-c.X, d.Y-c.Y
	cmpxs := cmpx*sy - cmpy*sx
	rxs := rx*sy - ry*sx
	if rxs == 0 {
		return false // segments are parallel.
	}
	rxsr := 1 / rxs
	t := cmpxs * rxsr
	u := cmpxr * rxsr
	if !((t >= 0) && (t <= 1) && (u >= 0) && (u <= 1)) {
		return false
	}
	return true
}

// pointOnLine returns true when a point is on a linestring
func pointOnLine(point Point, line Line) bool {
	if len(line) == 1 {
		return line[0] == point
	}
	for j := 0; j < len(line); j++ {
		if raycast(point, line[j], line[(j+1)%len(line)]).on {
			return true
		}
	}
	return false
}

func pointInRect(point Point, rect Rect) bool {
	if point.X < rect.Min.X || point.X > rect.Max.X {
		return false
	}
	if point.Y < rect.Min.Y || point.Y > rect.Max.Y {
		return false
	}
	return true
}

func pointInPolygon(point Point, polygon Polygon) bool {
	if !pointInRing(point, polygon.Exterior, true) {
		return false
	}
	for _, hole := range polygon.Holes {
		if pointInRing(point, hole, false) {
			return false
		}
	}
	return true
}

func rectInRect(a, b Rect) bool {
	if a.Min.X < b.Min.X || a.Max.X > b.Max.X {
		return false
	}
	if a.Min.Y < b.Min.Y || a.Max.Y > b.Max.Y {
		return false
	}
	return true
}

func rectIntersectsRect(a, b Rect) bool {
	if a.Min.Y > b.Max.Y || a.Max.Y < b.Min.Y {
		return false
	}
	if a.Min.X > b.Max.X || a.Max.X < b.Min.X {
		return false
	}
	return true
}

func ringInPolygon(ring Ring, polygon Polygon) bool {
	// all points in ring must be inside of the polygon
	for _, p := range ring {
		if !pointInPolygon(p, polygon) {
			return false
		}
	}
	// // all ring segments *must not* intersect with polygon segments
	// ext := polygon.Exterior
	// for i := 0; i < len(ring); i++ {
	// 	ringA, ringB := ring[i], ring[(i+1)%len(ring)]
	// 	for j := 0; j < len(ext); j++ {
	// 		extA, extB := ext[j], ext[(j+1)%len(ext)]
	// 		if segmentsIntersect(ringA, ringB, extA, extB) {
	// 			return false
	// 		}
	// 	}
	// }

	ok := true
	for _, hole := range polygon.Holes {
		if ringInPolygon(hole, Polygon{ring, nil}) {
			return false
		}
	}
	return ok
}

func doesIntersect(
	points []Point, pointsAreALineString bool, polygon Polygon,
) bool {
	switch len(points) {
	case 0:
		return false
	case 1:
		switch len(polygon.Exterior) {
		case 0:
			return false
		case 1:
			return points[0].X == polygon.Exterior[0].X &&
				points[0].Y == points[0].Y
		default:
			return pointInPolygon(points[0], polygon)
		}
	default:
		switch len(polygon.Exterior) {
		case 0:
			return false
		case 1:
			return pointInPolygon(polygon.Exterior[0],
				Polygon{points, polygon.Holes})
		}
	}
	if !rectIntersectsRect(Ring(points).Rect(), polygon.Exterior.Rect()) {
		return false
	}
	for i := 0; i < len(points); i++ {
		for j := 0; j < len(polygon.Exterior); j++ {
			if segmentsIntersect(
				points[i], points[(i+1)%len(points)],
				polygon.Exterior[j],
				polygon.Exterior[(j+1)%len(polygon.Exterior)],
			) {
				return true
			}
		}
	}
	for _, hole := range polygon.Holes {
		if ringInPolygon(Ring(points), Polygon{hole, nil}) {
			return false
		}
	}
	if ringInPolygon(Ring(points), Polygon{polygon.Exterior, nil}) {
		return true
	}
	if !pointsAreALineString {
		if ringInPolygon(polygon.Exterior, Polygon{points, nil}) {
			return true
		}
	}
	return false
}
