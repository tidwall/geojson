package poly

import "testing"

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
