package poly

import "testing"

func TestRect(t *testing.T) {
	if R(10, 10, 20, 20).InsidePoint(P(15, 15)) {
		t.Fatal("expected false")
	}
	if !R(15, 15, 15, 15).InsidePoint(P(15, 15)) {
		t.Fatal("expected true")
	}
	if !R(15, 15, 15, 15).InsideLine(Line{P(10, 10), P(20, 20)}) {
		t.Fatal("expected true")
	}
	if R(16, 15, 16, 15).InsideLine(Line{P(10, 10), P(20, 20)}) {
		t.Fatal("expected false")
	}
	if !R(10, 10, 20, 20).InsideRing(R(10, 10, 20, 20).Ring()) {
		t.Fatal("expected true")
	}
	if R(10, 10, 20, 20).InsideRing(R(11, 11, 20, 20).Ring()) {
		t.Fatal("expected false")
	}
	if !R(10, 10, 20, 20).IntersectsPoint(P(15, 15)) {
		t.Fatal("expected true")
	}
	if !R(10, 10, 20, 20).IntersectsPoint(P(20, 20)) {
		t.Fatal("expected true")
	}
	if R(10, 10, 20, 20).IntersectsPoint(P(21, 20)) {
		t.Fatal("expected false")
	}
	if !R(10, 10, 20, 20).IntersectsLine(Line{P(15, 0), P(15, 15), P(30, 15)}) {
		t.Fatal("expected true")
	}
	if R(10, 10, 20, 20).IntersectsLine(Line{P(15, 30), P(15, 45), P(30, 45)}) {
		t.Fatal("expected false")
	}
	if !R(18, 15, 22, 15).InsideLine(Line{P(15, 0), P(15, 15), P(30, 15)}) {
		t.Fatal("expected true")
	}
	if !R(15, 15, 22, 15).InsideLine(Line{P(15, 0), P(15, 15), P(30, 15)}) {
		t.Fatal("expected true")
	}
	if !R(15, 15, 30, 15).InsideLine(Line{P(15, 0), P(15, 15), P(30, 15)}) {
		t.Fatal("expected true")
	}
	if R(15, 15, 31, 15).InsideLine(Line{P(15, 0), P(15, 15), P(30, 15)}) {
		t.Fatal("expected false")
	}
	if R(15, 15, 31, 31).InsideLine(Line{P(15, 0), P(15, 15), P(30, 15)}) {
		t.Fatal("expected false")
	}
}
