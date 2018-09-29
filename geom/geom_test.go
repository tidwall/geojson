package geom

import "testing"

func TestRectVarious(t *testing.T) {
	expect(t, R(0, 0, 10, 10).ContainsRing(newRingSimple2(octagon)))
	expect(t, !R(5, 0, 15, 10).ContainsRing(newRingSimple2(octagon)))
	expect(t, R(5, 0, 15, 10).IntersectsRing(newRingSimple2(octagon)))
}
