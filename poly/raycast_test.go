package poly

import "testing"

func testRayRay(t *testing.T, p, a, b Point, expect rayres) {
	res := raycast(p, a, b)
	if res != expect {
		t.Fatalf("1) %v,%v,%v = %s, expect %s", p, a, b, res, expect)
	}
	res = raycast(p, a, b)
	if res != expect {
		t.Fatalf("1) %v,%v,%v = %s, expect %s", p, b, a, res, expect)
	}
}

func TestRayHorizontal(t *testing.T) {
	for x := float64(-1); x <= 4+1; x++ {
		expect := on
		if x < 0 {
			expect = left
		} else if x > 4 {
			expect = out
		}
		testRayRay(t, P(x, 0), P(0, 0), P(4, 0), expect)
	}
	for x := float64(-1); x <= 4+1; x++ {
		expect := out
		testRayRay(t, P(x, -1), P(0, 0), P(4, 0), expect)
	}
	for x := float64(-1); x <= 4+1; x++ {
		expect := out
		testRayRay(t, P(x, +1), P(0, 0), P(4, 0), expect)
	}
}

func TestRayVertical(t *testing.T) {
	for y := float64(-1); y <= 4+1; y++ {
		expect := on
		if y < 0 || y > 4 {
			expect = out
		}
		testRayRay(t, P(0, y), P(0, 0), P(0, 4), expect)
	}
	for y := float64(-1); y <= 4+1; y++ {
		expect := left
		if y < 0 || y > 4 {
			expect = out
		}
		testRayRay(t, P(-1, y), P(0, 0), P(0, 4), expect)
	}
	for y := float64(-1); y <= 4+1; y++ {
		expect := out
		testRayRay(t, P(+1, y), P(0, 0), P(0, 4), expect)
	}
}

func TestRayAngled(t *testing.T) {
	testRayRay(t, P(1, 3), P(0, 4), P(4, 0), on)
	testRayRay(t, P(0, 4), P(0, 4), P(4, 0), on)
	testRayRay(t, P(4, 0), P(0, 4), P(4, 0), on)
	testRayRay(t, P(1, 4), P(0, 4), P(4, 0), out)
	testRayRay(t, P(5, 0), P(0, 4), P(4, 0), out)
	testRayRay(t, P(-1, 4), P(0, 4), P(4, 0), left)
	testRayRay(t, P(3, 0), P(0, 4), P(4, 0), left)
	testRayRay(t, P(2, 2), P(0, 4), P(4, 0), on)
}
