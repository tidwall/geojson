package geom

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/tidwall/lotsa"
)

func init() {
	seed := time.Now().UnixNano()
	println(seed)
	rand.Seed(seed)
	if os.Getenv("PIPBENCH") != "1" {
		println("use PIPBENCH=1 for point-in-polygon benchmarks")
	}
}

func S(ax, ay, bx, by float64) Segment {
	return Segment{Point{ax, ay}, Point{bx, by}}
}
func R(minX, minY, maxX, maxY float64) Rect {
	return Rect{Point{minX, minY}, Point{maxX, maxY}}
}
func P(x, y float64) Point {
	return Point{x, y}
}

var (
	rectangle = []Point{{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0}}
	pentagon  = []Point{{2, 2}, {8, 0}, {10, 6}, {5, 10}, {0, 6}, {2, 2}}
	triangle  = []Point{{0, 0}, {10, 0}, {5, 10}, {0, 0}}
	trapezoid = []Point{{0, 0}, {10, 0}, {8, 10}, {2, 10}, {0, 0}}
	octagon   = []Point{
		{3, 0}, {7, 0}, {10, 3}, {10, 7},
		{7, 10}, {3, 10}, {0, 7}, {0, 3}, {3, 0},
	}
	concave1  = []Point{{5, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 5}, {5, 5}, {5, 0}}
	concave2  = []Point{{0, 0}, {5, 0}, {5, 5}, {10, 5}, {10, 10}, {0, 10}, {0, 0}}
	concave3  = []Point{{0, 0}, {10, 0}, {10, 5}, {5, 5}, {5, 10}, {0, 10}, {0, 0}}
	concave4  = []Point{{0, 0}, {10, 0}, {10, 10}, {5, 10}, {5, 5}, {0, 5}, {0, 0}}
	bowtie    = []Point{{0, 0}, {5, 4}, {10, 0}, {10, 10}, {5, 6}, {0, 10}, {0, 0}}
	notClosed = []Point{{0, 0}, {10, 0}, {10, 10}, {0, 10}}
)

func expect(t testing.TB, what bool) {
	t.Helper()
	if !what {
		t.Fatal("expection failure")
	}
}

func TestRingScan(t *testing.T) {
	test := func(t *testing.T, index bool) {
		rectangleRing := NewRing2(rectangle)
		if !index {
			rectangleRing.tree = nil
		} else {
			rectangleRing.buildTree()
		}
		var segs []Segment
		rectangleRing.Scan(func(seg Segment, idx int) bool {
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
		notClosedRing := NewRing2(rectangle)
		if !index {
			notClosedRing.tree = nil
		} else {
			notClosedRing.buildTree()
		}
		notClosedRing.Scan(func(seg Segment, idx int) bool {
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
		octagonRing := NewRing2(octagon)
		if !index {
			octagonRing.tree = nil
		} else {
			octagonRing.buildTree()
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
	simple := NewRing2(concave1)
	simple.tree = nil
	tree := NewRing2(concave1)
	tree.buildTree()

	expect(t, !simple.IntersectsSegment(S(0, 0, 3, 3), true))
	expect(t, !tree.IntersectsSegment(S(0, 0, 3, 3), true))
	expect(t, !simple.IntersectsSegment(S(0, 0, 3, 3), false))
	expect(t, !tree.IntersectsSegment(S(0, 0, 3, 3), false))

	expect(t, simple.IntersectsSegment(S(0, 0, 5, 5), true))
	expect(t, tree.IntersectsSegment(S(0, 0, 5, 5), true))
	expect(t, !simple.IntersectsSegment(S(0, 0, 5, 5), false))
	expect(t, !tree.IntersectsSegment(S(0, 0, 5, 5), false))

	expect(t, simple.IntersectsSegment(S(0, 0, 10, 10), true))
	expect(t, tree.IntersectsSegment(S(0, 0, 10, 10), true))
	expect(t, !simple.IntersectsSegment(S(0, 0, 10, 10), false))
	expect(t, !tree.IntersectsSegment(S(0, 0, 10, 10), false))

	expect(t, simple.IntersectsSegment(S(0, 0, 11, 11), true))
	expect(t, tree.IntersectsSegment(S(0, 0, 11, 11), true))
	expect(t, !simple.IntersectsSegment(S(0, 0, 11, 11), false))
	expect(t, !tree.IntersectsSegment(S(0, 0, 11, 11), false))

}

func TestRingIntersectsRing(t *testing.T) {
	simple := NewRing2(concave1)
	simple.tree = nil
	tree := NewRing2(concave1)
	tree.buildTree()
	small := NewRing2([]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}})
	small.tree = nil

	intersects := func(ring *Ring) bool {
		tt := simple.IntersectsRing(ring, true)
		if tree.IntersectsRing(ring, true) != tt {
			panic("structure mismatch")
		}
		return tt
	}

	intersectsOnEdgeNotAllowed := func(ring *Ring) bool {
		tt := simple.IntersectsRing(ring, false)
		if tree.IntersectsRing(ring, false) != tt {
			panic("structure mismatch")
		}
		return tt
	}

	expect(t, intersects(small))
	expect(t, intersects(small.move(-6, 0)))
	expect(t, intersects(small.move(6, 0)))
	expect(t, !intersects(small.move(-7, 0)))
	expect(t, !intersects(small.move(7, 0)))
	expect(t, intersects(small.move(1, 1)))
	expect(t, intersects(small.move(-1, -1)))
	expect(t, intersects(small.move(2, 2)))
	expect(t, !intersects(small.move(-2, -2)))
	expect(t, intersects(small.move(0, -6)))
	expect(t, intersects(small.move(0, 6)))
	expect(t, !intersects(small.move(0, -7)))
	expect(t, !intersects(small.move(0, 7)))

	expect(t, intersectsOnEdgeNotAllowed(small.move(-5, 0)))
	expect(t, intersectsOnEdgeNotAllowed(small.move(5, 0)))
	expect(t, intersectsOnEdgeNotAllowed(small.move(0, -5)))
	expect(t, intersectsOnEdgeNotAllowed(small.move(0, 5)))

	expect(t, !intersectsOnEdgeNotAllowed(small.move(-6, 0)))
	expect(t, !intersectsOnEdgeNotAllowed(small.move(6, 0)))
	expect(t, !intersectsOnEdgeNotAllowed(small.move(0, -6)))
	expect(t, !intersectsOnEdgeNotAllowed(small.move(0, 6)))

	expect(t, intersectsOnEdgeNotAllowed(small.move(1, 1)))
	expect(t, !intersectsOnEdgeNotAllowed(small.move(-1, -1)))

}

func TestBigRandomPIP(t *testing.T) {
	simple := NewRing2(az)
	simple.tree = nil
	tree := NewRing2(az)
	tree.buildTree()
	expect(t, simple.Rect() == tree.Rect())
	rect := tree.Rect()
	start := time.Now()
	for time.Since(start) < time.Second/4 {
		point := P(
			rand.Float64()*(rect.Max.X-rect.Min.X)+rect.Min.X,
			rand.Float64()*(rect.Max.Y-rect.Min.Y)+rect.Min.Y,
		)
		expect(t, tree.ContainsPoint(point, true) ==
			simple.ContainsPoint(point, true))
	}
}

func testBig(
	t *testing.T, label string, points []Point, pointIn, pointOut Point,
) {
	N := 10000
	simple := NewRing2(points)
	simple.tree = nil
	tree := NewRing2(points)
	tree.buildTree()
	pointOn := points[len(points)/2]

	expect(t, simple.ContainsPoint(pointIn, true))
	expect(t, tree.ContainsPoint(pointIn, true))

	expect(t, simple.ContainsPoint(pointOn, true))
	expect(t, tree.ContainsPoint(pointOn, true))

	expect(t, !simple.ContainsPoint(pointOn, false))
	expect(t, !tree.ContainsPoint(pointOn, false))

	expect(t, !simple.ContainsPoint(pointOut, true))
	expect(t, !tree.ContainsPoint(pointOut, true))
	if os.Getenv("PIPBENCH") == "1" {
		lotsa.Output = os.Stderr
		fmt.Printf(label + "/simp/in  ")
		lotsa.Ops(N, 1, func(_, _ int) {
			simple.ContainsPoint(pointIn, true)
		})
		fmt.Printf(label + "/tree/in  ")
		lotsa.Ops(N, 1, func(_, _ int) {
			tree.ContainsPoint(pointIn, true)
		})
		fmt.Printf(label + "/simp/on  ")
		lotsa.Ops(N, 1, func(_, _ int) {
			simple.ContainsPoint(pointOn, true)
		})
		fmt.Printf(label + "/tree/on  ")
		lotsa.Ops(N, 1, func(_, _ int) {
			tree.ContainsPoint(pointOn, true)
		})
		fmt.Printf(label + "/simp/out ")
		lotsa.Ops(N, 1, func(_, _ int) {
			simple.ContainsPoint(pointOut, true)
		})
		fmt.Printf(label + "/tree/out ")
		lotsa.Ops(N, 1, func(_, _ int) {
			tree.ContainsPoint(pointOut, true)
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
	circle := CircleRing(P(-100.1, 31.2), 660000, 10000).Points()
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
	circle = CircleRing(P(-100.1, 31.2), 660000, 2).Points()
	expect(t, len(circle) == 4)
}

func TestRingContainsRing(t *testing.T) {
	simple := NewRing2(concave1)
	simple.tree = nil
	tree := NewRing2(concave1)
	tree.buildTree()

	expect(t, simple.ContainsRing(simple, true))
	expect(t, simple.ContainsRing(tree, true))
	expect(t, tree.ContainsRing(simple, true))
	expect(t, tree.ContainsRing(tree, true))

	expect(t, !simple.ContainsRing(simple, false))
	expect(t, !simple.ContainsRing(tree, false))
	expect(t, !tree.ContainsRing(simple, false))
	expect(t, !tree.ContainsRing(tree, false))

	small := NewRing2([]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}})
	small.tree = nil

	expect(t, !simple.ContainsRing(small, true))
	expect(t, !tree.ContainsRing(small, true))

	for x := 1.0; x <= 4; x++ {
		expect(t, simple.ContainsRing(small.move(x, 0), true))
		expect(t, tree.ContainsRing(small.move(x, 0), true))
	}
	expect(t, !simple.ContainsRing(small.move(4, 0), false))
	expect(t, !tree.ContainsRing(small.move(4, 0), false))
	for y := 1.0; y <= 4; y++ {
		expect(t, simple.ContainsRing(small.move(0, y), true))
		expect(t, tree.ContainsRing(small.move(0, y), true))
	}
	expect(t, !simple.ContainsRing(small.move(0, 4), false))
	expect(t, !tree.ContainsRing(small.move(0, 4), false))

	for x := -1.0; x >= -4; x-- {
		expect(t, !simple.ContainsRing(small.move(x, 0), true))
		expect(t, !tree.ContainsRing(small.move(x, 0), true))
	}
	expect(t, !simple.ContainsRing(small.move(-4, 0), false))
	expect(t, !tree.ContainsRing(small.move(-4, 0), false))
	for y := -1.0; y >= -4; y-- {
		expect(t, !simple.ContainsRing(small.move(0, y), true))
		expect(t, !tree.ContainsRing(small.move(0, y), true))
	}
	expect(t, !simple.ContainsRing(small.move(0, -4), false))
	expect(t, !tree.ContainsRing(small.move(0, -4), false))

	expect(t, !simple.ContainsRing(small.move(1, 0), false))
	expect(t, !tree.ContainsRing(small.move(1, 0), false))
	expect(t, simple.ContainsRing(small.move(2, 0), false))
	expect(t, tree.ContainsRing(small.move(2, 0), false))
	expect(t, simple.ContainsRing(small.move(2, 2), false))
	expect(t, tree.ContainsRing(small.move(2, 2), false))
	expect(t, !simple.ContainsRing(small.move(-2, -2), false))
	expect(t, !tree.ContainsRing(small.move(-2, -2), false))

	expect(t, !simple.ContainsRing(small.move(5, 0), true))
	expect(t, !tree.ContainsRing(small.move(5, 0), true))
	expect(t, !simple.ContainsRing(small.move(-5, 0), true))
	expect(t, !tree.ContainsRing(small.move(-5, 0), true))

	expect(t, !simple.ContainsRing(small.move(0, 5), true))
	expect(t, !tree.ContainsRing(small.move(0, 5), true))
	expect(t, !simple.ContainsRing(small.move(0, -5), true))
	expect(t, !tree.ContainsRing(small.move(0, -5), true))

}
func TestBowtie(t *testing.T) {
	simple := NewRing2(bowtie)
	simple.tree = nil
	tree := NewRing2(bowtie)
	tree.buildTree()
	square := NewRing2([]Point{P(3, 3), P(7, 3), P(7, 7), P(3, 7), P(3, 3)})
	square.tree = nil

	expect(t, simple.IntersectsRing(square, true))
	expect(t, tree.IntersectsRing(square, true))
	expect(t, !simple.ContainsRing(square, true))
	expect(t, !tree.ContainsRing(square, true))

}

func TestRingVarious(t *testing.T) {
	ring := NewRing2(octagon[:len(octagon)-1])
	ring.buildTree()
	n := 0
	ring.Search(R(0, 0, 10, 10), func(seg Segment, index int) bool {
		n++
		return true
	})
	expect(t, n == 8)
	n = 0
	ring.Scan(func(seg Segment, idx int) bool {
		n++
		return true
	})
	expect(t, n == 8)
	n = 0
	ring.Scan(func(seg Segment, idx int) bool {
		n++
		return false
	})
	expect(t, n == 1)
	expect(t, ring.IntersectsSegment(S(0, 0, 4, 4), true))
	expect(t, !newRingSimple2([]Point{}).Convex())
	expect(t, newRingSimple2(octagon).Convex())
	expect(t, !newRingIndexed2([]Point{}).Convex())
	expect(t, newRingIndexed2(octagon).Convex())

	ring = NewRing2(octagon[:len(octagon)-1])
	ring.tree = nil
	n = 0
	ring.Search(R(0, 0, 10, 10), func(seg Segment, index int) bool {
		n++
		return false
	})
	expect(t, n == 1)
	n = 0
	ring.Scan(func(seg Segment, idx int) bool {
		n++
		return true
	})
	expect(t, n == 8)
	expect(t, ring.IntersectsSegment(S(0, 0, 4, 4), true))

	small := NewRing2([]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}})
	small.tree = nil

	expect(t, small.IntersectsRing(ring, true))
	expect(t, ring.IntersectsRing(small, true))

	expect(t, raycast(P(0, 0), P(0, 0), P(0, 0)).on)

	ring1 := NewRing2(octagon)
	ring1.tree = nil
	n1 := 0
	ring1.Scan(func(seg Segment, idx int) bool {
		n1++
		return true
	})
	expect(t, ring1.Closed())
	ring2 := NewRing2(octagon[:len(octagon)-1])
	ring2.tree = nil
	n2 := 0
	ring2.Scan(func(seg Segment, idx int) bool {
		n2++
		return true
	})
	expect(t, n1 == n2)
	expect(t, ring2.Closed())

	ring1 = NewRing2(octagon)
	ring1.buildTree()
	n1 = 0
	ring1.Scan(func(seg Segment, idx int) bool {
		n1++
		return true
	})
	expect(t, ring1.Closed())
	ring2 = NewRing2(octagon[:len(octagon)-1])
	ring2.buildTree()
	n2 = 0
	ring2.Scan(func(seg Segment, idx int) bool {
		n2++
		return true
	})
	expect(t, n1 == n2)
	expect(t, ring2.Closed())

	convex, rect := pointsConvexRect([]Point{P(0, 0)})
	expect(t, !convex)
	expect(t, rect == Rect{})
}

func newRingSimple2(points []Point) *Ring {
	ring := NewRing2(points)
	ring.tree = nil
	return ring
}
func newRingIndexed2(points []Point) *Ring {
	ring := NewRing2(points)
	ring.buildTree()
	return ring
}

func TestRingContainsPoint(t *testing.T) {
	expect(t, newRingSimple2(octagon).IntersectsPoint(P(4, 4), true))
	expect(t, newRingIndexed2(octagon).IntersectsPoint(P(4, 4), true))
}

func TestRingContainsSegment(t *testing.T) {
	expect(t, newRingSimple2(octagon).ContainsSegment(S(4, 4, 6, 6), true))
	expect(t, newRingIndexed2(octagon).ContainsSegment(S(4, 4, 6, 6), true))
	expect(t, !newRingSimple2(octagon).ContainsSegment(S(9, 4, 11, 6), true))
	expect(t, !newRingIndexed2(octagon).ContainsSegment(S(9, 4, 11, 6), true))
	expect(t, !newRingSimple2(octagon).ContainsSegment(S(11, 4, 9, 6), true))
	expect(t, !newRingIndexed2(octagon).ContainsSegment(S(11, 4, 9, 6), true))
	expect(t, !newRingSimple2(concave1).ContainsSegment(S(11, 4, 9, 6), true))
	expect(t, !newRingIndexed2(concave1).ContainsSegment(S(11, 4, 9, 6), true))
	expect(t, newRingSimple2(concave1).ContainsSegment(S(6, 6, 8, 8), true))
	expect(t, newRingIndexed2(concave1).ContainsSegment(S(6, 6, 8, 8), true))
	expect(t, !newRingSimple2(concave1).ContainsSegment(S(1, 6, 6, 1), true))
	expect(t, !newRingIndexed2(concave1).ContainsSegment(S(1, 6, 6, 1), true))
}
func TestRingContainsRect(t *testing.T) {
	expect(t, newRingSimple2(octagon).ContainsRect(R(4, 4, 6, 6), true))
	expect(t, newRingIndexed2(octagon).ContainsRect(R(4, 4, 6, 6), true))
	expect(t, newRingSimple2(octagon).ContainsRect(R(4, 4, 6, 6), false))
	expect(t, newRingIndexed2(octagon).ContainsRect(R(4, 4, 6, 6), false))
}
func TestRingIntersectsRect(t *testing.T) {
	expect(t, newRingSimple2(octagon).IntersectsRect(R(9, 4, 11, 6), true))
	expect(t, newRingIndexed2(octagon).IntersectsRect(R(9, 4, 11, 6), true))
	expect(t, !newRingSimple2(octagon).IntersectsRect(R(10, 4, 12, 6), false))
	expect(t, !newRingIndexed2(octagon).IntersectsRect(R(10, 4, 12, 6), false))
	expect(t, newRingSimple2(octagon).IntersectsRect(R(10, 4, 12, 6), true))
	expect(t, newRingIndexed2(octagon).IntersectsRect(R(10, 4, 12, 6), true))
	expect(t, !newRingSimple2(octagon).IntersectsRect(R(11, 4, 13, 6), true))
	expect(t, !newRingIndexed2(octagon).IntersectsRect(R(11, 4, 13, 6), true))
}
func TestRingContainsPoly(t *testing.T) {
	expect(t, newRingSimple2(octagon).ContainsPoly(
		NewPoly2(octagon, nil), true))
	expect(t, newRingIndexed2(octagon).ContainsPoly(
		NewPoly2(octagon, nil), true))
	expect(t, !newRingSimple2(octagon).ContainsPoly(
		NewPoly2(octagon, nil), false))
	expect(t, !newRingIndexed2(octagon).ContainsPoly(
		NewPoly2(octagon, nil), false))
}
func TestRingIntersectsPoly(t *testing.T) {
	expect(t, newRingSimple2(octagon).move(5, 0).IntersectsPoly(
		NewPoly2(octagon, nil), true))
	expect(t, newRingIndexed2(octagon).move(5, 0).IntersectsPoly(
		NewPoly2(octagon, nil), true))
	expect(t, newRingSimple2(octagon).move(10, 0).IntersectsPoly(
		NewPoly2(octagon, nil), true))
	expect(t, newRingIndexed2(octagon).move(10, 0).IntersectsPoly(
		NewPoly2(octagon, nil), true))
	expect(t, !newRingSimple2(octagon).move(10, 0).IntersectsPoly(
		NewPoly2(octagon, nil), false))
	expect(t, !newRingIndexed2(octagon).move(10, 0).IntersectsPoly(
		NewPoly2(octagon, nil), false))
}

func TestSegmentsIntersect(t *testing.T) {
	expect(t, !segmentsIntersect(P(0, 0), P(10, 10), P(11, 0), P(21, 10)))
	expect(t, !segmentsIntersect(P(0, 0), P(10, 10), P(-11, 0), P(-21, 10)))
	expect(t, !segmentsIntersect(P(10, 10), P(0, 10), P(11, 0), P(21, 10)))
	expect(t, !segmentsIntersect(P(10, 10), P(0, 10), P(-11, 0), P(-21, 10)))
	expect(t, !segmentsIntersect(P(0, 0), P(10, 10), P(0, 11), P(10, 21)))
	expect(t, !segmentsIntersect(P(0, 0), P(10, 10), P(0, -11), P(10, -21)))
	expect(t, !segmentsIntersect(P(10, 10), P(0, 0), P(0, 11), P(10, 21)))
	expect(t, !segmentsIntersect(P(10, 10), P(0, 0), P(0, -11), P(10, -21)))

	expect(t, !segmentsIntersect(P(0, 0), P(10, 0), P(11, 0), P(21, 0)))
	expect(t, !segmentsIntersect(P(0, 0), P(10, 0), P(0, 1), P(10, 1)))
	expect(t, !segmentsIntersect(P(0, 0), P(10, 0), P(0, -1), P(10, -1)))
	expect(t, !segmentsIntersect(P(0, 0), P(10, 10), P(1, 0), P(11, 10)))

}

func BenchmarkCircleRect(b *testing.B) {
	for i := 4; i < 256; i *= 2 {
		indexed := CircleRing(P(-112, 33), 1000, i)
		indexed.buildTree()
		simple := CircleRing(P(-112, 33), 1000, i)
		simple.tree = nil
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
