// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geo

import (
	"math"
)

const (
	earthRadius = 6371e3
	radians     = math.Pi / 180
	degrees     = 180 / math.Pi
)

// DistanceTo return the distance in meteres between two point.
func DistanceTo(latA, lonA, latB, lonB float64) (meters float64) {
	φ1 := latA * radians
	λ1 := lonA * radians
	φ2 := latB * radians
	λ2 := lonB * radians
	Δφ := φ2 - φ1
	Δλ := λ2 - λ1
	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

// DestinationPoint return the destination from a point based on a
// distance and bearing.
func DestinationPoint(lat, lon, meters, bearingDegrees float64) (
	destLat, destLon float64,
) {
	// see http://williams.best.vwh.net/avform.htm#LL
	δ := meters / earthRadius // angular distance in radians
	θ := bearingDegrees * radians
	φ1 := lat * radians
	λ1 := lon * radians
	φ2 := math.Asin(math.Sin(φ1)*math.Cos(δ) +
		math.Cos(φ1)*math.Sin(δ)*math.Cos(θ))
	λ2 := λ1 + math.Atan2(math.Sin(θ)*math.Sin(δ)*math.Cos(φ1),
		math.Cos(δ)-math.Sin(φ1)*math.Sin(φ2))
	λ2 = math.Mod(λ2+3*math.Pi, 2*math.Pi) - math.Pi // normalise to -180..+180°
	return φ2 * degrees, λ2 * degrees
}

// BearingTo returns the (initial) bearing from point 'A' to point 'B'.
func BearingTo(latA, lonA, latB, lonB float64) float64 {
	// tanθ = sinΔλ⋅cosφ2 / cosφ1⋅sinφ2 − sinφ1⋅cosφ2⋅cosΔλ
	// see mathforum.org/library/drmath/view/55417.html for derivation

	φ1 := latA * radians
	φ2 := latB * radians
	Δλ := (lonB - lonA) * radians
	y := math.Sin(Δλ) * math.Cos(φ2)
	x := math.Cos(φ1)*math.Sin(φ2) - math.Sin(φ1)*math.Cos(φ2)*math.Cos(Δλ)
	θ := math.Atan2(y, x)

	return math.Mod(θ*degrees+360, 360)
}

// // SegmentIntersectsCircle ...
// func SegmentIntersectsCircle(
// 	startLat, startLon, endLat, endLon, centerLat, centerLon, meters float64,
// ) bool {
// 	// These are faster checks.
// 	// If they succeed there's no need do complicate things.
// 	if DistanceTo(startLat, startLon, centerLat, centerLon) <= meters {
// 		return true
// 	}
// 	if DistanceTo(endLat, endLon, centerLat, centerLon) <= meters {
// 		return true
// 	}

// 	// Distance between start and end
// 	l := DistanceTo(startLat, startLon, endLat, endLon)

// 	// Unit direction vector
// 	dLat := (endLat - startLat) / l
// 	dLon := (endLon - startLon) / l

// 	// Point of the line closest to the center
// 	t := dLon*(centerLon-startLon) + dLat*(centerLat-startLat)
// 	pLat := t*dLat + startLat
// 	pLon := t*dLon + startLon
// 	if pLon < startLon || pLon > endLon || pLat < startLat || pLat > endLat {
// 		// closest point is outside the segment
// 		return false
// 	}

// 	// Distance from the closest point to the center
// 	return DistanceTo(centerLat, centerLon, pLat, pLon) <= meters
// }
