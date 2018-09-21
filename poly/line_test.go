package poly

import "testing"

func TestLine(t *testing.T) {
	if (Line{P(15, 0), P(15, 15), P(30, 15)}).InsidePoint(P(15, 15)) {
		t.Fatal("expected false")
	}
	if (Line{}).InsidePoint(P(15, 15)) {
		t.Fatal("expected false")
	}
	if !(Line{P(15, 15)}).InsidePoint(P(15, 15)) {
		t.Fatal("expected true")
	}
	if !(Line{P(15, 15), P(15, 15)}).InsidePoint(P(15, 15)) {
		t.Fatal("expected true")
	}
	if !(Line{P(15, 0), P(15, 15), P(30, 15)}).InsideRect(R(0, 0, 30, 30)) {
		t.Fatal("expected true")
	}
	if (Line{}).InsideRect(R(0, 0, 30, 30)) {
		t.Fatal("expected false")
	}
	if (Line{P(15, 0), P(15, 15), P(30, 15)}).InsideRect(R(0, 0, 20, 20)) {
		t.Fatal("expected false")
	}

	ln1 := Line{P(5, 0), P(5, 5), P(10, 5), P(10, 10), P(15, 10), P(15, 15)}
	lns := []Line{
		Line{P(7, 5), P(10, 5), P(10, 10), P(12, 10)},
		Line{P(7, 5), P(8, 5), P(10, 5), P(10, 10), P(12, 10)},
		Line{P(7, 5), P(8, 5), P(6, 5), P(10, 5), P(10, 8), P(10, 5),
			P(5, 5), P(10, 5), P(10, 10), P(12, 10)},
	}
	for i, ln := range lns {
		if !ln.InsideLine(ln1) {
			t.Fatalf("expected true for index: %d", i)
		}
	}
	if (Line{P(5, -1), P(5, 5), P(10, 5)}).InsideLine(ln1) {
		t.Fatal("expected false")
	}
	if (Line{P(5, 0), P(5, 5), P(5, 0), P(10, 0)}).InsideLine(ln1) {
		t.Fatal("expected false")
	}
	if (Line{
		P(5, 0), P(5, 5), P(10, 5), P(10, 10), P(15, 10), P(15, 15), P(20, 20),
	}).InsideLine(ln1) {
		t.Fatal("expected false")
	}
	if (Line{}).InsideLine(ln1) {
		t.Fatal("expected false")
	}
	if (Line{P(5, 0)}).InsideLine(Line{}) {
		t.Fatal("expected false")
	}
	if !(Line{P(5, 0)}).InsideLine(Line{P(5, 0), P(10, 0)}) {
		t.Fatal("expected true")
	}
	if !(Line{P(5, 0)}).InsideLine(Line{P(5, 0)}) {
		t.Fatal("expected true")
	}
	if !(Line{P(5, 0), P(5, 0)}).InsideLine(Line{P(5, 0)}) {
		t.Fatal("expected true")
	}
	if (Line{P(5, 0), P(5, 0), P(6, 0)}).InsideLine(Line{P(5, 0)}) {
		t.Fatal("expected false")
	}

	if !(Line{P(15, 0), P(15, 15), P(30, 15)}).InsideRing(R(0, 0, 30, 30).Ring()) {
		t.Fatal("expected true")
	}
	if (Line{}).InsideRing(R(0, 0, 30, 30).Ring()) {
		t.Fatal("expected false")
	}
	if (Line{P(15, 0), P(15, 15), P(30, 15)}).InsideRing(Ring{}) {
		t.Fatal("expected false")
	}
	if (Line{P(15, 0), P(15, 15), P(30, 15)}).InsideRing(R(0, 0, 20, 20).Ring()) {
		t.Fatal("expected false")
	}

}
