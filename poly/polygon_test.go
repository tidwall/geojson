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
