package poly

import "testing"

func testRayInside(t *testing.T, p Point, ps []Point, expect bool) {
	res := insideshpext(p, ps, true)
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

var texterior = Polygon{
	P(0, 0),
	P(0, 6),
	P(12, -6),
	P(12, 0),
	P(0, 0),
}
var tholeA = Polygon{
	P(1, 1),
	P(1, 2),
	P(2, 2),
	P(2, 1),
}
var tholeB = Polygon{
	P(11, -1),
	P(11, -3),
	P(9, -1),
}
var tholes = []Polygon{tholeA, tholeB}

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
		ok := points[i].p.Inside(texterior, tholes)
		if ok != points[i].ok {
			t.Fatalf("{%v,%v} = %t, expect %t", points[i].p.X, points[i].p.Y, ok, points[i].ok)
		}
	}
}

func TestInsideShapes(t *testing.T) {
	if texterior.Inside(texterior, nil) == false {
		t.Fatalf("expect true, got false")
	}
	if texterior.Inside(texterior, tholes) == true {
		t.Fatalf("expect false, got true")
	}
	if tholeA.Inside(texterior, nil) == false {
		t.Fatalf("expect true, got false")
	}
	if tholeB.Inside(texterior, nil) == false {
		t.Fatalf("expect true, got false")
	}
	if tholeA.Inside(tholeB, nil) == true {
		t.Fatalf("expect false, got true")
	}
}

func TestRectInsidePolygon(t *testing.T) {
	r1 := R(10, 10, 20, 20)
	r2 := R(0, 0, 30, 30)
	if !r1.Inside(r2.Polygon(), nil) {
		t.Fatalf("expected 'true'")
	}
	r3 := R(40, 40, 50, 50)
	if r1.Inside(r3.Polygon(), nil) {
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
	if !r1.Polygon().InsideRect(r2) {
		t.Fatalf("expected true")
	}
	r3 := R(40, 40, 50, 50)
	if r1.Polygon().InsideRect(r3) {
		t.Fatalf("expected false")
	}
	if (Polygon{}).InsideRect(r3) {
		t.Fatalf("expected false")
	}
}

func TestPolygonIntersectsRect(t *testing.T) {
	r1 := R(10, 10, 20, 20)
	r2 := R(0, 0, 30, 30)
	if !r1.Polygon().IntersectsRect(r2) {
		t.Fatalf("expected true")
	}
	r3 := R(40, 40, 50, 50)
	if r1.Polygon().IntersectsRect(r3) {
		t.Fatalf("expected false")
	}
	if (Polygon{}).IntersectsRect(r3) {
		t.Fatalf("expected false")
	}
}

func TestPolygonString(t *testing.T) {
	str := R(10, 10, 20, 20).Polygon().String()
	exp := "[[10,10],[20,10],[20,20],[10,20],[10,10]]"
	if str != exp {
		t.Fatalf("expected '%v', got '%v'", exp, str)
	}
}

func TestPolygonRect(t *testing.T) {
	p := Polygon{
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
