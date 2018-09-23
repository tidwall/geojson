package geojson

// Point is GeoJSON Point type
type Point struct {
	Coordinates Position
	bbox        *objBBox
	extra       []float64
	extraDims   int
}

// Center returns the center position of the point
func (point *Point) Center() Position {
	if point.bbox != nil {
		return point.bbox.rect.Center()
	}
	return point.Coordinates
}

// Rect returns the rectangle around the point
func (point *Point) Rect() Rect {
	if point.bbox != nil {
		return point.bbox.rect
	}
	return Rect{point.Coordinates, point.Coordinates}
}

// HasBBox returns true if theres a GeoJSON bbox member
func (point *Point) HasBBox() bool {
	return point.bbox != nil
}

// CachedRect returns a precaclulated rectangle, if any
func (point *Point) CachedRect() *Rect {
	return nil
}

// ForEachChild iterates over child objects.
func (point *Point) ForEachChild(iter func(child Object) bool) {}

// Contains returns true if object contains other object
func (point *Point) Contains(other Object) bool {
	return sharedContains(point, other,
		func(other Object) (accept, contains bool) {
			switch other := other.(type) {
			case *Position:
				return true, polyPoint(*other).InsidePoint(
					polyPoint(point.Coordinates),
				)
			case *Rect:
				return true, polyRect(*other).InsidePoint(
					polyPoint(point.Coordinates),
				)
			case *Point:
				return true, polyPoint(other.Coordinates).InsidePoint(
					polyPoint(point.Coordinates),
				)
			}
			return false, false
		},
	)
}
