// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geometry

import (
	"reflect"
	"testing"
)

func seriesForEachSegment(ring Ring, iter func(seg Segment) bool) {
	n := ring.NumSegments()
	for i := 0; i < n; i++ {
		if !iter(ring.SegmentAt(i)) {
			return
		}
	}
}

func seriesForEachPoint(ring Ring, iter func(point Point) bool) {
	n := ring.NumPoints()
	for i := 0; i < n; i++ {
		if !iter(ring.PointAt(i)) {
			return
		}
	}
}

func TestSeriesEmpty(t *testing.T) {
	var series *baseSeries
	expect(t, series.Empty())
	series2 := makeSeries(nil, false, false, &IndexOptions{Kind: None})
	expect(t, series2.Empty())
}

func TestSeriesIndex(t *testing.T) {
	series := makeSeries(nil, false, false, &IndexOptions{Kind: None})
	expect(t, series.Index() == nil)
	series = makeSeries([]Point{
		P(0, 0), P(10, 0), P(10, 10), P(0, 10), P(0, 0),
	}, true, true, &IndexOptions{Kind: RTree, MinPoints: 1})

	expect(t, series.Index() != nil)
}

func TestSeriesClockwise(t *testing.T) {
	var series baseSeries
	series = makeSeries([]Point{
		P(0, 0), P(10, 0), P(10, 10), P(0, 10), P(0, 0),
	}, true, true, DefaultIndexOptions)
	expect(t, !series.Clockwise())
	series = makeSeries([]Point{
		P(0, 0), P(10, 0), P(10, 10), P(0, 10),
	}, true, true, DefaultIndexOptions)
	expect(t, !series.Clockwise())
	series = makeSeries([]Point{
		P(0, 0), P(10, 0), P(10, 10),
	}, true, true, DefaultIndexOptions)
	expect(t, !series.Clockwise())
	series = makeSeries([]Point{
		P(0, 0), P(0, 10), P(10, 10), P(10, 0), P(0, 0),
	}, true, true, DefaultIndexOptions)
	expect(t, series.Clockwise())
	series = makeSeries([]Point{
		P(0, 0), P(0, 10), P(10, 10), P(10, 0),
	}, true, true, DefaultIndexOptions)
	expect(t, series.Clockwise())
	series = makeSeries([]Point{
		P(0, 0), P(0, 10), P(10, 10),
	}, true, true, DefaultIndexOptions)
	expect(t, series.Clockwise())
}

func TestSeriesConvex(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		series := makeSeries([]Point{
			P(0, 0), P(4, 0), P(4, 4), P(0, 4), P(0, 0),
		}, true, true, DefaultIndexOptions)
		expect(t, series.Convex())
	})
	t.Run("2", func(t *testing.T) {
		series := makeSeries([]Point{
			P(0, 0), P(4, 0), P(4, 4), P(3, 4), P(1, 4), P(0, 4), P(0, 0),
		}, true, true, DefaultIndexOptions)
		expect(t, series.Convex())
	})
	t.Run("3", func(t *testing.T) {
		series := makeSeries([]Point{
			P(0, 0), P(4, 0), P(4, 4), P(2, 5), P(0, 4), P(0, 0),
		}, true, true, DefaultIndexOptions)
		expect(t, series.Convex())
	})
	t.Run("4", func(t *testing.T) {
		series := makeSeries([]Point{
			P(0, 0), P(4, 0), P(4, 4),
			P(3, 4), P(2, 5), P(1, 4),
			P(0, 4), P(0, 0),
		}, true, true, DefaultIndexOptions)
		expect(t, !series.Convex())
	})
	t.Run("5", func(t *testing.T) {
		series := makeSeries([]Point{
			P(0, 0), P(4, 0), P(4, 4),
			P(3, 4), P(2, 3), P(1, 4),
			P(0, 4), P(0, 0),
		}, true, true, DefaultIndexOptions)
		expect(t, !series.Convex())
	})
	t.Run("5", func(t *testing.T) {
		series := makeSeries([]Point{
			P(0, 0), P(4, 0), P(4, 4), P(0, 4),
			P(-1, 2), P(0, 0),
		}, true, true, DefaultIndexOptions)
		expect(t, series.Convex())
	})
	t.Run("6", func(t *testing.T) {
		series := makeSeries([]Point{
			P(0, 0), P(4, 0), P(4, 4), P(0, 4), P(0, 3),
			P(-1, 2), P(0, 0),
		}, true, true, DefaultIndexOptions)
		expect(t, !series.Convex())
	})
	t.Run("6", func(t *testing.T) {
		series := makeSeries([]Point{
			P(0, 0), P(4, 0), P(4, 4), P(0, 4),
			P(-1, 2), P(0, 1), P(0, 0),
		}, true, true, DefaultIndexOptions)
		expect(t, !series.Convex())
	})
	t.Run("7", func(t *testing.T) {
		series := makeSeries([]Point{
			P(0, 0), P(4, 0), P(4, 4),
			P(3, 3), P(2, 5), P(1, 3), P(0, 4),
			P(0, 0),
		}, true, true, DefaultIndexOptions)
		expect(t, !series.Convex())
	})
	t.Run("8", func(t *testing.T) {
		series := makeSeries([]Point{
			P(0, 0), P(0, 4), P(1, 3), P(4, 4), P(2, 5), P(3, 3), P(4, 0),
			P(0, 0),
		}, true, true, DefaultIndexOptions)
		expect(t, !series.Convex())
	})
}

func testScanSeries(
	t *testing.T, series *baseSeries, index bool,
	expectSegmentCount int, expectedEmpty bool,
) []Segment {
	t.Helper()
	if index {
		series.buildIndex()
	} else {
		series.clearIndex()
	}
	var segs1 []Segment
	seriesForEachSegment(series, func(seg Segment) bool {
		segs1 = append(segs1, seg)
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
	series := makeSeries(octagon, true, true, DefaultIndexOptions)
	expect(t, reflect.DeepEqual(seriesCopyPoints(&series), octagon))
	expect(t, series.Convex())
	expect(t, series.Rect() == R(0, 0, 10, 10))
	expect(t, series.Closed())
	series = makeSeries(octagon, false, true, DefaultIndexOptions)
	expect(t, reflect.DeepEqual(seriesCopyPoints(&series), octagon))

	series = makeSeries(ri, true, true, DefaultIndexOptions)
	testScanSeries(t, &series, true, len(ri)-1, false)
	testScanSeries(t, &series, false, len(ri)-1, false)

	// small lines
	series = makeSeries([]Point{}, true, false, DefaultIndexOptions)
	testScanSeries(t, &series, true, 0, true)
	testScanSeries(t, &series, false, 0, true)

	series = makeSeries([]Point{P(5, 5)}, true, false, DefaultIndexOptions)
	testScanSeries(t, &series, true, 0, true)
	testScanSeries(t, &series, false, 0, true)

	series = makeSeries([]Point{P(5, 5), P(10, 10)}, true, false, DefaultIndexOptions)
	testScanSeries(t, &series, true, 1, false)
	testScanSeries(t, &series, false, 1, false)

	// small rings
	series = makeSeries([]Point{}, true, true, DefaultIndexOptions)
	testScanSeries(t, &series, true, 0, true)
	testScanSeries(t, &series, false, 0, true)

	series = makeSeries([]Point{P(5, 5)}, true, true, DefaultIndexOptions)
	testScanSeries(t, &series, true, 0, true)
	testScanSeries(t, &series, false, 0, true)

	series = makeSeries([]Point{P(5, 5), P(10, 10)}, true, true, DefaultIndexOptions)
	testScanSeries(t, &series, true, 0, true)
	testScanSeries(t, &series, false, 0, true)

	series = makeSeries([]Point{P(5, 5), P(10, 10), P(10, 5)}, true, true, DefaultIndexOptions)
	testScanSeries(t, &series, true, 3, false)
	testScanSeries(t, &series, false, 3, false)

}
func TestSeriesSearch(t *testing.T) {
	series := makeSeries(octagon, true, true, DefaultIndexOptions)
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
	seriesForEachSegment(&series, func(seg Segment) bool {
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
		series := makeSeries(ri[:len(ri)-1], true, true, DefaultIndexOptions)
		var seg2sA []Segment
		series.Search(series.Rect(), func(seg Segment, idx int) bool {
			seg2sA = append(seg2sA, seg)
			return true
		})
		var seg2sB []Segment
		seriesForEachSegment(&series, func(seg Segment) bool {
			seg2sB = append(seg2sB, seg)
			return true
		})

		expect(t, checkSegsDups(seg2sA, seg2sB))

		// use all points
		series2 := makeSeries(ri, true, true, DefaultIndexOptions)
		var seg2sC []Segment
		seriesForEachSegment(&series2, func(seg Segment) bool {
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
		series := makeSeries(az, true, false, DefaultIndexOptions)
		var seg2sA []Segment
		series.Search(series.Rect(), func(seg Segment, idx int) bool {
			seg2sA = append(seg2sA, seg)
			return true
		})
		var seg2sB []Segment
		seriesForEachSegment(&series, func(seg Segment) bool {
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
		series := makeSeries(rev, true, true, DefaultIndexOptions)
		var seg2sA []Segment
		series.Search(series.Rect(), func(seg Segment, idx int) bool {
			seg2sA = append(seg2sA, seg)
			return true
		})
		var seg2sB []Segment
		seriesForEachSegment(&series, func(seg Segment) bool {
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

func TestSeriesMove(t *testing.T) {
	shapes := [][]Point{ri, octagon}
	for _, shape := range shapes {
		series := makeSeries(shape, true, true, DefaultIndexOptions)
		series.clearIndex()
		series2 := series.Move(60, 70)
		expect(t, series2.NumPoints() == len(shape))
		for i := 0; i < len(shape); i++ {
			expect(t, series2.PointAt(i) == shape[i].Move(60, 70))
		}
		series.buildIndex()
		series2 = series.Move(60, 70)
		expect(t, series2.NumPoints() == len(shape))
		for i := 0; i < len(shape); i++ {
			expect(t, series2.PointAt(i) == shape[i].Move(60, 70))
		}
	}
}
