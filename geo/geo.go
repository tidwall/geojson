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
	piR         = math.Pi * earthRadius
	twoPiR      = 2 * piR
)

func Haversine(latA, lonA, latB, lonB float64) float64 {
	φ1 := latA * radians
	λ1 := lonA * radians
	φ2 := latB * radians
	λ2 := lonB * radians
	Δφ := φ2 - φ1
	Δλ := λ2 - λ1
	sΔφ2 := math.Sin(Δφ / 2)
	sΔλ2 := math.Sin(Δλ / 2)
	return sΔφ2*sΔφ2 + math.Cos(φ1)*math.Cos(φ2)*sΔλ2*sΔλ2
}

func NormalizeDistance(meters float64) float64 {
	return math.Mod(meters, twoPiR)
}

func DistanceToHaversine(meters float64) float64 {
	// convert the given distance to its haversine
	sin := math.Sin(0.5 * meters / earthRadius)
	return sin * sin
}

func DistanceFromHaversine(haversine float64) float64 {
	return earthRadius * 2 * math.Asin(math.Sqrt(haversine))
}

// DistanceTo return the distance in meters between two point.
func DistanceTo(latA, lonA, latB, lonB float64) (meters float64) {
	a := Haversine(latA, lonA, latB, lonB)
	return DistanceFromHaversine(a)
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

// RectFromCenter calculates the bounding box surrounding a circle.
func RectFromCenter(lat, lon, meters float64) (
	minLat, minLon, maxLat, maxLon float64,
) {
	// convert degrees to radians
	lat *= radians
	lon *= radians

	// Calculate ANGULAR RADIUS
	// see http://janmatuschek.de/LatitudeLongitudeBoundingCoordinates#UsingIndex
	r := meters / earthRadius

	// Calculate LATITUDE min and max
	// see http://janmatuschek.de/LatitudeLongitudeBoundingCoordinates#Latitude
	minLat = lat - r
	maxLat = lat + r

	// Calculate LONGITUDE min and max
	// see http://janmatuschek.de/LatitudeLongitudeBoundingCoordinates#Longitude
	rCos := math.Cos(r)
	if rCos > 0.999999999999999 {

		// This can occur when the meters is too miniscule to derive the outer
		// rectangle coordinates.
		minLat = lat
		minLon = lon
		maxLat = lat
		maxLon = lon

	} else {

		latSin, latCos := math.Sincos(lat)
		latT := math.Asin(latSin / rCos)
		latTSin, latTCos := math.Sincos(latT)
		lonΔ := math.Acos((rCos - latTSin*latSin) / (latTCos * latCos))

		minLon = lon - lonΔ
		maxLon = lon + lonΔ
	}

	// ADJUST mins and maxes for edge-of-map cases
	// see http://janmatuschek.de/LatitudeLongitudeBoundingCoordinates#PolesAnd180thMeridian

	// adjust for NORTH POLE
	if maxLat > math.Pi/2 {
		minLon = -math.Pi
		maxLat = math.Pi / 2
		maxLon = math.Pi
	}

	// adjust for SOUTH POLE
	if minLat < -math.Pi/2 {
		minLat = -math.Pi / 2
		minLon = -math.Pi
		maxLon = math.Pi
	}

	/* adjust for WRAPAROUND

	Creates a bounding box that wraps around the Earth like a belt, which
	results in returning false positive candidates (candidates that are
	farther away from the center than the distance of the search radius).

	An alternative method, possibly to be implemented in the future, would be
	to split the bounding box into two boxes. This would return fewer (or no)
	false positives, but will require significant changes to the API's of
	geoJSON and any of its dependents. */
	if minLon < -math.Pi || maxLon > math.Pi {
		minLon = -math.Pi
		maxLon = math.Pi
	}

	// convert radians to degrees
	minLat *= degrees
	minLon *= degrees
	maxLat *= degrees
	maxLon *= degrees
	return

}

func DegsToSemi(degs float64) int32 {
	return int32(degs * (math.Pow(2, 31) / 180.0))
}

func SemiToDegs(semi int32) float64 {
	return float64(semi) * (180.0 / math.Pow(2, 31))
}
