package geojson

import (
	"testing"

	"github.com/tidwall/geojson/geometry"
)

func TestClipLineString(t *testing.T) {
	ls := LO([]geometry.Point{
		{X: 1, Y: 1},
		{X: 2, Y: 2},
		{X: 3, Y: 1}})
	clipped := ls.Clipped(NewRect(1.5, 0.5, 2.5, 1.8))
	cl, ok := clipped.(*MultiLineString)
	if !ok {
		t.Fatal("wrong type")
	}
	if len(cl.children) != 2 {
		t.Fatal("result must have two parts in MultiString")
	}
}

func TestClipPolygon(t *testing.T) {
	exterior := []geometry.Point{
		{X: 2, Y: 2},
		{X: 1, Y: 2},
		{X: 1.5, Y: 1.5},
		{X: 1, Y: 1},
		{X: 2, Y: 1},
		{X: 2, Y: 2},
	}
	holes := [][]geometry.Point{
		[]geometry.Point{
			{X: 1.9, Y: 1.9},
			{X: 1.2, Y: 1.9},
			{X: 1.45, Y: 1.65},
			{X: 1.9, Y: 1.5},
			{X: 1.9, Y: 1.9},
		},
	}
	polygon := PPO(exterior, holes)
	clipped := polygon.Clipped(NewRect(1.3, 1.3, 1.4, 2.15))
	cp, ok := clipped.(*Polygon)
	if !ok {
		t.Fatal("wrong type")
	}
	if !cp.base.Exterior.Empty() && len(cp.base.Holes) != 1 {
		t.Fatal("result must have two parts in Polygon")
	}
}
