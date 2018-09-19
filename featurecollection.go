package geojson

import "github.com/tidwall/gjson"

// FeatureCollection is a GeoJSON FeatureCollection
type FeatureCollection struct {
	Features []Object
	BBox     BBox
}

// Rect returns a rectangle that contains the entire object
func (g FeatureCollection) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	return calcRectFromObjects(g.Features)
}

// Center is the center-most point of the object
func (g FeatureCollection) Center() Position {
	return g.Rect().Center()
}

// AppendJSON appends a json representation to destination
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

// loadJSONFeatureCollection will return a valid GeoJSON object.
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
		f, err = Load(value.Raw)
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
