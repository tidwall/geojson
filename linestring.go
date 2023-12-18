package geojson

import (
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/gjson"
)

type LineString struct {
	base  geometry.Line
	extra *extra
}

func NewLineString(line *geometry.Line) *LineString {
	return &LineString{base: *line}
}

func (g *LineString) Empty() bool {
	return g.base.Empty()
}

func (g *LineString) Valid() bool {
	return g.base.Valid()
}

func (g *LineString) Rect() geometry.Rect {
	return g.base.Rect()
}

func (g *LineString) Center() geometry.Point {
	return g.Rect().Center()
}

func (g *LineString) Base() *geometry.Line {
	return &g.base
}

func (g *LineString) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"LineString","coordinates":`...)
	dst, _ = appendJSONSeries(dst, &g.base, g.extra, 0)
	if g.extra != nil {
		dst = g.extra.appendJSONExtra(dst, false)
	}
	dst = append(dst, '}')
	return dst
}

func (g *LineString) String() string {
	return string(g.AppendJSON(nil))
}

func (g *LineString) JSON() string {
	return string(g.AppendJSON(nil))
}

func (g *LineString) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

func (g *LineString) Spatial() Spatial {
	return g
}

func (g *LineString) ForEach(iter func(geom Object) bool) bool {
	return iter(g)
}

func (g *LineString) Within(obj Object) bool {
	return obj.Contains(g)
}

func (g *LineString) Contains(obj Object) bool {
	return obj.Spatial().WithinLine(&g.base)
}

func (g *LineString) Intersects(obj Object) bool {
	return obj.Spatial().IntersectsLine(&g.base)
}

func (g *LineString) WithinRect(rect geometry.Rect) bool {
	return rect.ContainsLine(&g.base)
}

func (g *LineString) WithinPoint(point geometry.Point) bool {
	return point.ContainsLine(&g.base)
}

func (g *LineString) WithinLine(line *geometry.Line) bool {
	return line.ContainsLine(&g.base)
}

func (g *LineString) WithinPoly(poly *geometry.Poly) bool {
	return poly.ContainsLine(&g.base)
}

func (g *LineString) IntersectsPoint(point geometry.Point) bool {
	return g.base.IntersectsPoint(point)
}

func (g *LineString) IntersectsRect(rect geometry.Rect) bool {
	return g.base.IntersectsRect(rect)
}

func (g *LineString) IntersectsLine(line *geometry.Line) bool {
	return g.base.IntersectsLine(line)
}

func (g *LineString) IntersectsPoly(poly *geometry.Poly) bool {
	return g.base.IntersectsPoly(poly)
}

func (g *LineString) NumPoints() int {
	return g.base.NumPoints()
}

func parseJSONLineString(keys *parseKeys, opts *ParseOptions) (Object, error) {
	var g LineString
	points, ex, err := parseJSONLineStringCoords(keys, gjson.Result{}, opts)
	if err != nil {
		return nil, err
	}
	if len(points) < 2 {
		// Must have at least two points
		// https://tools.ietf.org/html/rfc7946#section-3.1.4
		return nil, errCoordinatesInvalid
	}
	gopts := toGeometryOpts(opts)
	line := geometry.NewLine(points, &gopts)
	g.base = *line
	g.extra = ex
	if err := parseBBoxAndExtras(&g.extra, keys, opts); err != nil {
		return nil, err
	}
	if opts.RequireValid {
		if !g.Valid() {
			return nil, errDataInvalid
		}
	}
	return &g, nil
}

func parseJSONLineStringCoords(
	keys *parseKeys, rcoords gjson.Result, opts *ParseOptions,
) ([]geometry.Point, *extra, error) {
	var err error
	var coords []geometry.Point
	var ex *extra
	var dims int
	if !rcoords.Exists() {
		rcoords = keys.rCoordinates
		if !rcoords.Exists() {
			return nil, nil, errCoordinatesMissing
		}
		if !rcoords.IsArray() {
			return nil, nil, errCoordinatesInvalid
		}
	}
	rcoords.ForEach(func(key, value gjson.Result) bool {
		if !value.IsArray() {
			err = errCoordinatesInvalid
			return false
		}
		var count int
		var nums [4]float64
		value.ForEach(func(key, value gjson.Result) bool {
			if count == 4 {
				return false
			}
			if value.Type != gjson.Number {
				err = errCoordinatesInvalid
				return false
			}
			nums[count] = value.Float()
			count++
			return true
		})
		if err != nil {
			return false
		}
		if count < 2 {
			err = errCoordinatesInvalid
			return false
		}
		coords = append(coords, geometry.Point{X: nums[0], Y: nums[1]})
		if ex == nil {
			if count > 2 {
				if len(coords) > 1 {
					err = errCoordinatesInvalid
					return false
				}
				ex = new(extra)
				if count > 3 {
					ex.dims = 2
				} else {
					ex.dims = 1
				}
				dims = int(ex.dims)
			}
		}
		if ex != nil {
			for i := 0; i < dims; i++ {
				ex.values = append(ex.values, nums[2+i])
			}
		}
		return true
	})
	if err != nil {
		return nil, nil, err
	}
	return coords, ex, err
}

func (g *LineString) Distance(obj Object) float64 {
	return obj.Spatial().DistanceLine(&g.base)
}

func (g *LineString) DistancePoint(point geometry.Point) float64 {
	return geoDistancePoints(g.Center(), point)
}

// DistanceRect ..
func (g *LineString) DistanceRect(rect geometry.Rect) float64 {
	return geoDistancePoints(g.Center(), rect.Center())
}

func (g *LineString) DistanceLine(line *geometry.Line) float64 {
	return geoDistancePoints(g.Center(), line.Rect().Center())
}

func (g *LineString) DistancePoly(poly *geometry.Poly) float64 {
	return geoDistancePoints(g.Center(), poly.Rect().Center())
}

func (g *LineString) Members() string {
	if g.extra != nil {
		return g.extra.members
	}
	return ""
}
