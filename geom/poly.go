package geom

// Poly ...
type Poly interface {
	Exterior() Ring
	Holes() []Ring

	ContainsPoint(point Point) bool

	ContainsRing(ring Ring) bool

	IntersectsRing(ring Ring) bool

	ContainsPoly(poly Poly) bool
	IntersectsPoly(poly Poly) bool
}

// NewPoly ...
func NewPoly(exterior []Point, holes [][]Point, index int) Poly {
	poly := new(sharedPoly)
	poly.exterior = NewRing(exterior, index)
	poly.holes = make([]Ring, len(holes))
	for i, hole := range holes {
		poly.holes[i] = NewRing(hole, index)
	}
	return poly
}

type sharedPoly struct {
	exterior Ring
	holes    []Ring
}

func (poly *sharedPoly) Exterior() Ring {
	return poly.exterior
}

func (poly *sharedPoly) Holes() []Ring {
	return poly.holes
}

func (poly *sharedPoly) ContainsPoint(point Point) bool {
	if !poly.Exterior().ContainsPoint(point, true) {
		return false
	}
	contains := true
	for _, hole := range poly.holes {
		if hole.ContainsPoint(point, false) {
			contains = false
			break
		}
	}
	return contains
}

func (poly *sharedPoly) ContainsRing(ring Ring) bool {
	// 1) other exterior must be fully contained inside of the ring.
	if !poly.Exterior().ContainsRing(ring, true) {
		return false
	}
	// 2) ring cannot intersect poly holes
	contains := true
	for _, polyHole := range poly.Holes() {
		if polyHole.IntersectsRing(ring, true) {
			contains = false
			break
		}
	}
	return contains
}

func (poly *sharedPoly) ContainsPoly(other Poly) bool {
	// 1) other exterior must be fully contained inside of the poly exterior.
	if !poly.Exterior().ContainsRing(other.Exterior(), true) {
		return false
	}
	// 2) ring cannot intersect poly holes
	contains := true
	for _, polyHole := range poly.Holes() {
		if polyHole.IntersectsRing(other.Exterior(), false) {
			contains = false
			// 3) unless the poly hole is contain inside of a other hole
			for _, otherHole := range other.Holes() {
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

func (poly *sharedPoly) IntersectsRing(ring Ring) bool {
	return ring.IntersectsPoly(poly, true)
}

func (poly *sharedPoly) IntersectsPoly(other Poly) bool {
	return poly.IntersectsRing(other.Exterior())
}
