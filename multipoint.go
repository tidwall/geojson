package geojson

import (
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/gjson"
)

type MultiPoint struct{ collection }

func NewMultiPoint(points []geometry.Point) *MultiPoint {
	g := new(MultiPoint)
	for _, point := range points {
		g.children = append(g.children, NewPoint(point))
	}
	g.parseInitRectIndex(DefaultParseOptions)
	return g
}

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
		dst = g.extra.appendJSONExtra(dst, false)
	}
	dst = append(dst, '}')
	return dst
}

func (g *MultiPoint) String() string {
	return string(g.AppendJSON(nil))
}

func (g *MultiPoint) JSON() string {
	return string(g.AppendJSON(nil))
}

func (g *MultiPoint) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

func (g *MultiPoint) AppendBinary(dst []byte) []byte {
	dst = append(dst, ':', binMultiPoint)
	return appendBinaryCollection(dst, g.collection)
}

func (g *MultiPoint) Binary() []byte {
	return g.AppendBinary(nil)
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
	var coords geometry.Point
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

func parseBinaryMultiPointObject(src []byte, opts *ParseOptions) (*MultiPoint, int) {
	mark := len(src)
	c, n := parseBinaryCollection(src, opts)
	if n <= 0 {
		return nil, 0
	}
	src = src[n:]
	g := &MultiPoint{collection: c}
	return g, mark - len(src)
}
