// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geometry

import "testing"

func TestLineNewLine(t *testing.T) {
	line := NewLine(u1, DefaultIndexOptions)
	expect(t, !line.Empty())
}

func TestLineGeometryDefaults(t *testing.T) {
	g := Geometry(&Line{})
	expect(t, g.Empty())
	expect(t, g.Rect() == R(0, 0, 0, 0))
	expect(t, !g.ContainsLine(nil))
	expect(t, !g.ContainsLine(&Line{}))
	expect(t, !g.ContainsPoint(Point{}))
	expect(t, !g.ContainsPoly(nil))
	expect(t, !g.ContainsPoly(&Poly{}))
	expect(t, !g.ContainsRect(Rect{}))
	expect(t, !g.IntersectsLine(nil))
	expect(t, !g.IntersectsLine(&Line{}))
	expect(t, !g.IntersectsPoint(Point{}))
	expect(t, !g.IntersectsPoly(nil))
	expect(t, !g.IntersectsPoly(&Poly{}))
	expect(t, !g.IntersectsRect(Rect{}))
}

func TestLineMove(t *testing.T) {
	ln1 := L(P(0, 1), P(2, 3), P(4, 5))
	ln2 := ln1.Move(7, 8)
	expect(t, ln1.NumPoints() == ln2.NumPoints())
	for i := 0; i < ln2.NumPoints(); i++ {
		expect(t, ln2.PointAt(i) == ln1.PointAt(i).Move(7, 8))
	}

	var line *Line
	expect(t, line.Move(0, 0) == nil)
	expect(t, (&Line{}).Move(0, 0) != nil)
}

func TestLineContainsPoint(t *testing.T) {
	line := NewLine(u1, DefaultIndexOptions)
	expect(t, line.ContainsPoint(P(0, 0)))
	expect(t, line.ContainsPoint(P(10, 10)))
	expect(t, line.ContainsPoint(P(0, 5)))
	expect(t, !line.ContainsPoint(P(5, 5)))
	line = NewLine(v1, DefaultIndexOptions)
	expect(t, line.ContainsPoint(P(0, 10)))
	expect(t, !line.ContainsPoint(P(0, 0)))
	expect(t, line.ContainsPoint(P(5, 0)))
	expect(t, line.ContainsPoint(P(2.5, 5)))
	line = nil
	expect(t, !line.ContainsPoint(Point{}))
	expect(t, !(&Line{}).ContainsPoint(Point{}))
}

func TestLineIntersectsPoint(t *testing.T) {
	line := NewLine(v1, DefaultIndexOptions)
	expect(t, line.IntersectsPoint(P(0, 10)))
	expect(t, !line.IntersectsPoint(P(0, 0)))
	expect(t, line.IntersectsPoint(P(5, 0)))
	expect(t, line.IntersectsPoint(P(2.5, 5)))
	line = nil
	expect(t, !line.IntersectsPoint(Point{}))
	expect(t, !(&Line{}).IntersectsPoint(Point{}))
}

func TestLineContainsRect(t *testing.T) {
	line := NewLine(v1, DefaultIndexOptions)
	expect(t, !line.ContainsRect(R(0, 0, 10, 10)))
	expect(t, line.ContainsRect(R(0, 10, 0, 10)))
	line = NewLine(u1, DefaultIndexOptions)
	expect(t, line.ContainsRect(R(0, 0, 0, 10)))
	line = nil
	expect(t, !line.ContainsRect(Rect{}))
	expect(t, !(&Line{}).ContainsRect(Rect{}))
}

func TestLineIntersectsRect(t *testing.T) {
	line := NewLine(v1, DefaultIndexOptions)
	expect(t, line.IntersectsRect(R(0, 0, 10, 10)))
	expect(t, line.IntersectsRect(R(0, 0, 2.5, 5)))
	expect(t, !line.IntersectsRect(R(0, 0, 2.4, 5)))
	line = nil
	expect(t, !line.IntersectsRect(Rect{}))
	expect(t, !(&Line{}).IntersectsRect(Rect{}))
}

func TestLineContainsLine(t *testing.T) {
	expect(t, !P(15, 15).ContainsLine(L(P(15, 0), P(15, 15), P(30, 15))))
	expect(t, !P(15, 15).ContainsLine(L()))
	expect(t, !P(15, 15).ContainsLine(L(P(15, 15))))
	expect(t, P(15, 15).ContainsLine(L(P(15, 15), P(15, 15))))
	expect(t, R(0, 0, 30, 30).ContainsLine(L(P(15, 0), P(15, 15), P(30, 15))))
	expect(t, !R(0, 0, 30, 30).ContainsLine(L()))
	expect(t, !R(0, 0, 20, 20).ContainsLine(L(P(15, 0), P(15, 15), P(30, 15))))
	ln1 := L(P(5, 0), P(5, 5), P(10, 5), P(10, 10), P(15, 10), P(15, 15))
	lns := []*Line{
		L(P(7, 5), P(10, 5), P(10, 10), P(12, 10)),
		L(P(7, 5), P(8, 5), P(10, 5), P(10, 10), P(12, 10)),
		L(P(7, 5), P(8, 5), P(6, 5), P(10, 5), P(10, 8), P(10, 5), P(5, 5),
			P(10, 5), P(10, 10), P(12, 10)),
	}
	for _, ln2 := range lns {
		expect(t, ln1.ContainsLine(ln2))
	}
	expect(t, !ln1.ContainsLine(L(P(5, -1), P(5, 5), P(10, 5))))
	expect(t, !ln1.ContainsLine(L(P(5, 0), P(5, 5), P(5, 0), P(10, 0))))
	expect(t, !ln1.ContainsLine(L(P(5, 0), P(5, 5), P(10, 5), P(10, 10),
		P(15, 10), P(15, 15), P(20, 20))))
	expect(t, !ln1.ContainsLine(L()))
	expect(t, !L().ContainsLine(L(P(5, 0))))
	expect(t, !L(P(5, 0), P(10, 0)).ContainsLine(L(P(5, 0))))
	expect(t, R(0, 0, 30, 30).ContainsLine(L(P(15, 0), P(15, 15), P(30, 15))))
	expect(t, !R(0, 0, 30, 30).ContainsLine(L()))
}

func TestLineClockwise(t *testing.T) {
	expect(t, L(P(0, 0), P(0, 10), P(10, 10), P(10, 0), P(0, 0)).Clockwise())
	expect(t, !L(P(0, 0), P(10, 0), P(10, 10), P(0, 10), P(0, 0)).Clockwise())
	expect(t, L(P(0, 0), P(0, 10), P(10, 10)).Clockwise())
	expect(t, !L(P(0, 0), P(10, 0), P(10, 10)).Clockwise())
}

func TestLineIntersectsLine(t *testing.T) {
	lns := [][]Point{u1, u2, u3, u4, v1, v2, v3, v4}
	for i := 0; i < len(lns); i++ {
		for j := 0; j < len(lns); j++ {
			expect(t, NewLine(lns[i], DefaultIndexOptions).IntersectsLine(
				NewLine(lns[j], DefaultIndexOptions),
			))
		}
	}
	line := NewLine(u1, DefaultIndexOptions)
	expect(t, !line.IntersectsLine(NewLine(nil, DefaultIndexOptions)))
	expect(t, !NewLine(nil, DefaultIndexOptions).IntersectsLine(NewLine(nil, DefaultIndexOptions)))
	expect(t, !NewLine(nil, DefaultIndexOptions).IntersectsLine(line))
	expect(t, line.IntersectsLine(line.Move(5, 0)))
	expect(t, line.IntersectsLine(line.Move(10, 0)))
	expect(t, !line.IntersectsLine(line.Move(11, 0)))
	expect(t, !L(v1...).IntersectsLine(L(v1...).Move(0, 1)))
	expect(t, !L(v1...).IntersectsLine(L(v1...).Move(0, -1)))
}

func TestLineContainsPoly(t *testing.T) {
	line := NewLine(u1, DefaultIndexOptions)
	poly := NewPoly(octagon, nil, DefaultIndexOptions)
	expect(t, !line.ContainsPoly(poly))
	expect(t, line.ContainsPoly(NewPoly(
		[]Point{P(0, 10), P(0, 0), P(0, 10)},
		nil, DefaultIndexOptions,
	)))
	expect(t, line.ContainsPoly(NewPoly(
		[]Point{P(0, 0), P(10, 0), P(0, 0)},
		nil, DefaultIndexOptions,
	)))
	expect(t, !L().ContainsPoly(NewPoly(
		[]Point{P(0, 0), P(10, 0), P(0, 0)},
		nil, DefaultIndexOptions,
	)))
	expect(t, !line.ContainsPoly(NewPoly(nil, nil, DefaultIndexOptions)))
}

func TestLineIntersectsPoly(t *testing.T) {
	line := NewLine(u1, DefaultIndexOptions)
	poly := NewPoly(octagon, nil, DefaultIndexOptions)
	expect(t, line.IntersectsPoly(poly))
	expect(t, line.IntersectsPoly(poly.Move(5, 0)))
	expect(t, line.IntersectsPoly(poly.Move(10, 0)))
	expect(t, !line.IntersectsPoly(poly.Move(11, 0)))
	expect(t, !line.IntersectsPoly(poly.Move(15, 0)))
}
