package geom

// Poly ...
type Poly struct {
	Exterior *Ring
	Holes    []*Ring
}

// NewPoly ...
func NewPoly(exterior []Point, holes [][]Point) *Poly {
	poly := new(Poly)
	poly.Exterior = NewRing(exterior)
	if len(holes) > 0 {
		poly.Holes = make([]*Ring, len(holes))
		for i := range holes {
			poly.Holes[i] = NewRing(holes[i])
		}
	}
	return poly
}

// Empty ...
func (poly *Poly) Empty() bool {
	return poly.Exterior.Empty()
}

// Rect ...
func (poly *Poly) Rect() Rect {
	return poly.Exterior.Rect()
}

// ContainsPoint ...
func (poly *Poly) ContainsPoint(point Point) bool {
	if !poly.Exterior.ContainsPoint(point, true) {
		return false
	}
	contains := true
	for _, hole := range poly.Holes {
		if hole.ContainsPoint(point, false) {
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
	panic("not ready")
}

// IntersectsRect ...
func (poly *Poly) IntersectsRect(rect Rect) bool {
	panic("not ready")
}

// ContainsLine ...
func (poly *Poly) ContainsLine(line *Line) bool {
	panic("not ready")
}

// IntersectsLine ...
func (poly *Poly) IntersectsLine(line *Line) bool {
	panic("not ready")
}

// ContainsPoly ...
func (poly *Poly) ContainsPoly(other *Poly) bool {
	// 1) other exterior must be fully contained inside of the poly exterior.
	if !poly.Exterior.ContainsRing(other.Exterior, true) {
		return false
	}
	// 2) ring cannot intersect poly holes
	contains := true
	for _, polyHole := range poly.Holes {
		if polyHole.IntersectsRing(other.Exterior, false) {
			contains = false
			// 3) unless the poly hole is contain inside of a other hole
			for _, otherHole := range other.Holes {
				if otherHole.ContainsRing(polyHole, true) {
					contains = true
					break
				}
			}
			if !contains {
				break
			}
		}
	}
	return contains
}

// IntersectsPoly ...
func (poly *Poly) IntersectsPoly(other *Poly) bool {
	return other.Exterior.IntersectsPoly(poly, true)
}
