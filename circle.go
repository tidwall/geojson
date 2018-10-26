package geojson

import (
	"strconv"

	"github.com/tidwall/geojson/geo"
	"github.com/tidwall/geojson/geometry"
)

// Circle ...
type Circle struct {
	Object
	center geometry.Point
	meters float64
	steps  int
	km     bool
	extra  *extra
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

func (g *Circle) Meters() float64 {
	return g.meters
}

func (g *Circle) Center() geometry.Point {
	return g.center
}

func (g *Circle) Contains(obj Object) bool {
	if p, ok := obj.(*Point); ok {
		return p.Distance(g) < g.Meters()
	}
	if c, ok := obj.(*Circle); ok {
		return c.Distance(g) < (c.Meters() + g.Meters())
	}
	if ls, ok := obj.(*LineString); ok {
		for i := 0; i < ls.base.NumPoints() ; i++ {
			if geoDistancePoints(ls.base.PointAt(i), g.Center()) > g.Meters() {
				return false
			}
		}
		return true
	}

	// Not sure if polygon already does this?
	if mp, ok := obj.(*MultiPoint); ok {
		for _, p := range mp.Children() {
			if !g.Contains(p) {
				return false
			}
		}
		return true
	}
	if mls, ok := obj.(*MultiLineString); ok {
		for _, p := range mls.Children() {
			if !g.Contains(p) {
				return false
			}
		}
		return true
	}

	// No simple cases, so using polygon approximation.
	return g.Object.Contains(g)
}

func (g *Circle) intersectsSegment(seg geometry.Segment) bool {
	start, end := seg.A, seg.B
	center := g.Center()
	meters := g.Meters()

	// These are faster checks.  If they succeed there's no need do complicate things.
	if geoDistancePoints(center, start) <= meters {
		return true
	}
	if geoDistancePoints(center, end) <= meters {
		return true
	}

	// Distance between start and end
	l := geo.DistanceTo(start.Y, start.X, end.Y, end.X)

	// Unit direction vector
	dx := (end.X - start.X) / l
	dy := (end.Y - start.Y) / l

	// Point of the line closest to the center
	t := dx * (center.X - start.X) + dy * (center.Y - start.Y)
	px := t * dx + start.X
	py := t * dy + start.Y
	if px < start.X || px > end.X || py < start.Y || py > end.Y {
		// closest point is outside the segment
		return false
	}

	// Distance from the closest point to the center
	return geo.DistanceTo(center.Y, center.X, py, px) <= meters
}

func (g *Circle) Intersects(obj Object) bool {
	if p, ok := obj.(*Point); ok {
		return p.Distance(g) <= g.Meters()
	}
	if c, ok := obj.(*Circle); ok {
		return c.Distance(g) <= (c.Meters() + g.Meters())
	}
	if ls, ok := obj.(*LineString); ok {
		for i := 0; i < ls.base.NumSegments() ; i++ {
			if g.intersectsSegment(ls.base.SegmentAt(i)) {
				return true
			}
		}
		return false
	}

	// Not sure if polygon already does this?
	if mp, ok := obj.(*MultiPoint); ok {
		for _, p := range mp.Children() {
			if !g.Intersects(p) {
				return false
			}
		}
		return true
	}
	if mls, ok := obj.(*MultiLineString); ok {
		for _, p := range mls.Children() {
			if !g.Intersects(p) {
				return false
			}
		}
		return true
	}

	// No simple cases, so using polygon approximation.
	return g.Object.Intersects(g)
}
