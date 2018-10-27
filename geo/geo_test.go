// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geo

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func init() {
	seed := time.Now().UnixNano()
	//seed = 1540656736244531000
	println(seed)
	rand.Seed(seed)
}

func TestGeoCalc(t *testing.T) {
	dist := 172853.26908429610193707048892974853515625
	bearing := 320.8560640269032546711969189345836639404296875
	latA, lonA := 33.112, -112.123
	latB, lonB := 34.312, -113.311
	// DistanceTo
	value := DistanceTo(latA, lonA, latB, lonB)
	if value != dist {
		t.Fatalf("expected '%v', got '%v'", dist, value)
	}
	// BearingTo
	value = BearingTo(latA, lonA, latB, lonB)
	if value != bearing {
		t.Fatalf("expected '%v', got '%v'", bearing, value)
	}
	// DestinationPoint
	value1, value2 := DestinationPoint(latA, lonA, dist, bearing)
	if value1 != latB {
		t.Fatalf("expected '%v', got '%v'", latB, value1)
	}
	if value2 != lonB {
		t.Fatalf("expected '%v', got '%v'", lonB, value2)
	}
}

func TestHaversine(t *testing.T) {
	latA := rand.Float64()*180 - 90
	lonA := rand.Float64()*360 - 180
	start := time.Now()
	for time.Since(start) < time.Second/4 {
		for i := 0; i < 1000; i++ {
			latB := rand.Float64()*180 - 90
			lonB := rand.Float64()*360 - 180
			latC := rand.Float64()*180 - 90
			lonC := rand.Float64()*360 - 180
			haver1 := Haversine(latA, lonA, latB, lonB)
			haver2 := Haversine(latA, lonA, latC, lonC)
			meters1 := DistanceTo(latA, lonA, latB, lonB)
			meters2 := DistanceTo(latA, lonA, latC, lonC)
			switch {
			case haver1 < haver2:
				if meters1 >= meters2 {
					t.Fatalf("failed")
				}
			case haver1 == haver2:
				if meters1 != meters2 {
					t.Fatalf("failed")
				}
			case haver1 > haver2:
				if meters1 <= meters2 {
					t.Fatalf("failed")
				}
			}
		}
	}
}

func TestNormalizeDistance(t *testing.T) {
	start := time.Now()
	for time.Since(start) < time.Second/4 {
		for i := 0; i < 1000; i++ {
			meters1 := rand.Float64() * 100000000
			meters2 := NormalizeDistance(meters1)
			dist1 := math.Floor(DistanceToHaversine(meters2) * 100000000.0)
			dist2 := math.Floor(DistanceToHaversine(meters1) * 100000000.0)
			if dist1 != dist2 {
				t.Fatalf("expected %f, got %f", dist2, dist1)
			}
		}
	}
}
