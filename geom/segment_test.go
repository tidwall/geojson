package geom

import (
	"fmt"
	"strings"
	"testing"
)

func TestSegmentRaycastPoint(t *testing.T) {
	// This is full coverage raycast test. It uses a 7x7 grid where the each
	// point is checked for a total of 49 tests per segment. There are 16
	// segments at 0º,30º,45º,90º,180º in both directions for a total of 784
	// tests.
	segs := []interface{}{
		S(1, 1, 5, 5), "A",
		S(5, 5, 1, 1), "B",
		[7]string{
			"-------",
			"-----*-",
			"++++*--",
			"+++*---",
			"++*----",
			"+*-----",
			"-------",
		},
		S(1, 5, 5, 1), "C",
		S(5, 1, 1, 5), "D",
		[7]string{
			"-------",
			"-*-----",
			"++*----",
			"+++*---",
			"++++*--",
			"+++++*-",
			"-------",
		},
		S(1, 3, 5, 3), "E",
		S(5, 3, 1, 3), "F",
		[7]string{
			"-------",
			"-------",
			"-------",
			"-*****-",
			"-------",
			"-------",
			"-------",
		},
		S(3, 5, 3, 1), "G",
		S(3, 1, 3, 5), "H",
		[7]string{
			"-------",
			"---*---",
			"+++*---",
			"+++*---",
			"+++*---",
			"+++*---",
			"-------",
		},
		S(1, 2, 5, 4), "I",
		S(5, 4, 1, 2), "J",
		[7]string{
			"-------",
			"-------",
			"-----*-",
			"+++*---",
			"+*-----",
			"-------",
			"-------",
		},
		S(1, 4, 5, 2), "K",
		S(5, 2, 1, 4), "L",
		[7]string{
			"-------",
			"-------",
			"-*-----",
			"+++*---",
			"+++++*-",
			"-------",
			"-------",
		},
		S(2, 1, 4, 5), "M",
		S(4, 5, 2, 1), "N",
		[7]string{
			"-------",
			"----*--",
			"++++---",
			"+++*---",
			"+++----",
			"++*----",
			"-------",
		},
		S(2, 5, 4, 1), "O",
		S(4, 1, 2, 5), "P",
		[7]string{
			"-------",
			"--*----",
			"+++----",
			"+++*---",
			"++++---",
			"++++*--",
			"-------",
		},
		S(3, 3, 3, 3), "Q",
		S(3, 3, 3, 3), "R",
		[7]string{
			"-------",
			"-------",
			"-------",
			"---*---",
			"-------",
			"-------",
			"-------",
		},
	}

	var ms string
	for i := 0; i < len(segs); i += 5 {
		segs := []interface{}{
			segs[i], segs[i+1], segs[i+4],
			segs[i+2], segs[i+3], segs[i+4],
		}
		for i := 0; i < len(segs); i += 3 {
			seg := segs[i].(Segment)
			label := segs[i+1].(string)
			grid := segs[i+2].([7]string)
			//
			var ngrid [7]string
			for y, sy := 0, 6; y < 7; y, sy = y+1, sy-1 {
				var nline string
				for x := 0; x < 7; x++ {
					// ch := grid[sy][x]
					pt := Point{float64(x), float64(y)}
					res := seg.Raycast(pt)
					if res.In {
						nline += "+"
					} else if res.On {
						nline += "*"
					} else {
						nline += "-"
					}
				}
				ngrid[sy] = nline
			}
			if grid != ngrid {
				ms += fmt.Sprintf("MISMATCH (%s) SEGMENT: %v\n", label, seg)
				ms += fmt.Sprintf("EXPECTED\n%s\n", strings.Join(grid[:], "\n"))
				ms += fmt.Sprintf("GOT\n%s\n", strings.Join(ngrid[:], "\n"))
			}
		}
	}
	if ms != "" {
		t.Fatalf("\n%s", ms)
	}
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
