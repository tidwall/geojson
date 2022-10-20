package geojson

import (
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/gjson"
)

type MultiPolygon struct{ collection }

func NewMultiPolygon(polys []*geometry.Poly) *MultiPolygon {
	g := new(MultiPolygon)
	for _, poly := range polys {
		g.children = append(g.children, NewPolygon(poly))
	}
	g.parseInitRectIndex(DefaultParseOptions)
	return g
}

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
		dst = g.extra.appendJSONExtra(dst, false)
	}
	dst = append(dst, '}')
	return dst
}

func (g *MultiPolygon) String() string {
	return string(g.AppendJSON(nil))
}

func (g *MultiPolygon) Valid() bool {
	valid := true
	for _, p := range g.children {
		if !p.Valid() {
			valid = false
		}
	}
	return valid
}

func (g *MultiPolygon) JSON() string {
	return string(g.AppendJSON(nil))
}

func (g *MultiPolygon) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

func parseJSONMultiPolygon(
	keys *parseKeys, opts *ParseOptions,
) (Object, error) {
	var g MultiPolygon
	var err error
	if !keys.rCoordinates.Exists() {
		return nil, errCoordinatesMissing
	}
	if !keys.rCoordinates.IsArray() {
		return nil, errCoordinatesInvalid
	}
	var coords [][]geometry.Point
	var ex *extra
	keys.rCoordinates.ForEach(func(_, value gjson.Result) bool {
		coords, ex, err = parseJSONPolygonCoords(keys, value, opts)
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
		var holes [][]geometry.Point
		if len(coords) > 1 {
			holes = coords[1:]
		}
		gopts := toGeometryOpts(opts)
		poly := geometry.NewPoly(exterior, holes, &gopts)
		g.children = append(g.children, &Polygon{base: *poly, extra: ex})
		return true
	})
	if err != nil {
		return nil, err
	}
	if err := parseBBoxAndExtras(&g.extra, keys, opts); err != nil {
		return nil, err
	}
	if opts.RequireValid {
		if !g.Valid() {
			return nil, errCoordinatesInvalid
		}
	}
	g.parseInitRectIndex(opts)
	return &g, nil
}

func (g *MultiPolygon) Members() string {
	if g.extra != nil {
		return g.extra.members
	}
	return ""
}
