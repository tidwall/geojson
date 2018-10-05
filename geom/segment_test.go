package geom

import (
	"testing"
)

func TestSegmentRaycastPoint(t *testing.T) {
	rr := func(in, on bool) RaycastResult {
		return RaycastResult{In: in, On: on}
	}
	t.Run("Angle", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			t.Run("LRBT", func(t *testing.T) {
				// angled segment from left to right, bottom to top
				s := S(0, 3, 1, 4)
				t.Run("CrossoverAbove", func(t *testing.T) {
					t.Run("Left", func(t *testing.T) {
						expect(t, s.Raycast(P(-0.1, 4.1)) == rr(false, false))
						expect(t, s.Raycast(P(0.0, 4.1)) == rr(false, false))
						expect(t, s.Raycast(P(0.1, 4.1)) == rr(false, false))
					})
					t.Run("Center", func(t *testing.T) {
						expect(t, s.Raycast(P(0.4, 4.1)) == rr(false, false))
						expect(t, s.Raycast(P(0.5, 4.1)) == rr(false, false))
						expect(t, s.Raycast(P(0.6, 4.1)) == rr(false, false))
					})
					t.Run("Right", func(t *testing.T) {
						expect(t, s.Raycast(P(0.9, 4.1)) == rr(false, false))
						expect(t, s.Raycast(P(1.0, 4.1)) == rr(false, false))
						expect(t, s.Raycast(P(1.1, 4.1)) == rr(false, false))
					})
				})
				t.Run("CrossoverTop", func(t *testing.T) {
					t.Run("Left", func(t *testing.T) {
						expect(t, s.Raycast(P(-0.1, 4)) == rr(false, false))
						expect(t, s.Raycast(P(0.0, 4)) == rr(false, false))
						expect(t, s.Raycast(P(0.1, 4)) == rr(false, false))
					})
					t.Run("Center", func(t *testing.T) {
						expect(t, s.Raycast(P(0.4, 4)) == rr(false, false))
						expect(t, s.Raycast(P(0.5, 4)) == rr(false, false))
						expect(t, s.Raycast(P(0.6, 4)) == rr(false, false))
					})
					t.Run("Right", func(t *testing.T) {
						expect(t, s.Raycast(P(0.9, 4)) == rr(false, false))
						expect(t, s.Raycast(P(1.0, 4)) == rr(false, true))
						expect(t, s.Raycast(P(1.1, 4)) == rr(false, false))
					})
				})
				t.Run("CrossoverMiddle", func(t *testing.T) {
					t.Run("Left", func(t *testing.T) {
						expect(t, s.Raycast(P(-0.1, 3.5)) == rr(true, false))
						expect(t, s.Raycast(P(0.0, 3.5)) == rr(true, false))
						expect(t, s.Raycast(P(0.1, 3.5)) == rr(true, false))
					})
					t.Run("Center", func(t *testing.T) {
						expect(t, s.Raycast(P(0.4, 3.5)) == rr(true, false))
						expect(t, s.Raycast(P(0.5, 3.5)) == rr(false, true))
						expect(t, s.Raycast(P(0.6, 3.5)) == rr(false, false))
					})
					t.Run("Right", func(t *testing.T) {
						expect(t, s.Raycast(P(0.9, 3.5)) == rr(false, false))
						expect(t, s.Raycast(P(1.0, 3.5)) == rr(false, false))
						expect(t, s.Raycast(P(1.1, 3.5)) == rr(false, false))
					})
				})
				t.Run("CrossoverBottom", func(t *testing.T) {
					t.Run("Left", func(t *testing.T) {
						expect(t, s.Raycast(P(-0.1, 3)) == rr(true, false))
						expect(t, s.Raycast(P(0.0, 3)) == rr(false, true))
						expect(t, s.Raycast(P(0.1, 3)) == rr(false, false))
					})
					t.Run("Center", func(t *testing.T) {
						expect(t, s.Raycast(P(0.4, 3)) == rr(false, false))
						expect(t, s.Raycast(P(0.5, 3)) == rr(false, false))
						expect(t, s.Raycast(P(0.6, 3)) == rr(false, false))
					})
					t.Run("Right", func(t *testing.T) {
						expect(t, s.Raycast(P(0.9, 3)) == rr(false, false))
						expect(t, s.Raycast(P(1.0, 3)) == rr(false, false))
						expect(t, s.Raycast(P(1.1, 3)) == rr(false, false))
					})
				})
				t.Run("CrossoverBelow", func(t *testing.T) {
					t.Run("Left", func(t *testing.T) {
						expect(t, s.Raycast(P(-0.1, 2.9)) == rr(false, false))
						expect(t, s.Raycast(P(0.0, 2.9)) == rr(false, false))
						expect(t, s.Raycast(P(0.1, 2.9)) == rr(false, false))
					})
					t.Run("Center", func(t *testing.T) {
						expect(t, s.Raycast(P(0.4, 2.9)) == rr(false, false))
						expect(t, s.Raycast(P(0.5, 2.9)) == rr(false, false))
						expect(t, s.Raycast(P(0.6, 2.9)) == rr(false, false))
					})
					t.Run("Right", func(t *testing.T) {
						expect(t, s.Raycast(P(0.9, 2.9)) == rr(false, false))
						expect(t, s.Raycast(P(1.0, 2.9)) == rr(false, false))
						expect(t, s.Raycast(P(1.1, 2.9)) == rr(false, false))
					})
				})
			})
			t.Run("RLBT", func(t *testing.T) {
				// angled segment from right to left, bottom to top
				s := S(1, 4, 0, 3)
				t.Run("CrossoverAbove", func(t *testing.T) {
					t.Run("Left", func(t *testing.T) {
						expect(t, s.Raycast(P(-0.1, 4.1)) == rr(false, false))
						expect(t, s.Raycast(P(0.0, 4.1)) == rr(false, false))
						expect(t, s.Raycast(P(0.1, 4.1)) == rr(false, false))
					})
					t.Run("Center", func(t *testing.T) {
						expect(t, s.Raycast(P(0.4, 4.1)) == rr(false, false))
						expect(t, s.Raycast(P(0.5, 4.1)) == rr(false, false))
						expect(t, s.Raycast(P(0.6, 4.1)) == rr(false, false))
					})
					t.Run("Right", func(t *testing.T) {
						expect(t, s.Raycast(P(0.9, 4.1)) == rr(false, false))
						expect(t, s.Raycast(P(1.0, 4.1)) == rr(false, false))
						expect(t, s.Raycast(P(1.1, 4.1)) == rr(false, false))
					})
				})
				t.Run("CrossoverTop", func(t *testing.T) {
					t.Run("Left", func(t *testing.T) {
						expect(t, s.Raycast(P(-0.1, 4)) == rr(false, false))
						expect(t, s.Raycast(P(0.0, 4)) == rr(false, false))
						expect(t, s.Raycast(P(0.1, 4)) == rr(false, false))
					})
					t.Run("Center", func(t *testing.T) {
						expect(t, s.Raycast(P(0.4, 4)) == rr(false, false))
						expect(t, s.Raycast(P(0.5, 4)) == rr(false, false))
						expect(t, s.Raycast(P(0.6, 4)) == rr(false, false))
					})
					t.Run("Right", func(t *testing.T) {
						expect(t, s.Raycast(P(0.9, 4)) == rr(false, false))
						expect(t, s.Raycast(P(1.0, 4)) == rr(false, true))
						expect(t, s.Raycast(P(1.1, 4)) == rr(false, false))
					})
				})
				t.Run("CrossoverMiddle", func(t *testing.T) {
					t.Run("Left", func(t *testing.T) {
						expect(t, s.Raycast(P(-0.1, 3.5)) == rr(true, false))
						expect(t, s.Raycast(P(0.0, 3.5)) == rr(true, false))
						expect(t, s.Raycast(P(0.1, 3.5)) == rr(true, false))
					})
					t.Run("Center", func(t *testing.T) {
						expect(t, s.Raycast(P(0.4, 3.5)) == rr(true, false))
						expect(t, s.Raycast(P(0.5, 3.5)) == rr(false, true))
						expect(t, s.Raycast(P(0.6, 3.5)) == rr(false, false))
					})
					t.Run("Right", func(t *testing.T) {
						expect(t, s.Raycast(P(0.9, 3.5)) == rr(false, false))
						expect(t, s.Raycast(P(1.0, 3.5)) == rr(false, false))
						expect(t, s.Raycast(P(1.1, 3.5)) == rr(false, false))
					})
				})
				t.Run("CrossoverBottom", func(t *testing.T) {
					t.Run("Left", func(t *testing.T) {
						expect(t, s.Raycast(P(-0.1, 3)) == rr(true, false))
						expect(t, s.Raycast(P(0.0, 3)) == rr(false, true))
						expect(t, s.Raycast(P(0.1, 3)) == rr(false, false))
					})
					t.Run("Center", func(t *testing.T) {
						expect(t, s.Raycast(P(0.4, 3)) == rr(false, false))
						expect(t, s.Raycast(P(0.5, 3)) == rr(false, false))
						expect(t, s.Raycast(P(0.6, 3)) == rr(false, false))
					})
					t.Run("Right", func(t *testing.T) {
						expect(t, s.Raycast(P(0.9, 3)) == rr(false, false))
						expect(t, s.Raycast(P(1.0, 3)) == rr(false, false))
						expect(t, s.Raycast(P(1.1, 3)) == rr(false, false))
					})
				})
				t.Run("CrossoverBelow", func(t *testing.T) {
					t.Run("Left", func(t *testing.T) {
						expect(t, s.Raycast(P(-0.1, 2.9)) == rr(false, false))
						expect(t, s.Raycast(P(0.0, 2.9)) == rr(false, false))
						expect(t, s.Raycast(P(0.1, 2.9)) == rr(false, false))
					})
					t.Run("Center", func(t *testing.T) {
						expect(t, s.Raycast(P(0.4, 2.9)) == rr(false, false))
						expect(t, s.Raycast(P(0.5, 2.9)) == rr(false, false))
						expect(t, s.Raycast(P(0.6, 2.9)) == rr(false, false))
					})
					t.Run("Right", func(t *testing.T) {
						expect(t, s.Raycast(P(0.9, 2.9)) == rr(false, false))
						expect(t, s.Raycast(P(1.0, 2.9)) == rr(false, false))
						expect(t, s.Raycast(P(1.1, 2.9)) == rr(false, false))
					})
				})

			})
		})
		t.Run("2", func(t *testing.T) {
			t.Run("LRTB", func(t *testing.T) {
				// angled segment from left to right, top to bottom
				s := S(3, 4, 4, 3)
				t.Run("CrossoverMiddle", func(t *testing.T) {
					t.Run("Left", func(t *testing.T) {
						expect(t, s.Raycast(P(2.9, 3.5)) == rr(true, false))
						expect(t, s.Raycast(P(3.0, 3.5)) == rr(true, false))
						expect(t, s.Raycast(P(3.1, 3.5)) == rr(true, false))
					})
					t.Run("Center", func(t *testing.T) {
						expect(t, s.Raycast(P(3.4, 3.5)) == rr(true, false))
						expect(t, s.Raycast(P(3.5, 3.5)) == rr(false, true))
						expect(t, s.Raycast(P(3.6, 3.5)) == rr(false, false))
					})
					t.Run("Right", func(t *testing.T) {
						expect(t, s.Raycast(P(3.9, 3.5)) == rr(false, false))
						expect(t, s.Raycast(P(4.0, 3.5)) == rr(false, false))
						expect(t, s.Raycast(P(4.1, 3.5)) == rr(false, false))
					})
				})
			})
			t.Run("LRBT", func(t *testing.T) {
				// angled segment from left to right, bottom to top
				s := S(4, 3, 3, 4)
				t.Run("CrossoverMiddle", func(t *testing.T) {
					t.Run("Left", func(t *testing.T) {
						expect(t, s.Raycast(P(2.9, 3.5)) == rr(true, false))
						expect(t, s.Raycast(P(3.0, 3.5)) == rr(true, false))
						expect(t, s.Raycast(P(3.1, 3.5)) == rr(true, false))
					})
					t.Run("Center", func(t *testing.T) {
						expect(t, s.Raycast(P(3.4, 3.5)) == rr(true, false))
						expect(t, s.Raycast(P(3.5, 3.5)) == rr(false, true))
						expect(t, s.Raycast(P(3.6, 3.5)) == rr(false, false))
					})
					t.Run("Right", func(t *testing.T) {
						expect(t, s.Raycast(P(3.9, 3.5)) == rr(false, false))
						expect(t, s.Raycast(P(4.0, 3.5)) == rr(false, false))
						expect(t, s.Raycast(P(4.1, 3.5)) == rr(false, false))
					})
				})
			})
		})
	})
}

func TestSegmentContainsPoint(t *testing.T) {
	expect(t, S(0, 0, 1, 1).ContainsPoint(P(0, 0)))
	expect(t, S(0, 0, 1, 1).ContainsPoint(P(0.5, 0.5)))
	expect(t, S(0, 0, 1, 1).ContainsPoint(P(1, 1)))
	expect(t, !S(0, 0, 1, 1).ContainsPoint(P(1.1, 1.1)))
	expect(t, !S(0, 0, 1, 1).ContainsPoint(P(0.5, 0.6)))
	expect(t, !S(0, 0, 1, 1).ContainsPoint(P(-0.1, -0.1)))
}

func TestSegmentCollinearPoint(t *testing.T) {
	expect(t, S(0, 0, 1, 1).CollinearPoint(P(-1, -1)))
	expect(t, S(0, 0, 1, 1).CollinearPoint(P(0.5, 0.5)))
	expect(t, S(0, 0, 1, 1).CollinearPoint(P(2, 2)))
	expect(t, S(1, 1, 0, 0).CollinearPoint(P(-1, -1)))
	expect(t, S(1, 1, 0, 0).CollinearPoint(P(0.5, 0.5)))
	expect(t, S(1, 1, 0, 0).CollinearPoint(P(2, 2)))
	expect(t, S(1, 0, 0, 1).CollinearPoint(P(2, -1)))
	expect(t, S(1, 0, 0, 1).CollinearPoint(P(0.5, 0.5)))
	expect(t, S(1, 0, 0, 1).CollinearPoint(P(-1, 2)))
	expect(t, S(0, 1, 1, 0).CollinearPoint(P(2, -1)))
	expect(t, S(0, 1, 1, 0).CollinearPoint(P(0.5, 0.5)))
	expect(t, S(0, 1, 1, 0).CollinearPoint(P(-1, 2)))
}

// func TestSegmentAngle(t *testing.T) {
// 	fmt.Printf("%v\n", math.Abs(math.Mod(S(0, 0, 1, 0).Angle(), math.Pi/2)))
// 	fmt.Printf("%v\n", math.Abs(math.Mod(S(1, 0, 0, 0).Angle(), math.Pi/2)))

// 	fmt.Printf("%v\n", S(0, 0, 0, 1).Angle()) //math.Abs(math.Mod(S(0, 0, 0, 1).Angle(), math.Pi/2)))
// 	fmt.Printf("%v\n", math.Abs(math.Mod(S(0, 1, 0, 0).Angle(), math.Pi/2)))
// 	fmt.Printf("%v\n", math.Abs(math.Mod(S(0, 0, 1, 1).Angle(), math.Pi/2)))
// 	fmt.Printf("%v\n", math.Abs(math.Mod(S(1, 1, 0, 0).Angle(), math.Pi/2)))
// 	fmt.Printf("%v\n", math.Abs(math.Mod(S(0, 1, 1, 0).Angle(), math.Pi/2)))
// 	fmt.Printf("%v\n", math.Abs(math.Mod(S(1, 0, 0, 1).Angle(), math.Pi/2)))
// }

func TestSegmentContainsSegment(t *testing.T) {
	expect(t, S(0, 0, 10, 10).ContainsSegment(S(0, 0, 10, 10)))
	expect(t, S(0, 0, 10, 10).ContainsSegment(S(2, 2, 10, 10)))
	expect(t, S(0, 0, 10, 10).ContainsSegment(S(2, 2, 8, 8)))
	expect(t, !S(0, 0, 10, 10).ContainsSegment(S(-1, -1, 8, 8)))
}

func TestSegmentIntersectsSegment(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		expect(t, S(0, 1, 1, 1).IntersectsSegment(S(0, 1, 1, 0)))
	})
	t.Run("2", func(t *testing.T) {
		expect(t, S(1, 1, 0, 1).IntersectsSegment(S(0, 1, 1, 0)))
	})
	return

	// expect(t, S(0, 0, 0, 1).IntersectsSegment(S(0, 1, 1, 1)))
	// expect(t, S(0, -1, 0, 1).IntersectsSegment(S(0, 1, 0, 1)))
	expect(t, S(0, 0, 1, 1).IntersectsSegment(S(1, 1, 0, 2)))

	// 'x'
	ab := S(0, 0, 1, 1)
	cd := S(0, 1, 1, 0)

	return
	expect(t, ab.IntersectsSegment(cd))
	// move AB diagonally
	expect(t, ab.Move(0.25, 0.25).IntersectsSegment(cd))
	expect(t, ab.Move(0.50, 0.50).IntersectsSegment(cd))
	expect(t, !ab.Move(0.75, 0.75).IntersectsSegment(cd))
	expect(t, !ab.Move(1, 1).IntersectsSegment(cd))
	expect(t, !ab.Move(1.25, 1.25).IntersectsSegment(cd))
	expect(t, ab.Move(-0.25, -0.25).IntersectsSegment(cd))
	expect(t, ab.Move(-0.50, -0.50).IntersectsSegment(cd))
	expect(t, !ab.Move(-0.75, -0.75).IntersectsSegment(cd))
	expect(t, !ab.Move(-1, -1).IntersectsSegment(cd))
	expect(t, !ab.Move(-1.25, -1.25).IntersectsSegment(cd))

	// move AB vertically
	expect(t, ab.Move(0, 0.25).IntersectsSegment(cd))
	expect(t, ab.Move(0, 0.50).IntersectsSegment(cd))
	expect(t, ab.Move(0, 0.75).IntersectsSegment(cd))
	expect(t, ab.Move(0, 1).IntersectsSegment(cd))
	expect(t, !ab.Move(0, 1.25).IntersectsSegment(cd))
	expect(t, ab.Move(0, -0.25).IntersectsSegment(cd))
	expect(t, ab.Move(0, -0.50).IntersectsSegment(cd))
	expect(t, ab.Move(0, -0.75).IntersectsSegment(cd))
	expect(t, ab.Move(0, -1).IntersectsSegment(cd))
	expect(t, !ab.Move(0, -1.25).IntersectsSegment(cd))

	expect(t, cd.IntersectsSegment(ab))
	expect(t, cd.Move(0.25, 0.25).IntersectsSegment(ab))
	expect(t, cd.Move(0.50, 0.50).IntersectsSegment(ab))
	expect(t, !cd.Move(0.75, 0.75).IntersectsSegment(ab))
	expect(t, !cd.Move(1, 1).IntersectsSegment(ab))
	expect(t, !cd.Move(1.25, 1.25).IntersectsSegment(ab))
	expect(t, cd.Move(-0.25, -0.25).IntersectsSegment(ab))
	expect(t, cd.Move(-0.50, -0.50).IntersectsSegment(ab))
	expect(t, !cd.Move(-0.75, -0.75).IntersectsSegment(ab))
	expect(t, !cd.Move(-1, -1).IntersectsSegment(ab))
	expect(t, !cd.Move(-1.25, -1.25).IntersectsSegment(ab))

}
