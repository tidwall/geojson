package poly

// anySegmentsIntersect return true when any segments intersect with each other.
// closeA and closeB will ensure that lines are closed thus creating rings.
func anySegmentsIntersect(
	lineA, lineB []Point,
	closeA, closeB bool,
	allowOn, exterior bool,
) bool {
	var a, b, c, d Point
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
			if segmentsIntersect(a, b, c, d) {
				if allowOn || (!raycast(a, c, d).on &&
					!raycast(b, c, d).on &&
					!raycast(c, a, b).on &&
					!raycast(d, a, b).on) {
					println(123)
					return true
				}
			}
		}
	}
	return false
}

// segmentsIntersect returns true when two line segments intersect
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

// ringConvex returns true if the ring is convex, a simple polygon where all
// points join using the same directional angle.
// Such a square, triangle, or octagon.
func ringConvex(ring Ring) bool {
	if len(ring) < 3 {
		return false
	}
	var dir int
	var a, b, c Point
	for i := 0; i < len(ring); i++ {
		a = ring[i]
		if i == len(ring)-1 {
			b = ring[0]
			c = ring[1]
		} else if i == len(ring)-2 {
			b = ring[i+1]
			c = ring[0]
		} else {
			b = ring[i+1]
			c = ring[i+2]
		}
		dx1 := b.X - a.X
		dy1 := b.Y - a.Y
		dx2 := c.X - b.X
		dy2 := c.Y - b.Y
		zCrossProduct := dx1*dy2 - dy1*dx2
		if dir == 0 {
			if zCrossProduct < 0 {
				dir = -1
			} else if zCrossProduct > 0 {
				dir = 1
			}
		} else if zCrossProduct < 0 {
			if dir == 1 {
				return false
			}
		} else if zCrossProduct > 0 {
			if dir == -1 {
				return false
			}
		}
	}
	return true
}

// ringConvexRect tests if the ring is convex and calculates the outer
// rectangle in one operation.
func ringConvexRect(ring Ring) (convex bool, rect Rect) {
	if len(ring) < 3 {
		return false, pointsRect(ring)
	}
	var concave bool
	var dir int
	var a, b, c Point
	for i := 0; i < len(ring); i++ {
		if i == 0 {
			rect = Rect{ring[i], ring[i]}
		} else {
			if ring[i].X < rect.Min.X {
				rect.Min.X = ring[i].X
			} else if ring[i].X > rect.Max.X {
				rect.Max.X = ring[i].X
			}
			if ring[i].Y < rect.Min.Y {
				rect.Min.Y = ring[i].Y
			} else if ring[i].Y > rect.Max.Y {
				rect.Max.Y = ring[i].Y
			}
		}
		if concave {
			continue
		}
		a = ring[i]
		if i == len(ring)-1 {
			b = ring[0]
			c = ring[1]
		} else if i == len(ring)-2 {
			b = ring[i+1]
			c = ring[0]
		} else {
			b = ring[i+1]
			c = ring[i+2]
		}
		dx1 := b.X - a.X
		dy1 := b.Y - a.Y
		dx2 := c.X - b.X
		dy2 := c.Y - b.Y
		zCrossProduct := dx1*dy2 - dy1*dx2
		if dir == 0 {
			if zCrossProduct < 0 {
				dir = -1
			} else if zCrossProduct > 0 {
				dir = 1
			}
		} else if zCrossProduct < 0 {
			if dir == 1 {
				concave = true
			}
		} else if zCrossProduct > 0 {
			if dir == -1 {
				concave = true
			}
		}
	}
	return !concave, rect
}

// pointsRect returns the outer rectangle for points. These point could be a
// line, ring, or just a bunch of points.
func pointsRect(points []Point) Rect {
	var rect Rect
	if len(points) > 0 {
		rect = Rect{points[0], points[0]}
		for i := 1; i < len(points); i++ {
			if points[i].X < rect.Min.X {
				rect.Min.X = points[i].X
			} else if points[i].X > rect.Max.X {
				rect.Max.X = points[i].X
			}
			if points[i].Y < rect.Min.Y {
				rect.Min.Y = points[i].Y
			} else if points[i].Y > rect.Max.Y {
				rect.Max.Y = points[i].Y
			}
		}
	}
	return rect
}
