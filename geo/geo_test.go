// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geo

import "testing"

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
