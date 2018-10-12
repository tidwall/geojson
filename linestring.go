package geojson

import (
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/gjson"
)

// LineString ...
type LineString struct {
	base  geometry.Line
	extra *extra
}

// NewLineString ...
func NewLineString(points []geometry.Point) *LineString {
	line := geometry.NewLine(points, geometry.DefaultIndex)
	return &LineString{base: *line}
}

// Empty ...
func (g *LineString) Empty() bool {
	return g.base.Empty()
}

// Rect ...
func (g *LineString) Rect() geometry.Rect {
	return g.base.Rect()
}

// Center ...
func (g *LineString) Center() geometry.Point {
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

// String ...
func (g *LineString) String() string {
	return string(g.AppendJSON(nil))
}

// JSON ...
func (g *LineString) JSON() string {
	return string(g.AppendJSON(nil))
}

// IsSpatial ...
func (g *LineString) IsSpatial() bool {
	return true
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
	return obj.withinLine(&g.base)
}

// Intersects ...
func (g *LineString) Intersects(obj Object) bool {
	return obj.intersectsLine(&g.base)
}

func (g *LineString) withinRect(rect geometry.Rect) bool {
	return rect.ContainsLine(&g.base)
}

func (g *LineString) withinPoint(point geometry.Point) bool {
	return point.ContainsLine(&g.base)
}

func (g *LineString) withinLine(line *geometry.Line) bool {
	return line.ContainsLine(&g.base)
}

func (g *LineString) withinPoly(poly *geometry.Poly) bool {
	return poly.ContainsLine(&g.base)
}

func (g *LineString) intersectsPoint(point geometry.Point) bool {
	return g.base.IntersectsPoint(point)
}

func (g *LineString) intersectsRect(rect geometry.Rect) bool {
	return g.base.IntersectsRect(rect)
}

func (g *LineString) intersectsLine(line *geometry.Line) bool {
	return g.base.IntersectsLine(line)
}

func (g *LineString) intersectsPoly(poly *geometry.Poly) bool {
	return g.base.IntersectsPoly(poly)
}

// NumPoints ...
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
	line := geometry.NewLine(points, opts.IndexGeometry)
	g.base = *line
	g.extra = ex
	if err := parseBBoxAndExtras(&g.extra, keys, opts); err != nil {
		return nil, err
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

// // Clipped ...
// func (g *LineString) Clipped(obj Object) Object {
// 	bbox := obj.Rect()
// 	var newPoints [][]geometry.Point
// 	var clipped geometry.Segment
// 	var rejected bool
// 	var line []geometry.Point
// 	nSegments := g.base.NumSegments()
// 	for i := 0; i < nSegments; i++ {
// 		clipped, rejected = ClipSegment(g.base.SegmentAt(i), bbox)
// 		if rejected {
// 			continue
// 		}
// 		if len(line) > 0 && line[len(line)-1] != clipped.A {
// 			newPoints = append(newPoints, line)
// 			line = []geometry.Point{clipped.A}
// 		} else if len(line) == 0 {
// 			line = append(line, clipped.A)
// 		}
// 		line = append(line, clipped.B)
// 	}
// 	if len(line) > 0 {
// 		newPoints = append(newPoints, line)
// 	}
// 	var children []Object
// 	for _, points := range newPoints {
// 		var lineString = new(LineString)
// 		line := geometry.NewLine(points, geometry.DefaultIndex)
// 		lineString.base = *line
// 		children = append(children, lineString)
// 	}
// 	if len(children) == 1 {
// 		return children[0]
// 	}
// 	multi := new(MultiLineString)
// 	multi.children = children
// 	multi.parseInitRectIndex(DefaultParseOptions)
// 	return multi
// }

// Distance ...
func (g *LineString) Distance(obj Object) float64 {
	return obj.distanceLine(&g.base)
}
func (g *LineString) distancePoint(point geometry.Point) float64 {
	return geoDistancePoints(g.Center(), point)
}
func (g *LineString) distanceRect(rect geometry.Rect) float64 {
	return geoDistancePoints(g.Center(), rect.Center())
}
func (g *LineString) distanceLine(line *geometry.Line) float64 {
	return geoDistancePoints(g.Center(), line.Rect().Center())
}
func (g *LineString) distancePoly(poly *geometry.Poly) float64 {
	return geoDistancePoints(g.Center(), poly.Rect().Center())
}
