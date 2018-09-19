package geojson

// Rect is a 2D rectangle
type Rect struct {
	// Min is the minimum point
	Min Position
	// Max is the maximum point
	Max Position
}

// ContainsRect test if rect contains other rect
func (rect Rect) ContainsRect(other Rect) bool {
	if other.Min.X < rect.Min.X || other.Min.Y < rect.Min.Y {
		return false
	}
	if other.Max.X > rect.Max.X || other.Max.Y < rect.Max.Y {
		return false
	}
	return true
}

// IntersectsRect test if rect contains other rect
func (rect Rect) IntersectsRect(other Rect) bool {
	if other.Min.X > rect.Max.X || other.Max.X < rect.Min.X {
		return false
	}
	if other.Min.Y > rect.Max.Y || other.Max.Y < rect.Min.Y {
		return false
	}
	return true
}

// IntersectsPolyLine test if rect intersects a polyline
func (rect Rect) IntersectsPolyLine(line []Position) bool {
	if rect.Min == rect.Max {
		return rect.Min.IntersectsPolyLine(line)
	}
	return polyLine(line).LineStringIntersects(polyRect(rect).Polygon(), nil)
}

// Intersects object
func (rect Rect) Intersects(o Object) bool {
	panic("turn rect into a polygon")
}

// Center returns the center point
func (rect Rect) Center() Position {
	return Position{
		X: (rect.Min.X + rect.Max.X) / 2,
		Y: (rect.Min.Y + rect.Max.Y) / 2,
	}
}

// BBox returns a conpatible bbox for rect
func (rect Rect) BBox() BBox {
	return rectBBox{xyBBox{rect}}
}

// rectBBox is a simple wrapper around the xyBBox that returns false for the
// Defined function
type rectBBox struct{ xyBBox }

// Defined return true when the BBox was defined by the parent geojson.
func (bbox rectBBox) Defined() bool {
	return false
}

// Rect returns the 2D rectangle
func (bbox rectBBox) Rect() Rect {
	return bbox.xyBBox.Rect()
}

// JSON is a geojson representation
func (bbox rectBBox) JSON() string {
	return bbox.xyBBox.JSON()
}

// AppendJSON appends the geojson representation to destination.
func (bbox rectBBox) AppendJSON(dst []byte) []byte {
	return bbox.xyBBox.AppendJSON(dst)
}

// calculateRect2 will calculate a rect from a list of positions.
func calculateRect2(coords []Position) Rect {
	var r Rect
	var n int
	for i := 0; i < len(coords); i++ {
		if n == 0 {
			r.Min = coords[i]
			r.Max = coords[i]
		} else {
			if coords[i].X < r.Min.X {
				r.Min.X = coords[i].X
			} else if coords[i].X > r.Max.X {
				r.Max.X = coords[i].X
			}
			if coords[i].Y < r.Min.Y {
				r.Min.Y = coords[i].Y
			} else if coords[i].Y > r.Max.Y {
				r.Max.Y = coords[i].Y
			}
		}
		n++
	}
	return r
}

// calculateRect3 will calculate a rect from a list of positions.
func calculateRect3(coords [][]Position) Rect {
	var r Rect
	var n int
	for i := 0; i < len(coords); i++ {
		for j := 0; j < len(coords[i]); j++ {
			if n == 0 {
				r.Min = coords[i][j]
				r.Max = coords[i][j]
			} else {
				if coords[i][j].X < r.Min.X {
					r.Min.X = coords[i][j].X
				} else if coords[i][j].X > r.Max.X {
					r.Max.X = coords[i][j].X
				}
				if coords[i][j].Y < r.Min.Y {
					r.Min.Y = coords[i][j].Y
				} else if coords[i][j].Y > r.Max.Y {
					r.Max.Y = coords[i][j].Y
				}
			}
			n++
		}
	}
	return r
}

// calculateRect4 will calculate a rect from a list of positions.
func calculateRect4(coords [][][]Position) Rect {
	var r Rect
	var n int
	for i := 0; i < len(coords); i++ {
		for j := 0; j < len(coords[i]); j++ {
			for k := 0; k < len(coords[i][j]); k++ {
				if n == 0 {
					r.Min = coords[i][j][k]
					r.Max = coords[i][j][k]
				} else {
					if coords[i][j][k].X < r.Min.X {
						r.Min.X = coords[i][j][k].X
					} else if coords[i][j][k].X > r.Max.X {
						r.Max.X = coords[i][j][k].X
					}
					if coords[i][j][k].Y < r.Min.Y {
						r.Min.Y = coords[i][j][k].Y
					} else if coords[i][j][k].Y > r.Max.Y {
						r.Max.Y = coords[i][j][k].Y
					}
				}
				n++
			}
		}
	}
	return r
}

// calculateRectObjs will calculate a rect from a list of objects.
func calculateRectObjs(objs []Object) Rect {
	var r Rect
	if len(objs) > 0 {
		r = objs[0].Rect()
		for i := 1; i < len(objs); i++ {
			gr := objs[i].Rect()
			if gr.Min.X < r.Min.X {
				r.Min.X = gr.Min.X
			} else if gr.Max.X > r.Max.X {
				r.Max.X = gr.Max.X
			}
			if gr.Min.Y < r.Min.Y {
				r.Min.Y = gr.Min.Y
			} else if gr.Max.Y > r.Max.Y {
				r.Max.Y = gr.Max.Y
			}
		}
	}
	return r
}
