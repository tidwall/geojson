// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geos

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestCircleNewCircle(t *testing.T) {
	circle := NewCircle(P(-112, 33), 1000, 2)
	expect(t, circle.ContainsPoint(P(-112, 33)))
}

func BenchmarkCircleContainsPoint(b *testing.B) {
	center := Point{-112, 33}
	meters := 1000.0
	var points []Point
	for j := 2; j <= 4096; j *= 2 {
		b.Run(fmt.Sprintf("%d", j), func(b *testing.B) {
			for i := 0; i < 2; i++ {
				name := "Simple"
				if i == 1 {
					name = "Indexed"
				}
				b.Run(name, func(b *testing.B) {
					poly := NewCircle(center, meters, j)
					if i == 0 {
						poly.Exterior.(*baseSeries).tree = nil
					} else {
						poly.Exterior.(*baseSeries).buildTree()
					}
					for len(points) < b.N {
						lat, lon := DestinationPoint(
							center.Y, center.X,
							meters*1.5*rand.Float64(),
							360*rand.Float64(),
						)
						points = append(points, Point{
							X: lon,
							Y: lat,
						})
					}
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						poly.ContainsPoint(points[i])
					}
				})
			}
		})
	}
}
