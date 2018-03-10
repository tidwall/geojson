package geojson

import (
	"github.com/tidwall/tile38/geojson/poly"
)

// withinObjectShared returns true if g is within o
func withinObjectShared(g, o Object) bool {
	bbp := o.bboxPtr()
	if bbp != nil {
		if !g.WithinBBox(*bbp) {
			return false
		}
		if o.IsBBoxDefined() {
			return true
		}
	}
	switch o := o.(type) {
	default:
		return false
	case SimplePoint:
		return g.WithinBBox(o.CalculatedBBox())
	case Point:
		return g.WithinBBox(o.CalculatedBBox())
	case MultiPoint:
		for i := range o.Coordinates {
			if g.Within(o.getPoint(i)) {
				return true
			}
		}
		return false
	case LineString:
		if len(o.Coordinates) == 0 {
			return false
		}
		switch g := g.(type) {
		default:
			return false
		case SimplePoint:
			return poly.Point(Position{X: g.X, Y: g.Y, Z: 0}).IntersectsLineString(polyPositions(o.Coordinates))
		case Point:
			return poly.Point(g.Coordinates).IntersectsLineString(polyPositions(o.Coordinates))
		case MultiPoint:
			if len(o.Coordinates) == 0 {
				return false
			}
			for _, p := range o.Coordinates {
				if !poly.Point(p).IntersectsLineString(polyPositions(o.Coordinates)) {
					return false
				}
			}
			return true
		}
	case MultiLineString:
		for i := range o.Coordinates {
			if g.Within(o.getLineString(i)) {
				return true
			}
		}
		return false
	case MultiPolygon:
		for i := range o.Coordinates {
			if g.Within(o.getPolygon(i)) {
				return true
			}
		}
		return false
	case Feature:
		return g.Within(o.Geometry)
	case FeatureCollection:
		for _, o := range o.Features {
			if g.Within(o) {
				return true
			}
		}
		return false
	case GeometryCollection:
		for _, o := range o.Geometries {
			if g.Within(o) {
				return true
			}
		}
		return false
	case Polygon:
		if len(o.Coordinates) == 0 {
			return false
		}
		exterior, holes := polyExteriorHoles(o.Coordinates)
		switch g := g.(type) {
		default:
			return false
		case SimplePoint:
			return poly.Point(Position{X: g.X, Y: g.Y, Z: 0}).Inside(exterior, holes)
		case Point:
			return poly.Point(g.Coordinates).Inside(exterior, holes)
		case MultiPoint:
			if len(g.Coordinates) == 0 {
				return false
			}
			for i := range g.Coordinates {
				if !g.getPoint(i).Within(o) {
					return false
				}
			}
			return true
		case LineString:
			return polyPositions(g.Coordinates).Inside(exterior, holes)
		case MultiLineString:
			if len(g.Coordinates) == 0 {
				return false
			}
			for i := range g.Coordinates {
				if !g.getLineString(i).Within(o) {
					return false
				}
			}
			return true
		case Polygon:
			if len(g.Coordinates) == 0 {
				return false
			}
			return polyPositions(g.Coordinates[0]).Inside(exterior, holes)
		case MultiPolygon:
			if len(g.Coordinates) == 0 {
				return false
			}
			for i := range g.Coordinates {
				if !g.getPolygon(i).Within(o) {
					return false
				}
			}
			return true
		case GeometryCollection:
			if len(g.Geometries) == 0 {
				return false
			}
			for _, g := range g.Geometries {
				if !g.Within(o) {
					return false
				}
			}
			return true
		case Feature:
			return g.Geometry.Within(o)
		case FeatureCollection:
			if len(g.Features) == 0 {
				return false
			}
			for _, g := range g.Features {
				if !g.Within(o) {
					return false
				}
			}
			return true
		}
	}
}

// intersectsObjectShared detects if g intersects with o
func intersectsObjectShared(g, o Object) bool {
	bbp := o.bboxPtr()
	if bbp != nil {
		if !g.IntersectsBBox(*bbp) {
			return false
		}
		if o.IsBBoxDefined() {
			return true
		}
	}
	switch o := o.(type) {
	default:
		return false
	case SimplePoint:
		return g.IntersectsBBox(o.CalculatedBBox())
	case Point:
		return g.IntersectsBBox(o.CalculatedBBox())
	case MultiPoint:
		for i := range o.Coordinates {
			if o.getPoint(i).Intersects(g) {
				return true
			}
		}
		return false
	case LineString:
		if g, ok := g.(LineString); ok {
			a := polyPositions(g.Coordinates)
			b := polyPositions(o.Coordinates)
			return a.LineStringIntersectsLineString(b)
		}
		return o.Intersects(g)
	case MultiLineString:
		for i := range o.Coordinates {
			if g.Intersects(o.getLineString(i)) {
				return true
			}
		}
		return false
	case MultiPolygon:
		for i := range o.Coordinates {
			if g.Intersects(o.getPolygon(i)) {
				return true
			}
		}
		return false
	case Feature:
		return g.Intersects(o.Geometry)
	case FeatureCollection:
		for _, f := range o.Features {
			if g.Intersects(f) {
				return true
			}
		}
		return false
	case GeometryCollection:
		for _, f := range o.Geometries {
			if g.Intersects(f) {
				return true
			}
		}
		return false
	case Polygon:
		if len(o.Coordinates) == 0 {
			return false
		}
		exterior, holes := polyExteriorHoles(o.Coordinates)
		switch g := g.(type) {
		default:
			return false
		case SimplePoint:
			return poly.Point(Position{X: g.X, Y: g.Y, Z: 0}).Intersects(exterior, holes)
		case Point:
			return poly.Point(g.Coordinates).Intersects(exterior, holes)
		case MultiPoint:
			for i := range g.Coordinates {
				if g.getPoint(i).Intersects(o) {
					return true
				}
			}
			return false
		case LineString:
			return polyPositions(g.Coordinates).LineStringIntersects(exterior, holes)
		case MultiLineString:
			for i := range g.Coordinates {
				if g.getLineString(i).Intersects(o) {
					return true
				}
			}
			return false
		case Polygon:
			if len(g.Coordinates) == 0 {
				return false
			}
			return polyPositions(g.Coordinates[0]).Intersects(exterior, holes)
		case MultiPolygon:
			for i := range g.Coordinates {
				if g.getPolygon(i).Intersects(o) {
					return true
				}
			}
			return false
		case GeometryCollection:
			for _, g := range g.Geometries {
				if g.Intersects(o) {
					return true
				}
			}
			return false
		case Feature:
			return g.Geometry.Intersects(o)
		case FeatureCollection:
			for _, g := range g.Features {
				if g.Intersects(o) {
					return true
				}
			}
			return false
		}
	}
}

// The object's calculated bounding box must intersect the radius of the circle to pass.
func nearbyObjectShared(g Object, x, y float64, meters float64) bool {
	if !g.hasPositions() {
		return false
	}
	center := Position{X: x, Y: y, Z: 0}
	bbox := g.CalculatedBBox()
	if bbox.Min.X == bbox.Max.X && bbox.Min.Y == bbox.Max.Y {
		// just a point, return is point is inside of the circle
		return center.DistanceTo(bbox.Min) <= meters
	}
	circlePoly := CirclePolygon(x, y, meters, 12)
	return g.Intersects(circlePoly)
}
