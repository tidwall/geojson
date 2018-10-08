package geojson

import (
	"strings"

	"github.com/tidwall/gjson"
)

// GeometryCollection ...
type GeometryCollection struct{ collection }

// AppendJSON appends the GeoJSON reprensentation to dst
func (g *GeometryCollection) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"GeometryCollection","geometries":[`...)
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

func parseJSONGeometryCollection(data string) (Object, error) {
	var g GeometryCollection
	rGeometries := gjson.Get(data, "geometries")
	if !rGeometries.Exists() {
		return nil, errGeometriesMissing
	}
	if !rGeometries.IsArray() {
		return nil, errGeometriesInvalid
	}
	var err error
	rGeometries.ForEach(func(key, value gjson.Result) bool {
		var f Object
		f, err = Parse(value.Raw)
		if err != nil {
			return false
		}
		g.children = append(g.children, f)
		return true
	})
	if err != nil {
		return nil, err
	}
	if err := parseBBoxAndFillExtra(data, &g.extra); err != nil {
		return nil, err
	}
	g.initRectIndex()
	return &g, nil
}
