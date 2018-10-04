package geom

import "testing"

func newPolyIndexed(exterior []Point, holes [][]Point) *Poly {
	poly := NewPoly(exterior, holes)
	poly.Exterior.(*baseSeries).buildTree()
	for _, hole := range poly.Holes {
		hole.(*baseSeries).buildTree()
	}
	return poly
}

func newPolySimple(exterior []Point, holes [][]Point) *Poly {
	poly := NewPoly(exterior, holes)
	poly.Exterior.(*baseSeries).tree = nil
	for _, hole := range poly.Holes {
		hole.(*baseSeries).tree = nil
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
	poly := &Poly{Exterior: R(0, 0, 10, 10), Holes: []RingX{R(2, 2, 8, 8)}}
	poly2 := poly.Move(5, 8)
	expect(t, poly2.Rect() == R(5, 8, 15, 18))
	expect(t, len(poly2.Holes) == 1)
	expect(t, poly2.Holes[0].Rect() == R(7, 10, 13, 16))
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
}
func TestPolyContainsRect(t *testing.T) {
	ring := []Point{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}
	hole := []Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}
	dualPolyTest(t, ring, [][]Point{hole}, func(t *testing.T, poly *Poly) {
		// expect(t, poly.ContainsRect(R(0, 0, 4, 4)))
		// expect(t, !poly.ContainsRect(R(0, 0, 5, 5)))
		expect(t, !poly.ContainsRect(R(2, 2, 6, 6)))
		// expect(t, !poly.ContainsRect(R(4.1, 4.1, 5.9, 5.9)))
		// expect(t, !poly.ContainsRect(R(4.1, 4.1, 5.9, 5.9)))
	})
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
	poly1 := NewPoly(octagon, holes1)
	poly2 := NewPoly(octagon, holes2)

	expect(t, !poly1.ContainsPoly(NewPoly(holes2[0], nil)))

	expect(t, !poly1.ContainsPoly(poly2))

	dualPolyTest(t, octagon, holes1, func(t *testing.T, poly *Poly) {
		expect(t, poly.ContainsPoly(poly1))
		expect(t, !poly.ContainsPoly(poly1.Move(1, 0)))
		expect(t, poly.ContainsPoly(NewPoly(holes1[0], nil)))
		expect(t, !poly.ContainsPoly(NewPoly(holes2[0], nil)))
		// expect(t, !poly.ContainsPoly(NewPoly(holes[0], nil).Move(1, 0)))
		// expect(t, !poly.ContainsPoly(NewPoly(holes[0], nil).Move(1, 0)))
	})

}

// func TestPolyVarious(t *testing.T) {
// 	exterior := octagon
// 	small := []Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}
// 	dualPolyTest(t, exterior, [][]Point{small},
// 		func(t *testing.T, poly *Poly) {
// 			expect(t, len(poly.Holes) == 1)
// 			expect(t, reflect.DeepEqual(
// 				seriesCopyPoints(poly.Holes[0]), small))
// 			expect(t, !poly.ContainsPoint(P(0, 0)))
// 			expect(t, poly.ContainsPoint(P(3, 3)))
// 			expect(t, poly.ContainsPoint(P(4, 4)))
// 			expect(t, !poly.ContainsPoint(P(5, 5)))
// 			expect(t, poly.ContainsPoint(P(6, 6)))
// 			expect(t, poly.ContainsPoint(P(7, 7)))
// 			expect(t, poly.IntersectsPoint(P(7, 7)))

// 			expect(t, poly.ContainsPoly(poly))
// 			// expect(t, !poly.ContainsRing(poly.Exterior))
// 			// expect(t, !poly.ContainsRing(poly.Holes[0]))

// 			// expect(t, !poly.ContainsRing(newRingSimple(small).move(10, 0)))
// 			expect(t, !poly.ContainsPoly(newPolySimple(
// 				seriesCopyPoints(newRingSimple2(small).(*baseSeries).Move(10, 0)),
// 				nil)))

// 		},
// 	)

// 	ex1 := newRingSimple2([]Point{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}})
// 	ex2 := newRingSimple2([]Point{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}})
// 	ex3 := newRingSimple2([]Point{{2, 2}, {8, 2}, {8, 8}, {2, 8}, {2, 2}})
// 	ex4 := newRingSimple2([]Point{{3, 3}, {7, 3}, {7, 7}, {3, 7}, {3, 3}})
// 	ex5 := newRingSimple2([]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}})
// 	// out5 := ex5.move(10, 0)

// 	p1 := newPolySimple(seriesCopyPoints(ex1),
// 		[][]Point{seriesCopyPoints(ex4)})
// 	p2 := newPolySimple(seriesCopyPoints(ex2),
// 		[][]Point{seriesCopyPoints(ex3)})
// 	p3 := newPolySimple(seriesCopyPoints(ex2),
// 		[][]Point{seriesCopyPoints(ex5)})

// 	expect(t, p1.ContainsPoly(p2))
// 	expect(t, !p1.ContainsPoly(p3))

// 	expect(t, p1.IntersectsPoly(p1))
// 	// expect(t, p1.IntersectsRing(ex1))
// 	// expect(t, !p1.IntersectsRing(out5))

// 	// expect(t, !p2.IntersectsRing(ex5))

// }
