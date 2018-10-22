// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geometry

import (
	"math/rand"
	"os"
	"testing"
	"time"
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
func L(points ...Point) *Line {
	return NewLine(points, DefaultIndexOptions)
}

var (
	// rings
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

	// lines
	u1 = []Point{{0, 10}, {0, 0}, {10, 0}, {10, 10}}
	u2 = []Point{{0, 0}, {10, 0}, {10, 10}, {0, 10}}
	u3 = []Point{{10, 0}, {10, 10}, {0, 10}, {0, 0}}
	u4 = []Point{{10, 10}, {0, 10}, {0, 0}, {10, 0}}

	v1 = []Point{{0, 10}, {5, 0}, {10, 10}}
	v2 = []Point{{0, 0}, {10, 5}, {0, 10}}
	v3 = []Point{{10, 0}, {5, 10}, {0, 0}}
	v4 = []Point{{10, 10}, {0, 5}, {10, 0}}
)

func expect(t testing.TB, what bool) {
	t.Helper()
	if !what {
		t.Fatal("expection failure")
	}
}

// func TestRectVarious(t *testing.T) {
// 	expect(t, R(0, 0, 10, 10).ContainsRing(newRingSimple2(octagon)))
// 	expect(t, !R(5, 0, 15, 10).ContainsRing(newRingSimple2(octagon)))
// 	expect(t, R(5, 0, 15, 10).IntersectsRing(newRingSimple2(octagon)))
// 	expect(t, R(0, 0, 10, 10).Center() == P(5, 5))
// }

func TestRaycastBounds(t *testing.T) {
	expect(t, S(0, 0, 10, 10).Raycast(P(20, -1)) == RaycastResult{false, false})
	expect(t, S(10, 10, 0, 0).Raycast(P(-1, 20)) == RaycastResult{false, false})
	expect(t, S(0, 0, 0, 0).Raycast(P(0, 0)) == RaycastResult{false, true})
	expect(t, S(0, 0, 0, 0).Raycast(P(0, 1)) == RaycastResult{false, false})
	expect(t, S(0, 0, 1, 0).Raycast(P(1, 0)) == RaycastResult{false, true})
	expect(t, S(1, 0, 0, 0).Raycast(P(1, 0)) == RaycastResult{false, true})
	expect(t, S(0, 1, 0, 0).Raycast(P(0, 1)) == RaycastResult{false, true})
	expect(t, S(0, 0, 0, 1).Raycast(P(0, 1)) == RaycastResult{false, true})
}
