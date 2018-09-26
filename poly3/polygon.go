package poly

// Polygon ...
type Polygon struct {
	Exterior Ring
	Holes    []Ring
}

// ContainsRing ...
func (polygon Polygon) ContainsRing(ring Ring) bool {
	if !polygon.Exterior.ContainsRing(ring) {
		return false
	}
	for _, hole := range polygon.Holes {
		if hole.intersectsRing(ring, false) {
			return false
		}
	}
	return true
}
