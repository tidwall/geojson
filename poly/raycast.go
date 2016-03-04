package poly

// This implementation of the raycast algorithm test if a point is
// to the left of a line, or on the segment line. Otherwise it is
// assumed that the point is outside of the segment line.

type rayres int

func (r rayres) String() string {
	switch r {
	default:
		return "unknown"
	case out:
		return "out"
	case left:
		return "left"
	case on:
		return "on"
	}
}

const (
	out  = rayres(0) // outside of the segment.
	left = rayres(1) // to the left of the segment
	on   = rayres(2) // on segment or vertex, special condition
)

func raycast(p, a, b Point) rayres {
	if a.Y == b.Y {
		// A and B share the same Y plane.
		if a.X == b.X {
			// AB is just a point.
			if p.X == a.X && p.Y == a.Y {
				return on
			}
			return out
		}
		// AB is a horizontal line.
		if p.Y != a.Y {
			// P is not on same Y plane as A and B.
			return out
		}
		// P is on same Y plane as A and B
		if a.X < b.X {
			if p.X >= a.X && p.X <= b.X {
				return on
			}
			if p.X < a.X {
				return left
			}
		} else {
			if p.X >= b.X && p.X <= a.X {
				return on
			}
			if p.X < b.X {
				return left
			}
		}
		return out
	}

	if a.X == b.X {
		// AB is a vertical line.
		if a.Y > b.Y {
			// A is above B
			if p.Y > a.Y || p.Y < b.Y {
				return out
			}
		} else {
			// B is above A
			if p.Y > b.Y || p.Y < a.Y {
				return out
			}
		}
		if p.X == a.X {
			return on
		}
		if p.X < a.X {
			return left
		}
		return out
	}

	// AB is an angled line
	if a.Y > b.Y {
		// swap A and B so that A is below B.
		a.X, a.Y, b.X, b.Y = b.X, b.Y, a.X, a.Y
	}
	if p.Y < a.Y || p.Y > b.Y {
		return out
	}
	if a.X < b.X {
		if p.X < a.X {
			return left
		}
		if p.X > b.X {
			return out
		}
	} else {
		if p.X < b.X {
			return left
		}
		if p.X > a.X {
			return out
		}
	}
	if (p.X == a.X && p.Y == a.Y) || (p.X == b.X && p.Y == b.Y) {
		// P is on a vertex.
		return on
	}
	v1 := (p.Y - a.Y) / (p.X - a.X)
	v2 := (b.Y - a.Y) / (b.X - a.X)
	if v1-v2 == 0 {
		// P is on a segment
		return on
	}
	if v1 >= v2 {
		return left
	}
	return out
}
