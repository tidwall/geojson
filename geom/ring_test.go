package geom

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/tidwall/lotsa"
)

func TestRingScan(t *testing.T) {
	test := func(t *testing.T, index bool) {
		rectangleRing := newRing(rectangle)
		if !index {
			rectangleRing.(*baseSeries).tree = nil
		} else {
			rectangleRing.(*baseSeries).buildTree()
		}
		var segs []Segment
		seriesForEachSegment(rectangleRing, func(seg Segment) bool {
			segs = append(segs, seg)
			return true
		})
		segsExpect := []Segment{
			S(0, 0, 10, 0),
			S(10, 0, 10, 10),
			S(10, 10, 0, 10),
			S(0, 10, 0, 0),
		}
		expect(t, len(segs) == len(segsExpect))
		for i := 0; i < len(segs); i++ {
			expect(t, segs[i] == segsExpect[i])
		}

		segs = nil
		notClosedRing := newRing(rectangle)
		if !index {
			notClosedRing.(*baseSeries).tree = nil
		} else {
			notClosedRing.(*baseSeries).buildTree()
		}
		seriesForEachSegment(notClosedRing, func(seg Segment) bool {
			segs = append(segs, seg)
			return true
		})
		expect(t, len(segs) == len(segsExpect))
		for i := 0; i < len(segs); i++ {
			expect(t, segs[i] == segsExpect[i])
		}
	}
	t.Run("Indexed", func(t *testing.T) {
		test(t, true)
	})
	t.Run("Simple", func(t *testing.T) {
		test(t, false)
	})
}

func TestRingSearch(t *testing.T) {
	test := func(t *testing.T, index bool) {
		octagonRing := newRing(octagon)
		if !index {
			octagonRing.(*baseSeries).tree = nil
		} else {
			octagonRing.(*baseSeries).buildTree()
		}
		var segs []Segment
		octagonRing.Search(R(0, 0, 0, 0), func(seg Segment, _ int) bool {
			segs = append(segs, seg)
			return true
		})
		segsExpect := []Segment{
			S(0, 3, 3, 0),
		}
		expect(t, checkSegsDups(segsExpect, segs))
		segs = nil
		octagonRing.Search(R(0, 0, 0, 10), func(seg Segment, _ int) bool {
			segs = append(segs, seg)
			return true
		})
		segsExpect = []Segment{
			S(3, 10, 0, 7),
			S(0, 7, 0, 3),
			S(0, 3, 3, 0),
		}
		expect(t, checkSegsDups(segsExpect, segs))
		segs = nil
		octagonRing.Search(R(0, 0, 5, 10), func(seg Segment, _ int) bool {
			segs = append(segs, seg)
			return true
		})
		segsExpect = []Segment{
			S(3, 0, 7, 0),
			S(7, 10, 3, 10),
			S(3, 10, 0, 7),
			S(0, 7, 0, 3),
			S(0, 3, 3, 0),
		}
		expect(t, checkSegsDups(segsExpect, segs))
	}
	t.Run("Indexed", func(t *testing.T) {
		test(t, true)
	})
	t.Run("Simple", func(t *testing.T) {
		test(t, false)
	})
}

func TestRingIntersectsSegment(t *testing.T) {
	simple := newRing(concave1)
	simple.(*baseSeries).tree = nil
	tree := newRing(concave1)
	tree.(*baseSeries).buildTree()

	expect(t, !ringIntersectsSegment(simple, S(0, 0, 3, 3), true))
	expect(t, !ringIntersectsSegment(tree, S(0, 0, 3, 3), true))
	expect(t, !ringIntersectsSegment(simple, S(0, 0, 3, 3), false))
	expect(t, !ringIntersectsSegment(tree, S(0, 0, 3, 3), false))

	expect(t, ringIntersectsSegment(simple, S(0, 0, 5, 5), true))
	expect(t, ringIntersectsSegment(tree, S(0, 0, 5, 5), true))
	expect(t, !ringIntersectsSegment(simple, S(0, 0, 5, 5), false))
	expect(t, !ringIntersectsSegment(tree, S(0, 0, 5, 5), false))

	expect(t, ringIntersectsSegment(simple, S(0, 0, 10, 10), true))
	expect(t, ringIntersectsSegment(tree, S(0, 0, 10, 10), true))
	expect(t, !ringIntersectsSegment(simple, S(0, 0, 10, 10), false))
	expect(t, !ringIntersectsSegment(tree, S(0, 0, 10, 10), false))

	expect(t, ringIntersectsSegment(simple, S(0, 0, 11, 11), true))
	expect(t, ringIntersectsSegment(tree, S(0, 0, 11, 11), true))
	expect(t, !ringIntersectsSegment(simple, S(0, 0, 11, 11), false))
	expect(t, !ringIntersectsSegment(tree, S(0, 0, 11, 11), false))

}

func TestRingIntersectsRing(t *testing.T) {
	small1 := newRing([]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}})
	small2 := newRing([]Point{{5, 4}, {7, 4}, {7, 6}, {5, 6}, {5, 4}})
	expect(t, ringIntersectsRing(small1, small2, true))
	expect(t, ringIntersectsRing(small1, small2, false))

	simple := newRing(concave1)
	simple.(*baseSeries).tree = nil
	tree := newRing(concave1)
	tree.(*baseSeries).buildTree()
	small := newRing([]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}})
	small.(*baseSeries).tree = nil

	intersects := func(ring Ring) bool {
		tt := ringIntersectsRing(simple, ring, true)
		if ringIntersectsRing(tree, ring, true) != tt {
			panic("structure mismatch")
		}
		return tt
	}

	intersectsOnEdgeNotAllowed := func(ring Ring) bool {
		tt := ringIntersectsRing(simple, ring, false)
		if ringIntersectsRing(tree, ring, false) != tt {
			panic("structure mismatch")
		}
		return tt
	}

	expect(t, intersects(small))
	expect(t, intersects(small.(*baseSeries).Move(-6, 0)))
	expect(t, intersects(small.(*baseSeries).Move(6, 0)))
	expect(t, !intersects(small.(*baseSeries).Move(-7, 0)))
	expect(t, !intersects(small.(*baseSeries).Move(7, 0)))
	expect(t, intersects(small.(*baseSeries).Move(1, 1)))
	expect(t, intersects(small.(*baseSeries).Move(-1, -1)))
	expect(t, intersects(small.(*baseSeries).Move(2, 2)))
	expect(t, !intersects(small.(*baseSeries).Move(-2, -2)))
	expect(t, intersects(small.(*baseSeries).Move(0, -6)))
	expect(t, intersects(small.(*baseSeries).Move(0, 6)))
	expect(t, !intersects(small.(*baseSeries).Move(0, -7)))
	expect(t, !intersects(small.(*baseSeries).Move(0, 7)))

	expect(t, intersectsOnEdgeNotAllowed(small.(*baseSeries).Move(-5, 0)))
	expect(t, intersectsOnEdgeNotAllowed(small.(*baseSeries).Move(5, 0)))
	expect(t, intersectsOnEdgeNotAllowed(small.(*baseSeries).Move(0, -5)))
	expect(t, intersectsOnEdgeNotAllowed(small.(*baseSeries).Move(0, 5)))

	expect(t, !intersectsOnEdgeNotAllowed(small.(*baseSeries).Move(-6, 0)))
	expect(t, !intersectsOnEdgeNotAllowed(small.(*baseSeries).Move(6, 0)))
	expect(t, !intersectsOnEdgeNotAllowed(small.(*baseSeries).Move(0, -6)))
	expect(t, !intersectsOnEdgeNotAllowed(small.(*baseSeries).Move(0, 6)))

	expect(t, intersectsOnEdgeNotAllowed(small.(*baseSeries).Move(1, 1)))
	expect(t, !intersectsOnEdgeNotAllowed(small.(*baseSeries).Move(-1, -1)))

}

func TestBigRandomPIP(t *testing.T) {
	simple := newRing(az)
	simple.(*baseSeries).tree = nil
	tree := newRing(az)
	tree.(*baseSeries).buildTree()
	expect(t, simple.Rect() == tree.Rect())
	rect := tree.Rect()
	start := time.Now()
	for time.Since(start) < time.Second/4 {
		point := P(
			rand.Float64()*(rect.Max.X-rect.Min.X)+rect.Min.X,
			rand.Float64()*(rect.Max.Y-rect.Min.Y)+rect.Min.Y,
		)
		expect(t, ringContainsPoint(tree, point, true) ==
			ringContainsPoint(simple, point, true))
	}
}

func testBig(
	t *testing.T, label string, points []Point, pointIn, pointOut Point,
) {
	N, T := 100000, 4

	simple := newRing(points)
	simple.(*baseSeries).tree = nil
	tree := newRing(points)
	tree.(*baseSeries).buildTree()
	pointOn := points[len(points)/2]

	// ioutil.WriteFile(label+".svg", []byte(tools.SVG(tree.(*baseSeries).tree)), 0666)

	expect(t, ringContainsPoint(simple, pointIn, true))
	expect(t, ringContainsPoint(tree, pointIn, true))

	expect(t, ringContainsPoint(simple, pointOn, true))
	expect(t, ringContainsPoint(tree, pointOn, true))

	expect(t, !ringContainsPoint(simple, pointOn, false))
	expect(t, !ringContainsPoint(tree, pointOn, false))

	expect(t, !ringContainsPoint(simple, pointOut, true))
	expect(t, !ringContainsPoint(tree, pointOut, true))
	if os.Getenv("PIPBENCH") == "1" {
		lotsa.Output = os.Stderr
		fmt.Printf(label + "/simp/in  ")
		lotsa.Ops(N, T, func(_, _ int) {
			ringContainsPoint(simple, pointIn, true)
		})
		fmt.Printf(label + "/tree/in  ")
		lotsa.Ops(N, T, func(_, _ int) {
			ringContainsPoint(tree, pointIn, true)
		})
		fmt.Printf(label + "/simp/on  ")
		lotsa.Ops(N, T, func(_, _ int) {
			ringContainsPoint(simple, pointOn, true)
		})
		fmt.Printf(label + "/tree/on  ")
		lotsa.Ops(N, T, func(_, _ int) {
			ringContainsPoint(tree, pointOn, true)
		})
		fmt.Printf(label + "/simp/out ")
		lotsa.Ops(N, T, func(_, _ int) {
			ringContainsPoint(simple, pointOut, true)
		})
		fmt.Printf(label + "/tree/out ")
		lotsa.Ops(N, T, func(_, _ int) {
			ringContainsPoint(tree, pointOut, true)
		})
	}
}

func TestBigArizona(t *testing.T) {
	testBig(t, "az", az, P(-112, 33), P(-114.477539062, 33.99802726))
}

func TestBigTexas(t *testing.T) {
	testBig(t, "tx", tx, P(-98.52539, 29.363027), P(-101.953125, 29.324720161))
}

func TestBigCircle(t *testing.T) {
	circle := seriesCopyPoints(CircleRing(P(-100.1, 31.2), 660000, 10000))
	if false {
		s := `{"type":"Polygon","coordinates":[[`
		for i, p := range circle {
			if i > 0 {
				s += ","
			}
			s += fmt.Sprintf("[%v,%v]", p.X, p.Y)
		}
		s += `]]}`
		println(s)
	}
	testBig(t, "circ", circle, P(-98.52, 29.363), P(-107.8857, 31.5410))
	circle = seriesCopyPoints(CircleRing(P(-100.1, 31.2), 660000, 2))
	expect(t, len(circle) == 4)
}

func TestRingContainsRing(t *testing.T) {
	simple := newRing(concave1)
	simple.(*baseSeries).tree = nil
	tree := newRing(concave1)
	tree.(*baseSeries).buildTree()

	expect(t, ringContainsRing(simple, simple, true))
	expect(t, ringContainsRing(simple, tree, true))
	expect(t, ringContainsRing(tree, simple, true))
	expect(t, ringContainsRing(tree, tree, true))

	expect(t, !ringContainsRing(simple, simple, false))
	expect(t, !ringContainsRing(simple, tree, false))
	expect(t, !ringContainsRing(tree, simple, false))
	expect(t, !ringContainsRing(tree, tree, false))

	small := newRing([]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}})
	small.(*baseSeries).tree = nil

	expect(t, !ringContainsRing(simple, small, true))
	expect(t, !ringContainsRing(tree, small, true))

	for x := 1.0; x <= 4; x++ {
		smallMoved := small.(*baseSeries).Move(x, 0)
		expect(t, ringContainsRing(simple, smallMoved, true))
		expect(t, ringContainsRing(tree, smallMoved, true))
	}
	expect(t, !ringContainsRing(simple, small.(*baseSeries).Move(4, 0), false))
	expect(t, !ringContainsRing(tree, small.(*baseSeries).Move(4, 0), false))
	for y := 1.0; y <= 4; y++ {
		expect(t, ringContainsRing(simple, small.(*baseSeries).Move(0, y), true))
		expect(t, ringContainsRing(tree, small.(*baseSeries).Move(0, y), true))
	}
	expect(t, !ringContainsRing(simple, small.(*baseSeries).Move(0, 4), false))
	expect(t, !ringContainsRing(tree, small.(*baseSeries).Move(0, 4), false))

	for x := -1.0; x >= -4; x-- {
		expect(t, !ringContainsRing(simple, small.(*baseSeries).Move(x, 0), true))
		expect(t, !ringContainsRing(tree, small.(*baseSeries).Move(x, 0), true))
	}
	expect(t, !ringContainsRing(simple, small.(*baseSeries).Move(-4, 0), false))
	expect(t, !ringContainsRing(tree, small.(*baseSeries).Move(-4, 0), false))
	for y := -1.0; y >= -4; y-- {
		expect(t, !ringContainsRing(simple, small.(*baseSeries).Move(0, y), true))
		expect(t, !ringContainsRing(tree, small.(*baseSeries).Move(0, y), true))
	}
	expect(t, !ringContainsRing(simple, small.(*baseSeries).Move(0, -4), false))
	expect(t, !ringContainsRing(tree, small.(*baseSeries).Move(0, -4), false))

	expect(t, !ringContainsRing(simple, small.(*baseSeries).Move(1, 0), false))
	expect(t, !ringContainsRing(tree, small.(*baseSeries).Move(1, 0), false))
	expect(t, ringContainsRing(simple, small.(*baseSeries).Move(2, 0), false))
	expect(t, ringContainsRing(tree, small.(*baseSeries).Move(2, 0), false))
	expect(t, ringContainsRing(simple, small.(*baseSeries).Move(2, 2), false))
	expect(t, ringContainsRing(tree, small.(*baseSeries).Move(2, 2), false))
	expect(t, !ringContainsRing(simple, small.(*baseSeries).Move(-2, -2), false))
	expect(t, !ringContainsRing(tree, small.(*baseSeries).Move(-2, -2), false))

	expect(t, !ringContainsRing(simple, small.(*baseSeries).Move(5, 0), true))
	expect(t, !ringContainsRing(tree, small.(*baseSeries).Move(5, 0), true))
	expect(t, !ringContainsRing(simple, small.(*baseSeries).Move(-5, 0), true))
	expect(t, !ringContainsRing(tree, small.(*baseSeries).Move(-5, 0), true))

	expect(t, !ringContainsRing(simple, small.(*baseSeries).Move(0, 5), true))
	expect(t, !ringContainsRing(tree, small.(*baseSeries).Move(0, 5), true))
	expect(t, !ringContainsRing(simple, small.(*baseSeries).Move(0, -5), true))
	expect(t, !ringContainsRing(tree, small.(*baseSeries).Move(0, -5), true))

}
func TestBowtie(t *testing.T) {
	simple := newRing(bowtie)
	simple.(*baseSeries).tree = nil
	tree := newRing(bowtie)
	tree.(*baseSeries).buildTree()
	square := newRing([]Point{P(3, 3), P(7, 3), P(7, 7), P(3, 7), P(3, 3)})
	square.(*baseSeries).tree = nil

	expect(t, ringIntersectsRing(simple, square, true))
	expect(t, ringIntersectsRing(tree, square, true))
	expect(t, !ringContainsRing(simple, square, true))
	expect(t, !ringContainsRing(tree, square, true))

}

func TestRingVarious(t *testing.T) {
	ring := newRing(octagon[:len(octagon)-1])
	ring.(*baseSeries).buildTree()
	n := 0
	ring.Search(R(0, 0, 10, 10), func(seg Segment, index int) bool {
		n++
		return true
	})
	expect(t, n == 8)
	n = 0
	seriesForEachSegment(ring, func(seg Segment) bool {
		n++
		return true
	})
	expect(t, n == 8)
	n = 0
	seriesForEachSegment(ring, func(seg Segment) bool {
		n++
		return false
	})
	expect(t, n == 1)
	expect(t, ringIntersectsSegment(ring, S(0, 0, 4, 4), true))
	expect(t, !newRingSimple2([]Point{}).Convex())
	expect(t, newRingSimple2(octagon).Convex())
	expect(t, !newRingIndexed2([]Point{}).Convex())
	expect(t, newRingIndexed2(octagon).Convex())

	ring = newRing(octagon[:len(octagon)-1])
	ring.(*baseSeries).tree = nil
	n = 0
	ring.Search(R(0, 0, 10, 10), func(seg Segment, index int) bool {
		n++
		return false
	})
	expect(t, n == 1)
	n = 0
	seriesForEachSegment(ring, func(seg Segment) bool {
		n++
		return true
	})
	expect(t, n == 8)
	expect(t, ringIntersectsSegment(ring, S(0, 0, 4, 4), true))

	small := newRing([]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}})
	small.(*baseSeries).tree = nil

	expect(t, ringIntersectsRing(small, ring, true))
	expect(t, ringIntersectsRing(ring, small, true))

	expect(t, S(0, 0, 0, 0).Raycast(P(0, 0)).On)

	ring1 := newRing(octagon)
	ring1.(*baseSeries).tree = nil
	n1 := 0
	seriesForEachSegment(ring1, func(seg Segment) bool {
		n1++
		return true
	})
	expect(t, ring1.(*baseSeries).Closed())
	ring2 := newRing(octagon[:len(octagon)-1])
	ring2.(*baseSeries).tree = nil
	n2 := 0
	seriesForEachSegment(ring2, func(seg Segment) bool {
		n2++
		return true
	})
	expect(t, n1 == n2)
	expect(t, ring2.(*baseSeries).Closed())

	ring1 = newRing(octagon)
	ring1.(*baseSeries).buildTree()
	n1 = 0
	seriesForEachSegment(ring1, func(seg Segment) bool {
		n1++
		return true
	})
	expect(t, ring1.(*baseSeries).Closed())
	ring2 = newRing(octagon[:len(octagon)-1])
	ring2.(*baseSeries).buildTree()
	n2 = 0
	seriesForEachSegment(ring2, func(seg Segment) bool {
		n2++
		return true
	})
	expect(t, n1 == n2)
	expect(t, ring2.(*baseSeries).Closed())

}

func newRingSimple2(points []Point) Ring {
	ring := newRing(points)
	ring.(*baseSeries).tree = nil
	return ring
}
func newRingIndexed2(points []Point) Ring {
	ring := newRing(points)
	ring.(*baseSeries).buildTree()
	return ring
}

func TestRingContainsPoint(t *testing.T) {
	expect(t, ringIntersectsPoint(newRingSimple2(octagon), P(4, 4), true))
	expect(t, ringIntersectsPoint(newRingIndexed2(octagon), P(4, 4), true))
}

func TestRingContainsSegment(t *testing.T) {
	expect(t, ringContainsSegment(newRingSimple2(octagon), S(4, 4, 6, 6), true))
	expect(t, ringContainsSegment(newRingIndexed2(octagon), S(4, 4, 6, 6), true))
	expect(t, !ringContainsSegment(newRingSimple2(octagon), S(9, 4, 11, 6), true))
	expect(t, !ringContainsSegment(newRingIndexed2(octagon), S(9, 4, 11, 6), true))
	expect(t, !ringContainsSegment(newRingSimple2(octagon), S(11, 4, 9, 6), true))
	expect(t, !ringContainsSegment(newRingIndexed2(octagon), S(11, 4, 9, 6), true))
	expect(t, !ringContainsSegment(newRingSimple2(concave1), S(11, 4, 9, 6), true))
	expect(t, !ringContainsSegment(newRingIndexed2(concave1), S(11, 4, 9, 6), true))
	expect(t, ringContainsSegment(newRingSimple2(concave1), S(6, 6, 8, 8), true))
	expect(t, ringContainsSegment(newRingIndexed2(concave1), S(6, 6, 8, 8), true))
	expect(t, !ringContainsSegment(newRingSimple2(concave1), S(1, 6, 6, 1), true))
	expect(t, !ringContainsSegment(newRingIndexed2(concave1), S(1, 6, 6, 1), true))
}
func TestRingContainsRect(t *testing.T) {
	expect(t, ringContainsRect(newRingSimple2(octagon), R(4, 4, 6, 6), true))
	expect(t, ringContainsRect(newRingIndexed2(octagon), R(4, 4, 6, 6), true))
	expect(t, ringContainsRect(newRingSimple2(octagon), R(4, 4, 6, 6), false))
	expect(t, ringContainsRect(newRingIndexed2(octagon), R(4, 4, 6, 6), false))
}
func TestRingIntersectsRect(t *testing.T) {
	expect(t, ringIntersectsRect(newRingSimple2(octagon), R(9, 4, 11, 6), true))
	expect(t, ringIntersectsRect(newRingIndexed2(octagon), R(9, 4, 11, 6), true))
	expect(t, !ringIntersectsRect(newRingSimple2(octagon), R(10, 4, 12, 6), false))
	expect(t, !ringIntersectsRect(newRingIndexed2(octagon), R(10, 4, 12, 6), false))
	expect(t, ringIntersectsRect(newRingSimple2(octagon), R(10, 4, 12, 6), true))
	expect(t, ringIntersectsRect(newRingIndexed2(octagon), R(10, 4, 12, 6), true))
	expect(t, !ringIntersectsRect(newRingSimple2(octagon), R(11, 4, 13, 6), true))
	expect(t, !ringIntersectsRect(newRingIndexed2(octagon), R(11, 4, 13, 6), true))
}
func TestRingContainsPoly(t *testing.T) {
	expect(t, ringContainsPoly(newRingSimple2(octagon), NewPoly(octagon, nil), true))
	expect(t, ringContainsPoly(newRingIndexed2(octagon), NewPoly(octagon, nil), true))
	expect(t, !ringContainsPoly(newRingSimple2(octagon), NewPoly(octagon, nil), false))
	expect(t, !ringContainsPoly(newRingIndexed2(octagon), NewPoly(octagon, nil), false))
}
func TestRingIntersectsPoly(t *testing.T) {
	expect(t, ringIntersectsPoly(newRingSimple2(octagon).(*baseSeries).Move(5, 0), NewPoly(octagon, nil), true))
	expect(t, ringIntersectsPoly(newRingIndexed2(octagon).(*baseSeries).Move(5, 0), NewPoly(octagon, nil), true))
	expect(t, ringIntersectsPoly(newRingSimple2(octagon).(*baseSeries).Move(10, 0), NewPoly(octagon, nil), true))
	expect(t, ringIntersectsPoly(newRingIndexed2(octagon).(*baseSeries).Move(10, 0), NewPoly(octagon, nil), true))
	expect(t, !ringIntersectsPoly(newRingSimple2(octagon).(*baseSeries).Move(10, 0), NewPoly(octagon, nil), false))
	expect(t, !ringIntersectsPoly(newRingIndexed2(octagon).(*baseSeries).Move(10, 0), NewPoly(octagon, nil), false))

	expect(t, !ringIntersectsPoly(newRingIndexed2(
		[]Point{P(4, 4), P(6, 4), P(6, 6), P(4, 6), P(4, 4)},
	), NewPoly(
		[]Point{P(0, 0), P(10, 0), P(10, 10), P(0, 10), P(0, 0)},
		[][]Point{{P(3, 3), P(7, 3), P(7, 7), P(3, 7), P(3, 3)}},
	), false))

}

func TestSegmentsIntersect(t *testing.T) {
	expect(t, !S(0, 0, 10, 10).IntersectsSegment(S(11, 0, 21, 10)))
	expect(t, !S(0, 0, 10, 10).IntersectsSegment(S(-11, 0, -21, 10)))
	expect(t, !S(10, 10, 0, 10).IntersectsSegment(S(11, 0, 21, 10)))
	expect(t, !S(10, 10, 0, 10).IntersectsSegment(S(-11, 0, -21, 10)))
	expect(t, !S(0, 0, 10, 10).IntersectsSegment(S(0, 11, 10, 21)))
	expect(t, !S(0, 0, 10, 10).IntersectsSegment(S(0, -11, 10, -21)))
	expect(t, !S(10, 10, 0, 0).IntersectsSegment(S(0, 11, 10, 21)))
	expect(t, !S(10, 10, 0, 0).IntersectsSegment(S(0, -11, 10, -21)))
	expect(t, !S(0, 0, 10, 0).IntersectsSegment(S(11, 0, 21, 0)))
	expect(t, !S(0, 0, 10, 0).IntersectsSegment(S(0, 1, 10, 1)))
	expect(t, !S(0, 0, 10, 0).IntersectsSegment(S(0, -1, 10, -1)))
	expect(t, !S(0, 0, 10, 10).IntersectsSegment(S(1, 0, 11, 10)))

}

func BenchmarkCircleRect(b *testing.B) {
	for i := 4; i < 256; i *= 2 {
		indexed := CircleRing(P(-112, 33), 1000, i)
		indexed.(*baseSeries).buildTree()
		simple := CircleRing(P(-112, 33), 1000, i)
		simple.(*baseSeries).tree = nil
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			b.Run("Simple", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					simple.Rect()
				}
			})
			b.Run("Indexed", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					indexed.Rect()
				}
			})
		})
	}
}
