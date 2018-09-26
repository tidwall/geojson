package poly

// Rect ...
type Rect struct {
	Min, Max Point
}

func expandRect(rect *Rect, point Point) {
	if point.X < rect.Min.X {
		rect.Min.X = point.X
	} else if point.X > rect.Max.X {
		rect.Max.X = point.X
	}
	if point.Y < rect.Min.Y {
		rect.Min.Y = point.Y
	} else if point.Y > rect.Max.Y {
		rect.Max.Y = point.Y
	}
}

func (rect Rect) ringPoints() [5]Point {
	return [5]Point{
		Point{rect.Min.X, rect.Min.Y},
		Point{rect.Max.X, rect.Min.Y},
		Point{rect.Max.X, rect.Max.Y},
		Point{rect.Min.X, rect.Max.Y},
		Point{rect.Min.X, rect.Min.Y},
	}
}

// ContainsPoint ...
func (rect Rect) ContainsPoint(point Point) bool {
	return point.IntersectsRect(rect)
}

// ContainsRect ...
func (rect Rect) ContainsRect(other Rect) bool {
	if other.Min.X < rect.Min.X || other.Max.X > rect.Max.X {
		return false
	}
	if other.Min.Y < rect.Min.Y || other.Max.Y > rect.Max.Y {
		return false
	}
	return true
}

// ContainsLine ...
func (rect Rect) ContainsLine(line Line) bool {
	if len(line) == 0 {
		return false
	}
	for _, point := range line {
		if !rect.ContainsPoint(point) {
			return false
		}
	}
	return true
}

// ContainsRing ...
func (rect Rect) ContainsRing(ring Ring) bool {
	return rect.ContainsLine(Line(ring))
}

// ContainsPolygon ...
func (rect Rect) ContainsPolygon(polygon Polygon) bool {
	return rect.ContainsRing(polygon.Exterior)
}

// IntersectsPoint ...
func (rect Rect) IntersectsPoint(point Point) bool {
	return point.IntersectsRect(rect)
}

// IntersectsRect ...
func (rect Rect) IntersectsRect(other Rect) bool {
	if rect.Min.Y > other.Max.Y || rect.Max.Y < other.Min.Y {
		return false
	}
	if rect.Min.X > other.Max.X || rect.Max.X < other.Min.X {
		return false
	}
	return true
}

func (rect Rect) intersectsPoints(points []Point) (accept, intersects bool) {
	if len(points) == 0 {
		return true, false
	}
	// test is rect contains any points and generate the rect along the way
	var pointsRect Rect
	for i, point := range points {
		if rect.ContainsPoint(point) {
			return true, true
		}
		if i == 0 {
			pointsRect = Rect{point, point}
		} else {
			expandRect(&pointsRect, point)
		}
	}
	// make sure that both rects intersect
	if !rect.IntersectsRect(pointsRect) {
		return true, false
	}
	// check if the points rect fully contains the points rect
	if pointsRect.ContainsRect(rect) {
		return true, true
	}
	// points intersect but we don't know about the underlying shape yet
	return false, false
}

// IntersectsLine ...
func (rect Rect) IntersectsLine(line Line) bool {
	accept, intersects := rect.intersectsPoints(line)
	if accept {
		return intersects
	}
	// do the slow segment edge checks
	rectRing := rect.ringPoints()
	return algoAnySegmentIntersects(rectRing[:], line, true, false)
}

// IntersectsRing ...
func (rect Rect) IntersectsRing(ring Ring) bool {
	accept, intersects := rect.intersectsPoints(ring)
	if accept {
		return intersects
	}
	if len(ring) == 0 {
		return false
	}
	rectRing := rect.ringPoints()
	return algoAnySegmentIntersects(rectRing[:], ring, true, true)
}

// // IntersectsPolygon ...
// func (rect Rect) IntersectsPolygon(polygon Polygon) bool {
// 	if rect.ContainsPolygon(polygon) {
// 		return true
// 	}
// 	rectRing := rect.ringPoints()
// 	return algoAnySegmentIntersects(rectRing[:], polygon.Exterior, true, true)
// }
