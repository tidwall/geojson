// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geometry

import "github.com/tidwall/boxtree/d2"

// DefaultIndex are the minumum number of points required before it makes
// sense to index the segments.
// 64 seems to be the sweet spot
const DefaultIndex = 64

// Series is just a series of points with utilities for efficiently accessing
// segments from rectangle queries, making stuff like point-in-polygon lookups
// very quick.
type Series interface {
	Rect() Rect
	Empty() bool
	Convex() bool
	Clockwise() bool
	NumPoints() int
	NumSegments() int
	PointAt(index int) Point
	SegmentAt(index int) Segment
	Search(rect Rect, iter func(seg Segment, index int) bool)
}

func seriesCopyPoints(series Series) []Point {
	points := make([]Point, series.NumPoints())
	for i := 0; i < len(points); i++ {
		points[i] = series.PointAt(i)
	}
	return points
}

// baseSeries is a concrete type containing all that is needed to make a Series.
type baseSeries struct {
	closed    bool        // points create a closed shape
	clockwise bool        // points move clockwise
	convex    bool        // points create a convex shape
	rect      Rect        // minumum bounding rectangle
	points    []Point     // original points
	tree      *d2.BoxTree // segment tree
}

// makeSeries returns a processed baseSeries.
func makeSeries(points []Point, copyPoints, closed bool, index int) baseSeries {
	var series baseSeries
	series.closed = closed
	if copyPoints {
		series.points = make([]Point, len(points))
		copy(series.points, points)
	} else {
		series.points = points
	}
	if index != 0 && len(points) >= int(index) {
		series.tree = new(d2.BoxTree)
	}
	series.convex, series.rect, series.clockwise =
		processPoints(points, closed, series.tree)
	return series
}

// Clockwise ...
func (series *baseSeries) Clockwise() bool {
	return series.clockwise
}

func (series *baseSeries) Move(deltaX, deltaY float64) Series {
	points := make([]Point, len(series.points))
	for i := 0; i < len(series.points); i++ {
		points[i].X = series.points[i].X + deltaX
		points[i].Y = series.points[i].Y + deltaY
	}
	nseries := makeSeries(points, false, series.closed, 0)
	if series.tree != nil {
		nseries.buildTree()
	}
	return &nseries
}

// Empty returns true if the series does not take up space.
func (series *baseSeries) Empty() bool {
	if series == nil {
		return true
	}
	return (series.closed && len(series.points) < 3) || len(series.points) < 2
}

// Rect returns the series rectangle
func (series *baseSeries) Rect() Rect {
	return series.rect
}

// Convex returns true if the points create a convex loop or linestring
func (series *baseSeries) Convex() bool {
	return series.convex
}

// Closed return true if the shape is closed
func (series *baseSeries) Closed() bool {
	return series.closed
}

// NumPoints returns the number of points in the series
func (series *baseSeries) NumPoints() int {
	return len(series.points)
}

// PointAt returns the point at index
func (series *baseSeries) PointAt(index int) Point {
	return series.points[index]
}

// Search finds a searches for segments that intersect the provided rectangle
func (series *baseSeries) Search(rect Rect, iter func(seg Segment, idx int) bool) {
	if series.tree == nil {
		n := series.NumSegments()
		for i := 0; i < n; i++ {
			seg := series.SegmentAt(i)
			if seg.Rect().IntersectsRect(rect) {
				if !iter(seg, i) {
					return
				}
			}
		}
	} else {
		series.tree.Search(
			[]float64{rect.Min.X, rect.Min.Y},
			[]float64{rect.Max.X, rect.Max.Y},
			func(_, _ []float64, value interface{}) bool {
				index := value.(int)
				var seg Segment
				seg.A = series.points[index]
				if series.closed && index == len(series.points)-1 {
					seg.B = series.points[0]
				} else {
					seg.B = series.points[index+1]
				}
				if !iter(seg, index) {
					return false
				}
				return true
			},
		)
	}
}

// NumSegments ...
func (series *baseSeries) NumSegments() int {
	if series.closed {
		if len(series.points) < 3 {
			return 0
		}
		if series.points[len(series.points)-1] == series.points[0] {
			return len(series.points) - 1
		}
		return len(series.points)
	}
	if len(series.points) < 2 {
		return 0
	}
	return len(series.points) - 1
}

// SegmentAt ...
func (series *baseSeries) SegmentAt(index int) Segment {
	var seg Segment
	seg.A = series.points[index]
	if index == len(series.points)-1 {
		seg.B = series.points[0]
	} else {
		seg.B = series.points[index+1]
	}
	return seg
}

func (series *baseSeries) buildTree() {
	if series.tree == nil {
		series.tree = new(d2.BoxTree)
		processPoints(series.points, series.closed, series.tree)
	}
}

// processPoints tests if the ring is convex, calculates the outer
// rectangle, and inserts segments into a boxtree in one pass.
func processPoints(points []Point, closed bool, tree *d2.BoxTree) (
	convex bool, rect Rect, clockwise bool,
) {
	if (closed && len(points) < 3) || len(points) < 2 {
		return
	}
	var concave bool
	var dir int
	var a, b, c Point
	var segCount int
	var cwc float64
	if closed {
		segCount = len(points)
	} else {
		segCount = len(points) - 1
	}

	for i := 0; i < len(points); i++ {
		// process the segments for tree insertion
		if tree != nil && i < segCount {
			var seg Segment
			seg.A = points[i]
			if closed && i == len(points)-1 {
				if seg.A == points[0] {
					break
				}
				seg.B = points[0]
			} else {
				seg.B = points[i+1]
			}
			rect := seg.Rect()
			tree.Insert(
				[]float64{rect.Min.X, rect.Min.Y},
				[]float64{rect.Max.X, rect.Max.Y}, i)
		}

		// process the rectangle inflation
		if i == 0 {
			rect = Rect{points[i], points[i]}
		} else {
			if points[i].X < rect.Min.X {
				rect.Min.X = points[i].X
			} else if points[i].X > rect.Max.X {
				rect.Max.X = points[i].X
			}
			if points[i].Y < rect.Min.Y {
				rect.Min.Y = points[i].Y
			} else if points[i].Y > rect.Max.Y {
				rect.Max.Y = points[i].Y
			}
		}

		// gather some point positions for concave and clockwise detection
		a = points[i]
		if i == len(points)-1 {
			b = points[0]
			c = points[1]
		} else if i == len(points)-2 {
			b = points[i+1]
			c = points[0]
		} else {
			b = points[i+1]
			c = points[i+2]
		}

		// process the clockwise detection
		cwc += (b.X - a.X) * (b.Y + a.Y)

		// process the convex calculation
		if concave {
			continue
		}

		zCrossProduct := (b.X-a.X)*(c.Y-b.Y) - (b.Y-a.Y)*(c.X-b.X)
		if dir == 0 {
			if zCrossProduct < 0 {
				dir = -1
			} else if zCrossProduct > 0 {
				dir = 1
			}
		} else if zCrossProduct < 0 {
			if dir == 1 {
				concave = true
			}
		} else if zCrossProduct > 0 {
			if dir == -1 {
				concave = true
			}
		}
	}
	return !concave, rect, cwc > 0
}
