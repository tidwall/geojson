package geojson

import (
	"strconv"

	"github.com/tidwall/geojson/geo"
	"github.com/tidwall/geojson/geometry"
)

// Circle ...
type Circle struct {
	Object
	center    geometry.Point
	meters    float64
	haversine float64
	steps     int
	km        bool
	extra     *extra
}

// NewCircle returns an circle object
func NewCircle(center geometry.Point, meters float64, steps int) *Circle {
	if steps < 3 {
		steps = 3
	}
	g := new(Circle)
	g.center = center
	g.meters = meters
	g.steps = steps
	if meters <= 0 {
		g.Object = NewPoint(center)
	} else {
		meters = geo.NormalizeDistance(meters)
		var points []geometry.Point
		step := 360.0 / float64(steps)
		i := 0
		for deg := 360.0; deg > 0; deg -= step {
			lat, lon := geo.DestinationPoint(center.Y, center.X, meters, deg)
			points = append(points, geometry.Point{X: lon, Y: lat})
			i++
		}
		// TODO: account for the pole and antimerdian. In most cases only a
		// polygon is needed, but when the circle bounds passes the 90/180
		// lines, we need to create a multipolygon
		points = append(points, points[0])
		g.Object = NewPolygon(
			geometry.NewPoly(points, nil, geometry.DefaultIndexOptions),
		)
		g.haversine = geo.DistanceToHaversine(meters)
	}
	return g
}

// AppendJSON ...
func (g *Circle) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Feature","geometry":`...)
	dst = append(dst, `{"type":"Point","coordinates":[`...)
	dst = strconv.AppendFloat(dst, g.center.X, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, g.center.Y, 'f', -1, 64)
	dst = append(dst, `]},"properties":{"type":"Circle","radius":`...)
	dst = strconv.AppendFloat(dst, g.meters, 'f', -1, 64)
	dst = append(dst, `,"radius_units":"m"}}`...)
	return dst
}

// JSON ...
func (g *Circle) JSON() string {
	return string(g.AppendJSON(nil))
}

// String ...
func (g *Circle) String() string {
	return string(g.AppendJSON(nil))
}

// Meters returns the circle's radius
func (g *Circle) Meters() float64 {
	return g.meters
}

// Center returns the circle's center point
func (g *Circle) Center() geometry.Point {
	return g.center
}

// Within returns true if circle is contained inside object
func (g *Circle) Within(obj Object) bool {
	return obj.Contains(g)
}

func (g *Circle) contains(p geometry.Point, allowOnEdge bool) bool {
	h := geo.Haversine(p.Y, p.X, g.center.Y, g.center.X)
	if allowOnEdge {
		return h <= g.haversine
	}
	return h < g.haversine
}

// Contains returns true if the circle contains other object
func (g *Circle) Contains(obj Object) bool {
	switch other := obj.(type) {
	case *Point:
		return g.contains(other.Center(), false)
	case *Circle:
		return other.Distance(g) < (other.meters + g.meters)
	case *LineString:
		for i := 0; i < other.base.NumPoints(); i++ {
			if geoDistancePoints(other.base.PointAt(i), g.center) > g.meters {
				return false
			}
		}
		return true
	case Collection:
		for _, p := range other.Children() {
			if !g.Contains(p) {
				return false
			}
		}
		return true
	default:
		// No simple cases, so using polygon approximation.
		return g.Object.Contains(other)
	}
}

func (g *Circle) intersectsSegment(seg geometry.Segment) bool {
	start, end := seg.A, seg.B

	// These are faster checks.  If they succeed there's no need do complicate things.
	if g.contains(start, true) {
		return true
	}
	if g.contains(end, true) {
		return true
	}

	// Distance between start and end
	l := geo.DistanceTo(start.Y, start.X, end.Y, end.X)

	// Unit direction vector
	dx := (end.X - start.X) / l
	dy := (end.Y - start.Y) / l

	// Point of the line closest to the center
	t := dx*(g.center.X-start.X) + dy*(g.center.Y-start.Y)
	px := t*dx + start.X
	py := t*dy + start.Y
	if px < start.X || px > end.X || py < start.Y || py > end.Y {
		// closest point is outside the segment
		return false
	}

	// Distance from the closest point to the center
	return g.contains(geometry.Point{X: px, Y: py}, true)
}

// Intersects returns true the circle intersects other object
func (g *Circle) Intersects(obj Object) bool {
	switch other := obj.(type) {
	case *Point:
		return g.contains(other.Center(), true)
	case *Circle:
		return other.Distance(g) <= (other.meters + g.meters)
	case *LineString:
		for i := 0; i < other.base.NumSegments(); i++ {
			if g.intersectsSegment(other.base.SegmentAt(i)) {
				return true
			}
		}
		return false
	case Collection:
		for _, p := range other.Children() {
			if g.Intersects(p) {
				return true
			}
		}
		return false
	default:
		// No simple cases, so using polygon approximation.
		return g.Object.Intersects(obj)
	}
}
