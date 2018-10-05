package geom

import (
	"testing"
)

func TestSegmentRaycastPoint(t *testing.T) {

	// // fmt.Printf("%+v\n", S(0, 0, 0, 1).Raycast(P(-2, 0.5)))
	// // fmt.Printf("%+v\n", S(0, 0, 0, 1).Raycast(P(2, 0.5)))
	// // fmt.Printf("%+v\n", S(0, 1, 0, 0).Raycast(P(-2, 0.5)))
	// // fmt.Printf("%+v\n", S(0, 1, 0, 0).Raycast(P(2, 0.5)))

	// // fmt.Printf("%+v\n", S(0, 0, 1, 1).Raycast(P(-2, 0.5)))
	// // fmt.Printf("%+v\n", S(0, 0, 1, 1).Raycast(P(2, 0.5)))
	// // fmt.Printf("%+v\n", S(1, 1, 0, 0).Raycast(P(-2, 0.5)))
	// // fmt.Printf("%+v\n", S(1, 1, 0, 0).Raycast(P(2, 0.5)))

	// // return
	// res := func(in, on bool) RaycastResult {
	// 	return RaycastResult{In: in, On: on}
	// }
	// t.Run("Vertical", func(t *testing.T) {
	// 	// vertial segment

	// 	seg := S(0, 0, 0, 1)
	// 	for i := 0; i < 2; i++ {
	// 		var name string
	// 		if i == 0 {
	// 			name = "Forwards"
	// 		} else {
	// 			name = "Backwards"
	// 			seg.A, seg.B = seg.B, seg.A
	// 		}
	// 		t.Run(name, func(t *testing.T) {
	// 			t.Run("Collinear", func(t *testing.T) {
	// 				expect(t, seg.Raycast(P(0, -0.5)) == res(false, false))
	// 				expect(t, seg.Raycast(P(0, 0)) == res(false, true))
	// 				expect(t, seg.Raycast(P(0, 0.5)) == res(false, true))
	// 				expect(t, seg.Raycast(P(0, 1)) == res(false, true))
	// 				expect(t, seg.Raycast(P(0, 1.5)) == res(false, false))
	// 			})
	// 			t.Run("Left", func(t *testing.T) {
	// 				expect(t, seg.Raycast(P(-0.5, -0.5)) == res(false, false))
	// 				expect(t, seg.Raycast(P(-0.5, 0)) == res(true, false))
	// 				expect(t, seg.Raycast(P(-0.5, 0.5)) == res(true, false))
	// 				expect(t, seg.Raycast(P(-0.5, 1)) == res(true, false))
	// 				expect(t, seg.Raycast(P(-0.5, 1.5)) == res(false, false))
	// 			})
	// 			t.Run("Right", func(t *testing.T) {
	// 				expect(t, seg.Raycast(P(0.5, -0.5)) == res(false, false))
	// 				expect(t, seg.Raycast(P(0.5, 0)) == res(false, false))
	// 				expect(t, seg.Raycast(P(0.5, 0.5)) == res(false, false))
	// 				expect(t, seg.Raycast(P(0.5, 1)) == res(false, false))
	// 				expect(t, seg.Raycast(P(0.5, 1.5)) == res(false, false))
	// 			})
	// 		})
	// 	}
	// })
	// // expect(t, S(1, 4, 0, 3).Raycast(P(4, 3.5)) == RaycastResult{true, false})
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
