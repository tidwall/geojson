package geojson

import (
	"github.com/tidwall/geojson/geos"
	"github.com/tidwall/gjson"
)

// LineString ...
type LineString struct {
	base  geos.Line
	extra *extra
}

// Empty ...
func (g *LineString) Empty() bool {
	if g.extra != nil && g.extra.bbox != nil {
		return false
	}
	return g.base.Empty()
}

// Rect ...
func (g *LineString) Rect() geos.Rect {
	if g.extra != nil && g.extra.bbox != nil {
		return *g.extra.bbox
	}
	return g.base.Rect()
}

// Center ...
func (g *LineString) Center() geos.Point {
	return g.Rect().Center()
}

// AppendJSON ...
func (g *LineString) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"LineString","coordinates":`...)
	dst, _ = appendJSONSeries(dst, &g.base, g.extra, 0)
	if g.extra != nil {
		dst = g.extra.appendJSONExtra(dst)
	}
	dst = append(dst, '}')
	return dst
}

// forEach ...
func (g *LineString) forEach(iter func(geom Object) bool) bool {
	return iter(g)
}

// Within ...
func (g *LineString) Within(obj Object) bool {
	return obj.Contains(g)
}

// Contains ...
func (g *LineString) Contains(obj Object) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return obj.withinRect(*g.extra.bbox)
	}
	return obj.withinLine(&g.base)
}

// Intersects ...
func (g *LineString) Intersects(obj Object) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return obj.intersectsRect(*g.extra.bbox)
	}
	return obj.intersectsLine(&g.base)
}

func (g *LineString) withinRect(rect geos.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return rect.ContainsRect(*g.extra.bbox)
	}
	return rect.ContainsLine(&g.base)
}

func (g *LineString) withinPoint(point geos.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return point.ContainsRect(*g.extra.bbox)
	}
	return point.ContainsLine(&g.base)
}

func (g *LineString) withinLine(line *geos.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return line.ContainsRect(*g.extra.bbox)
	}
	return line.ContainsLine(&g.base)
}

func (g *LineString) withinPoly(poly *geos.Poly) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return poly.ContainsRect(*g.extra.bbox)
	}
	return poly.ContainsLine(&g.base)
}

func (g *LineString) intersectsPoint(point geos.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoint(point)
	}
	return g.base.IntersectsPoint(point)
}

func (g *LineString) intersectsRect(rect geos.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsRect(rect)
	}
	return g.base.IntersectsRect(rect)
}

func (g *LineString) intersectsLine(line *geos.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsLine(line)
	}
	return g.base.IntersectsLine(line)
}

func (g *LineString) intersectsPoly(poly *geos.Poly) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoly(poly)
	}
	return g.base.IntersectsPoly(poly)
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
	line := geos.NewLine(points, opts.IndexGeometry)
	g.base = *line
	g.extra = ex
	if err := parseBBoxAndExtras(&g.extra, keys, opts); err != nil {
		return nil, err
	}
	return &g, nil
}

func parseJSONLineStringCoords(
	keys *parseKeys, rcoords gjson.Result, opts *ParseOptions,
) ([]geos.Point, *extra, error) {
	var err error
	var coords []geos.Point
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
		coords = append(coords, geos.Point{X: nums[0], Y: nums[1]})
		if ex == nil {
			if count > 2 {
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
