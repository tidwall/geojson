package geojson

import (
	"github.com/tidwall/geojson/geos"
	"github.com/tidwall/gjson"
)

// MultiPolygon ...
type MultiPolygon struct{ collection }

// AppendJSON ...
func (g *MultiPolygon) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"MultiPolygon","coordinates":[`...)
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

func parseJSONMultiPolygon(data string) (Object, error) {
	var g MultiPolygon
	var err error
	rcoords := gjson.Get(data, "coordinates")
	if !rcoords.Exists() {
		return nil, errCoordinatesMissing
	}
	if !rcoords.IsArray() {
		return nil, errCoordinatesInvalid
	}
	var coords [][]geos.Point
	var ex *extra
	rcoords.ForEach(func(_, value gjson.Result) bool {
		coords, ex, err = parseJSONPolygonCoords("", value)
		if err != nil {
			return false
		}
		if len(coords) == 0 {
			err = errCoordinatesInvalid // must be a linear ring
			return false
		}
		for _, p := range coords {
			if len(p) < 4 || p[0] != p[len(p)-1] {
				err = errCoordinatesInvalid // must be a linear ring
				return false
			}
		}
		exterior := coords[0]
		var holes [][]geos.Point
		if len(coords) > 1 {
			holes = coords[1:]
		}
		poly := geos.NewPoly(exterior, holes)

		g.children = append(g.children, &Polygon{base: *poly, extra: ex})
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
