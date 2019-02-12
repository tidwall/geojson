package geojson

import (
	"strings"

	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/tidwall/sjson"
)

// Feature ...
type Feature struct {
	base  Object
	extra *extra
}

// NewFeature returns a new GeoJSON Feature.
// The members must be a valid json object such as
// `{"id":"391","properties":{}}`, or it must be an empty string. It should not
// contain a "feature" member.
func NewFeature(geometry Object, members string) *Feature {
	g := new(Feature)
	g.base = geometry
	members = strings.TrimSpace(members)
	if members != "" && members != "{}" {
		if gjson.Valid(members) && gjson.Parse(members).IsObject() {
			if gjson.Get(members, "feature").Exists() {
				members, _ = sjson.Delete(members, "feature")
			}
			g.extra = new(extra)
			g.extra.members = string(pretty.UglyInPlace([]byte(members)))
		}
	}
	return g
}

// ForEach ...
func (g *Feature) ForEach(iter func(geom Object) bool) bool {
	return g.base.ForEach(iter)
}

// Empty ...
func (g *Feature) Empty() bool {
	return g.base.Empty()
}

// Valid ...
func (g *Feature) Valid() bool {
	return g.base.Valid()
}

// Rect ...
func (g *Feature) Rect() geometry.Rect {
	return g.base.Rect()
}

// Center ...
func (g *Feature) Center() geometry.Point {
	return g.Rect().Center()
}

// Base ...
func (g *Feature) Base() Object {
	return g.base
}

// Members ...
func (g *Feature) Members() string {
	if g.extra != nil {
		return g.extra.members
	}
	return ""
}

// AppendJSON ...
func (g *Feature) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Feature","geometry":`...)
	dst = g.base.AppendJSON(dst)
	dst = g.extra.appendJSONExtra(dst, true)
	dst = append(dst, '}')
	return dst

}

// String ...
func (g *Feature) String() string {
	return string(g.AppendJSON(nil))
}

// JSON ...
func (g *Feature) JSON() string {
	return string(g.AppendJSON(nil))
}

// MarshalJSON ...
func (g *Feature) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

// Spatial ...
func (g *Feature) Spatial() Spatial {
	return g
}

// Within ...
func (g *Feature) Within(obj Object) bool {
	return obj.Contains(g)
}

// Contains ...
func (g *Feature) Contains(obj Object) bool {
	return obj.Within(g.base)
}

// WithinRect ...
func (g *Feature) WithinRect(rect geometry.Rect) bool {
	return g.base.Spatial().WithinRect(rect)
}

// WithinPoint ...
func (g *Feature) WithinPoint(point geometry.Point) bool {
	return g.base.Spatial().WithinPoint(point)
}

// WithinLine ...
func (g *Feature) WithinLine(line *geometry.Line) bool {
	return g.base.Spatial().WithinLine(line)
}

// WithinPoly ...
func (g *Feature) WithinPoly(poly *geometry.Poly) bool {
	return g.base.Spatial().WithinPoly(poly)
}

// Intersects ...
func (g *Feature) Intersects(obj Object) bool {
	return obj.Intersects(g.base)
}

// IntersectsPoint ...
func (g *Feature) IntersectsPoint(point geometry.Point) bool {
	return g.base.Spatial().IntersectsPoint(point)
}

// IntersectsRect ...
func (g *Feature) IntersectsRect(rect geometry.Rect) bool {
	return g.base.Spatial().IntersectsRect(rect)
}

// IntersectsLine ...
func (g *Feature) IntersectsLine(line *geometry.Line) bool {
	return g.base.Spatial().IntersectsLine(line)
}

// IntersectsPoly ...
func (g *Feature) IntersectsPoly(poly *geometry.Poly) bool {
	return g.base.Spatial().IntersectsPoly(poly)
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
	if point, ok := g.base.(*Point); ok {
		if g.extra != nil {
			members := g.extra.members
			if gjson.Get(members, "properties.type").String() == "Circle" {
				// Circle
				radius := gjson.Get(members, "properties.radius").Float()
				units := gjson.Get(members, "properties.radius_units").String()
				switch units {
				case "", "m":
				case "km":
					radius *= 1000
				default:
					return nil, errCircleRadiusUnitsInvalid
				}
				return NewCircle(point.base, radius, 64), nil
			}
		}
	}
	return &g, nil
}

// Distance ...
func (g *Feature) Distance(obj Object) float64 {
	return g.base.Distance(obj)
}

// DistancePoint ...
func (g *Feature) DistancePoint(point geometry.Point) float64 {
	return g.base.Spatial().DistancePoint(point)
}

// DistanceRect ...
func (g *Feature) DistanceRect(rect geometry.Rect) float64 {
	return g.base.Spatial().DistanceRect(rect)
}

// DistanceLine ...
func (g *Feature) DistanceLine(line *geometry.Line) float64 {
	return g.base.Spatial().DistanceLine(line)
}

// DistancePoly ...
func (g *Feature) DistancePoly(poly *geometry.Poly) float64 {
	return g.base.Spatial().DistancePoly(poly)
}
