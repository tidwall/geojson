package poly

import (
	"testing"
)

func P(x, y float64) Point {
	return Point{x, y}
}

func R(minX, minY, maxX, maxY float64) Rect {
	return Rect{P(minX, minY), P(maxX, maxY)}
}

func TestRectIntersects(t *testing.T) {
	if !(Rect{P(0, 0), P(10, 10)}).IntersectsRect(Rect{P(-1, -1), P(1, 1)}) {
		t.Fatal("!")
	}
	if !(Rect{P(0, 0), P(10, 10)}).IntersectsRect(Rect{P(9, 9), P(11, 11)}) {
		t.Fatal("!")
	}
	if !(Rect{P(0, 0), P(10, 10)}).IntersectsRect(Rect{P(9, -1), P(11, 1)}) {
		t.Fatal("!")
	}
	if !(Rect{P(0, 0), P(10, 10)}).IntersectsRect(Rect{P(-1, 9), P(1, 11)}) {
		t.Fatal("!")
	}
	if !(Rect{P(0, 0), P(10, 10)}).IntersectsRect(Rect{P(-1, -1), P(0, 0)}) {
		t.Fatal("!")
	}
	if !(Rect{P(0, 0), P(10, 10)}).IntersectsRect(Rect{P(10, 10), P(11, 11)}) {
		t.Fatal("!")
	}
	if !(Rect{P(0, 0), P(10, 10)}).IntersectsRect(Rect{P(10, -1), P(11, 0)}) {
		t.Fatal("!")
	}
	if !(Rect{P(0, 0), P(10, 10)}).IntersectsRect(Rect{P(-1, 10), P(0, 11)}) {
		t.Fatal("!")
	}
	if !(Rect{P(0, 0), P(10, 10)}).IntersectsRect(Rect{P(1, 1), P(2, 2)}) {
		t.Fatal("!")
	}
}

func TestRectInside(t *testing.T) {
	if !(Rect{P(1, 1), P(9, 9)}).InsideRect(Rect{P(0, 0), P(10, 10)}) {
		t.Fatal("!")
	}
	if (Rect{P(-1, -1), P(9, 9)}).InsideRect(Rect{P(0, 0), P(10, 10)}) {
		t.Fatal("!")
	}
}

func TestRaycastAllMatch(t *testing.T) {
	res := raycast(P(0, 0), P(0, 0), P(0, 0))
	exp := rayres{false, true}
	if res != exp {
		t.Fatalf("expected '%v', got '%v'", exp, res)
	}
}

func TestIssue360(t *testing.T) {
	exterior := Ring{
		{-122.4408378, 37.7341129},
		{-122.4408378, 37.733},
		{-122.44, 37.733},
		{-122.44, 37.7341129},
		{-122.4408378, 37.7341129},
	}
	holes := []Ring{
		Ring{
			{-122.44060993194579, 37.73345766902749},
			{-122.44044363498686, 37.73345766902749},
			{-122.44044363498686, 37.73355524732416},
			{-122.44060993194579, 37.73355524732416},
			{-122.44060993194579, 37.73345766902749},
		},
		Ring{
			{-122.44060724973677, 37.7336888869566},
			{-122.4402102828026, 37.7336888869566},
			{-122.4402102828026, 37.7339752567853},
			{-122.44060724973677, 37.7339752567853},
			{-122.44060724973677, 37.7336888869566},
		},
	}
	_ = holes
	box := Ring{
		{-122.4434208869934, 37.73138728181471},
		{-122.43711233139038, 37.73138728181471},
		{-122.43711233139038, 37.73579951265516},
		{-122.4434208869934, 37.73579951265516},
		{-122.4434208869934, 37.73138728181471},
	}
	if !exterior.IntersectsRing(exterior) {
		t.Fatal("expected true")
	}
	if !exterior.InsideRing(exterior) {
		t.Fatal("expected true")
	}

	if !exterior.IntersectsRing(box) {
		t.Fatal("expected true")
	}
	if !exterior.InsideRing(box) {
		t.Fatal("expected true")
	}
	if !(Polygon{exterior, nil}).IntersectsRing(box) {
		t.Fatal("expected true")
	}
	if !(Polygon{exterior, nil}).InsideRing(box) {

		t.Fatal("expected true")
	}
	if !(Polygon{exterior, holes}).IntersectsRing(box) {
		t.Fatal("expected true")
	}
	if !(Polygon{exterior, holes}).InsideRing(box) {
		t.Fatal("expected true")
	}
	if !(Polygon{exterior, holes}).IntersectsPolygon(Polygon{box, nil}) {
		t.Fatal("expected true")
	}
	if !(Polygon{exterior, holes}).InsidePolygon(Polygon{box, nil}) {
		t.Fatal("expected true")
	}
}
