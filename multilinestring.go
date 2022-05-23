package geojson

import (
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/gjson"
)

type MultiLineString struct{ collection }

func NewMultiLineString(lines []*geometry.Line) *MultiLineString {
	g := new(MultiLineString)
	for _, line := range lines {
		g.children = append(g.children, NewLineString(line))
	}
	g.parseInitRectIndex(DefaultParseOptions)
	return g
}

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
		dst = g.extra.appendJSONExtra(dst, false)
	}
	dst = append(dst, '}')
	return dst

}

func (g *MultiLineString) String() string {
	return string(g.AppendJSON(nil))
}

func (g *MultiLineString) Valid() bool {
	valid := true
	for _, p := range g.children {
		if !p.Valid() {
			valid = false
		}
	}
	return valid
}

func (g *MultiLineString) JSON() string {
	return string(g.AppendJSON(nil))
}

func (g *MultiLineString) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}
func (g *MultiLineString) AppendBinary(dst []byte) []byte {
	dst = append(dst, ':', binMultiLineString)
	return appendBinaryCollection(dst, g.collection)
}

func (g *MultiLineString) Binary() []byte {
	return g.AppendBinary(nil)
}

func parseJSONMultiLineString(
	keys *parseKeys, opts *ParseOptions,
) (Object, error) {
	var g MultiLineString
	var err error
	if !keys.rCoordinates.Exists() {
		return nil, errCoordinatesMissing
	}
	if !keys.rCoordinates.IsArray() {
		return nil, errCoordinatesInvalid
	}
	var coords []geometry.Point
	var ex *extra
	keys.rCoordinates.ForEach(func(_, value gjson.Result) bool {
		coords, ex, err = parseJSONLineStringCoords(keys, value, opts)
		if err != nil {
			return false
		}
		if len(coords) < 2 {
			err = errCoordinatesInvalid
			return false
		}
		gopts := toGeometryOpts(opts)
		line := geometry.NewLine(coords, &gopts)
		g.children = append(g.children, &LineString{base: *line, extra: ex})
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

func parseBinaryMultiLineStringObject(src []byte, opts *ParseOptions) (*MultiLineString, int) {
	mark := len(src)
	c, n := parseBinaryCollection(src, opts)
	if n <= 0 {
		return nil, 0
	}
	src = src[n:]
	g := &MultiLineString{collection: c}
	return g, mark - len(src)
}
