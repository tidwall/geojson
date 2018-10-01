package geom

import "testing"

func TestPoint(t *testing.T) {
	expect(t, P(1, 2).ContainsPoint(P(1, 2)))
	expect(t, !P(1, 2).ContainsPoint(P(1, 3)))
	expect(t, P(1, 2).IntersectsPoint(P(1, 2)))
	expect(t, !P(1, 2).IntersectsPoint(P(1, 3)))
	expect(t, P(1, 2).ContainsRect(R(1, 2, 1, 2)))
	expect(t, !P(1, 2).ContainsRect(R(1, 2, 1, 3)))
	expect(t, P(0, 0).IntersectsRect(R(0, 0, 10, 10)))
	expect(t, P(5, 5).IntersectsRect(R(0, 0, 10, 10)))
	expect(t, P(10, 10).IntersectsRect(R(0, 0, 10, 10)))
	expect(t, !P(11, 11).IntersectsRect(R(0, 0, 10, 10)))
	expect(t, !P(5, 5).ContainsLine(NewLine([]Point{})))
	expect(t, !P(0, 0).ContainsLine(NewLine([]Point{})))
	expect(t, !P(5, 5).ContainsLine(NewLine([]Point{P(5, 5)})))
	expect(t, P(5, 5).ContainsLine(NewLine([]Point{P(5, 5), P(5, 5)})))
	expect(t, !P(5, 5).IntersectsLine(NewLine([]Point{P(5, 5)})))
	expect(t, !P(5, 5).IntersectsLine(NewLine([]Point{P(0, 0)})))
	expect(t, P(5, 5).IntersectsLine(NewLine([]Point{P(0, 0), P(5, 5)})))
	expect(t, P(5, 5).IntersectsLine(NewLine([]Point{P(0, 0), P(10, 10)})))
	expect(t, !P(6, 5).IntersectsLine(NewLine([]Point{P(0, 0), P(10, 10)})))
	expect(t, P(5, 5).ContainsPoly(
		NewPoly([]Point{P(5, 5), P(5, 5), P(5, 5)}, nil),
	))
	expect(t, !P(6, 5).ContainsPoly(
		NewPoly([]Point{P(5, 5), P(5, 5), P(5, 5)}, nil),
	))
	expect(t, P(5, 5).IntersectsPoly(NewPoly(octagon, nil)))
	expect(t, !P(5, 5).IntersectsPoly(NewPoly(octagon, [][]Point{
		{P(4, 4), P(6, 4), P(6, 6), P(4, 6)},
	})))

	expect(t, !P(0, 0).Empty())
}
