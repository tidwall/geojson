package geojson

import (
	"github.com/tidwall/geojson/geos"
	"github.com/tidwall/gjson"
)

// MultiLineString ...
type MultiLineString struct{ collection }

// AppendJSON ...
func (g *MultiLineString) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"MultiLineString","coordinates":[`...)
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

func parseJSONMultiLineString(data string, opts *ParseOptions) (Object, error) {
	var g MultiLineString
	var err error
	rcoords := gjson.Get(data, "coordinates")
	if !rcoords.Exists() {
		return nil, errCoordinatesMissing
	}
	if !rcoords.IsArray() {
		return nil, errCoordinatesInvalid
	}
	var coords []geos.Point
	var ex *extra
	rcoords.ForEach(func(_, value gjson.Result) bool {
		coords, ex, err = parseJSONLineStringCoords("", value, opts)
		if err != nil {
			return false
		}
		if len(coords) < 2 {
			err = errCoordinatesInvalid
			return false
		}
		line := geos.NewLine(coords, opts.IndexGeometry)
		g.children = append(g.children, &LineString{base: *line, extra: ex})
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
