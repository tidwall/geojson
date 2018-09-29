package geom

import (
	"reflect"
	"testing"
)

func TestABC(t *testing.T) {
	series2 := NewSeries(ri, true, true)
	_ = series2
}

func TestSeries(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		series := NewSeries(octagon, true, true)
		expect(t, reflect.DeepEqual(series.Points(), octagon))
		expect(t, series.Convex())
		expect(t, series.Rect() == R(0, 0, 10, 10))
		expect(t, series.Closed())
		series = NewSeries(octagon, false, true)
		expect(t, reflect.DeepEqual(series.Points(), octagon))
	})
	t.Run("Search", func(t *testing.T) {
		series := NewSeries(octagon, true, true)
		var segs []Segment
		series.Search(R(0, 0, 0, 0), func(seg Segment, _ int) bool {
			segs = append(segs, seg)
			return true
		})
		segsExpect := []Segment{
			S(0, 3, 3, 0),
		}
		expect(t, checkSegsDups(segsExpect, segs))
		segs = nil
		series.Search(R(0, 0, 0, 10), func(seg Segment, _ int) bool {
			segs = append(segs, seg)
			return true
		})
		segsExpect = []Segment{
			S(3, 10, 0, 7),
			S(0, 7, 0, 3),
			S(0, 3, 3, 0),
		}
		expect(t, checkSegsDups(segsExpect, segs))
		segs = nil
		series.Search(R(0, 0, 5, 10), func(seg Segment, _ int) bool {
			segs = append(segs, seg)
			return true
		})
		segsExpect = []Segment{
			S(3, 0, 7, 0),
			S(7, 10, 3, 10),
			S(3, 10, 0, 7),
			S(0, 7, 0, 3),
			S(0, 3, 3, 0),
		}
		expect(t, checkSegsDups(segsExpect, segs))
		var seg2sA []Segment
		series.Search(R(0, 0, 10, 10), func(seg Segment, idx int) bool {
			seg2sA = append(seg2sA, seg)
			return true
		})

		var seg2sB []Segment
		series.Scan(func(seg Segment, idx int) bool {
			seg2sB = append(seg2sB, seg)
			return true
		})

		expect(t, checkSegsDups(seg2sA, seg2sB))

		var first Segment
		var once bool
		series.Search(R(0, 0, 10, 10), func(seg Segment, idx int) bool {
			expect(t, !once)
			first = seg
			once = true
			return false
		})
		expect(t, first == seg2sA[0])

	})
	t.Run("Big", func(t *testing.T) {
		t.Run("Closed", func(t *testing.T) {

			// clip off the last point to force an auto closure
			series := NewSeries(ri[:len(ri)-1], true, true)
			var seg2sA []Segment
			series.Search(series.Rect(), func(seg Segment, idx int) bool {
				seg2sA = append(seg2sA, seg)
				return true
			})
			var seg2sB []Segment
			series.Scan(func(seg Segment, idx int) bool {
				seg2sB = append(seg2sB, seg)
				return true
			})

			expect(t, checkSegsDups(seg2sA, seg2sB))

			// use all points
			series2 := NewSeries(ri, true, true)
			var seg2sC []Segment
			series2.Scan(func(seg Segment, idx int) bool {
				seg2sC = append(seg2sC, seg)
				return true
			})
			expect(t, checkSegsDups(seg2sA, seg2sC))

			var first Segment
			var once bool
			series.Search(series.Rect(), func(seg Segment, idx int) bool {
				expect(t, !once)
				first = seg
				once = true
				return false
			})
			expect(t, first == seg2sA[0])
		})
		t.Run("Opened", func(t *testing.T) {
			series := NewSeries(az, true, false)
			var seg2sA []Segment
			series.Search(series.Rect(), func(seg Segment, idx int) bool {
				seg2sA = append(seg2sA, seg)
				return true
			})
			var seg2sB []Segment
			series.Scan(func(seg Segment, idx int) bool {
				seg2sB = append(seg2sB, seg)
				return true
			})
			expect(t, checkSegsDups(seg2sA, seg2sB))
		})

	})
	t.Run("Reverse", func(t *testing.T) {
		shapes := [][]Point{octagon, concave1, concave2, concave3, concave4}
		for _, shape := range shapes {
			var rev []Point
			for i := len(shape) - 1; i >= 0; i-- {
				rev = append(rev, shape[i])
			}
			series := NewSeries(rev, true, true)
			var seg2sA []Segment
			series.Search(series.Rect(), func(seg Segment, idx int) bool {
				seg2sA = append(seg2sA, seg)
				return true
			})
			var seg2sB []Segment
			series.Scan(func(seg Segment, idx int) bool {
				seg2sB = append(seg2sB, seg)
				return true
			})
		}
	})
}

func checkSegsDups(a, b []Segment) bool {
	if len(a) != len(b) {
		return false
	}
	for _, a := range a {
		var found bool
		for _, b := range b {
			if a == b {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
