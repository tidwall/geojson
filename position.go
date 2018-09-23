package geojson

import (
	"strconv"
)

type Position struct {
	X, Y float64
}

func (posn Position) BBoxDefined() bool {
	return false
}

func (posn Position) Rect() Rect {
	return Rect{Min: posn, Max: posn}
}

func (posn Position) Center() Position {
	return posn
}
func (posn Position) AppendJSON(dst []byte) []byte {
	return Point{Coordinates: posn}.AppendJSON(dst)
}
func (posn Position) Contains(other Object) bool {
	return objectContains(posn, other)
}
func (posn Position) Intersects(other Object) bool {
	return objectIntersects(posn, other)
}

func (posn Position) primativeContains(other Object) bool {
	ppoint := polyPoint(posn)
	switch other := other.(type) {
	case Position:
		return polyPoint(other).InsidePoint(ppoint)
	case Rect:
		return polyRect(other).InsidePoint(ppoint)
	case Point:
		return polyPoint(other.Coordinates).InsidePoint(ppoint)
	case LineString:
		return polyLine(other.Coordinates).InsidePoint(ppoint)
	case Polygon:
		return polyPolygon(other.Coordinates).InsidePoint(ppoint)
	}
	return false
}
func (posn Position) primativeIntersects(other Object) bool {
	ppoint := polyPoint(posn)
	switch other := other.(type) {
	case Position:
		return ppoint.IntersectsPoint(polyPoint(other))
	case Rect:
		return ppoint.IntersectsRect(polyRect(other))
	case Point:
		return ppoint.IntersectsPoint(polyPoint(other.Coordinates))
	case LineString:
		return ppoint.IntersectsLine(polyLine(other.Coordinates))
	case Polygon:
		return ppoint.IntersectsPolygon(polyPolygon(other.Coordinates))
	}
	return false
}

// ForEachChild iterates over child objects.
func (posn Position) ForEachChild(func(child Object) bool) {}

func appendJSONPosition(dst []byte, posn Position, ex *Extra, idx int) []byte {
	dst = append(dst, '[')
	dst = strconv.AppendFloat(dst, posn.X, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, posn.Y, 'f', -1, 64)
	if ex != nil {
		dims := int(ex.Dims)
		for i := 0; i < dims; i++ {
			dst = append(dst, ',')
			dst = strconv.AppendFloat(
				dst, ex.Positions[idx*dims+i], 'f', -1, 64,
			)
		}
	}
	dst = append(dst, ']')
	return dst
}
