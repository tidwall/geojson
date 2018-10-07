package geojson

import "github.com/tidwall/gjson"

// FeatureCollection GeoJSON type
type FeatureCollection struct {
	Features []Object
	BBox     BBox
}

// BBoxDefined return true if there is a defined GeoJSON "bbox" member
func (g FeatureCollection) BBoxDefined() bool {
	return g.BBox != nil && g.BBox.Defined()
}

// Rect returns the outer minimum bounding rectangle
func (g FeatureCollection) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	return calcRectFromObjects(g.Features)
}

// Center returns the center position of the object
func (g FeatureCollection) Center() Position {
	return g.Rect().Center()
}

// AppendJSON appends the GeoJSON reprensentation to dst
func (g FeatureCollection) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"FeatureCollection","features":[`...)
	for i := 0; i < len(g.Features); i++ {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = g.Features[i].AppendJSON(dst)
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
func (g FeatureCollection) ForEachChild(iter func(child Object) bool) {
	for _, child := range g.Features {
		if !iter(child) {
			return
		}
	}
}

// Within is the inverse of contains
func (g FeatureCollection) Within(other Object) bool {
	return other.Contains(g)
}

// Contains returns true if object contains other object
func (g FeatureCollection) Contains(other Object) bool {
	return collectionContains(g, other)
}

// Intersects returns true if object intersects with other object
func (g FeatureCollection) Intersects(other Object) bool {
	return collectionIntersects(g, other)
}

func loadJSONFeatureCollection(data string) (Object, error) {
	var g FeatureCollection
	rFeatures := gjson.Get(data, "features")
	if !rFeatures.Exists() {
		return nil, errFeaturesMissing
	}
	if !rFeatures.IsArray() {
		return nil, errFeaturesInvalid
	}
	var err error
	rFeatures.ForEach(func(key, value gjson.Result) bool {
		var f Object
		f, err = Parse(value.Raw)
		if err != nil {
			return false
		}
		g.Features = append(g.Features, f)
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
