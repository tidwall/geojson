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
func (posn Position) Within(other Object) bool {
	panic("unsupported")
}
func (posn Position) Intersects(other Object) bool {
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
