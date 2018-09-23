package geojson

// Rect is a 2D rectangle
type Rect struct {
	Min, Max Position
}

func calcRectFromChildren(obj Object) Rect {
	var r Rect
	var i int
	obj.ForEachChild(func(child Object) bool {
		r2 := child.Rect()
		if i == 0 {
			r = r2
		} else {
			if r2.Min.X < r.Min.X {
				r.Min.X = r2.Min.X
			} else if r2.Max.X > r.Max.X {
				r.Max.X = r2.Max.X
			}
			if r2.Min.Y < r.Min.Y {
				r.Min.Y = r2.Min.Y
			} else if r2.Max.Y > r.Max.Y {
				r.Max.Y = r2.Max.Y
			}
		}
		i++
		return true
	})
	return r
}

// Center returns the center position of the point
func (rect *Rect) Center() Position {
	return Position{
		(rect.Min.X + rect.Max.X) / 2,
		(rect.Min.Y + rect.Max.Y) / 2,
	}
}

// Rect returns the center position of the point
func (rect *Rect) Rect() Rect {
	return *rect
}

// HasBBox returns true if theres a GeoJSON bbox member
func (rect *Rect) HasBBox() bool {
	return false
}

// CachedRect returns a precaclulated rectangle, if any
func (rect *Rect) CachedRect() *Rect {
	return nil
}

// ForEachChild iterates over child objects.
func (rect *Rect) ForEachChild(iter func(child Object) bool) {}

// Contains returns true if object contains other object
func (rect *Rect) Contains(other Object) bool {
	return sharedContains(rect, other,
		func(other Object) (accept, contains bool) {
			switch other := other.(type) {
			case *Position:
				return true, polyPoint(*other).InsideRect(polyRect(*rect))
			case *Rect:
				return true, polyRect(*other).InsideRect(polyRect(*rect))
			case *Point:
				return true, polyPoint(other.Coordinates).InsideRect(
					polyRect(*rect),
				)
			}
			return false, false
		},
	)
}

// Intersects returns true if object intersects other object
func (rect *Rect) Intersects(other Object) bool {
	return sharedIntersects(rect, other,
		func(other Object) (accept, contains bool) {
			// switch other := other.(type) {
			// case *Position:
			// 	return true, polyPoint(*other).InsideRect(polyRect(*rect))
			// case *Rect:
			// 	return true, polyRect(*other).InsideRect(polyRect(*rect))
			// case *Point:
			// 	return true, polyPoint(other.Coordinates).InsideRect(
			// 		polyRect(*rect),
			// 	)
			// }
			return false, false
		},
	)
}
