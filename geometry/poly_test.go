// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geometry

import (
	"testing"
)

func newPolyIndexed(exterior []Point, holes [][]Point) *Poly {
	poly := NewPoly(exterior, holes, DefaultIndexOptions)
	poly.Exterior.(*baseSeries).buildIndex()
	for _, hole := range poly.Holes {
		hole.(*baseSeries).buildIndex()
	}
	return poly
}

func newPolySimple(exterior []Point, holes [][]Point) *Poly {
	poly := NewPoly(exterior, holes, DefaultIndexOptions)
	poly.Exterior.(*baseSeries).clearIndex()
	for _, hole := range poly.Holes {
		hole.(*baseSeries).clearIndex()
	}
	return poly
}

func dualPolyTest(
	t *testing.T, exterior []Point, holes [][]Point,
	do func(t *testing.T, poly *Poly),
) {
	t.Run("noindex", func(t *testing.T) {
		do(t, newPolySimple(exterior, holes))
	})
	t.Run("index", func(t *testing.T) {
		do(t, newPolyIndexed(exterior, holes))
	})
}

func TestPolyNewPoly(t *testing.T) {
	dualPolyTest(t, octagon, nil, func(t *testing.T, poly *Poly) {
		expect(t, !poly.Empty())
	})
	small := []Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}
	dualPolyTest(t, octagon, [][]Point{small}, func(t *testing.T, poly *Poly) {
		expect(t, !poly.Empty())
	})
}

func TestPolyGeometryDefaults(t *testing.T) {
	g := Geometry(&Poly{})
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

func TestPolyRect(t *testing.T) {
	dualPolyTest(t, octagon, nil, func(t *testing.T, poly *Poly) {
		expect(t, poly.Rect() == R(0, 0, 10, 10))
	})
	small := []Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}
	dualPolyTest(t, octagon, [][]Point{small}, func(t *testing.T, poly *Poly) {
		expect(t, poly.Rect() == R(0, 0, 10, 10))
	})
}

func TestPolyMove(t *testing.T) {
	small := []Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}
	dualPolyTest(t, octagon, [][]Point{small}, func(t *testing.T, poly *Poly) {
		poly2 := poly.Move(5, 8)
		expect(t, poly2.Rect() == R(5, 8, 15, 18))
	})
	poly := &Poly{Exterior: R(0, 0, 10, 10), Holes: []Ring{R(2, 2, 8, 8)}}
	poly2 := poly.Move(5, 8)
	expect(t, poly2.Rect() == R(5, 8, 15, 18))
	expect(t, len(poly2.Holes) == 1)
	expect(t, poly2.Holes[0].Rect() == R(7, 10, 13, 16))

	poly = nil
	expect(t, poly.Move(0, 0) == nil)
	expect(t, (&Poly{}).Move(0, 0) != nil)
}

func TestPolyContainsPoint(t *testing.T) {
	small := []Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}
	dualPolyTest(t, octagon, [][]Point{small}, func(t *testing.T, poly *Poly) {
		expect(t, !poly.ContainsPoint(P(0, 0)))
		expect(t, poly.ContainsPoint(P(0, 5)))
		expect(t, poly.ContainsPoint(P(3, 5)))
		expect(t, !poly.ContainsPoint(P(5, 5)))
	})
}

func TestPolyIntersectsPoint(t *testing.T) {
	small := []Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}
	dualPolyTest(t, octagon, [][]Point{small}, func(t *testing.T, poly *Poly) {
		expect(t, !poly.IntersectsPoint(P(0, 0)))
		expect(t, poly.IntersectsPoint(P(0, 5)))
		expect(t, poly.IntersectsPoint(P(3, 5)))
		expect(t, !poly.IntersectsPoint(P(5, 5)))
	})
	var poly *Poly
	expect(t, !poly.IntersectsPoint(Point{}))
	expect(t, !(&Poly{}).IntersectsPoint(Point{}))
}
func TestPolyContainsRect(t *testing.T) {
	ring := []Point{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}
	hole := []Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}
	dualPolyTest(t, ring, [][]Point{hole}, func(t *testing.T, poly *Poly) {
		expect(t, poly.ContainsRect(R(0, 0, 4, 4)))
		expect(t, !poly.ContainsRect(R(0, 0, 5, 5)))
		expect(t, !poly.ContainsRect(R(2, 2, 6, 6)))
		expect(t, !poly.ContainsRect(R(4.1, 4.1, 5.9, 5.9)))
		expect(t, !poly.ContainsRect(R(4.1, 4.1, 5.9, 5.9)))
	})

	var poly *Poly
	expect(t, !poly.ContainsRect(Rect{}))
	expect(t, !poly.IntersectsRect(Rect{}))
}

func TestPolyIntersectsRect(t *testing.T) {
	small := []Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}
	dualPolyTest(t, octagon, [][]Point{small}, func(t *testing.T, poly *Poly) {
		expect(t, poly.IntersectsRect(R(0, 4, 4, 6)))
		expect(t, poly.IntersectsRect(R(-1, 4, 4, 6)))
		expect(t, poly.IntersectsRect(R(4, 4, 6, 6)))
		expect(t, !poly.IntersectsRect(R(4.1, 4.1, 5.9, 5.9)))
		expect(t, !poly.IntersectsRect(R(0, 0, 1.4, 1.4)))
		expect(t, poly.IntersectsRect(R(0, 0, 1.5, 1.5)))
		expect(t, !poly.IntersectsRect(R(0, 0, 10, 10).Move(11, 0)))
	})
}

func TestPolyContainsLine(t *testing.T) {
	small := []Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}
	dualPolyTest(t, octagon, [][]Point{small}, func(t *testing.T, poly *Poly) {
		expect(t, poly.ContainsLine(L(P(3, 3), P(3, 7), P(7, 7), P(7, 3))))
		expect(t, !poly.ContainsLine(L(P(-1, 3), P(3, 7), P(7, 7), P(7, 3))))
		expect(t, poly.ContainsLine(L(P(4, 3), P(3, 7), P(7, 7), P(7, 3))))
		expect(t, !poly.ContainsLine(L(P(5, 3), P(3, 7), P(7, 7), P(7, 3))))
	})
}

func TestPolyIntersectsLine(t *testing.T) {
	holes := [][]Point{[]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}}
	dualPolyTest(t, octagon, holes, func(t *testing.T, poly *Poly) {
		expect(t, poly.IntersectsLine(L(P(3, 3), P(4, 4))))
		expect(t, poly.IntersectsLine(L(P(-1, 3), P(3, 7), P(7, 7), P(7, 3))))
		expect(t, poly.IntersectsLine(L(P(4, 3), P(3, 7), P(7, 7), P(7, 3))))
		expect(t, poly.IntersectsLine(L(P(5, 3), P(3, 7), P(7, 7), P(7, 3))))
		expect(t, !poly.IntersectsLine(
			L(P(5, 3), P(3, 7), P(7, 7), P(7, 3)).Move(11, 0),
		))
	})
}
func TestPolyContainsPoly(t *testing.T) {
	holes1 := [][]Point{[]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}}
	holes2 := [][]Point{[]Point{{5, 4}, {7, 4}, {7, 6}, {5, 6}, {5, 4}}}
	poly1 := NewPoly(octagon, holes1, DefaultIndexOptions)
	poly2 := NewPoly(octagon, holes2, DefaultIndexOptions)

	expect(t, !poly1.ContainsPoly(NewPoly(holes2[0], nil, DefaultIndexOptions)))
	expect(t, !poly1.ContainsPoly(poly2))

	dualPolyTest(t, octagon, holes1, func(t *testing.T, poly *Poly) {
		expect(t, poly.ContainsPoly(poly1))
		expect(t, !poly.ContainsPoly(poly1.Move(1, 0)))
		expect(t, poly.ContainsPoly(NewPoly(holes1[0], nil, DefaultIndexOptions)))
		expect(t, !poly.ContainsPoly(NewPoly(holes2[0], nil, DefaultIndexOptions)))
	})
}

func TestPolyClockwise(t *testing.T) {
	expect(t, !NewPoly(bowtie, nil, DefaultIndexOptions).Clockwise())
	var poly *Poly
	expect(t, !poly.Clockwise())
}

// https://github.com/tidwall/tile38/issues/369
func Test369(t *testing.T) {
	polyHoles := NewPoly([]Point{
		{-122.44154334068298, 37.73179457567642},
		{-122.43935465812682, 37.73179457567642},
		{-122.43935465812682, 37.7343740514423},
		{-122.44154334068298, 37.7343740514423},
		{-122.44154334068298, 37.73179457567642},
	}, [][]Point{
		[]Point{
			{-122.44104981422423, 37.73286371140448},
			{-122.44104981422423, 37.73424677678513},
			{-122.43990182876587, 37.73424677678513},
			{-122.43990182876587, 37.73286371140448},
			{-122.44104981422423, 37.73286371140448},
		},
		[]Point{
			{-122.44109272956847, 37.731870943026074},
			{-122.43976235389708, 37.731870943026074},
			{-122.43976235389708, 37.7326855231885},
			{-122.44109272956847, 37.7326855231885},
			{-122.44109272956847, 37.731870943026074},
		},
	}, DefaultIndexOptions)
	a := NewPoly([]Point{
		{-122.4408378, 37.7341129},
		{-122.4408378, 37.733},
		{-122.44, 37.733},
		{-122.44, 37.7343129},
		{-122.4408378, 37.7341129},
	}, nil, DefaultIndexOptions)
	b := NewPoly([]Point{
		{-122.44091033935547, 37.731981251280985},
		{-122.43994474411011, 37.731981251280985},
		{-122.43994474411011, 37.73254976045042},
		{-122.44091033935547, 37.73254976045042},
		{-122.44091033935547, 37.731981251280985},
	}, nil, DefaultIndexOptions)
	c := NewPoly([]Point{
		{-122.4408378, 37.7341129},
		{-122.4408378, 37.733},
		{-122.44, 37.733},
		{-122.44, 37.7341129},
		{-122.4408378, 37.7341129},
	}, nil, DefaultIndexOptions)
	d := NewPoly([]Point{
		{-182.4408378, 37.7341129},
		{-122.4408378, 37.733},
		{-122.44, 37.733},
		{-122.44, 37.7341129},
		{-122.4408378, 137.7341129},
	}, nil, DefaultIndexOptions)
	expect(t, polyHoles.IntersectsPoly(a))
	expect(t, !polyHoles.IntersectsPoly(b))
	expect(t, !polyHoles.IntersectsPoly(c))
	expect(t, a.IntersectsPoly(polyHoles))
	expect(t, !b.IntersectsPoly(polyHoles))
	expect(t, !c.IntersectsPoly(polyHoles))
	expect(t, a.Valid())
	expect(t, b.Valid())
	expect(t, c.Valid())
	expect(t, !d.Valid())
}
