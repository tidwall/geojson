package geojson

import "github.com/tidwall/geojson/geometry"

// Spatial ...
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

type emptySpatial int

func (s emptySpatial) WithinRect(rect geometry.Rect) bool         { return false }
func (s emptySpatial) WithinPoint(point geometry.Point) bool      { return false }
func (s emptySpatial) WithinLine(line *geometry.Line) bool        { return false }
func (s emptySpatial) WithinPoly(poly *geometry.Poly) bool        { return false }
func (s emptySpatial) IntersectsRect(rect geometry.Rect) bool     { return false }
func (s emptySpatial) IntersectsPoint(point geometry.Point) bool  { return false }
func (s emptySpatial) IntersectsLine(line *geometry.Line) bool    { return false }
func (s emptySpatial) IntersectsPoly(poly *geometry.Poly) bool    { return false }
func (s emptySpatial) DistanceRect(rect geometry.Rect) float64    { return 0 }
func (s emptySpatial) DistancePoint(point geometry.Point) float64 { return 0 }
func (s emptySpatial) DistanceLine(line *geometry.Line) float64   { return 0 }
func (s emptySpatial) DistancePoly(poly *geometry.Poly) float64   { return 0 }
