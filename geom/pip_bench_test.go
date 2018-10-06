package geom

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

	simple := newRingX(points)
	simple.(*baseSeries).tree = nil
	tree := newRingX(points)
	tree.(*baseSeries).buildTree()
	pointOn := points[len(points)/2]

	// ioutil.WriteFile(label+".svg", []byte(tools.SVG(tree.(*baseSeries).tree)), 0666)

	expect(t, ringxContainsPoint(simple, pointIn, true).hit)
	expect(t, ringxContainsPoint(tree, pointIn, true).hit)

	expect(t, ringxContainsPoint(simple, pointOn, true).hit)
	expect(t, ringxContainsPoint(tree, pointOn, true).hit)

	expect(t, !ringxContainsPoint(simple, pointOn, false).hit)
	expect(t, !ringxContainsPoint(tree, pointOn, false).hit)

	expect(t, !ringxContainsPoint(simple, pointOut, true).hit)
	expect(t, !ringxContainsPoint(tree, pointOut, true).hit)
	if os.Getenv("PIPBENCH") == "1" {
		lotsa.Output = os.Stderr
		fmt.Printf(label + "/simp/in  ")
		lotsa.Ops(N, T, func(_, _ int) {
			ringxContainsPoint(simple, pointIn, true)
		})
		fmt.Printf(label + "/tree/in  ")
		lotsa.Ops(N, T, func(_, _ int) {
			ringxContainsPoint(tree, pointIn, true)
		})
		fmt.Printf(label + "/simp/on  ")
		lotsa.Ops(N, T, func(_, _ int) {
			ringxContainsPoint(simple, pointOn, true)
		})
		fmt.Printf(label + "/tree/on  ")
		lotsa.Ops(N, T, func(_, _ int) {
			ringxContainsPoint(tree, pointOn, true)
		})
		fmt.Printf(label + "/simp/out ")
		lotsa.Ops(N, T, func(_, _ int) {
			ringxContainsPoint(simple, pointOut, true)
		})
		fmt.Printf(label + "/tree/out ")
		lotsa.Ops(N, T, func(_, _ int) {
			ringxContainsPoint(tree, pointOut, true)
		})
	}
}

func TestBigArizona(t *testing.T) {
	testBig(t, "az", az, P(-112, 33), P(-114.477539062, 33.99802726))
}

func TestBigTexas(t *testing.T) {
	testBig(t, "tx", tx, P(-98.52539, 29.363027), P(-101.953125, 29.324720161))
}
