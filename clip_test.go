package geojson

import "testing"

func TestClipLineString(t *testing.T) {
	ls, _ := fillLineString(
		[]Position{
			{X: 1, Y: 1},
			{X: 2, Y: 2},
			{X: 3, Y: 1},
		}, nil, nil)
	bbox := BBox{
		Min: Position{X: 1.5, Y: 0.5},
		Max: Position{X: 2.5, Y: 1.8},
	}
	clipped := ls.Clipped(bbox)
	cl, ok := clipped.(MultiLineString)
	if !ok {
		t.Fatal("wrong type")
	}
	if len(cl.Coordinates) != 2 {
		t.Fatal("result must have two parts in MultiString")
	}
}


func TestClipPolygon(t *testing.T) {
	outer := []Position{
		{X: 2, Y: 2},
		{X: 1, Y: 2},
		{X: 1.5, Y: 1.5},
		{X: 1, Y: 1},
		{X: 2, Y: 1},
		{X: 2, Y: 2},
	}
	inner := []Position{
		{X: 1.9, Y: 1.9},
		{X: 1.2, Y: 1.9},
		{X: 1.45, Y: 1.65},
		{X: 1.9, Y: 1.5},
		{X: 1.9, Y: 1.9},

	}
	polygon, _ := fillPolygon([][]Position{outer, inner}, nil, nil)
	bbox := BBox{
		Min: Position{X: 1.3, Y: 1.3},
		Max: Position{X: 1.4, Y: 2.15},
	}
	clipped := polygon.Clipped(bbox)
	cp, ok := clipped.(Polygon)
	if !ok {
		t.Fatal("wrong type")
	}
	if len(cp.Coordinates) != 2 {
		t.Fatal("result must have two parts in Polygon")
	}
}
