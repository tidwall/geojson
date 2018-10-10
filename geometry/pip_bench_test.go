// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geometry

import (
	"fmt"
	"os"
	"testing"

	"github.com/tidwall/lotsa"
)

func testBig(
	t *testing.T, label string, points []Point, pointIn, pointOut Point,
) {
	N, T := 100000, 4

	simple := newRing(points, DefaultIndex)
	simple.(*baseSeries).tree = nil
	tree := newRing(points, DefaultIndex)
	tree.(*baseSeries).buildTree()
	pointOn := points[len(points)/2]

	// ioutil.WriteFile(label+".svg", []byte(tools.SVG(tree.(*baseSeries).tree)), 0666)

	expect(t, ringContainsPoint(simple, pointIn, true).hit)
	expect(t, ringContainsPoint(tree, pointIn, true).hit)

	expect(t, ringContainsPoint(simple, pointOn, true).hit)
	expect(t, ringContainsPoint(tree, pointOn, true).hit)

	expect(t, !ringContainsPoint(simple, pointOn, false).hit)
	expect(t, !ringContainsPoint(tree, pointOn, false).hit)

	expect(t, !ringContainsPoint(simple, pointOut, true).hit)
	expect(t, !ringContainsPoint(tree, pointOut, true).hit)
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
