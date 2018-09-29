package geom

import (
	"github.com/tidwall/boxtree/d2"
)

// minTreePoints are the minumum number of points required before it makes
// sense to index an the segments in it's own rtree.
const minTreePoints = 32

// Series is just a series of points with utilities for efficiently accessing
// segments from rectangle queries, making stuff like point-in-polygon lookups
// very quick.
type Series struct {
	closed bool        // points create a closed shape
	convex bool        // points create a convex shape
	rect   Rect        // minumum bounding rectangle
	points []Point     // original points
	tree   *d2.BoxTree // segment tree. should be access though loadTree()
}

// NewSeries returns a new Series
func NewSeries(points []Point, copyPoints, closed bool) *Series {
	series := new(Series)
	series.closed = closed
	if copyPoints {
		series.points = make([]Point, len(points))
		copy(series.points, points)
	} else {
		series.points = points
	}
	if len(points) >= minTreePoints {
		series.tree = new(d2.BoxTree)
	}
	series.convex, series.rect = processPoints(points, closed, series.tree)
	return series
}

// Rect returns the series rectangle
func (series *Series) Rect() Rect {
	return series.rect
}

// Convex returns true if the points create a convex loop or linestring
func (series *Series) Convex() bool {
	return series.convex
}

// Closed return true if the shape is closed
func (series *Series) Closed() bool {
	return series.closed
}

// Points returns the original points
func (series *Series) Points() []Point {
	return series.points
}

// Search finds a searches for segments that intersect the provided rectangle
func (series *Series) Search(rect Rect, iter func(seg Segment, idx int) bool) {
	if series.tree == nil {
		series.Scan(func(seg Segment, idx int) bool {
			if seg.Rect().IntersectsRect(rect) {
				if !iter(seg, idx) {
					return false
				}
			}
			return true
		})
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

// Scan all segments in series
func (series *Series) Scan(iter func(seg Segment, idx int) bool) {
	var count int
	if series.closed {
		count = len(series.points)
	} else {
		count = len(series.points) - 1
	}
	for i := 0; i < count; i++ {
		var seg Segment
		seg.A = series.points[i]
		if series.closed && i == len(series.points)-1 {
			if seg.A == series.points[0] {
				break
			}
			seg.B = series.points[0]
		} else {
			seg.B = series.points[i+1]
		}
		if !iter(seg, i) {
			return
		}
	}
}

// processPoints tests if the ring is convex, calculates the outer
// rectangle, and inserts segments into a boxtree in one pass.
func processPoints(points []Point, closed bool, tree *d2.BoxTree) (
	convex bool, rect Rect,
) {
	var concave bool
	var dir int
	var a, b, c Point
	var segCount int
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

		// process the convex calculation
		if concave {
			continue
		}
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
		dx1 := b.X - a.X
		dy1 := b.Y - a.Y
		dx2 := c.X - b.X
		dy2 := c.Y - b.Y
		zCrossProduct := dx1*dy2 - dy1*dx2
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
	return !concave, rect
}
