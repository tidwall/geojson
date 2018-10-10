package geojson

import "github.com/tidwall/geojson/geometry"

// Feature ...
type Feature struct {
	base  Object
	extra *extra
}

// forEach ...
func (g *Feature) forEach(iter func(geom Object) bool) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return iter(g)
	}
	return g.base.forEach(iter)
}

// Empty ...
func (g *Feature) Empty() bool {
	if g.extra != nil && g.extra.bbox != nil {
		return false
	}
	return g.base.Empty()
}

// Rect ...
func (g *Feature) Rect() geometry.Rect {
	if g.extra != nil && g.extra.bbox != nil {
		return *g.extra.bbox
	}
	return g.base.Rect()
}

// Center ...
func (g *Feature) Center() geometry.Point {
	return g.Rect().Center()
}

// AppendJSON ...
func (g *Feature) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Feature","geometry":`...)
	dst = g.base.AppendJSON(dst)
	dst = g.extra.appendJSONExtra(dst)
	dst = append(dst, '}')
	return dst

}

// String ...
func (g *Feature) String() string {
	return string(g.AppendJSON(nil))
}

// Within ...
func (g *Feature) Within(obj Object) bool {
	return obj.Contains(g)
}

// Contains ...
func (g *Feature) Contains(obj Object) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return obj.withinRect(*g.extra.bbox)
	}
	return obj.Within(g.base)
}

func (g *Feature) withinRect(rect geometry.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return rect.ContainsRect(*g.extra.bbox)
	}
	return g.base.withinRect(rect)
}

func (g *Feature) withinPoint(point geometry.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return point.ContainsRect(*g.extra.bbox)
	}
	return g.base.withinPoint(point)
}

func (g *Feature) withinLine(line *geometry.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return line.ContainsRect(*g.extra.bbox)
	}
	return g.base.withinLine(line)
}

func (g *Feature) withinPoly(poly *geometry.Poly) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return poly.ContainsRect(*g.extra.bbox)
	}
	return g.base.withinPoly(poly)
}

// Intersects ...
func (g *Feature) Intersects(obj Object) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return obj.intersectsRect(*g.extra.bbox)
	}
	return obj.Intersects(g.base)
}

func (g *Feature) intersectsPoint(point geometry.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoint(point)
	}
	return g.base.intersectsPoint(point)
}

func (g *Feature) intersectsRect(rect geometry.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsRect(rect)
	}
	return g.base.intersectsRect(rect)
}

func (g *Feature) intersectsLine(line *geometry.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsLine(line)
	}
	return g.base.intersectsLine(line)
}

func (g *Feature) intersectsPoly(poly *geometry.Poly) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoly(poly)
	}
	return g.base.intersectsPoly(poly)
}

// NumPoints ...
func (g *Feature) NumPoints() int {
	return g.base.NumPoints()
}

// parseJSONFeature will return a valid GeoJSON object.
func parseJSONFeature(keys *parseKeys, opts *ParseOptions) (Object, error) {
	var g Feature
	if !keys.rGeometry.Exists() {
		return nil, errGeometryMissing
	}
	var err error
	g.base, err = Parse(keys.rGeometry.Raw, opts)
	if err != nil {
		return nil, err
	}
	if err := parseBBoxAndExtras(&g.extra, keys, opts); err != nil {
		return nil, err
	}
	return &g, nil
}
