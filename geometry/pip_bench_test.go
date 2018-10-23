// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geometry

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"runtime/debug"
	"testing"
	"time"

	"github.com/tidwall/lotsa"
)

func (kind IndexKind) shortString() string {
	switch kind {
	default:
		return "unkn"
	case None:
		return "none"
	case RTree:
		return "rtre"
	case QuadTree:
		return "quad"
	}
}

func testBig(
	t *testing.T, label string, points []Point, pointIn, pointOut Point,
) {
	N, T := 100000, 4

	opts := []IndexOptions{
		IndexOptions{Kind: None, MinPoints: 64},
		IndexOptions{Kind: QuadTree, MinPoints: 64},
		IndexOptions{Kind: RTree, MinPoints: 64},
	}
	for _, opts := range opts {
		var ms1, ms2 runtime.MemStats
		runtime.GC()
		debug.FreeOSMemory()
		runtime.ReadMemStats(&ms1)
		start := time.Now()
		ring := newRing(points, &opts)
		dur := time.Since(start)
		runtime.GC()
		debug.FreeOSMemory()
		runtime.ReadMemStats(&ms2)

		var randPoints []Point
		if os.Getenv("PIPBENCH") == "1" {
			rect := ring.Rect()
			randPoints = make([]Point, N)
			for i := 0; i < N; i++ {
				randPoints[i] = Point{
					X: (rect.Max.X-rect.Min.X)*rand.Float64() + rect.Min.X,
					Y: (rect.Max.Y-rect.Min.Y)*rand.Float64() + rect.Min.Y,
				}
			}
		}

		pointOn := points[len(points)/2]

		// tests
		expect(t, ringContainsPoint(ring, pointIn, true).hit)
		expect(t, ringContainsPoint(ring, pointOn, true).hit)
		expect(t, !ringContainsPoint(ring, pointOn, false).hit)
		expect(t, !ringContainsPoint(ring, pointOut, true).hit)
		if os.Getenv("PIPBENCH") == "1" {
			fmt.Printf("%s/%s     ", label, opts.Kind.shortString())
			mem := ms2.Alloc - ms1.Alloc
			fmt.Printf("%d points created in %s using %d bytes\n",
				ring.NumPoints(), dur, mem)
			lotsa.Output = os.Stdout
			fmt.Printf("%s/%s/in  ", label, opts.Kind.shortString())
			lotsa.Ops(N, T, func(_, _ int) {
				ringContainsPoint(ring, pointIn, true)
			})
			fmt.Printf("%s/%s/on  ", label, opts.Kind.shortString())
			lotsa.Ops(N, T, func(_, _ int) {
				ringContainsPoint(ring, pointOn, true)
			})
			fmt.Printf("%s/%s/out ", label, opts.Kind.shortString())
			lotsa.Ops(N, T, func(_, _ int) {
				ringContainsPoint(ring, pointOut, true)
			})
			fmt.Printf("%s/%s/rnd ", label, opts.Kind.shortString())
			lotsa.Ops(N, T, func(i, _ int) {
				ringContainsPoint(ring, randPoints[i], true)
			})
		}
	}
}

func TestBigArizona(t *testing.T) {
	testBig(t, "az", az, P(-112, 33), P(-114.477539062, 33.99802726))
}

func TestBigTexas(t *testing.T) {
	testBig(t, "tx", tx, P(-98.52539, 29.363027), P(-101.953125, 29.324720161))
}
