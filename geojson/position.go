package geojson

// Position is 2D position
type Position struct {
	X, Y float64
}

// Center returns itself
func (posn *Position) Center() Position {
	return *posn
}

// Rect returns a rectangle
func (posn *Position) Rect() Rect {
	return Rect{*posn, *posn}
}

// HasBBox returns true if theres a GeoJSON bbox member
func (posn *Position) HasBBox() bool {
	return false
}

// CachedRect returns a precaclulated rectangle, if any
func (posn *Position) CachedRect() *Rect {
	return nil
}

// ForEachChild iterates over child objects.
func (posn *Position) ForEachChild(iter func(child Object) bool) {}

// Contains returns true if object contains other object
func (posn *Position) Contains(other Object) bool {
	return sharedContains(posn, other,
		func(other Object) (accept, contains bool) {
			switch other := other.(type) {
			case *Position:
				return true, polyPoint(*other).InsidePoint(
					polyPoint(*posn),
				)
			case *Rect:
				return true, polyRect(*other).InsidePoint(
					polyPoint(*posn),
				)
			case *Point:
				return true, polyPoint(other.Coordinates).InsidePoint(
					polyPoint(*posn),
				)
			}
			return false, false
		},
	)
}
