package geojson

import (
	"github.com/tidwall/geojson/geos"
	"github.com/tidwall/gjson"
)

// MultiPoint ...
type MultiPoint struct{ collection }

// AppendJSON ...
func (g *MultiPoint) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"MultiPoint","coordinates":[`...)
	for i, g := range g.children {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = append(dst,
			gjson.GetBytes(g.AppendJSON(nil), "coordinates").String()...)
	}
	dst = append(dst, ']')
	if g.extra != nil {
		dst = g.extra.appendJSONExtra(dst)
	}
	dst = append(dst, '}')
	return dst
}

func parseJSONMultiPoint(keys *parseKeys, opts *ParseOptions) (Object, error) {
	var g MultiPoint
	var err error
	if !keys.rCoordinates.Exists() {
		return nil, errCoordinatesMissing
	}
	if !keys.rCoordinates.IsArray() {
		return nil, errCoordinatesInvalid
	}
	var coords geos.Point
	var ex *extra
	keys.rCoordinates.ForEach(func(_, value gjson.Result) bool {
		coords, ex, err = parseJSONPointCoords(keys, value, opts)
		if err != nil {
			return false
		}
		g.children = append(g.children, &Point{base: coords, extra: ex})
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
