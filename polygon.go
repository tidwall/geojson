package geojson

import (
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/gjson"
)

type Polygon struct {
	base  geometry.Poly
	extra *extra
}

func NewPolygon(poly *geometry.Poly) *Polygon {
	g := &Polygon{}
	if poly != nil {
		g.base = *poly
	}
	return g
}

func (g *Polygon) Empty() bool {
	return g.base.Empty()
}

func (g *Polygon) Valid() bool {
	return g.base.Valid()
}

func (g *Polygon) Rect() geometry.Rect {
	return g.base.Rect()
}

func (g *Polygon) Center() geometry.Point {
	return g.Rect().Center()
}

func (g *Polygon) Base() *geometry.Poly {
	return &g.base
}

func (g *Polygon) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Polygon","coordinates":[`...)
	if !g.base.Empty() {
		var pidx int
		dst, pidx = appendJSONSeries(dst, g.base.Exterior, g.extra, pidx)
		for _, hole := range g.base.Holes {
			dst = append(dst, ',')
			dst, pidx = appendJSONSeries(dst, hole, g.extra, pidx)
		}
	}
	dst = append(dst, ']')
	if g.extra != nil {
		dst = g.extra.appendJSONExtra(dst, false)
	}
	dst = append(dst, '}')
	return dst
}

func (g *Polygon) JSON() string {
	return string(g.AppendJSON(nil))
}

func (g *Polygon) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

func (g *Polygon) String() string {
	return string(g.AppendJSON(nil))
}

func (g *Polygon) Spatial() Spatial {
	return g
}

func (g *Polygon) ForEach(iter func(geom Object) bool) bool {
	return iter(g)
}

func (g *Polygon) Within(obj Object) bool {
	return obj.Contains(g)
}

func (g *Polygon) Contains(obj Object) bool {
	return obj.Spatial().WithinPoly(&g.base)
}

func (g *Polygon) WithinRect(rect geometry.Rect) bool {
	return rect.ContainsPoly(&g.base)
}

func (g *Polygon) WithinPoint(point geometry.Point) bool {
	return point.ContainsPoly(&g.base)
}

func (g *Polygon) WithinLine(line *geometry.Line) bool {
	return line.ContainsPoly(&g.base)
}

func (g *Polygon) WithinPoly(poly *geometry.Poly) bool {
	return poly.ContainsPoly(&g.base)
}

func (g *Polygon) Intersects(obj Object) bool {
	return obj.Spatial().IntersectsPoly(&g.base)
}

func (g *Polygon) IntersectsPoint(point geometry.Point) bool {
	return g.base.IntersectsPoint(point)
}

func (g *Polygon) IntersectsRect(rect geometry.Rect) bool {
	return g.base.IntersectsRect(rect)
}

func (g *Polygon) IntersectsLine(line *geometry.Line) bool {
	return g.base.IntersectsLine(line)
}

func (g *Polygon) IntersectsPoly(poly *geometry.Poly) bool {
	return g.base.IntersectsPoly(poly)
}

func (g *Polygon) NumPoints() int {
	n := g.base.Exterior.NumPoints()
	for _, hole := range g.base.Holes {
		n += hole.NumPoints()
	}
	return n
}

func parseJSONPolygon(keys *parseKeys, opts *ParseOptions) (Object, error) {
	var o Object
	coords, extra, err := parseJSONPolygonCoords(keys, gjson.Result{}, opts)
	if err != nil {
		return nil, err
	}
	if len(coords) == 0 {
		return nil, errCoordinatesInvalid // must be a linear ring
	}
	for _, p := range coords {
		if len(p) < 4 || p[0] != p[len(p)-1] {
			return nil, errCoordinatesInvalid // must be a linear ring
		}
	}
	exterior := coords[0]
	var holes [][]geometry.Point
	if len(coords) > 1 {
		holes = coords[1:]
	}
	gopts := toGeometryOpts(opts)
	if err := parseBBoxAndExtras(&extra, keys, opts); err != nil {
		return nil, err
	}
	if extra == nil && opts.AllowRects &&
		len(holes) == 0 && len(exterior) == 5 &&
		exterior[0].X < exterior[1].X &&
		exterior[0].Y == exterior[1].Y &&
		exterior[1].X == exterior[2].X &&
		exterior[1].Y < exterior[2].Y &&
		exterior[2].X > exterior[3].X &&
		exterior[2].Y == exterior[3].Y &&
		exterior[3].X == exterior[4].X &&
		exterior[3].Y > exterior[4].Y {
		// simple rectangle
		o = NewRect(geometry.Rect{
			Min: exterior[0],
			Max: exterior[2],
		})
	} else {
		g := Polygon{}
		poly := geometry.NewPoly(exterior, holes, &gopts)
		g.base = *poly
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

func parseJSONPolygonCoords(
	keys *parseKeys, rcoords gjson.Result, opts *ParseOptions,
) (
	[][]geometry.Point, *extra, error,
) {
	var err error
	var coords [][]geometry.Point
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
		coords = append(coords, []geometry.Point{})
		ii := len(coords) - 1
		value.ForEach(func(key, value gjson.Result) bool {
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
			coords[ii] = append(coords[ii], geometry.Point{X: nums[0], Y: nums[1]})
			if ex == nil {
				if count > 2 {
					if len(coords) > 1 || len(coords[ii]) > 1 {
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
		return err == nil
	})
	if err != nil {
		return nil, nil, err
	}
	return coords, ex, err
}

func (g *Polygon) Distance(obj Object) float64 {
	return obj.Spatial().DistancePoly(&g.base)
}

func (g *Polygon) DistancePoint(point geometry.Point) float64 {
	return geoDistancePoints(g.Center(), point)
}

func (g *Polygon) DistanceRect(rect geometry.Rect) float64 {
	return geoDistancePoints(g.Center(), rect.Center())
}

func (g *Polygon) DistanceLine(line *geometry.Line) float64 {
	return geoDistancePoints(g.Center(), line.Rect().Center())
}

func (g *Polygon) DistancePoly(poly *geometry.Poly) float64 {
	return geoDistancePoints(g.Center(), poly.Rect().Center())
}

func (g *Polygon) HasExtra() bool {
	return g.extra != nil
}

func (g *Polygon) Members() string {
	if g.extra != nil {
		return g.extra.members
	}
	return ""
}
