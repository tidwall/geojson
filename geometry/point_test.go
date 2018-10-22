// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geometry

import "testing"

func TestPointEmpty(t *testing.T) {
	expect(t, !P(0, 0).Empty())
}

func TestPointGeometryDefaults(t *testing.T) {
	g := Geometry(Point{})
	expect(t, !g.Empty())
	expect(t, g.Rect() == R(0, 0, 0, 0))
	expect(t, !g.ContainsLine(nil))
	expect(t, !g.ContainsLine(&Line{}))
	expect(t, g.ContainsPoint(Point{}))
	expect(t, !g.ContainsPoly(nil))
	expect(t, !g.ContainsPoly(&Poly{}))
	expect(t, g.ContainsRect(Rect{}))
	expect(t, !g.IntersectsLine(nil))
	expect(t, !g.IntersectsLine(&Line{}))
	expect(t, g.IntersectsPoint(Point{}))
	expect(t, !g.IntersectsPoly(nil))
	expect(t, !g.IntersectsPoly(&Poly{}))
	expect(t, g.IntersectsRect(Rect{}))
}

func TestPointRect(t *testing.T) {
	expect(t, P(5, 5).Rect() == R(5, 5, 5, 5))
}

func TestPointMove(t *testing.T) {
	expect(t, P(5, 6).Move(10, 10) == P(15, 16))
}

func TestPointContainsPoint(t *testing.T) {
	expect(t, P(5, 5).ContainsPoint(P(5, 5)))
	expect(t, !P(5, 5).ContainsPoint(P(6, 5)))
}

func TestPointIntersectsPoint(t *testing.T) {
	expect(t, P(5, 5).IntersectsPoint(P(5, 5)))
	expect(t, !P(5, 5).IntersectsPoint(P(6, 5)))
}

func TestPointContainsRect(t *testing.T) {
	expect(t, P(5, 5).ContainsRect(R(5, 5, 5, 5)))
	expect(t, !P(5, 5).ContainsRect(R(0, 0, 10, 10)))
}

func TestPointIntersectsRect(t *testing.T) {
	expect(t, P(5, 5).IntersectsRect(R(5, 5, 5, 5)))
	expect(t, P(5, 5).IntersectsRect(R(0, 0, 10, 10)))
	expect(t, P(0, 0).IntersectsRect(R(0, 0, 10, 10)))
	expect(t, !P(-1, 0).IntersectsRect(R(0, 0, 10, 10)))
}

func TestPointContainsLine(t *testing.T) {
	expect(t, !P(5, 5).ContainsLine(L()))
	expect(t, !P(5, 5).ContainsLine(L(P(5, 5))))
	expect(t, P(5, 5).ContainsLine(L(P(5, 5), P(5, 5))))
	expect(t, !P(5, 5).ContainsLine(L(P(5, 5), P(10, 10))))
}

func TestPointIntersectsLine(t *testing.T) {
	expect(t, !P(5, 5).IntersectsLine(L()))
	expect(t, !P(5, 5).IntersectsLine(L(P(5, 5))))
	expect(t, P(5, 5).IntersectsLine(L(P(5, 5), P(5, 5))))
	expect(t, P(5, 5).IntersectsLine(L(P(0, 0), P(10, 10))))
	expect(t, !P(6, 5).IntersectsLine(L(P(0, 0), P(10, 10))))
}

func TestPointContainsPoly(t *testing.T) {
	expect(t, !P(5, 5).ContainsPoly(NewPoly(nil, nil, DefaultIndexOptions)))
	expect(t, !P(5, 5).ContainsPoly(NewPoly([]Point{P(0, 0), P(10, 0)}, nil, DefaultIndexOptions)))
	expect(t, !P(5, 5).ContainsPoly(&Poly{Exterior: R(0, 0, 10, 10)}))
	expect(t, P(5, 5).ContainsPoly(&Poly{Exterior: R(5, 5, 5, 5)}))
}

func TestPointIntersectsPoly(t *testing.T) {
	octa := NewPoly(octagon, nil, DefaultIndexOptions)
	concave1 := NewPoly(concave1, nil, DefaultIndexOptions)
	expect(t, !P(5, 5).IntersectsPoly(NewPoly(nil, nil, DefaultIndexOptions)))
	expect(t, !P(5, 5).IntersectsPoly(NewPoly([]Point{P(0, 0), P(10, 0)}, nil, DefaultIndexOptions)))
	expect(t, P(5, 5).IntersectsPoly(octa))
	expect(t, P(0, 5).IntersectsPoly(octa))
	expect(t, !P(1, 1).IntersectsPoly(octa))
	expect(t, !P(4, 4).IntersectsPoly(concave1))
	expect(t, P(5, 5).IntersectsPoly(concave1))
	expect(t, P(6, 6).IntersectsPoly(concave1))
}
