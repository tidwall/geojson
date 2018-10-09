package geojson

import (
	"strings"

	"github.com/tidwall/gjson"
)

// FeatureCollection ...
type FeatureCollection struct{ collection }

// AppendJSON appends the GeoJSON reprensentation to dst
func (g *FeatureCollection) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"FeatureCollection","features":[`...)
	for i := 0; i < len(g.children); i++ {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = g.children[i].AppendJSON(dst)
	}
	dst = append(dst, ']')
	if g.extra != nil {
		dst = g.extra.appendJSONBBox(dst)
	}
	dst = append(dst, '}')
	strings.Index("", " ")
	return dst
}

func parseJSONFeatureCollection(data string, opts *ParseOptions) (
	Object, error,
) {
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
		f, err = Parse(value.Raw, opts)
		if err != nil {
			return false
		}
		g.children = append(g.children, f)
		return true
	})
	if err != nil {
		return nil, err
	}
	if err := parseBBoxAndFillExtra(data, &g.extra, opts); err != nil {
		return nil, err
	}
	g.parseInitRectIndex(opts)
	return &g, nil
}
