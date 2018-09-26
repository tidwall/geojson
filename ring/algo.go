package ring

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

// pointsConvex returns true if the ring is convex, a simple polygon where all
// points join using the same directional angle.
// Such a square, triangle, or octagon.
func pointsConvex(points []Point) bool {
	if len(points) < 3 {
		return false
	}
	var dir int
	var a, b, c Point
	for i := 0; i < len(points); i++ {
		a = points[i]
		if i == len(points)-1 {
			b = points[0]
			c = points[1]
		} else if i == len(points)-2 {
			b = points[i+1]
			c = points[0]
		} else {
			b = points[i+1]
			c = points[i+2]
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

// pointsConvexRect tests if the ring is convex and calculates the outer
// rectangle in one operation.
func pointsConvexRect(points []Point) (convex bool, rect Rect) {
	if len(points) < 3 {
		return false, pointsRect(points)
	}
	var concave bool
	var dir int
	var a, b, c Point
	for i := 0; i < len(points); i++ {
		if i == 0 {
			rect = Rect{points[i], points[i]}
		} else {
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
		if concave {
			continue
		}
		a = points[i]
		if i == len(points)-1 {
			b = points[0]
			c = points[1]
		} else if i == len(points)-2 {
			b = points[i+1]
			c = points[0]
		} else {
			b = points[i+1]
			c = points[i+2]
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

func pointInRing(point Point, ring []Point, allowOnEdge bool) bool {
	in := false
	var a, b Point
	for i := 0; i < len(ring); i++ {
		a = ring[i]
		if i == len(ring)-1 {
			b = ring[0]
		} else {
			b = ring[i+1]
		}
		res := raycast(point, a, b)
		if res.on {
			return allowOnEdge
		}
		if res.in {
			in = !in
		}
	}
	return in
}
