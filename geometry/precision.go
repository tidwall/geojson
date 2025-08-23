// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geometry

import "math"

// Epsilon for floating point precision across different platforms (ARM64, AMD64)
//
// We use 1e-8 instead of smaller values like 1e-10 for the following reasons:
// 1. Geometric calculations involve multiple floating point operations that accumulate errors
// 2. Cross-platform floating point differences between ARM64 and AMD64 can be in the 1e-9 range
// 3. For geographic coordinates (GeoJSON's primary use case), 1e-8 degrees is approximately 1mm precision
// 4. This epsilon provides a good balance between precision and robustness against computational noise
const epsilon = 1e-8

// FloatEqual checks if two float64 values are equal within epsilon tolerance
func FloatEqual(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

// FloatNotEqual checks if two float64 values are not equal within epsilon tolerance
func FloatNotEqual(a, b float64) bool {
	return !FloatEqual(a, b)
}

// FloatLess checks if a is less than b with epsilon tolerance
func FloatLess(a, b float64) bool {
	return (b - a) > epsilon
}

// FloatLessOrEqual checks if a is less than or equal to b with epsilon tolerance
func FloatLessOrEqual(a, b float64) bool {
	return FloatLess(a, b) || FloatEqual(a, b)
}

// FloatGreater checks if a is greater than b with epsilon tolerance
func FloatGreater(a, b float64) bool {
	return (a - b) > epsilon
}

// FloatGreaterOrEqual checks if a is greater than or equal to b with epsilon tolerance
func FloatGreaterOrEqual(a, b float64) bool {
	return FloatGreater(a, b) || FloatEqual(a, b)
}

// FloatZero checks if a float64 value is effectively zero
func FloatZero(x float64) bool {
	return math.Abs(x) < epsilon
}

// FloatNonZero checks if a float64 value is effectively non-zero
func FloatNonZero(x float64) bool {
	return !FloatZero(x)
}

// PointEqual checks if two points are equal within epsilon tolerance
func PointEqual(a, b Point) bool {
	return FloatEqual(a.X, b.X) && FloatEqual(a.Y, b.Y)
}

// PointNotEqual checks if two points are not equal within epsilon tolerance
func PointNotEqual(a, b Point) bool {
	return !PointEqual(a, b)
}

// RectEqual checks if two rectangles are equal within epsilon tolerance
func RectEqual(a, b Rect) bool {
	return PointEqual(a.Min, b.Min) && PointEqual(a.Max, b.Max)
}

// RectNotEqual checks if two rectangles are not equal within epsilon tolerance
func RectNotEqual(a, b Rect) bool {
	return !RectEqual(a, b)
}
