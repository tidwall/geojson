package geojson

import (
	"strings"

	"github.com/tidwall/gjson"
)

type GeometryCollection struct{ collection }

func NewGeometryCollection(geometries []Object) *GeometryCollection {
	g := new(GeometryCollection)
	g.children = geometries
	g.parseInitRectIndex(DefaultParseOptions)
	return g
}

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
		dst = g.extra.appendJSONExtra(dst, false)
	}
	dst = append(dst, '}')
	strings.Index("", " ")
	return dst
}

func (g *GeometryCollection) String() string {
	return string(g.AppendJSON(nil))
}

func (g *GeometryCollection) JSON() string {
	return string(g.AppendJSON(nil))
}

func (g *GeometryCollection) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

func parseJSONGeometryCollection(
	keys *parseKeys, opts *ParseOptions,
) (Object, error) {
	var g GeometryCollection
	if !keys.rGeometries.Exists() {
		return nil, errGeometriesMissing
	}
	if !keys.rGeometries.IsArray() {
		return nil, errGeometriesInvalid
	}
	var err error
	keys.rGeometries.ForEach(func(key, value gjson.Result) bool {
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
	if err := parseBBoxAndExtras(&g.extra, keys, opts); err != nil {
		return nil, err
	}
	g.parseInitRectIndex(opts)
	return &g, nil
}

func (g *GeometryCollection) Members() string {
	if g.extra != nil {
		return g.extra.members
	}
	return ""
}
