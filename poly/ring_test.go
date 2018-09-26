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

func TestIssue362(t *testing.T) {
	rect := Rect{
		Point{-122.4434208869934, 37.73138728181471},
		Point{-122.43711233139038, 37.73579951265516},
	}
	box := Ring{
		{-122.4434208869934, 37.73138728181471},
		{-122.43711233139038, 37.73138728181471},
		{-122.43711233139038, 37.73579951265516},
		{-122.4434208869934, 37.73579951265516},
		{-122.4434208869934, 37.73138728181471},
	}
	shape := Ring{
		{-122.4475622177124, 37.73590133026304},
		{-122.44369983673094, 37.72904529863455},
		{-122.44052410125731, 37.732778859926555},
		{-122.43713378906249, 37.729079240948714},
		{-122.43292808532715, 37.73865035274667},
		{-122.43962287902832, 37.735154664553534},
		{-122.44850635528563, 37.73885397998102},
		{-122.4475622177124, 37.7359013302630},
	}
	if box.InsideRing(shape) {
		t.Fatal("expected false")
	}
	if rect.InsideRing(shape) {
		t.Fatal("expected false")
	}
}
