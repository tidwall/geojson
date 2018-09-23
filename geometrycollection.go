package geojson

import "github.com/tidwall/gjson"

// GeometryCollection GeoJSON type
type GeometryCollection struct {
	Geometries []Object
	BBox       BBox
}

// BBoxDefined return true if there is a defined GeoJSON "bbox" member
func (g GeometryCollection) BBoxDefined() bool {
	return g.BBox != nil && g.BBox.Defined()
}

// Rect returns the outer minimum bounding rectangle
func (g GeometryCollection) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	return calcRectFromObjects(g.Geometries)
}

// Center returns the center position of the object
func (g GeometryCollection) Center() Position {
	return g.Rect().Center()
}

// AppendJSON appends the GeoJSON reprensentation to dst
func (g GeometryCollection) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"GeometryCollection","geometries":[`...)
	for i := 0; i < len(g.Geometries); i++ {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = g.Geometries[i].AppendJSON(dst)
	}
	dst = append(dst, ']')
	if g.BBox != nil && g.BBox.Defined() {
		dst = append(dst, `,"bbox":`...)
		dst = g.BBox.AppendJSON(dst)
	}
	dst = append(dst, '}')
	return dst
}

// ForEachChild iterates over child objects.
func (g GeometryCollection) ForEachChild(iter func(child Object) bool) {
	for _, child := range g.Geometries {
		if !iter(child) {
			return
		}
	}
}

// Contains returns true if object contains other object
func (g GeometryCollection) Contains(other Object) bool {
	return collectionContains(g, other)
}

// Intersects returns true if object intersects with other object
func (g GeometryCollection) Intersects(other Object) bool {
	return collectionIntersects(g, other)
}

func loadJSONGeometryCollection(data string) (Object, error) {
	var g GeometryCollection
	rgeometries := gjson.Get(data, "geometries")
	if !rgeometries.Exists() {
		return nil, errGeometriesMissing
	}
	if !rgeometries.IsArray() {
		return nil, errGeometriesInvalid
	}
	var err error
	rgeometries.ForEach(func(key, value gjson.Result) bool {
		var f Object
		f, err = Parse(value.Raw)
		if err != nil {
			return false
		}
		g.Geometries = append(g.Geometries, f)
		return true
	})
	if err != nil {
		return nil, err
	}
	g.BBox, err = loadBBox(data)
	if err != nil {
		return nil, err
	}
	if g.BBox == nil {
		g.BBox = bboxRect{g.Rect()}
	}
	return g, nil
}
