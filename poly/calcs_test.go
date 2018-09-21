package poly

import "testing"

func testRayInside(t *testing.T, p Point, ps []Point, expect bool) {
	res := pointInRing(p, ps, true)
	if res != expect {
		t.Fatalf("{%v,%v} = %t, expect %t", p.X, p.Y, res, expect)
	}
}

func TestRayInside(t *testing.T) {
	strange := []Point{P(0, 0), P(0, 3), P(4, -3), P(4, 0), P(0, 0)}

	// on the edge
	testRayInside(t, P(0, 0), strange, true)
	testRayInside(t, P(0, 3), strange, true)

	testRayInside(t, P(4, -3), strange, true)
	testRayInside(t, P(4, -2), strange, true)
	testRayInside(t, P(3, 0), strange, true)
	testRayInside(t, P(1, 0), strange, true)

	// ouside by just a tad
	testRayInside(t, P(-0.1, 0), strange, false)
	testRayInside(t, P(-0.1, -0.1), strange, false)
	testRayInside(t, P(0, 3.1), strange, false)
	testRayInside(t, P(0.1, 3.1), strange, false)
	testRayInside(t, P(-0.1, 3), strange, false)
	testRayInside(t, P(4, -3.1), strange, false)
	testRayInside(t, P(3.9, -3), strange, false)
	testRayInside(t, P(4.1, -2), strange, false)
	testRayInside(t, P(3, 0.1), strange, false)
	testRayInside(t, P(1, -0.1), strange, false)
}

func TestRayInside2(t *testing.T) {
	normal := []Point{P(0, 0), P(4, 3), P(5, 2), P(0, 0)}
	testRayInside(t, P(1, 2), normal, false)
	testRayInside(t, P(1, 3), normal, false)
	testRayInside(t, P(4, 2), normal, true)
	testRayInside(t, P(2, 1), normal, true)
}

var texterior = Ring{
	P(0, 0),
	P(0, 6),
	P(12, -6),
	P(12, 0),
	P(0, 0),
}
var tholeA = Ring{
	P(1, 1),
	P(1, 2),
	P(2, 2),
	P(2, 1),
}
var tholeB = Ring{
	P(11, -1),
	P(11, -3),
	P(9, -1),
}
var tholes = []Ring{tholeA, tholeB}

func TestRayExteriorHoles(t *testing.T) {

	type point struct {
		p  Point
		ok bool
	}

	points := []point{
		{P(.5, 3), true},
		{P(11.5, -4.5), true},
		{P(6, 0), true},

		{P(3.5, .5), true},
		{P(1.5, 1.5), false},
		{P(10.5, -1.5), false},
		{P(-2, 0), false},
		{P(0, -2), false},
		{P(8, -3), false},
		{P(8, 1), false},
		{P(14, -1), false},

		{P(8, -0.5), true},
		{P(8, -1.5), true},
		{P(8, -1), true},
	}
	// add the edges, all should be inside
	for i := 0; i < len(texterior); i++ {
		points = append(points, point{texterior[i], true})
	}
	for i := 0; i < len(tholeA); i++ {
		points = append(points, point{tholeA[i], true})
	}
	for i := 0; i < len(tholeB); i++ {
		points = append(points, point{tholeB[i], true})
	}

	for i := 0; i < len(points); i++ {
		ok := points[i].p.InsidePolygon(Polygon{texterior, tholes})
		if ok != points[i].ok {
			t.Fatalf("{%v,%v} = %t, expect %t", points[i].p.X, points[i].p.Y, ok, points[i].ok)
		}
	}
}

func TestInsideShapes(t *testing.T) {
	if texterior.InsidePolygon(Polygon{texterior, nil}) == false {
		t.Fatalf("expect true, got false")
	}
	if texterior.InsidePolygon(Polygon{texterior, tholes}) == true {
		t.Fatalf("expect false, got true")
	}
	if tholeA.InsidePolygon(Polygon{texterior, nil}) == false {
		t.Fatalf("expect true, got false")
	}
	if tholeB.InsidePolygon(Polygon{texterior, nil}) == false {
		t.Fatalf("expect true, got false")
	}
	if tholeA.InsidePolygon(Polygon{tholeB, nil}) == true {
		t.Fatalf("expect false, got true")
	}
}

func TestRectInsidePolygon(t *testing.T) {
	r1 := R(10, 10, 20, 20)
	r2 := R(0, 0, 30, 30)
	if !r1.InsidePolygon(r2.Polygon()) {
		t.Fatalf("expected 'true'")
	}
	r3 := R(40, 40, 50, 50)
	if r1.InsidePolygon(r3.Polygon()) {
		t.Fatalf("expected 'false'")
	}
}

func TestPointInsideRect(t *testing.T) {
	if !P(10, 10).InsideRect(R(0, 0, 20, 20)) {
		t.Fatalf("expected true")
	}
	if P(10, -1).InsideRect(R(0, 0, 20, 20)) {
		t.Fatalf("expected false")
	}
	if P(21, 10).InsideRect(R(0, 0, 20, 20)) {
		t.Fatalf("expected false")
	}
}

func TestPolygonInsideRect(t *testing.T) {
	r1 := R(10, 10, 20, 20)
	r2 := R(0, 0, 30, 30)
	if !r1.Ring().InsideRect(r2) {
		t.Fatalf("expected true")
	}
	r3 := R(40, 40, 50, 50)
	if r1.Ring().InsideRect(r3) {
		t.Fatalf("expected false")
	}
	if (Ring{}).InsideRect(r3) {
		t.Fatalf("expected false")
	}
}

func TestPolygonIntersectsRect(t *testing.T) {
	r1 := R(10, 10, 20, 20)
	r2 := R(0, 0, 30, 30)
	if !r1.Ring().IntersectsRect(r2) {
		t.Fatalf("expected true")
	}
	r3 := R(40, 40, 50, 50)
	if r1.Ring().IntersectsRect(r3) {
		t.Fatalf("expected false")
	}
	if (Ring{}).IntersectsRect(r3) {
		t.Fatalf("expected false")
	}
}

func TestPolygonString(t *testing.T) {
	str := R(10, 10, 20, 20).Ring().String()
	exp := "[[10,10],[20,10],[20,20],[10,20],[10,10]]"
	if str != exp {
		t.Fatalf("expected '%v', got '%v'", exp, str)
	}
}

func TestPolygonRect(t *testing.T) {
	p := Ring{
		P(10, 10), P(20, 10), P(30, 20), P(40, 0),
		P(50, 50), P(40, 30), P(30, 20), P(0, 0),
		P(10, 10),
	}
	r := p.Rect()
	exp := R(0, 0, 50, 50)
	if r != exp {
		t.Fatalf("expected '%v', got '%v'", exp, r)
	}
}

func TestInsideRect(t *testing.T) {
	if !R(10, 10, 20, 20).InsideRect(R(0, 0, 30, 30)) {
		t.Fatal("expected true")
	}
	if R(10, 10, 20, 20).InsideRect(R(20, 20, 30, 30)) {
		t.Fatal("expected false")
	}
	if R(10, 10, 20, 20).InsideRect(R(0, 0, 15, 15)) {
		t.Fatal("expected false")
	}
	if R(10, 10, 20, 20).InsideRect(R(0, 20, 30, 50)) {
		t.Fatal("expected false")
	}
}

func testIntersectsLinesA(t *testing.T, a, b, c, d Point, expect bool) {
	res := segmentsIntersect(a, b, c, d)
	if res != expect {
		t.Fatalf("{%v,%v}, {%v,%v} = %t, expect %t", a, b, c, d, res, expect)
	}
	res = segmentsIntersect(b, a, c, d)
	if res != expect {
		t.Fatalf("{%v,%v}, {%v,%v} = %t, expect %t", b, a, c, d, res, expect)
	}
	res = segmentsIntersect(a, b, d, c)
	if res != expect {
		t.Fatalf("{%v,%v}, {%v,%v} = %t, expect %t", a, b, d, c, res, expect)
	}
	res = segmentsIntersect(b, a, d, c)
	if res != expect {
		t.Fatalf("{%v,%v}, {%v,%v} = %t, expect %t", b, a, d, c, res, expect)
	}
}

func testIntersectsLines(t *testing.T, a, b, c, d Point, expect bool) {
	testIntersectsLinesA(t, a, b, c, d, expect)
	testIntersectsLinesA(t, c, d, a, b, expect)
}

func TestIntersectsLines(t *testing.T) {
	testIntersectsLines(t, P(0, 6), P(12, -6), P(0, 0), P(12, 0), true)
	testIntersectsLines(t, P(0, 0), P(5, 5), P(5, 5), P(0, 10), true)
	testIntersectsLines(t, P(0, 0), P(5, 5), P(5, 6), P(0, 10), false)
	testIntersectsLines(t, P(0, 0), P(5, 5), P(5, 4), P(0, 10), true)
	testIntersectsLines(t, P(0, 0), P(2, 2), P(0, 2), P(2, 0), true)
	testIntersectsLines(t, P(0, 0), P(2, 2), P(0, 2), P(1, 1), true)
	testIntersectsLines(t, P(0, 0), P(2, 2), P(2, 0), P(1, 1), true)
	testIntersectsLines(t, P(0, 0), P(0, 4), P(1, 4), P(4, 1), false)
	testIntersectsLines(t, P(0, 0), P(0, 4), P(1, 4), P(4, 4), false)
	testIntersectsLines(t, P(0, 0), P(0, 4), P(4, 1), P(4, 4), false)
	testIntersectsLines(t, P(0, 0), P(4, 0), P(1, 4), P(4, 1), false)
	testIntersectsLines(t, P(0, 0), P(4, 0), P(1, 4), P(4, 4), false)
	testIntersectsLines(t, P(0, 0), P(4, 0), P(4, 1), P(4, 4), false)
	testIntersectsLines(t, P(0, 4), P(4, 0), P(1, 4), P(4, 1), false)
	testIntersectsLines(t, P(0, 4), P(4, 0), P(1, 4), P(4, 4), false)
	testIntersectsLines(t, P(0, 4), P(4, 0), P(4, 1), P(4, 4), false)
}

func testIntersectsShapes(t *testing.T, exterior Ring, holes []Ring, shape Ring, expect bool) {
	got := shape.IntersectsPolygon(Polygon{exterior, holes})
	if got != expect {
		t.Fatalf("%v intersects %v = %v, expect %v", shape, exterior, got, expect)
	}
	got = exterior.IntersectsPolygon(Polygon{shape, nil})
	if got != expect {
		t.Fatalf("%v intersects %v = %v, expect %v", exterior, shape, got, expect)
	}
}

func TestIntersectsShapes(t *testing.T) {

	testIntersectsShapes(t,
		Ring{P(6, 0), P(12, 0), P(12, -6), P(6, 0)},
		nil,
		Ring{P(0, 0), P(0, 6), P(6, 0), P(0, 0)},
		true)

	testIntersectsShapes(t,
		Ring{P(7, 0), P(12, 0), P(12, -6), P(7, 0)},
		nil,
		Ring{P(0, 0), P(0, 6), P(6, 0), P(0, 0)},
		false)

	testIntersectsShapes(t,
		Ring{P(0.5, 0.5), P(0.5, 4.5), P(4.5, 0.5), P(0.5, 0.5)},
		nil,
		Ring{P(0, 0), P(0, 6), P(6, 0), P(0, 0)},
		true)

	testIntersectsShapes(t,
		Ring{P(0, 0), P(0, 6), P(6, 0), P(0, 0)},
		[]Ring{{P(1, 1), P(1, 2), P(2, 2), P(2, 1), P(1, 1)}},
		Ring{P(0.5, 0.5), P(0.5, 4.5), P(4.5, 0.5), P(0.5, 0.5)},
		true)

	testIntersectsShapes(t,
		Ring{P(0, 0), P(0, 10), P(10, 10), P(10, 0), P(0, 0)},
		[]Ring{{P(2, 2), P(2, 6), P(6, 6), P(6, 2), P(2, 2)}},
		Ring{P(1, 1), P(1, 9), P(9, 9), P(9, 1), P(1, 1)},
		true)
}

func TestPointIntersectsLine(t *testing.T) {
	poly := Line{P(0, 0), P(10, 10), P(20, 0)}
	if !P(0, 0).IntersectsLine(poly) {
		t.Fatal("expected true")
	}
	if !P(10, 10).IntersectsLine(poly) {
		t.Fatal("expected true")
	}
	if !P(20, 0).IntersectsLine(poly) {
		t.Fatal("expected true")
	}
	if !P(5, 5).IntersectsLine(poly) {
		t.Fatal("expected true")
	}
	if !P(15, 5).IntersectsLine(poly) {
		t.Fatal("expected true")
	}
	if P(20, 5).IntersectsLine(poly) {
		t.Fatal("expected false")
	}
}

func TestPointIntersects(t *testing.T) {
	poly := Ring{P(0, 0), P(10, 0), P(10, 10), P(0, 10), P(0, 0)}

	if !P(0, 0).IntersectsPolygon(Polygon{poly, nil}) {
		t.Fatal("expected true")
	}
	if !P(10, 0).IntersectsPolygon(Polygon{poly, nil}) {
		t.Fatal("expected true")
	}
	if !P(0, 10).IntersectsPolygon(Polygon{poly, nil}) {
		t.Fatal("expected true")
	}
	if !P(10, 10).IntersectsPolygon(Polygon{poly, nil}) {
		t.Fatal("expected true")
	}
}
func TestRectIntersectsPolygon(t *testing.T) {
	poly := Ring{P(0, 0), P(10, 0), P(10, 10), P(0, 10), P(0, 0)}
	if !R(0, 0, 5, 5).IntersectsPolygon(Polygon{poly, nil}) {
		t.Fatal("expected true")
	}
	if R(15, 15, 20, 20).IntersectsPolygon(Polygon{poly, nil}) {
		t.Fatal("expected true")
	}
}

func TestLineIntersectsLine(t *testing.T) {
	line1 := Line{P(0, 0), P(10, 10), P(20, 0)}
	line2 := Line{P(0, 1), P(10, 11), P(20, 1)}
	line3 := Line{P(0, -1), P(10, 11), P(20, -1)}
	if line1.IntersectsLine(line2) {
		t.Fatal("expected false")
	}
	if !line1.IntersectsLine(line3) {
		t.Fatal("expected true")
	}
}

func TestLineIntersects(t *testing.T) {
	poly := Ring{P(0, 0), P(10, 0), P(10, 10), P(0, 10), P(0, 0)}
	line1 := Line{P(0, 0), P(10, 10), P(20, 0)}

	if !line1.IntersectsPolygon(Polygon{poly, nil}) {
		t.Fatal("expected true")
	}
}

func TestDoesIntersect(t *testing.T) {
	exterior := Ring{P(0, 0), P(10, 0), P(10, 10), P(0, 10), P(0, 0)}
	holes := []Ring{{P(2, 2), P(2, 8), P(8, 8), P(8, 2), P(2, 2)}}

	if doesIntersect(nil, false, Polygon{exterior, holes}) {
		t.Fatal("expected false")
	}
	if doesIntersect(Ring{P(5, 5)}, false, Polygon{exterior, holes}) {
		t.Fatal("expected false")
	}
	if !doesIntersect(Ring{P(1, 1)}, false, Polygon{exterior, holes}) {
		t.Fatal("expected true")
	}
	if doesIntersect(Ring{P(1, 1)}, false, Polygon{Ring{}, nil}) {
		t.Fatal("expected false")
	}
	if !doesIntersect(Ring{P(1, 1)}, false, Polygon{Ring{P(1, 1)}, nil}) {
		t.Fatal("expected true")
	}
	if doesIntersect(Ring{P(1, 1), P(2, 2)}, false, Polygon{Ring{}, nil}) {
		t.Fatal("expected false")
	}
	if !doesIntersect(exterior, false, Polygon{Ring{P(1, 1)}, nil}) {
		t.Fatal("expected true")
	}

	inner := Ring{P(3, 3), P(7, 3), P(7, 7), P(3, 7), P(3, 3)}
	if doesIntersect(inner, false, Polygon{exterior, holes}) {
		t.Fatal("expected false")
	}

	outside := Ring{P(30, 30), P(60, 30), P(60, 60), P(30, 60), P(30, 30)}
	if doesIntersect(outside, false, Polygon{exterior, holes}) {
		t.Fatal("expected false")
	}

	// triangles

	tri1 := Ring{P(0, 0), P(10, 0), P(5, 10), P(0, 0)}
	tri2 := Ring{P(7, 9), P(17, 9), P(12, 19), P(7, 9)}

	if doesIntersect(tri1, false, Polygon{tri2, nil}) {
		t.Fatal("expected false")
	}

	// if !line1.doesIntersect(poly, nil) {
	// 	t.Fatal("expected true")
	// }
}

func TestLineIntersectsRect(t *testing.T) {
	if !(Line{P(0, 0), P(30, 30)}).IntersectsRect(R(10, 10, 20, 20)) {
		t.Fatal("expected true")
	}
	if (Line{P(100, 100), P(300, 300)}).IntersectsRect(R(10, 10, 20, 20)) {
		t.Fatal("expected false")
	}
}

func TestLineIntersectsPoint(t *testing.T) {
	if !(Line{P(0, 0), P(30, 30)}).IntersectsPoint(P(15, 15)) {
		t.Fatal("expected true")
	}
	if (Line{P(0, 0), P(30, 30)}).IntersectsPoint(P(15, 10)) {
		t.Fatal("expected false")
	}
}
