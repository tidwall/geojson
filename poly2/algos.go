package poly

import "math"

func algoPointInRing(point Point, ring Ring, exterior bool) bool {
	in := false
	var a, b Point
	for i := 0; i < len(ring); i++ {
		a = ring[i]
		if i == len(ring)-1 {
			b = ring[0]
		} else {
			b = ring[i+1]
		}
		res := algoRaycast(point, a, b)
		if res.on {
			return exterior
		}
		if res.in {
			in = !in
		}
	}
	return in
}

func algoPointInPolygon(point Point, polygon Polygon) bool {
	if !algoPointInRing(point, polygon.Exterior, true) {
		return false
	}
	for _, hole := range polygon.Holes {
		if algoPointInRing(point, hole, false) {
			return false
		}
	}
	return true
}

// algoAnySegmentsIntersect return true when any segment intersects from two
// line series.
func algoAnySegmentIntersects(lineA, lineB []Point, closeA, closeB bool) bool {
	var a, b Point
	var c, d Point
	for i := 0; i < len(lineA); i++ {
		a = lineA[i]
		if i == len(lineA)-1 {
			if !closeA {
				break
			}
			b = lineA[0]
		} else {
			b = lineA[i+1]
		}
		for j := 0; j < len(lineB); j++ {
			c = lineB[j]
			if j == len(lineB)-1 {
				if !closeB {
					break
				}
				d = lineB[0]
			} else {
				d = lineB[j+1]
			}
			if algoSegmentsIntersect(a, b, c, d) {
				return true
			}
		}
	}
	return false
}

// algoPointOnLine returns true when a point is on a linestring
func algoPointOnLine(point Point, line Line) bool {
	if len(line) == 1 {
		return line[0] == point
	}
	var a, b Point
	for i := 0; i < len(line); i++ {
		a = line[i]
		if i == len(line)-1 {
			b = line[0]
		} else {
			b = line[i+1]
		}
		if algoRaycast(point, a, b).on {
			return true
		}
	}
	return false
}

// algoSegmentsIntersect returns true when two line segments intersect
func algoSegmentsIntersect(
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

type rayres struct {
	in, on bool
}

func algoRaycast(p, a, b Point) rayres {
	// make sure that the point is inside the segment bounds
	if a.Y < b.Y && (p.Y < a.Y || p.Y > b.Y) {
		return rayres{false, false}
	} else if a.Y > b.Y && (p.Y < b.Y || p.Y > a.Y) {
		return rayres{false, false}
	}

	// test if point is in on the segment
	if a.Y == b.Y {
		if a.X == b.X {
			if p == a {
				return rayres{false, true}
			}
			return rayres{false, false}
		}
		if p.Y == b.Y {
			// horizontal segment
			// check if the point in on the line
			if a.X < b.X {
				if p.X >= a.X && p.X <= b.X {
					return rayres{false, true}
				}
			} else {
				if p.X >= b.X && p.X <= a.X {
					return rayres{false, true}
				}
			}
		}
	}
	if a.X == b.X && p.X == b.X {
		// vertical segment
		// check if the point in on the line
		if a.Y < b.Y {
			if p.Y >= a.Y && p.Y <= b.Y {
				return rayres{false, true}
			}
		} else {
			if p.Y >= b.Y && p.Y <= a.Y {
				return rayres{false, true}
			}
		}
	}
	if (p.X-a.X)/(b.X-a.X) == (p.Y-a.Y)/(b.Y-a.Y) {
		return rayres{false, true}
	}

	// do the actual raycast here.
	for p.Y == a.Y || p.Y == b.Y {
		p.Y = math.Nextafter(p.Y, math.Inf(1))
	}
	if a.Y < b.Y {
		if p.Y < a.Y || p.Y > b.Y {
			return rayres{false, false}
		}
	} else {
		if p.Y < b.Y || p.Y > a.Y {
			return rayres{false, false}
		}
	}
	if a.X > b.X {
		if p.X > a.X {
			return rayres{false, false}
		}
		if p.X < b.X {
			return rayres{true, false}
		}
	} else {
		if p.X > b.X {
			return rayres{false, false}
		}
		if p.X < a.X {
			return rayres{true, false}
		}
	}
	if a.Y < b.Y {
		if (p.Y-a.Y)/(p.X-a.X) >= (b.Y-a.Y)/(b.X-a.X) {
			return rayres{true, false}
		}
	} else {
		if (p.Y-b.Y)/(p.X-b.X) >= (a.Y-b.Y)/(a.X-b.X) {
			return rayres{true, false}
		}
	}
	return rayres{false, false}
}

// func rectForPoints(points []Point) Rect {
// 	var bbox Rect
// 	for i, p := range points {
// 		if i == 0 {
// 			bbox.Min = p
// 			bbox.Max = p
// 		} else {
// 			if p.X < bbox.Min.X {
// 				bbox.Min.X = p.X
// 			} else if p.X > bbox.Max.X {
// 				bbox.Max.X = p.X
// 			}
// 			if p.Y < bbox.Min.Y {
// 				bbox.Min.Y = p.Y
// 			} else if p.Y > bbox.Max.Y {
// 				bbox.Max.Y = p.Y
// 			}
// 		}
// 	}
// 	return bbox
// }

// func algoDoesIntersect(
// 	points []Point, pointsAreALineString bool, polygon Polygon,
// ) bool {
// 	// 	switch len(points) {
// 	// 	case 0:
// 	// 		return false
// 	// 	case 1:
// 	// 		switch len(polygon.Exterior) {
// 	// 		case 0:
// 	// 			return false
// 	// 		case 1:
// 	// 			return points[0].X == polygon.Exterior[0].X &&
// 	// 				points[0].Y == points[0].Y
// 	// 		default:
// 	// 			return algoPointInPolygon(points[0], polygon)
// 	// 		}
// 	// 	default:
// 	// 		switch len(polygon.Exterior) {
// 	// 		case 0:
// 	// 			return false
// 	// 		case 1:
// 	// 			return algoPointInPolygon(
// 	// 				polygon.Exterior[0], Polygon{points, polygon.Holes},
// 	// 			)
// 	// 		}
// 	// 	}
// 	// 	if !rectForPoints(points).IntersectsRect(rectForPoints(polygon.Exterior)) {
// 	// 		return false
// 	// 	}
// 	// 	for i := 0; i < len(points); i++ {
// 	// 		for j := 0; j < len(polygon.Exterior); j++ {
// 	// 			if algoSegmentsIntersect(
// 	// 				points[i], points[(i+1)%len(points)],
// 	// 				polygon.Exterior[j],
// 	// 				polygon.Exterior[(j+1)%len(polygon.Exterior)],
// 	// 			) {
// 	// 				return true
// 	// 			}
// 	// 		}
// 	// 	}
// 	// 	for _, hole := range polygon.Holes {
// 	// 		if ringInPolygon(Ring(points), Polygon{hole, nil}) {
// 	// 			return false
// 	// 		}
// 	// 	}
// 	// 	if ringInPolygon(Ring(points), Polygon{polygon.Exterior, nil}) {
// 	// 		return true
// 	// 	}
// 	// 	if !pointsAreALineString {
// 	// 		if ringInPolygon(polygon.Exterior, Polygon{points, nil}) {
// 	// 			return true
// 	// 		}
// 	// 	}
// 	return false
// }
