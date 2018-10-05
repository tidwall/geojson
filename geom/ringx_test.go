package geom

import (
	"fmt"
	"testing"
)

func newRingXSimple(points []Point) RingX {
	ring := newRingX(points)
	if ring.(*baseSeries).tree != nil {
		ring.(*baseSeries).tree = nil
	}
	return ring
}
func newRingXIndexed(points []Point) RingX {
	ring := newRingX(points)
	if ring.(*baseSeries).tree == nil {
		ring.(*baseSeries).buildTree()
	}
	return ring
}

func testDualRingX(t *testing.T, points []Point, do func(t *testing.T, ring RingX)) {
	t.Run("index", func(t *testing.T) {
		do(t, newRingXSimple(points))
	})
	t.Run("noindex", func(t *testing.T) {
		do(t, newRingXIndexed(points))
	})
}

func TestRingXContainsPoint(t *testing.T) {
	// testDualRingX(t, octagon, func(t *testing.T, ring RingX) {
	// 	expect(t, !ringxContainsPoint(ring, P(0, 0), true).hit)
	// 	expect(t, ringxContainsPoint(ring, P(0, 5), true).hit)
	// 	expect(t, !ringxContainsPoint(ring, P(0, 5), false).hit)
	// 	expect(t, ringxContainsPoint(ring, P(4, 4), true).hit)
	// 	expect(t, !ringxContainsPoint(ring, P(1.4, 1.4), true).hit)
	// 	expect(t, ringxContainsPoint(ring, P(1.5, 1.5), true).hit)
	// 	expect(t, !ringxContainsPoint(ring, P(1.5, 1.5), false).hit)
	// })
	shape := []Point{
		P(0, 0), P(4, 0), P(4, 3),
		P(3, 4), P(1, 4),
		P(0, 3), P(0, 0),
	}
	ring := newRingX(shape)
	//expect(t, !insideshpext(P(4, 3.5), shape))
	expect(t, !ringxContainsPoint(ring, P(4, 3.5), true).hit)

	// testDualRingX(t, ring, func(t *testing.T, ring RingX) {
	// 	expect(t, !ringxContainsPoint(ring, P(4, 3.5), true).hit)
	// })
}

func TestRingXIntersectsPoint(t *testing.T) {
	testDualRingX(t, octagon, func(t *testing.T, ring RingX) {
		expect(t, !ringxIntersectsPoint(ring, P(0, 0), true).hit)
		expect(t, ringxIntersectsPoint(ring, P(0, 5), true).hit)
		expect(t, !ringxIntersectsPoint(ring, P(0, 5), false).hit)
		expect(t, ringxIntersectsPoint(ring, P(4, 4), true).hit)
		expect(t, !ringxIntersectsPoint(ring, P(1.4, 1.4), true).hit)
		expect(t, ringxIntersectsPoint(ring, P(1.5, 1.5), true).hit)
		expect(t, !ringxIntersectsPoint(ring, P(1.5, 1.5), false).hit)
	})
}

func TestRingXIntersectsSegment(t *testing.T) {
	t.Run("Cases", func(t *testing.T) {
		ring := newRingX([]Point{
			P(0, 0), P(4, 0), P(4, 4), P(0, 4), P(0, 0),
		})
		t.Run("1", func(t *testing.T) {
			expect(t, ringxIntersectsSegment(ring, S(2, 2, 4, 4), true))
			expect(t, ringxIntersectsSegment(ring, S(2, 2, 4, 4), false))
		})
		t.Run("2", func(t *testing.T) {
			expect(t, ringxIntersectsSegment(ring, S(2, 0, 4, 2), true))
			expect(t, ringxIntersectsSegment(ring, S(2, 0, 4, 2), false))
		})
		t.Run("3", func(t *testing.T) {
			expect(t, ringxIntersectsSegment(ring, S(1, 0, 3, 0), true))
			expect(t, !ringxIntersectsSegment(ring, S(1, 0, 3, 0), false))
		})
		t.Run("4", func(t *testing.T) {
			expect(t, ringxIntersectsSegment(ring, S(2, 4, 2, 0), true))
			expect(t, ringxIntersectsSegment(ring, S(2, 4, 2, 0), false))
		})
		t.Run("5", func(t *testing.T) {
			expect(t, ringxIntersectsSegment(ring, S(2, 0, 2, -3), true))
			expect(t, !ringxIntersectsSegment(ring, S(2, 0, 2, -3), false))
		})
		t.Run("6", func(t *testing.T) {
			expect(t, !ringxIntersectsSegment(ring, S(2, -1, 2, -3), true))
			expect(t, !ringxIntersectsSegment(ring, S(2, -1, 2, -3), false))
		})
		t.Run("7", func(t *testing.T) {
			expect(t, ringxIntersectsSegment(ring, S(-1, 2, 3, 2), true))
			expect(t, ringxIntersectsSegment(ring, S(-1, 2, 3, 2), false))
		})
		t.Run("8", func(t *testing.T) {
			expect(t, ringxIntersectsSegment(ring, S(-1, 3, 2, 1), true))
			expect(t, ringxIntersectsSegment(ring, S(-1, 3, 2, 1), false))
		})
		t.Run("9", func(t *testing.T) {
			expect(t, ringxIntersectsSegment(ring, S(-1, 2, 5, 1), true))
			expect(t, ringxIntersectsSegment(ring, S(-1, 2, 5, 1), false))
		})
		t.Run("10", func(t *testing.T) {
			expect(t, ringxIntersectsSegment(ring, S(2, 5, 1, -1), true))
			expect(t, ringxIntersectsSegment(ring, S(2, 5, 1, -1), false))
		})
		t.Run("11", func(t *testing.T) {
			expect(t, ringxIntersectsSegment(ring, S(0, 4, 4, 0), true))
			expect(t, ringxIntersectsSegment(ring, S(0, 4, 4, 0), false))
		})
		t.Run("12", func(t *testing.T) {
			ring := newRingX([]Point{
				{0, 0}, {2, 0}, {4, 0}, {4, 4}, {0, 4}, {0, 0},
			})
			expect(t, ringxIntersectsSegment(ring, S(1, 0, 3, 0), true))
			expect(t, !ringxIntersectsSegment(ring, S(1, 0, 3, 0), false))
		})
		t.Run("13", func(t *testing.T) {
			ring := newRingX([]Point{
				{0, 0}, {4, 0}, {4, 4}, {2, 4}, {0, 4}, {0, 0},
			})
			expect(t, ringxIntersectsSegment(ring, S(0, 4, 4, 4), true))
			expect(t, !ringxIntersectsSegment(ring, S(0, 4, 4, 4), false))
		})
		t.Run("14", func(t *testing.T) {
			ring := newRingX([]Point{
				{0, 0}, {4, 0}, {4, 4}, {2, 4}, {0, 4}, {0, 0},
			})
			expect(t, ringxIntersectsSegment(ring, S(-1, -2, 0, 0), true))
			expect(t, !ringxIntersectsSegment(ring, S(-1, -2, 0, 0), false))
		})
		t.Run("15", func(t *testing.T) {
			ring := newRingX([]Point{
				{0, 0}, {4, 0}, {4, 4}, {2, 4}, {0, 4}, {0, 0},
			})
			expect(t, ringxIntersectsSegment(ring, S(0, 4, 5, 4), true))
			expect(t, !ringxIntersectsSegment(ring, S(0, 4, 5, 4), false))
		})
		t.Run("16", func(t *testing.T) {
			ring := newRingX([]Point{
				{0, 0}, {4, 0}, {4, 4}, {2, 4}, {0, 4}, {0, 0},
			})
			expect(t, ringxIntersectsSegment(ring, S(1, 4, 5, 4), true))
			expect(t, !ringxIntersectsSegment(ring, S(1, 4, 5, 4), false))
		})
		t.Run("17", func(t *testing.T) {
			ring := newRingX([]Point{
				{0, 0}, {4, 0}, {4, 4}, {2, 4}, {0, 4}, {0, 0},
			})
			expect(t, ringxIntersectsSegment(ring, S(1, 4, 4, 4), true))
			expect(t, !ringxIntersectsSegment(ring, S(1, 4, 4, 4), false))
		})
	})

	// convex shape
	t.Run("Octagon", func(t *testing.T) {
		shape := octagon
		t.Run("Outside", func(t *testing.T) {
			// test coming off of the edges
			segs := []Segment{
				S(0, 5, -5, 5).Move(-1, 0),           // left
				S(1.5, 1.5, -3.5, -3.5).Move(-1, -1), // bottom-left corner
				S(5, 0, 5, -5).Move(0, -1),           // bottom
				S(8.5, 1.5, 13.5, -3.5).Move(1, -1),  // bottom-right corner
				S(10, 5, 15, 5).Move(1, 0),           // right
				S(8.5, 8.5, 13.5, 13.5).Move(1, 1),   // top-right corner
				S(5, 10, 5, 15).Move(0, 5),           // top
				S(1.5, 8.5, -3.5, 13.5).Move(-1, 1),  // top-left corner
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, !ringxIntersectsSegment(ring, seg, true))
						expect(t, !ringxIntersectsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, !ringxIntersectsSegment(ring, seg, true))
						expect(t, !ringxIntersectsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("OutsideEdgeTouch", func(t *testing.T) {
			// test coming off of the edges
			segs := []Segment{
				S(0, 5, -5, 5),          // left
				S(1.5, 1.5, -3.5, -3.5), // bottom-left corner
				S(5, 0, 5, -5),          // bottom
				S(8.5, 1.5, 13.5, -3.5), // bottom-right corner
				S(10, 5, 15, 5),         // right
				S(8.5, 8.5, 13.5, 13.5), // top-right corner
				S(5, 10, 5, 15),         // top
				S(1.5, 8.5, -3.5, 13.5), // top-left corner
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, !ringxIntersectsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, !ringxIntersectsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("EdgeCross", func(t *testing.T) {
			segs := []Segment{
				S(0, 5, -5, 5).Move(2.5, 0),              // left
				S(1.5, 1.5, -3.5, -3.5).Move(2.5, 2.5),   // bottom-left corner
				S(5, 0, 5, -5).Move(0, 2.5),              // bottom
				S(8.5, 1.5, 13.5, -3.5).Move(-2.5, 2.5),  // bottom-right corner
				S(10, 5, 15, 5).Move(-2.5, 0),            // right
				S(8.5, 8.5, 13.5, 13.5).Move(-2.5, -2.5), // top-right corner
				S(5, 10, 5, 15).Move(0, -2.5),            // top
				S(1.5, 8.5, -3.5, 13.5).Move(2.5, -2.5),  // top-left corner
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("InsideEdgeTouch", func(t *testing.T) {
			segs := []Segment{
				S(0, 5, -5, 5).Move(5, 0),            // left
				S(1.5, 1.5, -3.5, -3.5).Move(5, 5),   // bottom-left corner
				S(5, 0, 5, -5).Move(0, 5),            // bottom
				S(8.5, 1.5, 13.5, -3.5).Move(-5, 5),  // bottom-right corner
				S(10, 5, 15, 5).Move(-5, 0),          // right
				S(8.5, 8.5, 13.5, 13.5).Move(-5, -5), // top-right corner
				S(5, 10, 5, 15).Move(0, -5),          // top
				S(1.5, 8.5, -3.5, 13.5).Move(5, -5),  // top-left corner
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("Inside", func(t *testing.T) {
			segs := []Segment{
				S(0, 5, -5, 5).Move(6, 0),            // left
				S(1.5, 1.5, -3.5, -3.5).Move(6, 6),   // bottom-left corner
				S(5, 0, 5, -5).Move(0, 6),            // bottom
				S(8.5, 1.5, 13.5, -3.5).Move(-6, 6),  // bottom-right corner
				S(10, 5, 15, 5).Move(-6, 0),          // right
				S(8.5, 8.5, 13.5, 13.5).Move(-6, -6), // top-right corner
				S(5, 10, 5, 15).Move(0, -6),          // top
				S(1.5, 8.5, -3.5, 13.5).Move(6, -6),  // top-left corner
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("InsideEdgeToEdge", func(t *testing.T) {
			segs := []Segment{
				S(0, 5, 10, 5),        // left to right
				S(1.5, 1.5, 8.5, 8.5), // bottom-left to top-right
				S(5, 0, 5, 10),        // bottom to top
				S(8.5, 1.5, 1.5, 8.5), // bottom-right to top-left
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("PassoverEdgeToEdge", func(t *testing.T) {
			segs := []Segment{
				S(-1, 5, 11, 5),       // left to right
				S(0.5, 0.5, 9.5, 9.5), // bottom-left to top-right
				S(5, -1, 5, 11),       // bottom to top
				S(9.5, 0.5, 0.5, 9.5), // bottom-right to top-left
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
					})
				})
			}
		})
	})

	// concave shape
	t.Run("Concave", func(t *testing.T) {
		shape := concave1
		t.Run("Outside", func(t *testing.T) {
			segs := []Segment{
				S(0, 7.5, -5, 7.5).Move(-1, 0), // left
				S(2.5, 5, 2.5, 0).Move(0, -1),  // bottom-left-1
				S(5, 2.5, 0, 2.5).Move(-1, 0),  // bottom-left-2
				S(7.5, 0, 7.5, -5).Move(0, -1), // bottom
				S(10, 5, 15, 5).Move(1, 0),     // right
				S(5, 10, 5, 15).Move(0, 1),     // top
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, !ringxIntersectsSegment(ring, seg, true))
						expect(t, !ringxIntersectsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, !ringxIntersectsSegment(ring, seg, true))
						expect(t, !ringxIntersectsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("OutsideEdgeTouch", func(t *testing.T) {
			segs := []Segment{
				S(0, 7.5, -5, 7.5), // left
				S(2.5, 5, 2.5, 0),  // bottom-left-1
				S(5, 2.5, 0, 2.5),  // bottom-left-2
				S(7.5, 0, 7.5, -5), // bottom
				S(10, 5, 15, 5),    // right
				S(5, 10, 5, 15),    // top
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, !ringxIntersectsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, !ringxIntersectsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("EdgeCross", func(t *testing.T) {
			segs := []Segment{
				S(0, 7.5, -5, 7.5).Move(2.5, 0), // left
				S(2.5, 5, 2.5, 0).Move(0, 2.5),  // bottom-left-1
				S(5, 2.5, 0, 2.5).Move(2.5, 0),  // bottom-left-2
				S(7.5, 0, 7.5, -5).Move(0, 2.5), // bottom
				S(10, 5, 15, 5).Move(-2.5, 0),   // right
				S(5, 10, 5, 15).Move(0, -2.5),   // top
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("InsideEdgeTouch", func(t *testing.T) {
			segs := []Segment{
				S(0, 7.5, -3, 7.5).Move(3, 0), // 0:left
				S(2.5, 5, 2.5, 2).Move(0, 3),  // 1:bottom-left-1
				S(5, 2.5, 2, 2.5).Move(3, 0),  // 2:bottom-left-2
				S(7.5, 0, 7.5, -3).Move(0, 3), // 3:bottom
				S(10, 5, 13, 5).Move(-3, 0),   // 4:right
				S(5, 10, 5, 13).Move(0, -3),   // 5:top
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("Inside", func(t *testing.T) {
			segs := []Segment{
				S(0, 7.5, -3, 7.5).Move(4, 0), // 0:left
				S(2.5, 5, 2.5, 2).Move(0, 4),  // 1:bottom-left-1
				S(5, 2.5, 2, 2.5).Move(4, 0),  // 2:bottom-left-2
				S(7.5, 0, 7.5, -3).Move(0, 4), // 3:bottom
				S(10, 5, 13, 5).Move(-4, 0),   // 4:right
				S(5, 10, 5, 13).Move(0, -4),   // 5:top
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("InsideEdgeToEdge", func(t *testing.T) {
			segs := []Segment{
				S(0, 7.5, 10, 7.5), // 0:left
				S(2.5, 5, 2.5, 10), // 1:bottom-left-1
				S(5, 2.5, 10, 2.5), // 2:bottom-left-2
				S(7.5, 0, 7.5, 10), // 3:bottom
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("PassoverEdgeToEdge", func(t *testing.T) {
			segs := []Segment{
				S(-1, 7.5, 11, 7.5), // 0:left
				S(2.5, 4, 2.5, 11),  // 1:bottom-left-1
				S(4, 2.5, 11, 2.5),  // 2:bottom-left-2
				S(7.5, -1, 7.5, 11), // 3:bottom
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxIntersectsSegment(ring, seg, true))
						expect(t, ringxIntersectsSegment(ring, seg, false))
					})
				})
			}
		})
	})
}

func TestRingXContainsSegment(t *testing.T) {
	// convex shape
	t.Run("Octagon", func(t *testing.T) {
		shape := octagon
		t.Run("Outside", func(t *testing.T) {
			// test coming off of the edges
			segs := []Segment{
				S(0, 5, -5, 5).Move(-1, 0),           // left
				S(1.5, 1.5, -3.5, -3.5).Move(-1, -1), // bottom-left corner
				S(5, 0, 5, -5).Move(0, -1),           // bottom
				S(8.5, 1.5, 13.5, -3.5).Move(1, -1),  // bottom-right corner
				S(10, 5, 15, 5).Move(1, 0),           // right
				S(8.5, 8.5, 13.5, 13.5).Move(1, 1),   // top-right corner
				S(5, 10, 5, 15).Move(0, 5),           // top
				S(1.5, 8.5, -3.5, 13.5).Move(-1, 1),  // top-left corner
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("OutsideEdgeTouch", func(t *testing.T) {
			// test coming off of the edges
			segs := []Segment{
				S(0, 5, -5, 5),          // left
				S(1.5, 1.5, -3.5, -3.5), // bottom-left corner
				S(5, 0, 5, -5),          // bottom
				S(8.5, 1.5, 13.5, -3.5), // bottom-right corner
				S(10, 5, 15, 5),         // right
				S(8.5, 8.5, 13.5, 13.5), // top-right corner
				S(5, 10, 5, 15),         // top
				S(1.5, 8.5, -3.5, 13.5), // top-left corner
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("EdgeCross", func(t *testing.T) {
			segs := []Segment{
				S(0, 5, -5, 5).Move(2.5, 0),              // left
				S(1.5, 1.5, -3.5, -3.5).Move(2.5, 2.5),   // bottom-left corner
				S(5, 0, 5, -5).Move(0, 2.5),              // bottom
				S(8.5, 1.5, 13.5, -3.5).Move(-2.5, 2.5),  // bottom-right corner
				S(10, 5, 15, 5).Move(-2.5, 0),            // right
				S(8.5, 8.5, 13.5, 13.5).Move(-2.5, -2.5), // top-right corner
				S(5, 10, 5, 15).Move(0, -2.5),            // top
				S(1.5, 8.5, -3.5, 13.5).Move(2.5, -2.5),  // top-left corner
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("InsideEdgeTouch", func(t *testing.T) {
			segs := []Segment{
				S(0, 5, -5, 5).Move(5, 0),            // left
				S(1.5, 1.5, -3.5, -3.5).Move(5, 5),   // bottom-left corner
				S(5, 0, 5, -5).Move(0, 5),            // bottom
				S(8.5, 1.5, 13.5, -3.5).Move(-5, 5),  // bottom-right corner
				S(10, 5, 15, 5).Move(-5, 0),          // right
				S(8.5, 8.5, 13.5, 13.5).Move(-5, -5), // top-right corner
				S(5, 10, 5, 15).Move(0, -5),          // top
				S(1.5, 8.5, -3.5, 13.5).Move(5, -5),  // top-left corner
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("Inside", func(t *testing.T) {
			segs := []Segment{
				S(0, 5, -5, 5).Move(6, 0),            // left
				S(1.5, 1.5, -3.5, -3.5).Move(6, 6),   // bottom-left corner
				S(5, 0, 5, -5).Move(0, 6),            // bottom
				S(8.5, 1.5, 13.5, -3.5).Move(-6, 6),  // bottom-right corner
				S(10, 5, 15, 5).Move(-6, 0),          // right
				S(8.5, 8.5, 13.5, 13.5).Move(-6, -6), // top-right corner
				S(5, 10, 5, 15).Move(0, -6),          // top
				S(1.5, 8.5, -3.5, 13.5).Move(6, -6),  // top-left corner
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxContainsSegment(ring, seg, true))
						expect(t, ringxContainsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxContainsSegment(ring, seg, true))
						expect(t, ringxContainsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("InsideEdgeToEdge", func(t *testing.T) {
			segs := []Segment{
				S(0, 5, 10, 5),        // left to right
				S(1.5, 1.5, 8.5, 8.5), // bottom-left to top-right
				S(5, 0, 5, 10),        // bottom to top
				S(8.5, 1.5, 1.5, 8.5), // bottom-right to top-left
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("PassoverEdgeToEdge", func(t *testing.T) {
			segs := []Segment{
				S(-1, 5, 11, 5),       // left to right
				S(0.5, 0.5, 9.5, 9.5), // bottom-left to top-right
				S(5, -1, 5, 11),       // bottom to top
				S(9.5, 0.5, 0.5, 9.5), // bottom-right to top-left
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
					})
				})
			}
		})
	})

	// concave shape
	t.Run("Concave", func(t *testing.T) {
		shape := concave1
		t.Run("Outside", func(t *testing.T) {
			segs := []Segment{
				S(0, 7.5, -5, 7.5).Move(-1, 0), // left
				S(2.5, 5, 2.5, 0).Move(0, -1),  // bottom-left-1
				S(5, 2.5, 0, 2.5).Move(-1, 0),  // bottom-left-2
				S(7.5, 0, 7.5, -5).Move(0, -1), // bottom
				S(10, 5, 15, 5).Move(1, 0),     // right
				S(5, 10, 5, 15).Move(0, 1),     // top
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("OutsideEdgeTouch", func(t *testing.T) {
			segs := []Segment{
				S(0, 7.5, -5, 7.5), // left
				S(2.5, 5, 2.5, 0),  // bottom-left-1
				S(5, 2.5, 0, 2.5),  // bottom-left-2
				S(7.5, 0, 7.5, -5), // bottom
				S(10, 5, 15, 5),    // right
				S(5, 10, 5, 15),    // top
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("EdgeCross", func(t *testing.T) {
			segs := []Segment{
				S(0, 7.5, -5, 7.5).Move(2.5, 0), // left
				S(2.5, 5, 2.5, 0).Move(0, 2.5),  // bottom-left-1
				S(5, 2.5, 0, 2.5).Move(2.5, 0),  // bottom-left-2
				S(7.5, 0, 7.5, -5).Move(0, 2.5), // bottom
				S(10, 5, 15, 5).Move(-2.5, 0),   // right
				S(5, 10, 5, 15).Move(0, -2.5),   // top
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("InsideEdgeTouch", func(t *testing.T) {
			segs := []Segment{
				S(0, 7.5, -3, 7.5).Move(3, 0), // 0:left
				S(2.5, 5, 2.5, 2).Move(0, 3),  // 1:bottom-left-1
				S(5, 2.5, 2, 2.5).Move(3, 0),  // 2:bottom-left-2
				S(7.5, 0, 7.5, -3).Move(0, 3), // 3:bottom
				S(10, 5, 13, 5).Move(-3, 0),   // 4:right
				S(5, 10, 5, 13).Move(0, -3),   // 5:top
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("Inside", func(t *testing.T) {
			segs := []Segment{
				S(0, 7.5, -3, 7.5).Move(4, 0), // 0:left
				S(2.5, 5, 2.5, 2).Move(0, 4),  // 1:bottom-left-1
				S(5, 2.5, 2, 2.5).Move(4, 0),  // 2:bottom-left-2
				S(7.5, 0, 7.5, -3).Move(0, 4), // 3:bottom
				S(10, 5, 13, 5).Move(-4, 0),   // 4:right
				S(5, 10, 5, 13).Move(0, -4),   // 5:top
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxContainsSegment(ring, seg, true))
						expect(t, ringxContainsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxContainsSegment(ring, seg, true))
						expect(t, ringxContainsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("InsideEdgeToEdge", func(t *testing.T) {
			segs := []Segment{
				S(0, 7.5, 10, 7.5), // 0:left
				S(2.5, 5, 2.5, 10), // 1:bottom-left-1
				S(5, 2.5, 10, 2.5), // 2:bottom-left-2
				S(7.5, 0, 7.5, 10), // 3:bottom
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
					})
				})
			}
		})
		t.Run("PassoverEdgeToEdge", func(t *testing.T) {
			segs := []Segment{
				S(-1, 7.5, 11, 7.5), // 0:left
				S(2.5, 4, 2.5, 11),  // 1:bottom-left-1
				S(4, 2.5, 11, 2.5),  // 2:bottom-left-2
				S(7.5, -1, 7.5, 11), // 3:bottom
			}
			for i, seg := range segs {
				t.Run(fmt.Sprintf("Segment-%d", i), func(t *testing.T) {
					testDualRingX(t, shape, func(t *testing.T, ring RingX) {
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
						seg.A, seg.B = seg.B, seg.A
						expect(t, !ringxContainsSegment(ring, seg, true))
						expect(t, !ringxContainsSegment(ring, seg, false))
					})
				})
			}
		})
	})

	t.Run("Cases", func(t *testing.T) {
		ring := newRingX([]Point{
			P(0, 0), P(4, 0), P(4, 4), P(3, 4),
			P(2, 3), P(1, 4), P(0, 4), P(0, 0),
		})
		t.Run("1", func(t *testing.T) {
			expect(t, !ringxContainsSegment(ring, S(1.5, 3.5, 2.5, 3.5), true))
			expect(t, !ringxContainsSegment(ring, S(1.5, 3.5, 2.5, 3.5), false))
		})
		t.Run("2", func(t *testing.T) {
			expect(t, !ringxContainsSegment(ring, S(1.0, 3.5, 2.5, 3.5), true))
			expect(t, !ringxContainsSegment(ring, S(1.0, 3.5, 2.5, 3.5), false))
		})
		t.Run("3", func(t *testing.T) {
			expect(t, ringxContainsSegment(ring, S(1.25, 3.75, 1.75, 3.25), true))
			expect(t, !ringxContainsSegment(ring, S(1.25, 3.75, 1.75, 3.25), false))
		})
		t.Run("4", func(t *testing.T) {
			expect(t, !ringxContainsSegment(ring, S(1.5, 3.5, 3, 3.5), true))
			expect(t, !ringxContainsSegment(ring, S(1.5, 3.5, 3, 3.5), false))
		})
		t.Run("5", func(t *testing.T) {
			expect(t, !ringxContainsSegment(ring, S(1.0, 3.5, 3, 3.5), true))
			expect(t, !ringxContainsSegment(ring, S(1.0, 3.5, 3, 3.5), false))
		})
		t.Run("6", func(t *testing.T) {
			ring := newRingX([]Point{
				P(0, 0), P(4, 0), P(4, 4), P(3, 4),
				P(2, 5), P(1, 4), P(0, 4), P(0, 0),
			})
			expect(t, ringxContainsSegment(ring, S(1.5, 4.5, 2.5, 4.5), true))
			expect(t, !ringxContainsSegment(ring, S(1.5, 4.5, 2.5, 4.5), false))
		})
		t.Run("7", func(t *testing.T) {
			ring := newRingX([]Point{
				P(0, 0), P(4, 0), P(4, 4), P(3, 4),
				P(2.5, 3), P(2, 4), P(1.5, 3),
				P(1, 4), P(0, 4), P(0, 0),
			})
			expect(t, !ringxContainsSegment(ring, S(1.25, 3.5, 2.75, 3.5), true))
			expect(t, !ringxContainsSegment(ring, S(1.25, 3.5, 2.75, 3.5), false))
		})
		t.Run("8", func(t *testing.T) {
			ring := newRingX([]Point{
				P(0, 0), P(4, 0), P(4, 4), P(3, 4),
				P(2.5, 5), P(2, 4), P(1.5, 5),
				P(1, 4), P(0, 4), P(0, 0),
			})
			expect(t, !ringxContainsSegment(ring, S(1.25, 4.5, 2.75, 4.5), true))
			expect(t, !ringxContainsSegment(ring, S(1.25, 4.5, 2.75, 4.5), false))
		})
		t.Run("9", func(t *testing.T) {
			expect(t, ringxContainsSegment(ring, S(1, 2, 3, 2), true))
			expect(t, ringxContainsSegment(ring, S(1, 2, 3, 2), false))
		})
		t.Run("10", func(t *testing.T) {
			expect(t, ringxContainsSegment(ring, S(1, 3, 3, 2), true))
			expect(t, ringxContainsSegment(ring, S(1, 3, 3, 2), false))
		})
		t.Run("11", func(t *testing.T) {
			expect(t, ringxContainsSegment(ring, S(1, 2, 3, 3), true))
			expect(t, ringxContainsSegment(ring, S(1, 2, 3, 3), false))
		})
		t.Run("12", func(t *testing.T) {
			expect(t, ringxContainsSegment(ring, S(1.5, 3.5, 1.5, 2), true))
			expect(t, !ringxContainsSegment(ring, S(1.5, 3.5, 1.5, 2), false))
		})
		t.Run("13", func(t *testing.T) {
			ring := newRingX([]Point{
				P(0, 0), P(2, 0), P(4, 0), P(4, 4), P(3, 4),
				P(2, 3), P(1, 4), P(0, 4), P(0, 0),
			})
			expect(t, ringxContainsSegment(ring, S(1, 0, 3, 0), true))
			expect(t, !ringxContainsSegment(ring, S(1, 0, 3, 0), false))
		})
		t.Run("14", func(t *testing.T) {
			ring := newRingX([]Point{
				P(0, 0), P(4, 0), P(2, 2), P(0, 4), P(0, 0),
			})
			expect(t, ringxContainsSegment(ring, S(1, 3, 3, 1), true))
			expect(t, ringxContainsSegment(ring, S(3, 1, 1, 3), true))
			expect(t, !ringxContainsSegment(ring, S(1, 3, 3, 1), false))
			expect(t, !ringxContainsSegment(ring, S(3, 1, 1, 3), false))
		})
		t.Run("15", func(t *testing.T) {
			expect(t, ringxContainsSegment(ring, S(1, 3, 3, 3), true))
			expect(t, !ringxContainsSegment(ring, S(1, 3, 3, 3), false))
		})
		t.Run("16", func(t *testing.T) {
			expect(t, ringxContainsSegment(ring, S(1, 3, 2, 3), true))
			expect(t, !ringxContainsSegment(ring, S(1, 3, 2, 3), false))
		})
	})
}

func TestRingXContainsRing(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		expect(t, !ringxContainsRing(newRingX(nil), R(0, 0, 1, 1), true))
		expect(t, !ringxContainsRing(R(0, 0, 1, 1), newRingX(nil), true))
	})
	t.Run("Cases", func(t *testing.T) {
		// concave
		ring := newRingX([]Point{
			P(0, 0), P(4, 0), P(4, 4), P(3, 4),
			P(2, 3), P(1, 4), P(0, 4), P(0, 0),
		})
		t.Run("1", func(t *testing.T) {
			expect(t, ringxContainsRing(ring, R(1, 1, 3, 2), true))
			expect(t, ringxContainsRing(ring, R(1, 1, 3, 2), false))
		})
		t.Run("2", func(t *testing.T) {
			expect(t, ringxContainsRing(ring, R(0, 0, 2, 1), true))
			expect(t, !ringxContainsRing(ring, R(0, 0, 2, 1), false))
		})
		t.Run("3", func(t *testing.T) {
			expect(t, !ringxContainsRing(ring, R(-1.5, 1, 1.5, 2), true))
			expect(t, !ringxContainsRing(ring, R(-1.5, 1, 1.5, 2), false))
		})
		t.Run("4", func(t *testing.T) {
			expect(t, !ringxContainsRing(ring, R(1, 2.5, 3, 3.5), true))
			expect(t, !ringxContainsRing(ring, R(1, 2.5, 3, 3.5), false))
		})
		t.Run("5", func(t *testing.T) {
			expect(t, ringxContainsRing(ring, R(1, 2, 3, 3), true))
			expect(t, !ringxContainsRing(ring, R(1, 2, 3, 3), false))
		})
		// convex
		ring = newRingX([]Point{
			P(0, 0), P(4, 0), P(4, 4),
			P(3, 4), P(2, 5), P(1, 4),
			P(0, 4), P(0, 0),
		})
		t.Run("6", func(t *testing.T) {
			expect(t, ringxContainsRing(ring, R(1, 2, 3, 3), true))
			expect(t, ringxContainsRing(ring, R(1, 2, 3, 3), false))
		})
		t.Run("7", func(t *testing.T) {
			expect(t, ringxContainsRing(ring, R(1, 3, 3, 4), true))
			expect(t, !ringxContainsRing(ring, R(1, 3, 3, 4), false))
		})
		t.Run("8", func(t *testing.T) {
			expect(t, !ringxContainsRing(ring, R(1, 3.5, 3, 4.5), true))
			expect(t, !ringxContainsRing(ring, R(1, 3.5, 3, 4.5), false))
		})
		t.Run("9", func(t *testing.T) {
			ring = newRingX([]Point{
				P(0, 0), P(2, 0), P(4, 0), P(4, 4), P(3, 4),
				P(2, 5), P(1, 4), P(0, 4), P(0, 0),
			})
			expect(t, ringxContainsRing(ring, R(1, 0, 3, 1), true))
			expect(t, !ringxContainsRing(ring, R(1, 0, 3, 1), false))
		})
		t.Run("10", func(t *testing.T) {
			ring = newRingX([]Point{
				P(0, 0), P(4, 0), P(4, 3), P(2, 4),
				P(0, 3), P(0, 0),
			})
			expect(t, ringxContainsRing(ring, R(1, 1, 3, 2), true))
			expect(t, ringxContainsRing(ring, R(1, 1, 3, 2), false))
		})
		ring = newRingX([]Point{
			P(0, 0), P(4, 0), P(4, 3), P(3, 4), P(1, 4), P(0, 3), P(0, 0),
		})
		t.Run("11", func(t *testing.T) {
			expect(t, ringxContainsRing(ring, R(1, 1, 3, 2), true))
			expect(t, ringxContainsRing(ring, R(1, 1, 3, 2), false))
		})
		t.Run("12", func(t *testing.T) {
			expect(t, ringxContainsRing(ring, R(0, 1, 2, 2), true))
			expect(t, !ringxContainsRing(ring, R(0, 1, 2, 2), false))
		})
		t.Run("13", func(t *testing.T) {
			expect(t, !ringxContainsRing(ring, R(-1, 1, 1, 2), true))
			expect(t, !ringxContainsRing(ring, R(-1, 1, 1, 2), false))
		})
		t.Run("14", func(t *testing.T) {
			expect(t, ringxContainsRing(ring, R(0.5, 2.5, 2.5, 3.5), true))
			expect(t, !ringxContainsRing(ring, R(0.5, 2.5, 2.5, 3.5), false))
		})
		t.Run("15", func(t *testing.T) {
			expect(t, !ringxContainsRing(ring, R(0.25, 2.75, 2.25, 3.75), true))
			expect(t, !ringxContainsRing(ring, R(0.25, 2.75, 2.25, 3.75), false))
		})
		t.Run("16", func(t *testing.T) {
			expect(t, !ringxContainsRing(ring, R(-2, 1, -1, 2), true))
			expect(t, !ringxContainsRing(ring, R(-2, 1, -1, 2), false))
		})
		t.Run("17", func(t *testing.T) {
			expect(t, !ringxContainsRing(ring, R(-0.5, 3.5, 0.5, 4.5), true))
			expect(t, !ringxContainsRing(ring, R(-0.5, 3.5, 0.5, 4.5), false))
		})
		t.Run("18", func(t *testing.T) {
			expect(t, !ringxContainsRing(ring, R(-0.75, 3.75, 0.25, 4.75), true))
			expect(t, !ringxContainsRing(ring, R(-0.74, 3.75, 0.25, 4.75), false))
		})
		t.Run("19", func(t *testing.T) {
			expect(t, !ringxContainsRing(ring, R(-1, -1, 5, 5), true))
			expect(t, !ringxContainsRing(ring, R(-1, -1, 5, 5), false))
		})

	})
}

func TestRingXIntersectsRing(t *testing.T) {
	intersectsBothWays := func(ringA, ringB RingX, allowOnEdge bool) bool {
		t1 := ringxIntersectsRing(ringA, ringB, allowOnEdge)
		t2 := ringxIntersectsRing(ringB, ringA, allowOnEdge)
		if t1 != t2 {
			panic("mismatch")
		}
		return t1
	}
	t.Run("Empty", func(t *testing.T) {
		expect(t, !intersectsBothWays(newRingX(nil), R(0, 0, 1, 1), true))
		expect(t, !intersectsBothWays(R(0, 0, 1, 1), newRingX(nil), true))
	})
	t.Run("Cases", func(t *testing.T) {
		// concave
		ring := newRingX([]Point{
			P(0, 0), P(4, 0), P(4, 4), P(3, 4),
			P(2, 3), P(1, 4), P(0, 4), P(0, 0),
		})
		t.Run("1", func(t *testing.T) {
			expect(t, intersectsBothWays(ring, R(1, 1, 3, 2), true))
			expect(t, intersectsBothWays(ring, R(1, 1, 3, 2), false))
		})
		t.Run("2", func(t *testing.T) {
			expect(t, intersectsBothWays(ring, R(0, 0, 2, 1), true))
			expect(t, intersectsBothWays(ring, R(0, 0, 2, 1), false))
		})
		t.Run("3", func(t *testing.T) {
			expect(t, intersectsBothWays(ring, R(-1.5, 1, 1.5, 2), true))
			expect(t, intersectsBothWays(ring, R(-1.5, 1, 1.5, 2), false))
		})
		t.Run("4", func(t *testing.T) {
			expect(t, intersectsBothWays(ring, R(1, 2.5, 3, 3.5), true))
			expect(t, intersectsBothWays(ring, R(1, 2.5, 3, 3.5), false))
		})
		t.Run("5", func(t *testing.T) {
			expect(t, intersectsBothWays(ring, R(1, 2, 3, 3), true))
			expect(t, intersectsBothWays(ring, R(1, 2, 3, 3), false))
		})
		// convex
		ring = newRingX([]Point{
			P(0, 0), P(4, 0), P(4, 4),
			P(3, 4), P(2, 5), P(1, 4),
			P(0, 4), P(0, 0),
		})
		t.Run("6", func(t *testing.T) {
			expect(t, intersectsBothWays(ring, R(1, 2, 3, 3), true))
			expect(t, intersectsBothWays(ring, R(1, 2, 3, 3), false))
		})
		t.Run("7", func(t *testing.T) {
			expect(t, intersectsBothWays(ring, R(1, 3, 3, 4), true))
			expect(t, intersectsBothWays(ring, R(1, 3, 3, 4), false))
		})
		t.Run("8", func(t *testing.T) {
			expect(t, intersectsBothWays(ring, R(1, 3.5, 3, 4.5), true))
			expect(t, intersectsBothWays(ring, R(1, 3.5, 3, 4.5), false))
		})
		t.Run("9", func(t *testing.T) {
			ring = newRingX([]Point{
				P(0, 0), P(2, 0), P(4, 0), P(4, 4), P(3, 4),
				P(2, 5), P(1, 4), P(0, 4), P(0, 0),
			})
			expect(t, intersectsBothWays(ring, R(1, 0, 3, 1), true))
			expect(t, intersectsBothWays(ring, R(1, 0, 3, 1), false))
		})
		t.Run("10", func(t *testing.T) {
			ring = newRingX([]Point{
				P(0, 0), P(4, 0), P(4, 3), P(2, 4),
				P(0, 3), P(0, 0),
			})
			expect(t, intersectsBothWays(ring, R(1, 1, 3, 2), true))
			expect(t, intersectsBothWays(ring, R(1, 1, 3, 2), false))
		})
		ring = newRingX([]Point{
			P(0, 0), P(4, 0), P(4, 3), P(3, 4), P(1, 4), P(0, 3), P(0, 0),
		})
		t.Run("11", func(t *testing.T) {
			expect(t, intersectsBothWays(ring, R(1, 1, 3, 2), true))
			expect(t, intersectsBothWays(ring, R(1, 1, 3, 2), false))
		})
		t.Run("12", func(t *testing.T) {
			expect(t, intersectsBothWays(ring, R(0, 1, 2, 2), true))
			expect(t, intersectsBothWays(ring, R(0, 1, 2, 2), false))
		})
		t.Run("13", func(t *testing.T) {
			expect(t, intersectsBothWays(ring, R(-1, 1, 1, 2), true))
			expect(t, intersectsBothWays(ring, R(-1, 1, 1, 2), false))
		})
		t.Run("14", func(t *testing.T) {
			expect(t, intersectsBothWays(ring, R(0.5, 2.5, 2.5, 3.5), true))
			expect(t, intersectsBothWays(ring, R(0.5, 2.5, 2.5, 3.5), false))
		})
		t.Run("15", func(t *testing.T) {
			expect(t, intersectsBothWays(ring, R(0.25, 2.75, 2.25, 3.75), true))
			expect(t, intersectsBothWays(ring, R(0.25, 2.75, 2.25, 3.75), false))
		})
		t.Run("16", func(t *testing.T) {
			expect(t, !intersectsBothWays(ring, R(-2, 1, -1, 2), true))
			expect(t, !intersectsBothWays(ring, R(-2, 1, -1, 2), false))
		})
		t.Run("17", func(t *testing.T) {
			expect(t, intersectsBothWays(ring, R(-0.5, 3.5, 0.5, 4.5), true))
			expect(t, !intersectsBothWays(ring, R(-0.5, 3.5, 0.5, 4.5), false))
		})
		t.Run("18", func(t *testing.T) {
			expect(t, !intersectsBothWays(ring, R(-0.75, 3.75, 0.25, 4.75), true))
			expect(t, !intersectsBothWays(ring, R(-0.74, 3.75, 0.25, 4.75), false))
		})
		t.Run("19", func(t *testing.T) {
			expect(t, ringxIntersectsRing(ring, R(-1, -1, 5, 5), true))
			expect(t, ringxIntersectsRing(ring, R(-1, -1, 5, 5), false))
		})
	})
}

func TestRingXContainsLine(t *testing.T) {
	t.Run("Cases", func(t *testing.T) {
		// convex
		ring := newRingX([]Point{
			P(0, 0), P(4, 0), P(4, 3),
			P(3, 4), P(1, 4),
			P(0, 3), P(0, 0),
		})
		makeLine := func(start Point) *Line {
			return NewLine([]Point{
				start, start.Move(0, 1), start.Move(1, 1),
				start.Move(1, 2), start.Move(2, 2),
			})
		}
		t.Run("1", func(t *testing.T) {
			expect(t, ringxContainsLine(ring, makeLine(P(1, 1)), true))
			expect(t, ringxContainsLine(ring, makeLine(P(1, 1)), false))
		})
		t.Run("2", func(t *testing.T) {
			expect(t, ringxContainsLine(ring, makeLine(P(0, 1)), true))
			expect(t, !ringxContainsLine(ring, makeLine(P(0, 1)), false))
		})
		t.Run("3", func(t *testing.T) {
			expect(t, !ringxContainsLine(ring, makeLine(P(-0.5, 1)), true))
			expect(t, !ringxContainsLine(ring, makeLine(P(-0.5, 1)), false))
		})
		t.Run("4", func(t *testing.T) {
			expect(t, ringxContainsLine(ring, makeLine(P(0, 2)), true))
			expect(t, !ringxContainsLine(ring, makeLine(P(0, 2)), false))
		})
		t.Run("5", func(t *testing.T) {
			expect(t, !ringxContainsLine(ring, makeLine(P(0, 2.5)), true))
			expect(t, !ringxContainsLine(ring, makeLine(P(0, 2.5)), false))
		})
		t.Run("6", func(t *testing.T) {
			expect(t, !ringxContainsLine(ring, makeLine(P(2, 2)), true))
			expect(t, !ringxContainsLine(ring, makeLine(P(2, 2)), false))
		})
		t.Run("7", func(t *testing.T) {
			expect(t, !ringxContainsLine(ring, makeLine(P(2, 1.5)), true))
			//expect(t, !ringxContainsLine(ring, makeLine(P(2, 1.5)), false))
		})
		t.Run("8", func(t *testing.T) {
			expect(t, ringxContainsLine(ring, makeLine(P(2, 1)), true))
			expect(t, !ringxContainsLine(ring, makeLine(P(2, 1)), false))
		})
		t.Run("9", func(t *testing.T) {
			expect(t, ringxContainsLine(ring, makeLine(P(2, 0)), true))
			expect(t, !ringxContainsLine(ring, makeLine(P(2, 0)), false))
		})
		t.Run("10", func(t *testing.T) {
			expect(t, ringxContainsLine(ring, makeLine(P(1.5, 0)), true))
			expect(t, !ringxContainsLine(ring, makeLine(P(1.5, 0)), false))
		})
		t.Run("11", func(t *testing.T) {
			expect(t, !ringxContainsLine(ring, makeLine(P(1, -1)), true))
			expect(t, !ringxContainsLine(ring, makeLine(P(1, -1)), false))
		})
		t.Run("12", func(t *testing.T) {
			expect(t, !ringxContainsLine(ring, makeLine(P(1.5, -0.5)), true))
			expect(t, !ringxContainsLine(ring, makeLine(P(1.5, -0.5)), false))
		})
		t.Run("13", func(t *testing.T) {
			expect(t, ringxContainsLine(ring, makeLine(P(0, 0)), true))
			expect(t, !ringxContainsLine(ring, makeLine(P(0, 0)), false))
		})
		t.Run("14", func(t *testing.T) {
			expect(t, !ringxContainsLine(ring, makeLine(P(3, 1)), true))
			expect(t, !ringxContainsLine(ring, makeLine(P(3, 1)), false))
		})
		// t.Run("15", func(t *testing.T) {
		// 	line := L(P(-1,-1),P(-1,2),P(2,2),P(2,
		// 	expect(t, !ringxContainsLine(ring, line, true))
		// 	expect(t, !ringxContainsLine(ring, line, false))
		// })
	})
}
