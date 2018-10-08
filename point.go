package geojson

import (
	"github.com/tidwall/geojson/geos"
	"github.com/tidwall/gjson"
)

// Point ...
type Point struct {
	base  geos.Point
	extra *extra
}

// // ForEach ...
// func (g *Point) ForEach(iter func(geom geos.Geometry) bool) bool {
// 	if g.extra != nil && g.extra.bbox != nil {
// 		return iter(*g.extra.bbox)
// 	}
// 	return iter(g.base)
// }

// forEach ...
func (g *Point) forEach(iter func(geom Object) bool) bool {
	return iter(g)
}

// NewPoint ...
func NewPoint(x, y float64) *Point {
	return &Point{base: geos.Point{X: x, Y: y}}
}

// Empty ...
func (g *Point) Empty() bool {
	return g.base.Empty()
}

// Rect ...
func (g *Point) Rect() geos.Rect {
	if g.extra != nil && g.extra.bbox != nil {
		return *g.extra.bbox
	}
	return g.base.Rect()
}

// Center ...
func (g *Point) Center() geos.Point {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.Center()
	}
	return g.base
}

// AppendJSON ...
func (g *Point) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Point","coordinates":`...)
	dst = appendJSONPoint(dst, g.base, g.extra, 0)
	if g.extra != nil {
		dst = g.extra.appendJSONBBox(dst)
	}
	dst = append(dst, '}')
	return dst
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

func (g *Point) withinRect(rect geos.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return rect.ContainsRect(*g.extra.bbox)
	}
	return rect.ContainsPoint(g.base)
}

func (g *Point) withinPoint(point geos.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return point.ContainsRect(*g.extra.bbox)
	}
	return point.ContainsPoint(g.base)
}

func (g *Point) withinLine(line *geos.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return line.ContainsRect(*g.extra.bbox)
	}
	return line.ContainsPoint(g.base)
}

func (g *Point) withinPoly(poly *geos.Poly) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return poly.ContainsRect(*g.extra.bbox)
	}
	return poly.ContainsPoint(g.base)
}

func (g *Point) intersectsPoint(point geos.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoint(point)
	}
	return g.base.IntersectsPoint(point)
}

func (g *Point) intersectsRect(rect geos.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsRect(rect)
	}
	return g.base.IntersectsRect(rect)
}

func (g *Point) intersectsLine(line *geos.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsLine(line)
	}
	return g.base.IntersectsLine(line)
}

func (g *Point) intersectsPoly(poly *geos.Poly) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoly(poly)
	}
	return g.base.IntersectsPoly(poly)
}

func parseJSONPoint(data string) (Object, error) {
	var g Point
	var err error
	g.base, g.extra, err = parseJSONPointCoords(data, gjson.Result{})
	if err != nil {
		return nil, err
	}
	if err := parseBBoxAndFillExtra(data, &g.extra); err != nil {
		return nil, err
	}
	return &g, nil
}

func parseJSONPointCoords(data string, rcoords gjson.Result) (
	geos.Point, *extra, error,
) {
	var coords geos.Point
	var ex *extra
	if !rcoords.Exists() {
		rcoords = gjson.Get(data, "coordinates")
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
	coords = geos.Point{X: nums[0], Y: nums[1]}
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
