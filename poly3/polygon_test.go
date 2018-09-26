package poly

import "testing"

func TestPolygonContainsRing(t *testing.T) {
	small := Ring{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}
	polyNoHole := Polygon{octagon, nil}
	polyWithHole := Polygon{octagon, []Ring{small}}
	expect(t, polyNoHole.ContainsRing(small.move(0, 0)))
	expect(t, polyNoHole.ContainsRing(small.move(1, 0)))

	//expect(t, !polyWithHole.ContainsRing(small.move(0, 0)))
	expect(t, !polyWithHole.ContainsRing(small.move(1, 0)))
	expect(t, polyWithHole.ContainsRing(small.move(2, 0)))

	// //poly :=

	// // expect(t, poly.ContainsRing(octagon))
	// // expect(t, !poly.ContainsRing(octagon.move(1, 0)))
	// small := Ring{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}

	// poly = Polygon{octagon, []Ring{small}}
	// // expect(t, !poly.ContainsRing(octagon))
	// // expect(t, !poly.ContainsRing(small))

}
