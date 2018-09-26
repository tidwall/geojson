package ring

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
	test := func(t *testing.T, indexed bool) {
		rectangleRing := NewRing(rectangle, indexed)
		var segs []Segment
		rectangleRing.Scan(func(seg Segment) bool {
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
		notClosedRing := NewRing(rectangle, indexed)
		notClosedRing.Scan(func(seg Segment) bool {
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
	test := func(t *testing.T, indexed bool) {
		octagonRing := NewRing(octagon, indexed)
		var segs []Segment
		octagonRing.Search(R(0, 0, 0, 0), func(seg Segment, _ int) bool {
			segs = append(segs, seg)
			return true
		})
		segsExpect := []Segment{
			S(0, 3, 3, 0),
		}
		for i := 0; i < len(segs); i++ {
			expect(t, segs[i] == segsExpect[i])
		}
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
		for i := 0; i < len(segs); i++ {
			expect(t, segs[i] == segsExpect[i])
		}
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

func TestRingIntersectsSegment(t *testing.T) {
	simple := NewRing(concave1, false)
	tree := NewRing(concave1, true)

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
	simple := NewRing(concave1, false)
	tree := NewRing(concave1, true)
	small := NewRing([]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}, false).(*simpleRing)

	intersects := func(ring Ring) bool {
		tt := simple.IntersectsRing(ring, true)
		if tree.IntersectsRing(ring, true) != tt {
			panic("structure mismatch")
		}
		return tt
	}

	intersectsOnEdgeNotAllowed := func(ring Ring) bool {
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
	simple := NewRing(az, false)
	tree := NewRing(az, true)
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

func TestBigArizona(t *testing.T) {
	simple := NewRing(az, false)
	tree := NewRing(az, true)
	pointIn := P(-112, 33)
	pointOut := P(-114.47753906249999, 33.99802726234877)
	pointOn := P(-114.604715, 35.061744)

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
		fmt.Printf("az/tree/in  ")
		lotsa.Ops(1000, 1, func(_, _ int) {
			simple.ContainsPoint(pointIn, true)
		})
		fmt.Printf("az/simp/in  ")
		lotsa.Ops(1000, 1, func(_, _ int) {
			tree.ContainsPoint(pointIn, true)
		})
		fmt.Printf("az/simp/on  ")
		lotsa.Ops(1000, 1, func(_, _ int) {
			simple.ContainsPoint(pointOn, true)
		})
		fmt.Printf("az/tree/on  ")
		lotsa.Ops(1000, 1, func(_, _ int) {
			tree.ContainsPoint(pointOn, true)
		})
		fmt.Printf("az/simp/out ")
		lotsa.Ops(1000, 1, func(_, _ int) {
			simple.ContainsPoint(pointOut, true)
		})
		fmt.Printf("az/tree/out ")
		lotsa.Ops(1000, 1, func(_, _ int) {
			tree.ContainsPoint(pointOut, true)
		})
	}
}

func TestBigTexas(t *testing.T) {
	simple := NewRing(tx, false)
	tree := NewRing(tx, true)
	pointIn := P(-98.525390625, 29.36302703778376)
	pointOut := P(-101.953125, 29.32472016151103)
	pointOn := P(-100.402214, 28.532657)

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
		fmt.Printf("tx/simp/in  ")
		lotsa.Ops(1000, 1, func(_, _ int) {
			simple.ContainsPoint(pointIn, true)
		})
		fmt.Printf("tx/tree/in  ")
		lotsa.Ops(1000, 1, func(_, _ int) {
			tree.ContainsPoint(pointIn, true)
		})
		fmt.Printf("tx/simp/on  ")
		lotsa.Ops(1000, 1, func(_, _ int) {
			simple.ContainsPoint(pointOn, true)
		})
		fmt.Printf("tx/tree/on  ")
		lotsa.Ops(1000, 1, func(_, _ int) {
			tree.ContainsPoint(pointOn, true)
		})
		fmt.Printf("tx/simp/out ")
		lotsa.Ops(1000, 1, func(_, _ int) {
			simple.ContainsPoint(pointOut, true)
		})
		fmt.Printf("tx/tree/out ")
		lotsa.Ops(1000, 1, func(_, _ int) {
			tree.ContainsPoint(pointOut, true)
		})
	}
}

func TestRingContainsRing(t *testing.T) {
	simple := NewRing(concave1, false)
	tree := NewRing(concave1, true)

	expect(t, simple.ContainsRing(simple, true))
	expect(t, simple.ContainsRing(tree, true))
	expect(t, tree.ContainsRing(simple, true))
	expect(t, tree.ContainsRing(tree, true))

	expect(t, !simple.ContainsRing(simple, false))
	expect(t, !simple.ContainsRing(tree, false))
	expect(t, !tree.ContainsRing(simple, false))
	expect(t, !tree.ContainsRing(tree, false))

	small := NewRing([]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}, false).(*simpleRing)

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
	simple := NewRing(bowtie, false)
	tree := NewRing(bowtie, true)
	square := NewRing([]Point{P(3, 3), P(7, 3), P(7, 7), P(3, 7), P(3, 3)}, false)

	expect(t, simple.IntersectsRing(square, true))
	expect(t, tree.IntersectsRing(square, true))
	expect(t, !simple.ContainsRing(square, true))
	expect(t, !tree.ContainsRing(square, true))

}

func TestVarious(t *testing.T) {
	ring := NewRing(octagon[:len(octagon)-1], true)
	n := 0
	ring.Search(R(0, 0, 10, 10), func(seg Segment, index int) bool {
		n++
		return true
	})
	expect(t, n == 8)
	n = 0
	ring.Scan(func(seg Segment) bool {
		n++
		return true
	})
	expect(t, n == 8)
	n = 0
	ring.Scan(func(seg Segment) bool {
		n++
		return false
	})
	expect(t, n == 1)
	expect(t, ring.IntersectsSegment(S(0, 0, 4, 4), true))
	expect(t, !NewRing([]Point{}, false).Convex())
	expect(t, NewRing(octagon, false).Convex())
	expect(t, !NewRing([]Point{}, true).Convex())
	expect(t, NewRing(octagon, true).Convex())

	ring = NewRing(octagon[:len(octagon)-1], false)
	n = 0
	ring.Search(R(0, 0, 10, 10), func(seg Segment, index int) bool {
		n++
		return false
	})
	expect(t, n == 1)
	n = 0
	ring.Scan(func(seg Segment) bool {
		n++
		return true
	})
	expect(t, n == 8)
	expect(t, ring.IntersectsSegment(S(0, 0, 4, 4), true))

	small := NewRing([]Point{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}, false).(*simpleRing)
	expect(t, small.IntersectsRing(ring, true))
	expect(t, ring.IntersectsRing(small, true))

	expect(t, raycast(P(0, 0), P(0, 0), P(0, 0)).on)
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
