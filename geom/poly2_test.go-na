package geom

import (
	"reflect"
	"testing"
)

func dualPolyTest(
	t *testing.T, exterior []Point, holes [][]Point,
	do func(t *testing.T, poly Poly),
) {
	t.Run("noindex", func(t *testing.T) {
		do(t, NewPoly(exterior, holes, NoIndex))
	})
	t.Run("index", func(t *testing.T) {
		do(t, NewPoly(exterior, holes, Index))
	})
}

func TestPolyVarious(t *testing.T) {
	exterior := octagon
	small := []Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}
	dualPolyTest(t, exterior, [][]Point{small},
		func(t *testing.T, poly Poly) {
			expect(t, len(poly.Holes()) == 1)
			expect(t, reflect.DeepEqual(poly.Holes()[0].Points(), small))
			expect(t, !poly.ContainsPoint(P(0, 0)))
			expect(t, poly.ContainsPoint(P(3, 3)))
			expect(t, poly.ContainsPoint(P(4, 4)))
			expect(t, !poly.ContainsPoint(P(5, 5)))
			expect(t, poly.ContainsPoint(P(6, 6)))
			expect(t, poly.ContainsPoint(P(7, 7)))

			expect(t, poly.ContainsPoly(poly))
			expect(t, !poly.ContainsRing(poly.Exterior()))
			expect(t, !poly.ContainsRing(poly.Holes()[0]))

			expect(t, !poly.ContainsRing(newRingSimple(small).move(10, 0)))
			expect(t, !poly.ContainsPoly(NewPoly(
				newRingSimple(small).move(10, 0).Points(), nil, NoIndex)))

		},
	)

	ex1 := newRingSimple([]Point{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}})
	ex2 := newRingSimple([]Point{{1, 1}, {9, 1}, {9, 9}, {1, 9}, {1, 1}})
	ex3 := newRingSimple([]Point{{2, 2}, {8, 2}, {8, 8}, {2, 8}, {2, 2}})
	ex4 := newRingSimple([]Point{{3, 3}, {7, 3}, {7, 7}, {3, 7}, {3, 3}})
	ex5 := newRingSimple([]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}})
	out5 := ex5.move(10, 0)

	p1 := NewPoly(ex1.Points(), [][]Point{ex4.Points()}, NoIndex)
	p2 := NewPoly(ex2.Points(), [][]Point{ex3.Points()}, NoIndex)
	p3 := NewPoly(ex2.Points(), [][]Point{ex5.Points()}, NoIndex)

	expect(t, p1.ContainsPoly(p2))
	expect(t, !p1.ContainsPoly(p3))

	expect(t, p1.IntersectsPoly(p1))
	expect(t, p1.IntersectsRing(ex1))
	expect(t, !p1.IntersectsRing(out5))

	expect(t, !p2.IntersectsRing(ex5))

}
