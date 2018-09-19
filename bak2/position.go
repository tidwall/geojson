package geojson

import (
	"strconv"

	"github.com/tidwall/geojson/geohash"
	"github.com/tidwall/geojson/poly"
)

// Position is a simple X,Y point
type Position poly.Point

func (posn Position) appendJSON(dst []byte, extra *extra, index int) []byte {
	dst = append(dst, '[')
	dst = strconv.AppendFloat(dst, posn.X, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, posn.Y, 'f', -1, 64)
	dims := extra.dims()
	for i := 0; i < dims; i++ {
		dst = append(dst, ',')
		dst = strconv.AppendFloat(dst, extra.coords[index*dims+i], 'f', -1, 64)
	}
	dst = append(dst, ']')
	return dst
}

// Geohash converts the object to a geohash value
func (posn Position) Geohash(precision int) (string, error) {
	return geohash.Encode(posn.Y, posn.X, precision)
}

// IntersectsPolyLine test if point intersects polyline
func (posn Position) IntersectsPolyLine(line []Position) bool {
	return polyPoint(posn).IntersectsLineString(polyLine(line))
}
