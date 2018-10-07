// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geos

// NewCircle returns a circle polygon
func NewCircle(center Point, meters float64, segments int) *Poly {
	if segments < 3 {
		segments = 3
	}
	points := make([]Point, segments+1)
	step := 360.0 / float64(segments)
	i := 0
	for deg := 360.0; deg > 0; deg -= step {
		lat, lon := DestinationPoint(center.Y, center.X, meters, deg)
		points[i] = Point{X: lon, Y: lat}
		i++
	}
	points[i] = points[0]
	return NewPoly(points, nil)
}
