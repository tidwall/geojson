package geojson

import (
	"math"

	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/gjson"
)

type Point struct {
	base  geometry.Point
	extra *extra
}

func NewPoint(point geometry.Point) *Point {
	return &Point{base: point}
}

func NewPointZ(point geometry.Point, z float64) *Point {
	return &Point{
		base:  point,
		extra: &extra{dims: 1, values: []float64{z}},
	}
}

func (g *Point) ForEach(iter func(geom Object) bool) bool {
	return iter(g)
}

func (g *Point) Empty() bool {
	return g.base.Empty()
}

func (g *Point) Valid() bool {
	return g.base.Valid()
}

func (g *Point) Rect() geometry.Rect {
	return g.base.Rect()
}

func (g *Point) Spatial() Spatial {
	return g
}

func (g *Point) Center() geometry.Point {
	return g.base
}

func (g *Point) Base() geometry.Point {
	return g.base
}

func (g *Point) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Point","coordinates":`...)
	dst = appendJSONPoint(dst, g.base, g.extra, 0)
	dst = g.extra.appendJSONExtra(dst, false)
	dst = append(dst, '}')
	return dst
}

func (g *Point) JSON() string {
	return string(g.AppendJSON(nil))
}

func (g *Point) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

func (g *Point) String() string {
	return string(g.AppendJSON(nil))
}

func (g *Point) Within(obj Object) bool {
	return obj.Contains(g)
}

func (g *Point) Contains(obj Object) bool {
	return obj.Spatial().WithinPoint(g.base)
}

func (g *Point) Intersects(obj Object) bool {
	if obj, ok := obj.(*Circle); ok {
		return obj.Contains(g)
	}
	return obj.Spatial().IntersectsPoint(g.base)
}

func (g *Point) WithinRect(rect geometry.Rect) bool {
	return rect.ContainsPoint(g.base)
}

func (g *Point) WithinPoint(point geometry.Point) bool {
	return point.ContainsPoint(g.base)
}

func (g *Point) WithinLine(line *geometry.Line) bool {
	return line.ContainsPoint(g.base)
}

func (g *Point) WithinPoly(poly *geometry.Poly) bool {
	return poly.ContainsPoint(g.base)
}

func (g *Point) IntersectsPoint(point geometry.Point) bool {
	return g.base.IntersectsPoint(point)
}

func (g *Point) IntersectsRect(rect geometry.Rect) bool {
	return g.base.IntersectsRect(rect)
}

func (g *Point) IntersectsLine(line *geometry.Line) bool {
	return g.base.IntersectsLine(line)
}

func (g *Point) IntersectsPoly(poly *geometry.Poly) bool {
	return g.base.IntersectsPoly(poly)
}

func (g *Point) NumPoints() int {
	return 1
}

func (g *Point) Z() float64 {
	if g.extra != nil && len(g.extra.values) > 0 {
		return g.extra.values[0]
	}
	return 0
}

func parseJSONPoint(keys *parseKeys, opts *ParseOptions) (Object, error) {
	var o Object
	base, extra, err := parseJSONPointCoords(keys, gjson.Result{}, opts)
	if err != nil {
		return nil, err
	}
	if err := parseBBoxAndExtras(&extra, keys, opts); err != nil {
		return nil, err
	}
	if extra == nil && opts.AllowSimplePoints {
		var g SimplePoint
		g.Point = base
		o = &g
	} else {
		var g Point
		g.base = base
		g.extra = extra
		o = &g
	}
	if opts.RequireValid {
		if !o.Valid() {
			return nil, errCoordinatesInvalid
		}
	}
	return o, nil
}

func parseJSONPointCoords(
	keys *parseKeys, rcoords gjson.Result, opts *ParseOptions,
) (geometry.Point, *extra, error) {
	var coords geometry.Point
	var ex *extra
	if !rcoords.Exists() {
		rcoords = keys.rCoordinates
		if !rcoords.Exists() {
			return coords, nil, errCoordinatesMissing
		}
		if !rcoords.IsArray() {
			return coords, nil, errCoordinatesInvalid
		}
	}
	var err error
	var count int
	var nums [4]float64
	rcoords.ForEach(func(key, value gjson.Result) bool {
		if count == 4 {
			return false
		}
		if value.Type != gjson.Number {
			if value.Type == gjson.Null {
				// convert "null" to "NaN"
				nums[count] = math.NaN()
				count++
				return true
			}
			err = errCoordinatesInvalid
			return false
		}
		nums[count] = value.Float()
		count++
		return true
	})
	if err != nil {
		return coords, nil, err
	}
	if count < 2 {
		return coords, nil, errCoordinatesInvalid
	}
	coords = geometry.Point{X: nums[0], Y: nums[1]}
	if count > 2 {
		ex = new(extra)
		if count > 3 {
			ex.dims = 2
		} else {
			ex.dims = 1
		}
		ex.values = make([]float64, count-2)
		for i := 2; i < count; i++ {
			ex.values[i-2] = nums[i]
		}
	}
	return coords, ex, nil
}

func (g *Point) Distance(obj Object) float64 {
	return obj.Spatial().DistancePoint(g.base)
}

func (g *Point) DistancePoint(point geometry.Point) float64 {
	return geoDistancePoints(g.Center(), point)
}

func (g *Point) DistanceRect(rect geometry.Rect) float64 {
	return geoDistancePoints(g.Center(), rect.Center())
}

func (g *Point) DistanceLine(line *geometry.Line) float64 {
	return geoDistancePoints(g.Center(), line.Rect().Center())
}

func (g *Point) DistancePoly(poly *geometry.Poly) float64 {
	return geoDistancePoints(g.Center(), poly.Rect().Center())
}

// IsSimple returns true if the Point can be converted to a SimplePoint
func (g *Point) IsSimple() bool {
	return g.extra == nil
}

// IsPoint returns true if the object is a {"type":"Point"}
func IsPoint(obj Object) (z float64, ok bool) {
	switch pt := obj.(type) {
	case *SimplePoint:
		return 0, true
	case *Point:
		return pt.Z(), true
	}
	return 0, false
}

func (g *Point) Members() string {
	if g.extra != nil {
		return g.extra.members
	}
	return ""
}
