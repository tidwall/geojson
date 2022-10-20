package geojson

import (
	"strings"

	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"github.com/tidwall/sjson"
)

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

func (g *Feature) ForEach(iter func(geom Object) bool) bool {
	return iter(g)
}

func (g *Feature) Empty() bool {
	return g.base.Empty()
}

func (g *Feature) Valid() bool {
	return g.base.Valid()
}

func (g *Feature) Rect() geometry.Rect {
	return g.base.Rect()
}

func (g *Feature) Center() geometry.Point {
	return g.Rect().Center()
}

func (g *Feature) Base() Object {
	return g.base
}

func (g *Feature) Members() string {
	if g.extra != nil {
		return g.extra.members
	}
	return ""
}

func (g *Feature) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Feature","geometry":`...)
	dst = g.base.AppendJSON(dst)
	dst = g.extra.appendJSONExtra(dst, true)
	dst = append(dst, '}')
	return dst

}

func (g *Feature) String() string {
	return string(g.AppendJSON(nil))
}

func (g *Feature) JSON() string {
	return string(g.AppendJSON(nil))
}

func (g *Feature) MarshalJSON() ([]byte, error) {
	return g.AppendJSON(nil), nil
}

func (g *Feature) Spatial() Spatial {
	return g
}

func (g *Feature) Within(obj Object) bool {
	return obj.Contains(g)
}

func (g *Feature) Contains(obj Object) bool {
	return g.base.Contains(obj)
}

func (g *Feature) WithinRect(rect geometry.Rect) bool {
	return g.base.Spatial().WithinRect(rect)
}

func (g *Feature) WithinPoint(point geometry.Point) bool {
	return g.base.Spatial().WithinPoint(point)
}

func (g *Feature) WithinLine(line *geometry.Line) bool {
	return g.base.Spatial().WithinLine(line)
}

func (g *Feature) WithinPoly(poly *geometry.Poly) bool {
	return g.base.Spatial().WithinPoly(poly)
}

func (g *Feature) Intersects(obj Object) bool {
	return g.base.Intersects(obj)
}

func (g *Feature) IntersectsPoint(point geometry.Point) bool {
	return g.base.Spatial().IntersectsPoint(point)
}

func (g *Feature) IntersectsRect(rect geometry.Rect) bool {
	return g.base.Spatial().IntersectsRect(rect)
}

func (g *Feature) IntersectsLine(line *geometry.Line) bool {
	return g.base.Spatial().IntersectsLine(line)
}

func (g *Feature) IntersectsPoly(poly *geometry.Poly) bool {
	return g.base.Spatial().IntersectsPoly(poly)
}

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
			if !opts.DisableCircleType &&
				gjson.Get(members, "properties.type").String() == "Circle" {
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

func (g *Feature) Distance(obj Object) float64 {
	return g.base.Distance(obj)
}

func (g *Feature) DistancePoint(point geometry.Point) float64 {
	return g.base.Spatial().DistancePoint(point)
}

func (g *Feature) DistanceRect(rect geometry.Rect) float64 {
	return g.base.Spatial().DistanceRect(rect)
}

func (g *Feature) DistanceLine(line *geometry.Line) float64 {
	return g.base.Spatial().DistanceLine(line)
}

func (g *Feature) DistancePoly(poly *geometry.Poly) float64 {
	return g.base.Spatial().DistancePoly(poly)
}
