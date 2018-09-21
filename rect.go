package geojson

type Rect struct {
	Min, Max Position
}

func (rect Rect) HasBBox() bool {
	return false
}

func (rect Rect) Rect() Rect {
	return rect
}

// Center returns the center point
func (rect Rect) Center() Position {
	return Position{
		X: (rect.Min.X + rect.Max.X) / 2,
		Y: (rect.Min.Y + rect.Max.Y) / 2,
	}
}

func (rect Rect) AppendJSON(dst []byte) []byte {
	if rect.Min == rect.Max {
		return rect.Min.AppendJSON(dst)
	}
	return (Polygon{Coordinates: [][]Position{[]Position{
		rect.Min,
		Position{rect.Max.X, rect.Min.Y},
		rect.Max,
		Position{rect.Min.X, rect.Max.Y},
		rect.Min,
	}}}).AppendJSON(dst)
}
func (rect Rect) ForEach(func(child Object) bool) {

}

func (rect Rect) Expand(posn Position) Rect {
	if posn.X < rect.Min.X {
		rect.Min.X = posn.X
	} else if posn.X > rect.Max.X {
		rect.Max.X = posn.X
	}
	if posn.Y < rect.Min.Y {
		rect.Min.Y = posn.Y
	} else if posn.Y > rect.Max.Y {
		rect.Max.Y = posn.Y
	}
	return rect
}

func (rect Rect) Union(rect2 Rect) Rect {
	if rect2.Min.X < rect.Min.X {
		rect.Min.X = rect2.Min.X
	} else if rect2.Max.X > rect.Max.X {
		rect.Max.X = rect2.Max.X
	}
	if rect2.Min.Y < rect.Min.Y {
		rect.Min.Y = rect2.Min.Y
	} else if rect2.Max.Y > rect.Max.Y {
		rect.Max.Y = rect2.Max.Y
	}
	return rect
}

func (rect Rect) ContainsRect(other Rect) bool {
	if other.Min.X < rect.Min.X || other.Max.X > rect.Max.X {
		return false
	}
	if other.Min.Y < rect.Min.Y || other.Max.Y > rect.Max.Y {
		return false
	}
	return true
}

func (rect Rect) IntersectsRect(other Rect) bool {
	if other.Min.X > rect.Max.X || other.Max.X < rect.Min.X {
		return false
	}
	if other.Min.Y > rect.Max.Y || other.Max.Y < rect.Min.Y {
		return false
	}
	return true
}

func (rect Rect) ContainsPosition(posn Position) bool {
	return posn.X >= rect.Min.X && posn.X <= rect.Max.X &&
		posn.Y >= rect.Min.Y && posn.Y <= rect.Max.Y
}

func (rect Rect) Contains(other Object) bool {
	// basic types
	switch other := other.(type) {
	case Position:
		return rect.ContainsPosition(other)
	case Rect:
		return rect.ContainsRect(other)
	}
	// entire inner bounds must be fully contained inside of rect.
	return rect.ContainsRect(other.Rect())
}

func (rect Rect) Intersects(other Object) bool {
	// simple types
	switch other := other.(type) {
	case Position:
		return rect.ContainsPosition(other)
	case Rect:
		return rect.IntersectsRect(other)
	}
	// bbox types
	if !rect.IntersectsRect(other.Rect()) {
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
		return polyPoint(other.Coordinates).InsideRect(polyRect(rect))
	case LineString:
		return polyLine(other.Coordinates).LineStringIntersectsPolygon(
			polyRect(rect).Polygon(), nil,
		)
	case Polygon:
		return polyRect(rect).IntersectsPolygon(polyPolygon(other.Coordinates))
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
}

func calcRectFromObjects(objs []Object) Rect {
	var r Rect
	for i := 0; i < len(objs); i++ {
		r2 := objs[i].Rect()
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
	}
	return r
}

// bboxRect is a simple wrapper around the xyBBox that returns false for the
// Defined function
type bboxRect struct{ rect Rect }

// Defined return true when the BBox was defined by the parent geojson.
func (bbox bboxRect) Defined() bool {
	return false
}

// Rect returns the 2D rectangle
func (bbox bboxRect) Rect() Rect {
	return bbox.rect.Rect()
}

// AppendJSON appends the geojson representation to destination.
func (bbox bboxRect) AppendJSON(dst []byte) []byte {
	return bbox.rect.AppendJSON(dst)
}
