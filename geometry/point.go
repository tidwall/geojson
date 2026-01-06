// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geometry

import "math"

type Point struct {
	X, Y float64
}

func (point Point) Move(deltaX, deltaY float64) Point {
	return Point{X: point.X + deltaX, Y: point.Y + deltaY}
}

func (point Point) Empty() bool {
	return math.IsNaN(point.X) && math.IsNaN(point.Y)
}

func (point Point) Valid() bool {
	return point.X >= -180 && point.X <= 180 && point.Y >= -90 && point.Y <= 90
}

func (point Point) Rect() Rect {
	return Rect{point, point}
}

func (point Point) ContainsPoint(other Point) bool {
	return PointEqual(point, other)
}

func (point Point) IntersectsPoint(other Point) bool {
	return PointEqual(point, other)
}

func (point Point) ContainsRect(rect Rect) bool {
	return RectEqual(point.Rect(), rect)
}

func (point Point) IntersectsRect(rect Rect) bool {
	return rect.ContainsPoint(point)
}

func (point Point) ContainsLine(line *Line) bool {
	if line == nil {
		return false
	}
	return !line.Empty() && RectEqual(line.Rect(), point.Rect())
}

func (point Point) IntersectsLine(line *Line) bool {
	if line == nil {
		return false
	}
	return line.IntersectsPoint(point)
}

func (point Point) ContainsPoly(poly *Poly) bool {
	if poly == nil {
		return false
	}
	return !poly.Empty() && RectEqual(poly.Rect(), point.Rect())
}

func (point Point) IntersectsPoly(poly *Poly) bool {
	if poly == nil {
		return false
	}
	return poly.IntersectsPoint(point)
}
