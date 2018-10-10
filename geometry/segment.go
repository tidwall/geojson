// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geometry

import (
	"math"
)

// Segment is a two point line
type Segment struct {
	A, B Point
}

// Move a segment by delta
func (seg Segment) Move(deltaX, deltaY float64) Segment {
	return Segment{
		A: Point{X: seg.A.X + deltaX, Y: seg.A.Y + deltaY},
		B: Point{X: seg.B.X + deltaX, Y: seg.B.Y + deltaY},
	}
}

// Rect is the outer boundaries of the segment.
func (seg Segment) Rect() Rect {
	var rect Rect
	rect.Min = seg.A
	rect.Max = seg.B
	if rect.Min.X > rect.Max.X {
		rect.Min.X, rect.Max.X = rect.Max.X, rect.Min.X
	}
	if rect.Min.Y > rect.Max.Y {
		rect.Min.Y, rect.Max.Y = rect.Max.Y, rect.Min.Y
	}
	return rect
}

// CollinearPoint ...
func (seg Segment) CollinearPoint(point Point) bool {
	cmpx, cmpy := point.X-seg.A.X, point.Y-seg.A.Y
	rx, ry := seg.B.X-seg.A.X, seg.B.Y-seg.A.Y
	cmpxr := cmpx*ry - cmpy*rx
	return cmpxr == 0
}

// ContainsPoint ...
func (seg Segment) ContainsPoint(point Point) bool {
	return seg.Raycast(point).On
}

// // Angle ...
// func (seg Segment) Angle() float64 {
// 	return math.Atan2(seg.B.Y-seg.A.Y, seg.B.X-seg.A.X)
// }

// RaycastResult holds the results of the Raycast operation
type RaycastResult struct {
	In bool // point on the left
	On bool // point is directly on top of
}

// Raycast performs the raycast operation
func (seg Segment) Raycast(point Point) RaycastResult {

	p, a, b := point, seg.A, seg.B
	// make sure that the point is inside the segment bounds
	if a.Y < b.Y && (p.Y < a.Y || p.Y > b.Y) {
		return RaycastResult{false, false}
	} else if a.Y > b.Y && (p.Y < b.Y || p.Y > a.Y) {
		return RaycastResult{false, false}
	}

	// test if point is in on the segment
	if a.Y == b.Y {
		if a.X == b.X {
			if p == a {
				return RaycastResult{false, true}
			}
			return RaycastResult{false, false}
		}
		if p.Y == b.Y {
			// horizontal segment
			// check if the point in on the line
			if a.X < b.X {
				if p.X >= a.X && p.X <= b.X {
					return RaycastResult{false, true}
				}
			} else {
				if p.X >= b.X && p.X <= a.X {
					return RaycastResult{false, true}
				}
			}
		}
	}
	if a.X == b.X && p.X == b.X {
		// vertical segment
		// check if the point in on the line
		if a.Y < b.Y {
			if p.Y >= a.Y && p.Y <= b.Y {
				return RaycastResult{false, true}
			}
		} else {
			if p.Y >= b.Y && p.Y <= a.Y {
				return RaycastResult{false, true}
			}
		}
	}
	if (p.X-a.X)/(b.X-a.X) == (p.Y-a.Y)/(b.Y-a.Y) {
		return RaycastResult{false, true}
	}

	// do the actual raycast here.
	for p.Y == a.Y || p.Y == b.Y {
		p.Y = math.Nextafter(p.Y, math.Inf(1))
	}
	if a.Y < b.Y {
		if p.Y < a.Y || p.Y > b.Y {
			return RaycastResult{false, false}
		}
	} else {
		if p.Y < b.Y || p.Y > a.Y {
			return RaycastResult{false, false}
		}
	}
	if a.X > b.X {
		if p.X >= a.X {
			return RaycastResult{false, false}
		}
		if p.X <= b.X {
			return RaycastResult{true, false}
		}
	} else {
		if p.X >= b.X {
			return RaycastResult{false, false}
		}
		if p.X <= a.X {
			return RaycastResult{true, false}
		}
	}
	if a.Y < b.Y {
		if (p.Y-a.Y)/(p.X-a.X) >= (b.Y-a.Y)/(b.X-a.X) {
			return RaycastResult{true, false}
		}
	} else {
		if (p.Y-b.Y)/(p.X-b.X) >= (a.Y-b.Y)/(a.X-b.X) {
			return RaycastResult{true, false}
		}
	}
	return RaycastResult{false, false}
}

// IntersectsSegment detects if segment intersects with other segement
func (seg Segment) IntersectsSegment(other Segment) bool {
	a, b, c, d := seg.A, seg.B, other.A, other.B
	// do the bounding boxes intersect?
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
	if seg.A == other.A || seg.A == other.B ||
		seg.B == other.A || seg.B == other.B {
		return true
	}

	// the following code is from http://ideone.com/PnPJgb
	cmpx, cmpy := c.X-a.X, c.Y-a.Y
	rx, ry := b.X-a.X, b.Y-a.Y
	cmpxr := cmpx*ry - cmpy*rx
	if cmpxr == 0 {
		// Lines are collinear, and so intersect if they have any overlap
		if !(((c.X-a.X <= 0) != (c.X-b.X <= 0)) ||
			((c.Y-a.Y <= 0) != (c.Y-b.Y <= 0))) {
			return seg.Raycast(other.A).On || seg.Raycast(other.B).On
			//return false
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

// ContainsSegment returns true if segment contains other segment
func (seg Segment) ContainsSegment(other Segment) bool {
	return seg.Raycast(other.A).On && seg.Raycast(other.B).On
}
