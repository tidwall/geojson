package geojson

type Rect struct {
	Min, Max Position
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
