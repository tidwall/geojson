package geometry

import "math"

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
	if FloatEqual(a.Y, b.Y) {
		if FloatEqual(a.X, b.X) {
			if PointEqual(p, a) {
				return RaycastResult{false, true}
			}
			return RaycastResult{false, false}
		}
		if FloatEqual(p.Y, b.Y) {
			// horizontal segment
			// check if the point in on the line
			if FloatLess(a.X, b.X) {
				if FloatGreaterOrEqual(p.X, a.X) && FloatLessOrEqual(p.X, b.X) {
					return RaycastResult{false, true}
				}
			} else {
				if FloatGreaterOrEqual(p.X, b.X) && FloatLessOrEqual(p.X, a.X) {
					return RaycastResult{false, true}
				}
			}
		}
	}
	if FloatEqual(a.X, b.X) && FloatEqual(p.X, b.X) {
		// vertical segment
		// check if the point in on the line
		if FloatLess(a.Y, b.Y) {
			if FloatGreaterOrEqual(p.Y, a.Y) && FloatLessOrEqual(p.Y, b.Y) {
				return RaycastResult{false, true}
			}
		} else {
			if FloatGreaterOrEqual(p.Y, b.Y) && FloatLessOrEqual(p.Y, a.Y) {
				return RaycastResult{false, true}
			}
		}
	}
	// Use epsilon comparison for slope equality instead of direct division comparison
	if FloatNonZero(b.X-a.X) && FloatNonZero(b.Y-a.Y) {
		slope1 := (p.X - a.X) / (b.X - a.X)
		slope2 := (p.Y - a.Y) / (b.Y - a.Y)
		if FloatEqual(slope1, slope2) {
			return RaycastResult{false, true}
		}
	}

	// do the actual raycast here.
	for FloatEqual(p.Y, a.Y) || FloatEqual(p.Y, b.Y) {
		p.Y = math.Nextafter(p.Y, math.Inf(1))
	}
	if FloatLess(a.Y, b.Y) {
		if FloatLess(p.Y, a.Y) || FloatGreater(p.Y, b.Y) {
			return RaycastResult{false, false}
		}
	} else {
		if FloatLess(p.Y, b.Y) || FloatGreater(p.Y, a.Y) {
			return RaycastResult{false, false}
		}
	}
	if FloatGreater(a.X, b.X) {
		if FloatGreaterOrEqual(p.X, a.X) {
			return RaycastResult{false, false}
		}
		if FloatLessOrEqual(p.X, b.X) {
			return RaycastResult{true, false}
		}
	} else {
		if FloatGreaterOrEqual(p.X, b.X) {
			return RaycastResult{false, false}
		}
		if FloatLessOrEqual(p.X, a.X) {
			return RaycastResult{true, false}
		}
	}
	if FloatLess(a.Y, b.Y) {
		// Use epsilon-safe division comparison
		if FloatNonZero(p.X-a.X) && FloatNonZero(b.X-a.X) {
			slope1 := (p.Y - a.Y) / (p.X - a.X)
			slope2 := (b.Y - a.Y) / (b.X - a.X)
			if FloatGreaterOrEqual(slope1, slope2) {
				return RaycastResult{true, false}
			}
		}
	} else {
		// Use epsilon-safe division comparison
		if FloatNonZero(p.X-b.X) && FloatNonZero(a.X-b.X) {
			slope1 := (p.Y - b.Y) / (p.X - b.X)
			slope2 := (a.Y - b.Y) / (a.X - b.X)
			if FloatGreaterOrEqual(slope1, slope2) {
				return RaycastResult{true, false}
			}
		}
	}
	return RaycastResult{false, false}
}
