package geojson

import (
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/gjson"
)

// Point ...
type Point struct {
	base  geometry.Point
	extra *extra
}

// forEach ...
func (g *Point) forEach(iter func(geom Object) bool) bool {
	return iter(g)
}

// NewPoint ...
func NewPoint(x, y float64) *Point {
	return &Point{base: geometry.Point{X: x, Y: y}}
}

// Empty ...
func (g *Point) Empty() bool {
	return g.base.Empty()
}

// Rect ...
func (g *Point) Rect() geometry.Rect {
	if g.extra != nil && g.extra.bbox != nil {
		return *g.extra.bbox
	}
	return g.base.Rect()
}

// Center ...
func (g *Point) Center() geometry.Point {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.Center()
	}
	return g.base
}

// AppendJSON ...
func (g *Point) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Point","coordinates":`...)
	dst = appendJSONPoint(dst, g.base, g.extra, 0)
	dst = g.extra.appendJSONExtra(dst)
	dst = append(dst, '}')
	return dst
}

// String ...
func (g *Point) String() string {
	return string(g.AppendJSON(nil))
}

// Within ...
func (g *Point) Within(obj Object) bool {
	return obj.Contains(g)
}

// Contains ...
func (g *Point) Contains(obj Object) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return obj.withinRect(*g.extra.bbox)
	}
	return obj.withinPoint(g.base)
}

// Intersects ...
func (g *Point) Intersects(obj Object) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return obj.intersectsRect(*g.extra.bbox)
	}
	return obj.intersectsPoint(g.base)
}

func (g *Point) withinRect(rect geometry.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return rect.ContainsRect(*g.extra.bbox)
	}
	return rect.ContainsPoint(g.base)
}

func (g *Point) withinPoint(point geometry.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return point.ContainsRect(*g.extra.bbox)
	}
	return point.ContainsPoint(g.base)
}

func (g *Point) withinLine(line *geometry.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return line.ContainsRect(*g.extra.bbox)
	}
	return line.ContainsPoint(g.base)
}

func (g *Point) withinPoly(poly *geometry.Poly) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return poly.ContainsRect(*g.extra.bbox)
	}
	return poly.ContainsPoint(g.base)
}

func (g *Point) intersectsPoint(point geometry.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoint(point)
	}
	return g.base.IntersectsPoint(point)
}

func (g *Point) intersectsRect(rect geometry.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsRect(rect)
	}
	return g.base.IntersectsRect(rect)
}

func (g *Point) intersectsLine(line *geometry.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsLine(line)
	}
	return g.base.IntersectsLine(line)
}

func (g *Point) intersectsPoly(poly *geometry.Poly) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoly(poly)
	}
	return g.base.IntersectsPoly(poly)
}

// NumPoints ...
func (g *Point) NumPoints() int {
	return 0
}

// Nearby ...
func (g *Point) Nearby(center geometry.Point, meters float64) bool {
	panic("not ready")
}

func parseJSONPoint(keys *parseKeys, opts *ParseOptions) (Object, error) {
	var g Point
	var err error
	g.base, g.extra, err = parseJSONPointCoords(keys, gjson.Result{}, opts)
	if err != nil {
		return nil, err
	}
	if err := parseBBoxAndExtras(&g.extra, keys, opts); err != nil {
		return nil, err
	}
	return &g, nil
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
