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
		})

}
