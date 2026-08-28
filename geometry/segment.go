// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geometry

func eqZero(x float64) bool {
	return FloatZero(x)
}

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

func (seg Segment) CollinearPoint(point Point) bool {
	cmpx, cmpy := point.X-seg.A.X, point.Y-seg.A.Y
	rx, ry := seg.B.X-seg.A.X, seg.B.Y-seg.A.Y
	cmpxr := cmpx*ry - cmpy*rx
	return eqZero(cmpxr)
}

func (seg Segment) ContainsPoint(point Point) bool {
	return seg.Raycast(point).On
}

// func (seg Segment) Angle() float64 {
// 	return math.Atan2(seg.B.Y-seg.A.Y, seg.B.X-seg.A.X)
// }

// IntersectsSegment detects if segment intersects with other segement
func (seg Segment) IntersectsSegment(other Segment) bool {
	a, b, c, d := seg.A, seg.B, other.A, other.B
	// do the bounding boxes intersect?
	if FloatGreater(a.Y, b.Y) {
		if FloatGreater(c.Y, d.Y) {
			if FloatGreater(b.Y, c.Y) || FloatLess(a.Y, d.Y) {
				return false
			}
		} else {
			if FloatGreater(b.Y, d.Y) || FloatLess(a.Y, c.Y) {
				return false
			}
		}
	} else {
		if FloatGreater(c.Y, d.Y) {
			if FloatGreater(a.Y, c.Y) || FloatLess(b.Y, d.Y) {
				return false
			}
		} else {
			if FloatGreater(a.Y, d.Y) || FloatLess(b.Y, c.Y) {
				return false
			}
		}
	}
	if FloatGreater(a.X, b.X) {
		if FloatGreater(c.X, d.X) {
			if FloatGreater(b.X, c.X) || FloatLess(a.X, d.X) {
				return false
			}
		} else {
			if FloatGreater(b.X, d.X) || FloatLess(a.X, c.X) {
				return false
			}
		}
	} else {
		if FloatGreater(c.X, d.X) {
			if FloatGreater(a.X, c.X) || FloatLess(b.X, d.X) {
				return false
			}
		} else {
			if FloatGreater(a.X, d.X) || FloatLess(b.X, c.X) {
				return false
			}
		}
	}
	if PointEqual(seg.A, other.A) || PointEqual(seg.A, other.B) ||
		PointEqual(seg.B, other.A) || PointEqual(seg.B, other.B) {
		return true
	}

	// the following code is from http://ideone.com/PnPJgb
	cmpx, cmpy := c.X-a.X, c.Y-a.Y
	rx, ry := b.X-a.X, b.Y-a.Y
	cmpxr := cmpx*ry - cmpy*rx
	if eqZero(cmpxr) {
		// Lines are collinear, and so intersect if they have any overlap
		if !(((FloatLessOrEqual(c.X-a.X, 0)) != (FloatLessOrEqual(c.X-b.X, 0))) ||
			((FloatLessOrEqual(c.Y-a.Y, 0)) != (FloatLessOrEqual(c.Y-b.Y, 0)))) {
			return seg.Raycast(other.A).On || seg.Raycast(other.B).On
			//return false
		}
		return true
	}
	sx, sy := d.X-c.X, d.Y-c.Y
	cmpxs := cmpx*sy - cmpy*sx
	rxs := rx*sy - ry*sx
	if eqZero(rxs) {
		return false // segments are parallel.
	}
	rxsr := 1 / rxs
	t := cmpxs * rxsr
	u := cmpxr * rxsr
	if !((FloatGreaterOrEqual(t, 0)) && (FloatLessOrEqual(t, 1)) && (FloatGreaterOrEqual(u, 0)) && (FloatLessOrEqual(u, 1))) {
		return false
	}
	return true
}

// ContainsSegment returns true if segment contains other segment
func (seg Segment) ContainsSegment(other Segment) bool {
	return seg.Raycast(other.A).On && seg.Raycast(other.B).On
}
