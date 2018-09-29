package geom

// Poly2 ...
type Poly2 struct {
	Exterior *Ring2
	Holes    []*Ring2
}

// NewPoly2 ...
func NewPoly2(exterior []Point, holes [][]Point) *Poly2 {
	poly := new(Poly2)
	poly.Exterior = NewRing2(exterior)
	if len(holes) > 0 {
		poly.Holes = make([]*Ring2, len(holes))
		for i := range holes {
			poly.Holes[i] = NewRing2(holes[i])
		}
	}
	return poly
}

// ContainsPoint ...
func (poly *Poly2) ContainsPoint(point Point) bool {
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
func (poly *Poly2) IntersectsPoint(point Point) bool {
	return poly.ContainsPoint(point)
}

// ContainsPoly ...
func (poly *Poly2) ContainsPoly(other *Poly2) bool {
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
func (poly *Poly2) IntersectsPoly(other *Poly2) bool {
	return other.Exterior.IntersectsPoly(poly, true)
}
