package geojson

import (
	"github.com/tidwall/geojson/geos"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
)

// Feature ...
type Feature struct {
	base       Object
	extra      *extra
	ID         string
	Properties string
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
func (g *Feature) Rect() geos.Rect {
	if g.extra != nil && g.extra.bbox != nil {
		return *g.extra.bbox
	}
	return g.base.Rect()
}

// Center ...
func (g *Feature) Center() geos.Point {
	return g.Rect().Center()
}

// AppendJSON ...
func (g *Feature) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Feature","geometry":`...)
	dst = g.base.AppendJSON(dst)
	if g.extra != nil {
		dst = g.extra.appendJSONBBox(dst)
	}
	if g.ID != "" {
		dst = append(dst, `,"id":`...)
		dst = append(dst, g.ID...)
	}
	if g.Properties != "" {
		dst = append(dst, `,"properties":`...)
		dst = append(dst, g.Properties...)
	}
	dst = append(dst, '}')
	return dst

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

func (g *Feature) withinRect(rect geos.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return rect.ContainsRect(*g.extra.bbox)
	}
	return g.base.withinRect(rect)
}

func (g *Feature) withinPoint(point geos.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return point.ContainsRect(*g.extra.bbox)
	}
	return g.base.withinPoint(point)
}

func (g *Feature) withinLine(line *geos.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return line.ContainsRect(*g.extra.bbox)
	}
	return g.base.withinLine(line)
}

func (g *Feature) withinPoly(poly *geos.Poly) bool {
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

func (g *Feature) intersectsPoint(point geos.Point) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoint(point)
	}
	return g.base.intersectsPoint(point)
}

func (g *Feature) intersectsRect(rect geos.Rect) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsRect(rect)
	}
	return g.base.intersectsRect(rect)
}

func (g *Feature) intersectsLine(line *geos.Line) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsLine(line)
	}
	return g.base.intersectsLine(line)
}

func (g *Feature) intersectsPoly(poly *geos.Poly) bool {
	if g.extra != nil && g.extra.bbox != nil {
		return g.extra.bbox.IntersectsPoly(poly)
	}
	return g.base.intersectsPoly(poly)
}

// parseJSONFeature will return a valid GeoJSON object.
func parseJSONFeature(data string, opts *ParseOptions) (Object, error) {
	var g Feature
	rgeometry := gjson.Get(data, "geometry")
	if !rgeometry.Exists() {
		return nil, errGeometryMissing
	}
	var err error
	g.base, err = Parse(rgeometry.Raw, opts)
	if err != nil {
		return nil, err
	}
	if err := parseBBoxAndFillExtra(data, &g.extra, opts); err != nil {
		return nil, err
	}
	id := gjson.Get(data, "id").Raw
	if len(id) > 0 {
		g.ID = string(pretty.UglyInPlace([]byte(id)))
	}
	properties := gjson.Get(data, "properties").Raw
	if len(properties) > 0 {
		g.Properties = string(pretty.UglyInPlace([]byte(properties)))
	}
	return &g, nil
}
