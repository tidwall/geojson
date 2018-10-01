package geom

import "testing"

func TestSegment(t *testing.T) {
	expect(t, S(0, 0, 10, 10).ContainsSegment(S(0, 0, 10, 10)))
	expect(t, S(0, 0, 10, 10).ContainsSegment(S(2, 2, 10, 10)))
	expect(t, S(0, 0, 10, 10).ContainsSegment(S(2, 2, 8, 8)))
	expect(t, !S(0, 0, 10, 10).ContainsSegment(S(-1, -1, 8, 8)))
}
