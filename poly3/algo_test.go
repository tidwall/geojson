package poly

import "testing"

func TestConvexAlgo(t *testing.T) {
	expect(t, !ringConvex(nil))
	expect(t, ringConvex(rectangle))
	expect(t, ringConvex(pentagon))
	expect(t, ringConvex(triangle))
	expect(t, ringConvex(trapezoid))
	expect(t, ringConvex(octagon))
	expect(t, !ringConvex(concave1))
	expect(t, !ringConvex(concave2))
	expect(t, !ringConvex(concave3))
	expect(t, !ringConvex(concave4))
	expect(t, !ringConvex(bowtie))
}
