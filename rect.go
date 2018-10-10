package geojson

import (
	"github.com/tidwall/geojson/geometry"
)

// Rect ...
type Rect struct {
	base geometry.Rect
}

// NewRect ...
func NewRect(minX, minY, maxX, maxY float64) *Rect {
	return &Rect{base: geometry.Rect{
		Min: geometry.Point{X: minX, Y: minY},
		Max: geometry.Point{X: maxX, Y: maxY},
	}}
}

// forEach ...
func (g *Rect) forEach(iter func(geom Object) bool) bool {
	return iter(g)
}

// Empty ...
func (g *Rect) Empty() bool {
	return g.base.Empty()
}

// Rect ...
func (g *Rect) Rect() geometry.Rect {
	return g.base
}

// Center ...
func (g *Rect) Center() geometry.Point {
	return g.base.Center()
}

// AppendJSON ...
func (g *Rect) AppendJSON(dst []byte) []byte {
	panic("not ready")
}

// String ...
func (g *Rect) String() string {
	return string(g.AppendJSON(nil))
}

// Within ...
func (g *Rect) Within(obj Object) bool {
	return obj.Contains(g)
}

// Contains ...
func (g *Rect) Contains(obj Object) bool {
	return obj.withinRect(g.base)
}

// Intersects ...
func (g *Rect) Intersects(obj Object) bool {
	return obj.intersectsRect(g.base)
}

func (g *Rect) withinRect(rect geometry.Rect) bool {
	return rect.ContainsRect(g.base)
}

func (g *Rect) withinPoint(point geometry.Point) bool {
	return point.ContainsRect(g.base)
}

func (g *Rect) withinLine(line *geometry.Line) bool {
	return line.ContainsRect(g.base)
}

func (g *Rect) withinPoly(poly *geometry.Poly) bool {
	return poly.ContainsRect(g.base)
}

func (g *Rect) intersectsPoint(point geometry.Point) bool {
	return g.base.IntersectsPoint(point)
}

func (g *Rect) intersectsRect(rect geometry.Rect) bool {
	return g.base.IntersectsRect(rect)
}

func (g *Rect) intersectsLine(line *geometry.Line) bool {
	return g.base.IntersectsLine(line)
}

func (g *Rect) intersectsPoly(poly *geometry.Poly) bool {
	return g.base.IntersectsPoly(poly)
}

// NumPoints ...
func (g *Rect) NumPoints() int {
	return 2
}

// Nearby ...
func (g *Rect) Nearby(center geometry.Point, meters float64) bool {
	panic("not ready")
}

// RectFromCenter returns a geospatial rect
func RectFromCenter(center geometry.Point, meters float64) geometry.Rect {
	panic("not ready")
	// var outer geometry.Rect
	// outer.Min.Y, outer.Min.X, outer.Max.Y, outer.Max.X =
	// 	BoundsFromCenter(center.Y, center.X, meters)
	// if outer.Min.X == outer.Max.X {
	// 	switch outer.Min.X {
	// 	case -180:
	// 		outer.Max.X = 180
	// 	case 180:
	// 		outer.Min.X = -180
	// 	}
	// }
	// return outer
}

// const (
// 	earthRadius = 6371e3
// 	radians     = math.Pi / 180
// 	degrees     = 180 / math.Pi
// )

// // RectFromCenter calculates the bounding box surrounding .
// func RectFromCenter(lat, lon, meters float64) (
// 	latMin, lonMin, latMax, lonMax float64,
// ) {

// 	// see http://janmatuschek.de/LatitudeLongitudeBoundingCoordinates#Latitude
// 	lat = lat * radians
// 	lon = lon * radians

// 	r := meters / earthRadius // angular radius

// 	latMin = lat - r
// 	latMax = lat + r

// 	latT := math.Asin(math.Sin(lat) / math.Cos(r))
// 	lonΔ := math.Acos((math.Cos(r) - math.Sin(latT)*math.Sin(lat)) /
// 		(math.Cos(latT) * math.Cos(lat)))

// 	lonMin = lon - lonΔ
// 	lonMax = lon + lonΔ

// 	// Adjust for north poll
// 	if latMax > math.Pi/2 {
// 		lonMin = -math.Pi
// 		latMax = math.Pi / 2
// 		lonMax = math.Pi
// 	}

// 	// Adjust for south poll
// 	if latMin < -math.Pi/2 {
// 		latMin = -math.Pi / 2
// 		lonMin = -math.Pi
// 		lonMax = math.Pi
// 	}

// 	// Adjust for wraparound. Remove this if the commented-out condition below
// 	// this block is added.
// 	if lonMin < -math.Pi || lonMax > math.Pi {
// 		lonMin = -math.Pi
// 		lonMax = math.Pi
// 	}

// 	// // Consider splitting area into two bboxes, using the below checks, and
// 	// // erasing above block for performance. See
// 	// http://janmatuschek.de/LatitudeLongitudeBoundingCoordinates#PolesAnd180thMeridian

// 	// // Adjust for wraparound if minimum longitude is less than -180 degrees.
// 	// if lonMin < -math.Pi {
// 	// // box 1:
// 	// 	latMin = latMin
// 	// 	latMax = latMax
// 	// 	lonMin += 2*math.Pi
// 	// 	lonMax = math.Pi
// 	// // box 2:
// 	// 	latMin = latMin
// 	// 	latMax = latMax
// 	// 	lonMin = -math.Pi
// 	// 	lonMax = lonMax
// 	// }

// 	// // Adjust for wraparound if maximum longitude is greater than 180 degrees.
// 	// if lonMax > math.Pi {
// 	// // box 1:
// 	// 	latMin = latMin
// 	// 	latMax = latMax
// 	// 	lonMin = lonMin
// 	// 	lonMax = -math.Pi
// 	// // box 2:
// 	// 	latMin = latMin
// 	// 	latMax = latMax
// 	// 	lonMin = -math.Pi
// 	// 	lonMax -= 2*math.Pi
// 	// }

// 	// normalise to -180..+180°
// 	lonMin = math.Mod(lonMin+3*math.Pi, 2*math.Pi) - math.Pi
// 	lonMax = math.Mod(lonMax+3*math.Pi, 2*math.Pi) - math.Pi

// 	return latMin * degrees, lonMin * degrees, latMax * degrees, lonMax * degrees
// }
