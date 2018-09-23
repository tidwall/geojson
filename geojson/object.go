package geojson

// Object is a GeoJSON object
type Object interface {
	Center() Position
	Rect() Rect
	CachedRect() *Rect
	Contains(other Object) bool
	HasBBox() bool
	ForEachChild(iter func(child Object) bool)
}

var _ = []Object{
	&Position{},
	&Rect{},
	&Point{},
	&GeometryCollection{},
}

// BBox is a GeoJSON bbox member
type objBBox struct {
	rect  Rect
	extra []float64
}

func sharedContains(
	g, other Object,
	primativeContains func(other Object) (accept, contains bool),
) bool {
	if g.HasBBox() {
		// g has a bbox so use that
		gRect := g.Rect()
		return gRect.Contains(other)
	}
	cachedRect := g.CachedRect()
	if cachedRect != nil {
		if !cachedRect.Contains(other) {
			return false
		}
	}
	if other.HasBBox() {
		// other has a bbox so use that
		otherRect := other.Rect()
		other = &otherRect
	} else {
		cachedRect := other.CachedRect()
		if cachedRect != nil {
			if !g.Contains(cachedRect) {
				return false
			}
		}
	}

	accept, contains := primativeContains(other)
	if accept {
		return contains
	}

	contains = true
	other.ForEachChild(func(child Object) bool {
		accept = true
		if !g.Contains(child) {
			contains = false
			return false
		}
		return true
	})
	return accept && contains
}

func sharedIntersects(
	g, other Object,
	primativeIntersects func(other Object) (accept, contains bool),
) bool {
	panic(123)
	// var cachedRect *Rect
	// if g.HasBBox() {
	// 	gRect := g.Rect()
	// 	return gRect.Intersects(other)
	// }
	// cachedRect = g.CachedRect()
	// if cachedRect != nil {
	// 	return cachedRect.Intersects(other)
	// }
	// if other.HasBBox() {
	// 	otherRect := other.Rect()
	// 	other = &otherRect
	// } else {
	// 	cachedRect = other.CachedRect()
	// 	if cachedRect != nil {
	// 		if !g.Contains(cachedRect) {
	// 			return false
	// 		}
	// 	}
	// }

	// accept, contains := primativeContains(other)
	// if accept {
	// 	return contains
	// }

	// contains = true
	// other.ForEachChild(func(child Object) bool {
	// 	accept = true
	// 	if !g.Contains(child) {
	// 		contains = false
	// 		return false
	// 	}
	// 	return true
	// })
	// return accept && contains
}

func collectionContains(col, other Object, testBounds bool) bool {
	if testBounds {
		if col.HasBBox() {
			colRect := col.Rect()
			return colRect.Contains(other)
		}
		cachedRect := col.CachedRect()
		if cachedRect != nil {
			if !cachedRect.Contains(other) {
				return false
			}
		}
	}
	var contains bool
	col.ForEachChild(func(child Object) bool {
		if child.Contains(other) {
			contains = true
			return false
		}
		return true
	})
	return contains
}
