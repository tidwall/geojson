package poly

import "testing"

func TestPolygon(t *testing.T) {
	poly := R(10, 10, 20, 20).Polygon()
	if poly.InsidePoint(P(15, 15)) {
		t.Fatal("expected false")
	}
	if (Polygon{}).InsidePoint(P(15, 15)) {
		t.Fatal("expected false")
	}
	if !R(15, 15, 15, 15).Polygon().InsidePoint(P(15, 15)) {
		t.Fatal("expected true")
	}
	if !R(15, 15, 15, 15).Polygon().InsideLine(Line{P(10, 15), P(20, 15)}) {
		t.Fatal("expected true")
	}
	if !R(15, 15, 15, 20).Polygon().InsideLine(Line{P(15, 10), P(15, 30)}) {
		t.Fatal("expected true")
	}
	if (Ring{}).InsideLine(Line{}) {
		t.Fatal("expected false")
	}
	if R(10, 10, 20, 20).Polygon().InsideRing(R(12, 12, 18, 18).Ring()) {
		t.Fatal("expected false")
	}
	if !R(10, 10, 20, 20).Polygon().InsideRing(R(8, 8, 22, 22).Ring()) {
		t.Fatal("expected true")
	}
	if !R(10, 10, 20, 20).Polygon().IntersectsLine(Line{P(15, 10), P(15, 30)}) {
		t.Fatal("expected true")
	}
	if !R(10, 10, 20, 20).Polygon().IntersectsPoint(P(15, 15)) {
		t.Fatal("expected true")
	}
	if R(10, 10, 20, 20).Polygon().IntersectsPoint(P(30, 30)) {
		t.Fatal("expected false")
	}
	if !R(10, 10, 20, 20).Polygon().IntersectsRing(R(0, 0, 15, 15).Ring()) {
		t.Fatal("expected true")
	}
	if !poly.InsideRect(R(9, 9, 21, 21)) {
		t.Fatal("expected true")
	}

	exterior := R(10, 10, 20, 20).Ring()
	//holes := []Ring{R(13, 13, 17, 17).Ring()}

	if !(Polygon{exterior, nil}).InsidePolygon(Polygon{exterior, nil}) {
		t.Fatal("expected false")
	}
	if !(Polygon{exterior, nil}).IntersectsRect(R(10, 10, 20, 20)) {
		t.Fatal("expected false")
	}
	if !(Polygon{exterior, nil}).IntersectsPolygon(Polygon{exterior, nil}) {
		t.Fatal("expected false")
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
