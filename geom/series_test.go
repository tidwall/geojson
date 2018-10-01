package geom

import (
	"reflect"
	"testing"
)

func testScanSeries(
	t *testing.T, series *Series, index bool,
	expectSegmentCount int, expectedEmpty bool,
) []Segment {
	t.Helper()
	if index {
		series.buildTree()
	} else {
		series.tree = nil
	}
	var segs1 []Segment
	lastIdx := -1
	series.ForEachSegment(func(seg Segment, idx int) bool {
		expect(t, idx == lastIdx+1)
		segs1 = append(segs1, seg)
		lastIdx = idx
		return true
	})
	var segs2Count int
	segs2 := make([]Segment, len(segs1))
	series.Search(series.Rect(), func(seg Segment, idx int) bool {
		segs2[idx] = seg
		segs2Count++
		return true
	})
	expect(t, segs2Count == len(segs2))
	expect(t, len(segs1) == len(segs2))
	if len(segs1) != 0 {
		expect(t, reflect.DeepEqual(segs1, segs2))
	}
	expect(t, expectSegmentCount == len(segs1))
	expect(t, series.Empty() == expectedEmpty)
	return segs1
}

func TestSeriesBasic(t *testing.T) {
	series := MakeSeries(octagon, true, true)
	expect(t, reflect.DeepEqual(ringCopyPoints(&series), octagon))
	expect(t, series.Convex())
	expect(t, series.Rect() == R(0, 0, 10, 10))
	expect(t, series.Closed())
	series = MakeSeries(octagon, false, true)
	expect(t, reflect.DeepEqual(ringCopyPoints(&series), octagon))

	series = MakeSeries(ri, true, true)
	testScanSeries(t, &series, true, len(ri)-1, false)
	testScanSeries(t, &series, false, len(ri)-1, false)

	// small lines
	series = MakeSeries([]Point{}, true, false)
	testScanSeries(t, &series, true, 0, true)
	testScanSeries(t, &series, false, 0, true)

	series = MakeSeries([]Point{P(5, 5)}, true, false)
	testScanSeries(t, &series, true, 0, true)
	testScanSeries(t, &series, false, 0, true)

	series = MakeSeries([]Point{P(5, 5), P(10, 10)}, true, false)
	testScanSeries(t, &series, true, 1, false)
	testScanSeries(t, &series, false, 1, false)

	// small rings
	series = MakeSeries([]Point{}, true, true)
	testScanSeries(t, &series, true, 0, true)
	testScanSeries(t, &series, false, 0, true)

	series = MakeSeries([]Point{P(5, 5)}, true, true)
	testScanSeries(t, &series, true, 0, true)
	testScanSeries(t, &series, false, 0, true)

	series = MakeSeries([]Point{P(5, 5), P(10, 10)}, true, true)
	testScanSeries(t, &series, true, 0, true)
	testScanSeries(t, &series, false, 0, true)

	series = MakeSeries([]Point{P(5, 5), P(10, 10), P(10, 5)}, true, true)
	testScanSeries(t, &series, true, 3, false)
	testScanSeries(t, &series, false, 3, false)

}
func TestSeriesSearch(t *testing.T) {
	series := MakeSeries(octagon, true, true)
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
	series.ForEachSegment(func(seg Segment, idx int) bool {
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

}
func TestSeriesBig(t *testing.T) {
	t.Run("Closed", func(t *testing.T) {
		// clip off the last point to force an auto closure
		series := MakeSeries(ri[:len(ri)-1], true, true)
		var seg2sA []Segment
		series.Search(series.Rect(), func(seg Segment, idx int) bool {
			seg2sA = append(seg2sA, seg)
			return true
		})
		var seg2sB []Segment
		series.ForEachSegment(func(seg Segment, idx int) bool {
			seg2sB = append(seg2sB, seg)
			return true
		})

		expect(t, checkSegsDups(seg2sA, seg2sB))

		// use all points
		series2 := MakeSeries(ri, true, true)
		var seg2sC []Segment
		series2.ForEachSegment(func(seg Segment, idx int) bool {
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
		series := MakeSeries(az, true, false)
		var seg2sA []Segment
		series.Search(series.Rect(), func(seg Segment, idx int) bool {
			seg2sA = append(seg2sA, seg)
			return true
		})
		var seg2sB []Segment
		series.ForEachSegment(func(seg Segment, idx int) bool {
			seg2sB = append(seg2sB, seg)
			return true
		})
		expect(t, checkSegsDups(seg2sA, seg2sB))
	})
}
func TestSeriesReverse(t *testing.T) {
	shapes := [][]Point{octagon, concave1, concave2, concave3, concave4}
	for _, shape := range shapes {
		var rev []Point
		for i := len(shape) - 1; i >= 0; i-- {
			rev = append(rev, shape[i])
		}
		series := MakeSeries(rev, true, true)
		var seg2sA []Segment
		series.Search(series.Rect(), func(seg Segment, idx int) bool {
			seg2sA = append(seg2sA, seg)
			return true
		})
		var seg2sB []Segment
		series.ForEachSegment(func(seg Segment, idx int) bool {
			seg2sB = append(seg2sB, seg)
			return true
		})
	}
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
