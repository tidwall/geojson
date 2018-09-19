package geojson

import (
	"unsafe"

	"github.com/tidwall/gjson"
)

// GeometryCollection is a GeoJSON GeometryCollection
type GeometryCollection struct {
	Geometries []Object
	BBox       BBox
}

// Rect returns a rectangle that contains the entire object
func (g GeometryCollection) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	return calculateRectObjs(g.Geometries)
}

// Center is the center-most point of the object
func (g GeometryCollection) Center() Position {
	return g.Rect().Center()
}

// AppendJSON appends a json representation to destination
func (g GeometryCollection) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"GeometryCollection","geometries":[`...)
	for i := 0; i < len(g.Geometries); i++ {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = g.Geometries[i].AppendJSON(dst)
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
func (g GeometryCollection) JSON() string {
	return string(g.AppendJSON(nil))
}

// String returns a string representation of the object
func (g GeometryCollection) String() string {
	return g.JSON()
}

// Stats of the object
func (g GeometryCollection) Stats() Stats {
	stats := Stats{
		Weight:        int(unsafe.Sizeof(g)) + bboxWeight(g.BBox),
		PositionCount: bboxPositionCount(g.BBox),
	}
	for _, g := range g.Geometries {
		s := g.Stats()
		stats.Weight += s.Weight
		stats.PositionCount += s.PositionCount
	}
	return stats
}

// Contains another object
func (g GeometryCollection) Contains(o Object) bool {
	if contains, certain := objectContainsRect(g.BBox, g, o); certain {
		return contains
	}
	for _, g := range g.Geometries {
		if g.Contains(o) {
			return true
		}
	}
	return false
}

// Intersects another object
func (g GeometryCollection) Intersects(o Object) bool {
	if intersects, certain := objectIntersectsRect(g.BBox, g, o); certain {
		return intersects
	}
	for _, g := range g.Geometries {
		if g.Intersects(o) {
			return true
		}
	}
	return false
}

// IntersectsPolyLine test if object intersect a polyline
func (g GeometryCollection) IntersectsPolyLine(line []Position) bool {
	if g.BBox != nil && g.BBox.Defined() {
		return g.BBox.Rect().IntersectsPolyLine(line)
	}
	for _, g := range g.Geometries {
		if g.IntersectsPolyLine(line) {
			return true
		}
	}
	return false
}

// parseGeoJSONGeometryCollection will return a valid GeoJSON object.
func parseGeoJSONGeometryCollection(data string) (Object, error) {
	var g GeometryCollection
	rgeometries := gjson.Get(data, "geometries")
	if !rgeometries.Exists() {
		return nil, errGeometriesMissing
	}
	if !rgeometries.IsArray() {
		return nil, errGeometriesInvalid
	}
	var err error
	rgeometries.ForEach(func(key, value gjson.Result) bool {
		var f Object
		f, err = Parse(value.Raw)
		if err != nil {
			return false
		}
		g.Geometries = append(g.Geometries, f)
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
