package geojson

import "github.com/tidwall/geojson/geometry"

type Spatial interface {
	WithinRect(rect geometry.Rect) bool
	WithinPoint(point geometry.Point) bool
	WithinLine(line *geometry.Line) bool
	WithinPoly(poly *geometry.Poly) bool
	IntersectsRect(rect geometry.Rect) bool
	IntersectsPoint(point geometry.Point) bool
	IntersectsLine(line *geometry.Line) bool
	IntersectsPoly(poly *geometry.Poly) bool
	DistanceRect(rect geometry.Rect) float64
	DistancePoint(point geometry.Point) float64
	DistanceLine(line *geometry.Line) float64
	DistancePoly(poly *geometry.Poly) float64
}

var _ = []Spatial{
	&Point{}, &LineString{}, &Polygon{}, &Feature{},
	&MultiPoint{}, &MultiLineString{}, &MultiPolygon{},
	&GeometryCollection{}, &FeatureCollection{}, &Rect{},
	EmptySpatial{},
}

type EmptySpatial struct{}

func (s EmptySpatial) WithinRect(rect geometry.Rect) bool {
	return false
}

func (s EmptySpatial) WithinPoint(point geometry.Point) bool {
	return false
}

func (s EmptySpatial) WithinLine(line *geometry.Line) bool {
	return false
}

func (s EmptySpatial) WithinPoly(poly *geometry.Poly) bool {
	return false
}

func (s EmptySpatial) IntersectsRect(rect geometry.Rect) bool {
	return false
}

func (s EmptySpatial) IntersectsPoint(point geometry.Point) bool {
	return false
}

func (s EmptySpatial) IntersectsLine(line *geometry.Line) bool {
	return false
}

func (s EmptySpatial) IntersectsPoly(poly *geometry.Poly) bool {
	return false
}

func (s EmptySpatial) DistanceRect(rect geometry.Rect) float64 {
	return 0
}

func (s EmptySpatial) DistancePoint(point geometry.Point) float64 {
	return 0
}

func (s EmptySpatial) DistanceLine(line *geometry.Line) float64 {
	return 0
}

func (s EmptySpatial) DistancePoly(poly *geometry.Poly) float64 {
	return 0
}
