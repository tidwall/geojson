package geojson

import "strconv"

type Position struct {
	X, Y float64
}

func (posn Position) HasBBox() bool {
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
	rect := other.Rect()
	return rect.Min == rect.Max && rect.Min == posn
}
func (posn Position) Intersects(other Object) bool {
	switch other := other.(type) {
	case Position:
		return posn == other
	case Rect:
		return other.ContainsPosition(posn)
	}
	// bbox types
	if !other.Rect().ContainsPosition(posn) {
		// no intersection
		return false
	}
	// yes they intersect
	if other.HasBBox() {
		// nothing more to check
		return true
	}
	// geometry types
	switch other := other.(type) {
	case Point:
		return polyPoint(other.Coordinates) == polyPoint(posn)
	case LineString:
		return polyLine(other.Coordinates).LineStringIntersects()
	case Polygon:
		return polyRect(rect).Intersects(polyPolygon(other.Coordinates))
	}
	// check types with children
	var intersects bool
	other.ForEach(func(child Object) bool {
		if rect.Intersects(child) {
			intersects = true
			return false
		}
		return true
	})
	return intersects

	panic("unsupported")
}

func (posn Position) ForEach(func(child Object) bool) {}

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
