package geojson

import (
	"unsafe"

	"github.com/tidwall/gjson"
)

// FeatureCollection is a GeoJSON FeatureCollection
type FeatureCollection struct {
	Features []Object
	BBox     BBox
}

// Rect returns a rectangle that contains the entire object
func (g FeatureCollection) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	return calculateRectObjs(g.Features)
}

// Center is the center-most point of the object
func (g FeatureCollection) Center() Position {
	return g.Rect().Center()
}

// AppendJSON appends a json representation to destination
func (g FeatureCollection) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"FeatureCollection","features":[`...)
	for i := 0; i < len(g.Features); i++ {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = g.Features[i].AppendJSON(dst)
	}
	dst = append(dst, ']')
	if g.BBox != nil && g.BBox.Defined() {
		dst = append(dst, `,"bbox":`...)
		dst = g.BBox.AppendJSON(dst)
	}
	dst = append(dst, '}')
	return dst
}

// JSON returns a json representation of the object
func (g FeatureCollection) JSON() string {
	return string(g.AppendJSON(nil))
}

// String returns a string representation of the object
func (g FeatureCollection) String() string {
	return g.JSON()
}

// Stats of the object
func (g FeatureCollection) Stats() Stats {
	stats := Stats{
		Weight:        int(unsafe.Sizeof(g)) + bboxWeight(g.BBox),
		PositionCount: bboxPositionCount(g.BBox),
	}
	for _, g := range g.Features {
		s := g.Stats()
		stats.Weight += s.Weight
		stats.PositionCount += s.PositionCount
	}
	return stats
}

// Contains another object
func (g FeatureCollection) Contains(o Object) bool {
	if contains, certain := objectContainsRect(g.BBox, g, o); certain {
		return contains
	}
	for _, g := range g.Features {
		if g.Contains(o) {
			return true
		}
	}
	return false
}

// Intersects another object
func (g FeatureCollection) Intersects(o Object) bool {
	if intersects, certain := objectIntersectsRect(g.BBox, g, o); certain {
		return intersects
	}
	for _, g := range g.Features {
		if g.Intersects(o) {
			return true
		}
	}
	return false
}

// IntersectsPolyLine test if object intersect a polyline
func (g FeatureCollection) IntersectsPolyLine(line []Position) bool {
	if g.BBox != nil && g.BBox.Defined() {
		return g.BBox.Rect().IntersectsPolyLine(line)
	}
	for _, g := range g.Features {
		if g.IntersectsPolyLine(line) {
			return true
		}
	}
	return false
}

// parseGeoJSONFeatureCollection will return a valid GeoJSON object.
func parseGeoJSONFeatureCollection(data string) (Object, error) {
	var g FeatureCollection
	rfeatures := gjson.Get(data, "features")
	if !rfeatures.Exists() {
		return nil, errFeaturesMissing
	}
	if !rfeatures.IsArray() {
		return nil, errFeaturesInvalid
	}
	var err error
	rfeatures.ForEach(func(key, value gjson.Result) bool {
		var f Object
		f, err = Parse(value.Raw)
		if err != nil {
			return false
		}
		g.Features = append(g.Features, f)
		return true
	})
	if err != nil {
		return nil, err
	}
	g.BBox, err = parseBBox(data)
	if err != nil {
		return nil, err
	}
	if g.BBox == nil {
		g.BBox = g.Rect().BBox()
	}
	return g, nil
}
