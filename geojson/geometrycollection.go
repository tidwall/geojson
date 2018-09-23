package geojson

// GeometryCollection GeoJSON type
type GeometryCollection struct {
	Geometries []Object
	bbox       *objBBox
	cachedRect *Rect
}

// Center returns the center position of the point
func (gcol *GeometryCollection) Center() Position {
	rect := gcol.Rect()
	return rect.Center()
}

// Rect returns the boundary
func (gcol *GeometryCollection) Rect() Rect {
	if gcol.bbox != nil {
		return gcol.bbox.rect
	}
	if gcol.cachedRect != nil {
		return *gcol.cachedRect
	}
	return calcRectFromChildren(gcol)
}

// HasBBox returns true if theres a GeoJSON bbox member
func (gcol *GeometryCollection) HasBBox() bool {
	return gcol.bbox != nil
}

// CachedRect returns a precaclulated rectangle, if any
func (gcol *GeometryCollection) CachedRect() *Rect {
	return gcol.cachedRect
}

// ForEachChild iterates over child objects.
func (gcol *GeometryCollection) ForEachChild(iter func(child Object) bool) {
	for _, child := range gcol.Geometries {
		if !iter(child) {
			return
		}
	}
}

// Contains returns true if object contains other object
func (gcol *GeometryCollection) Contains(other Object) bool {
	return collectionContains(gcol, other, true)
}
