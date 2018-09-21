package poly

import "testing"

func TestPoint(t *testing.T) {
	if !P(15, 15).InsidePoint(P(15, 15)) {
		t.Fatal("expected true")
	}
	if P(15, 15).InsidePoint(P(16, 15)) {
		t.Fatal("expected false")
	}
	if !P(15, 15).InsideRing(R(10, 10, 20, 20).Ring()) {
		t.Fatal("expected true")
	}
	if !P(10, 10).InsideRing(R(10, 10, 20, 20).Ring()) {
		t.Fatal("expected true")
	}
	if P(9, 10).InsideRing(R(10, 10, 20, 20).Ring()) {
		t.Fatal("expected false")
	}
	if !P(15, 15).IntersectsPoint(P(15, 15)) {
		t.Fatal("expected true")
	}
	if !P(15, 15).IntersectsRect(R(10, 10, 20, 20)) {
		t.Fatal("expected true")
	}
	if !P(15, 15).IntersectsRing(R(10, 10, 20, 20).Ring()) {
		t.Fatal("expected true")
	}
	if P(9, 15).IntersectsRing(R(10, 10, 20, 20).Ring()) {
		t.Fatal("expected false")
	}
	if P(9, 15).InsideLine(Line{}) {
		t.Fatal("expected false")
	}
	if !P(9, 15).InsideLine(Line{P(9, 15)}) {
		t.Fatal("expected true")
	}
	if !P(9, 15).InsideRing(Ring{P(9, 15)}) {
		t.Fatal("expected true")
	}
	if P(9, 15).InsideRing(Ring{}) {
		t.Fatal("expected false")
	}

}
