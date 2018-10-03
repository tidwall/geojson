package geom

// Poly ...
type Poly struct {
	Exterior Ring
	Holes    []Ring
}

// NewPoly ...
func NewPoly(exterior []Point, holes [][]Point) *Poly {
	poly := new(Poly)
	poly.Exterior = newRing(exterior)
	if len(holes) > 0 {
		poly.Holes = make([]Ring, len(holes))
		for i := range holes {
			poly.Holes[i] = newRing(holes[i])
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
		npoly.Exterior = Ring(series.Move(deltaX, deltaY))
	} else {
		nseries := makeSeries(seriesCopyPoints(poly.Exterior), false, true)
		npoly.Exterior = Ring(nseries.Move(deltaX, deltaY))
	}
	if len(poly.Holes) > 0 {
		npoly.Holes = make([]Ring, len(poly.Holes))
		for i, hole := range poly.Holes {
			if series, ok := hole.(*baseSeries); ok {
				npoly.Holes[i] = Ring(series.Move(deltaX, deltaY))
			} else {
				nseries := makeSeries(seriesCopyPoints(hole), false, true)
				npoly.Holes[i] = Ring(nseries.Move(deltaX, deltaY))
			}
		}
	}
	return npoly
}

// ContainsPoint ...
func (poly *Poly) ContainsPoint(point Point) bool {
	if !ringContainsPoint(poly.Exterior, point, true) {
		return false
	}
	contains := true
	for _, hole := range poly.Holes {
		if ringContainsPoint(hole, point, false) {
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
	if !ringContainsLine(poly.Exterior, line, true) {
		return false
	}
	for _, polyHole := range poly.Holes {
		if ringIntersectsLine(polyHole, line, false) {
			return false
		}
	}
	return true
}

// IntersectsLine ...
func (poly *Poly) IntersectsLine(line *Line) bool {
	return ringIntersectsLine(poly.Exterior, line, true)
}

// ContainsPoly ...
func (poly *Poly) ContainsPoly(other *Poly) bool {
	// 1) other exterior must be fully contained inside of the poly exterior.
	if !ringContainsRing(poly.Exterior, other.Exterior, true) {
		return false
	}
	// 2) ring cannot intersect poly holes
	println(1)
	contains := true
	for _, polyHole := range poly.Holes {
		println(2)

		println(ringString(polyHole))
		println(ringString(other.Exterior))

		println("--", ringIntersectsRing(polyHole, other.Exterior, false))
		println("--", ringIntersectsRing(polyHole, other.Exterior, true))
		if ringIntersectsRing(polyHole, other.Exterior, false) {
			contains = false
			println(3)
			// 3) unless the poly hole is contain inside of a other hole
			for _, otherHole := range other.Holes {
				if ringContainsRing(otherHole, polyHole, true) {
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
	return ringIntersectsPoly(other.Exterior, poly, true)
}
