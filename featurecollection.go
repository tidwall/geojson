package geojson

import "github.com/tidwall/gjson"

type FeatureCollection struct {
	Features []Object
	BBox     BBox
}

func (g FeatureCollection) BBoxDefined() bool {
	return g.BBox != nil && g.BBox.Defined()
}

func (g FeatureCollection) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	return calcRectFromObjects(g.Features)
}

func (g FeatureCollection) Center() Position {
	return g.Rect().Center()
}

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
func (g FeatureCollection) ForEach(iter func(child Object) bool) {
	for _, child := range g.Features {
		if !iter(child) {
			return
		}
	}
}

func (g FeatureCollection) Contains(other Object) bool {
	return collectionObjectContains(g, other)
}

func (g FeatureCollection) Intersects(other Object) bool {
	return collectionObjectIntersects(g, other)
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
