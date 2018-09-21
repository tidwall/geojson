package poly

import "testing"

func TestRing(t *testing.T) {
	ring := R(10, 10, 20, 20).Ring()
	if ring.InsidePoint(P(15, 15)) {
		t.Fatal("expected false")
	}
	if (Ring{}).InsidePoint(P(15, 15)) {
		t.Fatal("expected false")
	}
	if !R(15, 15, 15, 15).Ring().InsidePoint(P(15, 15)) {
		t.Fatal("expected true")
	}
	if !R(15, 15, 15, 15).Ring().InsideLine(Line{P(10, 15), P(20, 15)}) {
		t.Fatal("expected true")
	}
	if !R(15, 15, 15, 20).Ring().InsideLine(Line{P(15, 10), P(15, 30)}) {
		t.Fatal("expected true")
	}
	if (Ring{}).InsideLine(Line{}) {
		t.Fatal("expected false")
	}
	if R(10, 10, 20, 20).Ring().InsideRing(R(12, 12, 18, 18).Ring()) {
		t.Fatal("expected false")
	}
	if !R(10, 10, 20, 20).Ring().InsideRing(R(8, 8, 22, 22).Ring()) {
		t.Fatal("expected true")
	}
	if !R(10, 10, 20, 20).Ring().IntersectsLine(Line{P(15, 10), P(15, 30)}) {
		t.Fatal("expected true")
	}
	if !R(10, 10, 20, 20).Ring().IntersectsPoint(P(15, 15)) {
		t.Fatal("expected true")
	}
	if R(10, 10, 20, 20).Ring().IntersectsPoint(P(30, 30)) {
		t.Fatal("expected false")
	}
	if !R(10, 10, 20, 20).Ring().IntersectsRing(R(0, 0, 15, 15).Ring()) {
		t.Fatal("expected true")
	}
}
