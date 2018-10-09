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
		dst = g.extra.appendJSONBBox(dst)
	}
	dst = append(dst, '}')
	return dst
}

func parseJSONMultiPoint(data string, opts *ParseOptions) (Object, error) {
	var g MultiPoint
	var err error
	rcoords := gjson.Get(data, "coordinates")
	if !rcoords.Exists() {
		return nil, errCoordinatesMissing
	}
	if !rcoords.IsArray() {
		return nil, errCoordinatesInvalid
	}
	var coords geos.Point
	var ex *extra
	rcoords.ForEach(func(_, value gjson.Result) bool {
		coords, ex, err = parseJSONPointCoords("", value, opts)
		if err != nil {
			return false
		}
		g.children = append(g.children, &Point{base: coords, extra: ex})
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
