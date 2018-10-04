package geom

import "testing"

func TestSegmentContainsSegment(t *testing.T) {
	expect(t, S(0, 0, 10, 10).ContainsSegment(S(0, 0, 10, 10)))
	expect(t, S(0, 0, 10, 10).ContainsSegment(S(2, 2, 10, 10)))
	expect(t, S(0, 0, 10, 10).ContainsSegment(S(2, 2, 8, 8)))
	expect(t, !S(0, 0, 10, 10).ContainsSegment(S(-1, -1, 8, 8)))
}

func TestSegmentIntersectsSegment(t *testing.T) {
	expect(t, S(1, 4, 2, 3).IntersectsSegment(S(1, 3, 3, 3)))
}
