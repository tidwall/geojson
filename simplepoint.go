package geojson

import "github.com/tidwall/geojson/geometry"

type SimplePoint struct {
	geometry.Point
}

// NewSimplePoint returns a new SimplePoint object.
func NewSimplePoint(point geometry.Point) *SimplePoint {
	return &SimplePoint{Point: point}
}

func (g *SimplePoint) ForEach(iter func(geom Object) bool) bool {
	return iter(g)
}

func (g *SimplePoint) Empty() bool {
	return g.Point.Empty()
}

func (g *SimplePoint) Valid() bool {
	return g.Point.Valid()
}

func (g *SimplePoint) Rect() geometry.Rect {
	return g.Point.Rect()
}

func (g *SimplePoint) Spatial() Spatial {
	return g
}

func (g *SimplePoint) Center() geometry.Point {
	return g.Point
}

func (g *SimplePoint) Base() geometry.Point {
	return g.Point
}

func (g *SimplePoint) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Point","coordinates":`...)
	dst = appendJSONPoint(dst, g.Point, nil, 0)
	dst = append(dst, '}')
	return dst
}

func (g *SimplePoint) JSON() string {
	return string(g.AppendJSON(nil))
}

func (g *SimplePoint) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

func (g *SimplePoint) String() string {
	return string(g.AppendJSON(nil))
}

func (g *SimplePoint) Within(obj Object) bool {
	return obj.Contains(g)
}

func (g *SimplePoint) Contains(obj Object) bool {
	return obj.Spatial().WithinPoint(g.Point)
}

func (g *SimplePoint) Intersects(obj Object) bool {
	if obj, ok := obj.(*Circle); ok {
		return obj.Contains(g)
	}
	return obj.Spatial().IntersectsPoint(g.Point)
}

func (g *SimplePoint) WithinRect(rect geometry.Rect) bool {
	return rect.ContainsPoint(g.Point)
}

func (g *SimplePoint) WithinPoint(point geometry.Point) bool {
	return point.ContainsPoint(g.Point)
}

func (g *SimplePoint) WithinLine(line *geometry.Line) bool {
	return line.ContainsPoint(g.Point)
}

func (g *SimplePoint) WithinPoly(poly *geometry.Poly) bool {
	return poly.ContainsPoint(g.Point)
}

func (g *SimplePoint) IntersectsPoint(point geometry.Point) bool {
	return g.Point.IntersectsPoint(point)
}

func (g *SimplePoint) IntersectsRect(rect geometry.Rect) bool {
	return g.Point.IntersectsRect(rect)
}

func (g *SimplePoint) IntersectsLine(line *geometry.Line) bool {
	return g.Point.IntersectsLine(line)
}

func (g *SimplePoint) IntersectsPoly(poly *geometry.Poly) bool {
	return g.Point.IntersectsPoly(poly)
}

func (g *SimplePoint) NumPoints() int {
	return 1
}

func (g *SimplePoint) Distance(obj Object) float64 {
	return obj.Spatial().DistancePoint(g.Point)
}

func (g *SimplePoint) DistancePoint(point geometry.Point) float64 {
	return geoDistancePoints(g.Center(), point)
}

func (g *SimplePoint) DistanceRect(rect geometry.Rect) float64 {
	return geoDistancePoints(g.Center(), rect.Center())
}

func (g *SimplePoint) DistanceLine(line *geometry.Line) float64 {
	return geoDistancePoints(g.Center(), line.Rect().Center())
}

func (g *SimplePoint) DistancePoly(poly *geometry.Poly) float64 {
	return geoDistancePoints(g.Center(), poly.Rect().Center())
}

func (g *SimplePoint) Members() string {
	return ""
}
