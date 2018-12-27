// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geometry

import (
	"testing"
)

func TestRectCenter(t *testing.T) {
	expect(t, R(0, 0, 10, 10).Center() == P(5, 5))
	expect(t, R(0, 0, 0, 0).Center() == P(0, 0))
}

func TestRectArea(t *testing.T) {
	expect(t, R(0, 0, 10, 10).Area() == 100)
}

func TestRectMove(t *testing.T) {
	expect(t, R(1, 2, 3, 4).Move(5, 6) == R(6, 8, 8, 10))
}

func TestRectIndex(t *testing.T) {
	expect(t, (Rect{}).Index() == nil)
}

func TestRectNumPoints(t *testing.T) {
	expect(t, R(0, 0, 10, 10).NumPoints() == 5)
}

func TestRectNumSegments(t *testing.T) {
	expect(t, R(0, 0, 10, 10).NumSegments() == 4)
}

func TestRectPointAt(t *testing.T) {
	expect(t, R(0, 0, 10, 10).PointAt(0) == P(0, 0))
	expect(t, R(0, 0, 10, 10).PointAt(1) == P(10, 0))
	expect(t, R(0, 0, 10, 10).PointAt(2) == P(10, 10))
	expect(t, R(0, 0, 10, 10).PointAt(3) == P(0, 10))
	expect(t, R(0, 0, 10, 10).PointAt(4) == P(0, 0))
	defer func() { expect(t, recover() != nil) }()
	R(0, 0, 10, 10).PointAt(5)
}

func TestRectSegmentAt(t *testing.T) {
	expect(t, R(0, 0, 10, 10).SegmentAt(0) == S(0, 0, 10, 0))
	expect(t, R(0, 0, 10, 10).SegmentAt(1) == S(10, 0, 10, 10))
	expect(t, R(0, 0, 10, 10).SegmentAt(2) == S(10, 10, 0, 10))
	expect(t, R(0, 0, 10, 10).SegmentAt(3) == S(0, 10, 0, 0))
	defer func() { expect(t, recover() != nil) }()
	R(0, 0, 10, 10).SegmentAt(4)
}

func TestRectSearch(t *testing.T) {
	rect := R(0, 0, 10, 10)
	var count int
	rect.Search(R(0, 0, 10, 10), func(seg Segment, idx int) bool {
		expect(t, rect.PointAt(idx) == seg.A)
		count++
		return true
	})
	expect(t, count == 4)
	count = 0
	rect.Search(R(0, 4, 10, 5), func(seg Segment, idx int) bool {
		expect(t, rect.PointAt(idx) == seg.A)
		count++
		return true
	})
	expect(t, count == 2)
	count = 0
	rect.Search(R(0, 4, 10, 5), func(seg Segment, idx int) bool {
		expect(t, rect.PointAt(idx) == seg.A)
		count++
		return false
	})
	expect(t, count == 1)
}

func TestRectEmpty(t *testing.T) {
	expect(t, !R(0, 0, 10, 10).Empty())
}

func TestRectValid(t *testing.T) {
	expect(t, R(0, 0, 10, 10).Valid())
}

func TestRectRect(t *testing.T) {
	expect(t, R(0, 0, 10, 10).Rect() == R(0, 0, 10, 10))
}

func TestRectConvex(t *testing.T) {
	expect(t, R(0, 0, 10, 10).Convex())
}

func TestRectContainsPoint(t *testing.T) {
	for x := 0.0; x <= 10; x++ {
		for y := 0.0; y <= 10; y++ {
			expect(t, R(0, 0, 10, 10).ContainsPoint(P(x, y)))
		}
	}
	expect(t, !R(0, 0, 10, 10).ContainsPoint(P(-15, -15)))
	expect(t, !R(0, 0, 10, 10).ContainsPoint(P(-15, 5)))
	expect(t, !R(0, 0, 10, 10).ContainsPoint(P(-15, 15)))
	expect(t, !R(0, 0, 10, 10).ContainsPoint(P(0, -15)))
}

func TestRectIntersectsPoint(t *testing.T) {
	expect(t, R(0, 0, 10, 10).IntersectsPoint(P(5, 5)))
	expect(t, !R(0, 0, 10, 10).IntersectsPoint(P(15, 15)))
}

func TestRectContainsRect(t *testing.T) {
	expect(t, R(0, 0, 10, 10).ContainsRect(R(0, 0, 10, 10)))
	expect(t, R(0, 0, 10, 10).ContainsRect(R(2, 2, 10, 10)))
	expect(t, R(0, 0, 10, 10).ContainsRect(R(2, 2, 8, 8)))
	expect(t, R(0, 0, 10, 10).ContainsRect(R(0, 0, 8, 8)))
	expect(t, R(0, 0, 10, 10).ContainsRect(R(0, 2, 8, 8)))
	expect(t, R(0, 0, 10, 10).ContainsRect(R(2, 0, 8, 8)))
	expect(t, R(0, 0, 10, 10).ContainsRect(R(2, 2, 10, 8)))
	expect(t, R(0, 0, 10, 10).ContainsRect(R(2, 2, 8, 10)))
	expect(t, !R(0, 0, 10, 10).ContainsRect(R(-1, 0, 10, 10)))
	expect(t, !R(0, 0, 10, 10).ContainsRect(R(0, -1, 10, 10)))
	expect(t, !R(0, 0, 10, 10).ContainsRect(R(0, 0, 11, 10)))
	expect(t, !R(0, 0, 10, 10).ContainsRect(R(0, 0, 10, 11)))
}

func TestRectIntersectsRect(t *testing.T) {
	expect(t, R(0, 0, 10, 10).IntersectsRect(R(0, 0, 10, 10)))
	expect(t, R(0, 0, 10, 10).IntersectsRect(R(2, 2, 8, 8)))
	expect(t, R(0, 0, 10, 10).IntersectsRect(R(-1, 0, 10, 10)))
	expect(t, R(0, 0, 10, 10).IntersectsRect(R(0, -1, 10, 10)))
	expect(t, R(0, 0, 10, 10).IntersectsRect(R(0, 0, 11, 10)))
	expect(t, R(0, 0, 10, 10).IntersectsRect(R(0, 0, 10, 11)))
	expect(t, !R(0, 0, 10, 10).IntersectsRect(R(11, 0, 21, 10)))
	expect(t, !R(0, 0, 10, 10).IntersectsRect(R(0, 11, 10, 21)))
	expect(t, !R(0, 0, 10, 10).IntersectsRect(R(11, 0, 21, 10)))
	expect(t, !R(0, 0, 10, 10).IntersectsRect(R(11, 11, 21, 21)))
	expect(t, !R(0, 0, 10, 10).IntersectsRect(R(-11, 11, 1, 21)))
	expect(t, !R(0, 0, 10, 10).IntersectsRect(R(-11, -11, -1, -1)))
}

func TestRectContainsLine(t *testing.T) {
	expect(t, R(0, 0, 10, 10).ContainsLine(L(P(0, 0), P(10, 10))))
	expect(t, !R(0, 0, 10, 10).ContainsLine(L(P(0, 0), P(11, 11))))
	expect(t, !R(0, 0, 10, 10).ContainsLine(L()))
}

func TestRectIntersectsLine(t *testing.T) {
	expect(t, R(0, 0, 10, 10).IntersectsLine(L(P(0, 0), P(10, 10))))
	expect(t, R(0, 0, 10, 10).IntersectsLine(L(P(0, 0), P(11, 11))))
	expect(t, !R(0, 0, 10, 10).IntersectsLine(L()))
	expect(t, !R(0, 0, 10, 10).IntersectsLine(L(P(11, 11), P(12, 12))))
}

func TestRectContainsPoly(t *testing.T) {
	oct := NewPoly(octagon, nil, DefaultIndexOptions)
	expect(t, R(0, 0, 10, 10).ContainsPoly(oct))
	expect(t, !R(0, 0, 10, 10).ContainsPoly(oct.Move(1, 0)))
	expect(t, !R(0, 0, 10, 10).ContainsPoly(oct.Move(1, 1)))
	expect(t, !R(0, 0, 10, 10).ContainsPoly(NewPoly(nil, nil, DefaultIndexOptions)))
}

func TestRectIntersectsPoly(t *testing.T) {
	oct := NewPoly(octagon, nil, DefaultIndexOptions)
	expect(t, R(0, 0, 10, 10).IntersectsPoly(oct))
	expect(t, R(0, 0, 10, 10).IntersectsPoly(oct.Move(1, 0)))
	expect(t, R(0, 0, 10, 10).IntersectsPoly(oct.Move(0, 1)))
	expect(t, !R(0, 0, 10, 10).IntersectsPoly(oct.Move(10, 10)))
	expect(t, !R(0, 0, 10, 10).IntersectsPoly(oct.Move(11, 10)))
	expect(t, !R(0, 0, 10, 10).IntersectsPoly(oct.Move(-11, 0)))
	expect(t, !R(0, 0, 10, 10).IntersectsPoly(NewPoly(nil, nil, DefaultIndexOptions)))
}

func TestRectClockwise(t *testing.T) {
	expect(t, !R(10, 11, 12, 13).Clockwise())
}
