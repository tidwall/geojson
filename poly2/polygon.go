package poly

// Polygon ...
type Polygon struct {
	Exterior Ring
	Holes    []Ring
}

// // ContainsPoint ...
// func (polygon Polygon) ContainsPoint(point Point) bool {
// 	if !polygon.Exterior.ContainsPoint(point) {
// 		return false
// 	}
// 	for _, hole := range polygon.Holes {
// 		if algoPointInPolygon(point, hole, false) {
// 			return false
// 		}
// 	}
// 	return true
// }

// // ContainsRing ...
// func (polygon Polygon) ContainsRing(ring Ring) bool {
// 	if !polygon.Exterior.ContainsRing(ring) {
// 		return false
// 	}
// 	for _, hole := range polygon.Holes {
// 		if algoPointInPolygon(point, hole, false) {
// 			return false
// 		}
// 	}
// 	return true
// }
