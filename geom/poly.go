package geom

// Poly ...
type Poly interface {
	Exterior() Ring
	Holes() []Ring
	ContainsPoint(point Point) bool
	ContainsRing(ring Ring) bool
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
	if !poly.Exterior().ContainsRing(ring, true) {
		return false
	}
	contains := true
	for _, hole := range poly.holes {
		if hole.IntersectsRing(ring, false) {
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
	// 2) all poly holes must be fully contained in a
	return true
	// contains := true
	// otherHoles := other.Holes()
	// polyHoles := poly.Holes()
	// for _, polyHole := range polyHoles {
	// 	for _, hole := range otherHoles {
	// 		if poly.Contains(hole) {
	// 			contains = false
	// 			break
	// 		}
	// 	}
	// }
	// return contains
}
