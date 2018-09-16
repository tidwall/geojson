// Derived from javascript at http://www.movable-type.co.uk/scripts/geohash.html
//
// Original copyright states...
// "Geohash encoding/decoding and associated functions
// (c) Chris Veness 2014 / MIT Licence"

package geohash

import "errors"

var (
	errInvalidPrecision = errors.New("invalid precision")
	errEncodingError    = errors.New("encoding error")
	errInvalidGeohash   = errors.New("invalid geohash")
)
var base32R = [...]int8{
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, -1,
	-1, -1, -1, -1, -1, -1, -1, 10, 11, 12, 13, 14, 15, 16, -1, 17, 18, -1, 19,
	20, -1, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, -1, -1, -1, -1, -1, -1,
	-1, 10, 11, 12, 13, 14, 15, 16, -1, 17, 18, -1, 19, 20, -1, 21, 22, 23, 24,
	25, 26, 27, 28, 29, 30, 31, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1,
}

var base32F = [...]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'b',
	'c', 'd', 'e', 'f', 'g', 'h', 'j', 'k', 'm', 'n', 'p',
	'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
}

// Encode latitude/longitude to geohash, either to specified precision or to
// automatically evaluated precision.
func Encode(lat, lon float64, precision int) (string, error) {
	var idx = 0 // index into base32 map
	var bit = 0 // each char holds 5 bits
	var evenBit = true
	var latMin = -90.0
	var latMax = 90.0
	var lonMin = -180.0
	var lonMax = 180.0
	if precision < 1 {
		return "", errInvalidPrecision
	}
	geohash := make([]byte, 0, precision)
	for len(geohash) < precision {
		if evenBit {
			// bisect E-W longitude
			var lonMid = (lonMin + lonMax) / 2
			if lon > lonMid {
				idx = idx*2 + 1
				lonMin = lonMid
			} else {
				idx = idx * 2
				lonMax = lonMid
			}
		} else {
			// bisect N-S latitude
			var latMid = (latMin + latMax) / 2
			if lat > latMid {
				idx = idx*2 + 1
				latMin = latMid
			} else {
				idx = idx * 2
				latMax = latMid
			}
		}
		evenBit = !evenBit

		bit = bit + 1
		if bit == 5 {
			// 5 bits gives us a character: append it and start over
			geohash = append(geohash, base32F[idx])
			bit = 0
			idx = 0
		}
	}
	return string(geohash), nil
}

// Decode geohash to latitude/longitude (location is approximate centre of
//  geohash cell, to reasonable precision).
func Decode(geohash string) (lat, lon float64, err error) {
	swLat, swLon, neLat, neLon, err1 := Bounds(geohash) // <-- the hard work
	if err1 != nil {
		return 0, 0, err1
	}
	return (neLat-swLat)/2 + swLat, (neLon-swLon)/2 + swLon, nil
}

// Bounds returns SW/NE latitude/longitude bounds of specified geohash.
func Bounds(geohash string) (swLat, swLon, neLat, neLon float64, err error) {
	var evenBit = true
	var latMin = -90.0
	var latMax = 90.0
	var lonMin = -180.0
	var lonMax = 180.0
	for i := 0; i < len(geohash); i++ {
		var chr = geohash[i]
		var idx = base32R[chr]
		if idx == -1 {
			return 0, 0, 0, 0, errInvalidGeohash
		}
		for n := uint(4); ; n-- {
			var bitN = idx >> n & 1
			if evenBit {
				// longitude
				var lonMid = (lonMin + lonMax) / 2
				if bitN == 1 {
					lonMin = lonMid
				} else {
					lonMax = lonMid
				}
			} else {
				// latitude
				var latMid = (latMin + latMax) / 2
				if bitN == 1 {
					latMin = latMid
				} else {
					latMax = latMid
				}
			}
			evenBit = !evenBit
			if n == 0 {
				break
			}
		}
	}
	return latMin, lonMin, latMax, lonMax, nil
}
