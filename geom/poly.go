package geom

// Poly ...
type Poly struct {
	Exterior RingX
	Holes    []RingX
}

// NewPoly ...
func NewPoly(exterior []Point, holes [][]Point) *Poly {
	poly := new(Poly)
	poly.Exterior = newRingX(exterior)
	if len(holes) > 0 {
		poly.Holes = make([]RingX, len(holes))
		for i := range holes {
			poly.Holes[i] = newRingX(holes[i])
		}
	}
	return poly
}

// Clockwise ...
func (poly *Poly) Clockwise() bool {
	return poly.Exterior.Clockwise()
}

// Empty ...
func (poly *Poly) Empty() bool {
	return poly.Exterior.Empty()
}

// Rect ...
func (poly *Poly) Rect() Rect {
	return poly.Exterior.Rect()
}

// Move the polygon by delta. Returns a new polygon
func (poly *Poly) Move(deltaX, deltaY float64) *Poly {
	npoly := new(Poly)
	if series, ok := poly.Exterior.(*baseSeries); ok {
		npoly.Exterior = RingX(series.Move(deltaX, deltaY))
	} else {
		nseries := makeSeries(seriesCopyPoints(poly.Exterior), false, true)
		npoly.Exterior = RingX(nseries.Move(deltaX, deltaY))
	}
	if len(poly.Holes) > 0 {
		npoly.Holes = make([]RingX, len(poly.Holes))
		for i, hole := range poly.Holes {
			if series, ok := hole.(*baseSeries); ok {
				npoly.Holes[i] = RingX(series.Move(deltaX, deltaY))
			} else {
				nseries := makeSeries(seriesCopyPoints(hole), false, true)
				npoly.Holes[i] = RingX(nseries.Move(deltaX, deltaY))
			}
		}
	}
	return npoly
}

// ContainsPoint ...
func (poly *Poly) ContainsPoint(point Point) bool {
	if !ringxContainsPoint(poly.Exterior, point, true).hit {
		return false
	}
	contains := true
	for _, hole := range poly.Holes {
		if ringxContainsPoint(hole, point, false).hit {
			contains = false
			break
		}
	}
	return contains
}

// IntersectsPoint ...
func (poly *Poly) IntersectsPoint(point Point) bool {
	return poly.ContainsPoint(point)
}

// ContainsRect ...
func (poly *Poly) ContainsRect(rect Rect) bool {
	// convert rect into a polygon
	return poly.ContainsPoly(&Poly{Exterior: rect})
}

// IntersectsRect ...
func (poly *Poly) IntersectsRect(rect Rect) bool {
	// convert rect into a polygon
	return poly.IntersectsPoly(&Poly{Exterior: rect})
}

// ContainsLine ...
func (poly *Poly) ContainsLine(line *Line) bool {
	if !ringxContainsLine(poly.Exterior, line, true) {
		return false
	}
	for _, polyHole := range poly.Holes {
		if ringxIntersectsLine(polyHole, line, false) {
			return false
		}
	}
	return true
}

// IntersectsLine ...
func (poly *Poly) IntersectsLine(line *Line) bool {
	return ringxIntersectsLine(poly.Exterior, line, true)
}

// ContainsPoly ...
func (poly *Poly) ContainsPoly(other *Poly) bool {
	println(0)
	// 1) other exterior must be fully contained inside of the poly exterior.
	if !ringxContainsRing(poly.Exterior, other.Exterior, true) {
		return false
	}
	// 2) ring cannot intersect poly holes
	println(1)
	contains := true
	for _, polyHole := range poly.Holes {
		println(2)

		// println(ringxString(polyHole))
		// println(ringxString(other.Exterior))

		println("--", ringxIntersectsRing(polyHole, other.Exterior, false))
		println("--", ringxIntersectsRing(polyHole, other.Exterior, true))
		if ringxIntersectsRing(polyHole, other.Exterior, false) {
			contains = false
			println(3)
			// 3) unless the poly hole is contain inside of a other hole
			for _, otherHole := range other.Holes {
				if ringxContainsRing(otherHole, polyHole, true) {
					contains = true
					println(4)
					break
				}
			}
			if !contains {
				println(5)
				break
			}
		}
	}
	println(6)
	return contains
}

// IntersectsPoly ...
func (poly *Poly) IntersectsPoly(other *Poly) bool {
	return ringxIntersectsPoly(other.Exterior, poly, true)
}
