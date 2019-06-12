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
	// RectFromCenter
	// expected mins and maxes for wraparound and pole tests
	expMinLat := -90.0
	expMinLon := -180.0
	expMaxLat := 90.0
	expMaxLon := 180.0

	dist1 := 1600000.0

	wraparoundTests := []struct {
		name                   string
		lat, lon, searchRadius float64
	}{
		{name: "Wraparound E", lat: 0.0, lon: 179.0, searchRadius: dist1},  // at equator near 180th meridian East
		{name: "Wraparound W", lat: 0.0, lon: -179.0, searchRadius: dist1}, // at equator near 180th meridian West
	}

	for _, tt := range wraparoundTests {
		t.Run(tt.name, func(t *testing.T) {
			_, minLon, _, maxLon := RectFromCenter(tt.lat, tt.lon, tt.searchRadius)
			if !(minLon == expMinLon && maxLon == expMaxLon) {
				t.Errorf("\nexpected minLon = '%v', maxLon = '%v'"+"\ngot minLon = '%v', maxLon = '%v'\n",
					expMinLon, expMaxLon,
					minLon, maxLon)
			}
		})
	}

	northPoleTests := []struct {
		name                   string
		lat, lon, searchRadius float64
	}{
		{name: "North Pole", lat: 89.0, lon: 90.0, searchRadius: dist1}, // near North Pole
		{name: "North Pole: Tile38 iss422", lat: 13.0257553, lon: 77.6672509, searchRadius: 9000000.0},
	}
	for _, tt := range northPoleTests {
		t.Run(tt.name, func(t *testing.T) {
			_, minLon, maxLat, maxLon := RectFromCenter(tt.lat, tt.lon, tt.searchRadius)
			if !(minLon == expMinLon && maxLat == expMaxLat && maxLon == expMaxLon) {
				t.Errorf("\nexpected minLon = '%v', maxLat = '%v', maxLon = '%v'"+"\ngot minLon = '%v', maxLat = '%v', maxLon = '%v'",
					expMinLon, expMaxLat, expMaxLon,
					minLon, maxLat, maxLon)
			}
		})
	}

	southPoleTests := []struct {
		name                   string
		lat, lon, searchRadius float64
	}{
		{name: "South Pole", lat: -89.0, lon: 90.0, searchRadius: dist1}, // near South Pole
	}
	for _, tt := range southPoleTests {
		t.Run(tt.name, func(t *testing.T) {
			minLat, minLon, _, maxLon := RectFromCenter(tt.lat, tt.lon, tt.searchRadius)
			if !(minLat == expMinLat && minLon == expMinLon && maxLon == expMaxLon) {
				t.Errorf("\nexpected minLat = '%v', minLon = '%v', maxLon = '%v'"+"\ngot minLat = '%v', minLon = '%v', maxLon = '%v'",
					expMinLat, expMinLon, expMaxLon,
					minLat, minLon, maxLon)
			}
		})
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
	for time.Since(start) < time.Second {
		for i := 0; i < 1000; i++ {
			meters1 := rand.Float64() * earthRadius * 3 // wrap three times
			meters2 := NormalizeDistance(meters1)
			dist1 := math.Floor(DistanceToHaversine(meters2) * 1e8)
			dist2 := math.Floor(DistanceToHaversine(meters1) * 1e8)
			if dist1 != dist2 {
				t.Fatalf("expected %f, got %f", dist2, dist1)
			}
		}
	}
}

type point struct {
	lat, lon float64
}

func BenchmarkHaversine(b *testing.B) {
	pointA := point{
		lat: rand.Float64()*180 - 90,
		lon: rand.Float64()*360 - 180,
	}
	points := make([]point, b.N)
	for i := 0; i < b.N; i++ {
		points[i].lat = rand.Float64()*180 - 90
		points[i].lon = rand.Float64()*360 - 180
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Haversine(pointA.lat, pointA.lon, points[i].lat, points[i].lon)
	}
}

func BenchmarkDistanceTo(b *testing.B) {
	pointA := point{
		lat: rand.Float64()*180 - 90,
		lon: rand.Float64()*360 - 180,
	}
	points := make([]point, b.N)
	for i := 0; i < b.N; i++ {
		points[i].lat = rand.Float64()*180 - 90
		points[i].lon = rand.Float64()*360 - 180
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DistanceTo(pointA.lat, pointA.lon, points[i].lat, points[i].lon)
	}
}
