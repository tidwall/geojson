package geojson

import (
	"unsafe"

	"github.com/tidwall/gjson"
)

// Feature is a GeoJSON Feature
type Feature struct {
	Geometry   Object
	BBox       BBox
	ID         gjson.Result
	Properties gjson.Result
}

// Rect returns a rectangle that contains the entire object
func (g Feature) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	return g.Geometry.Rect()
}

// Center is the center-most point of the object
func (g Feature) Center() Position {
	return g.Rect().Center()
}

// AppendJSON appends a json representation to destination
func (g Feature) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Feature","geometry":`...)
	dst = g.Geometry.AppendJSON(dst)
	if g.BBox != nil && g.BBox.Defined() {
		dst = append(dst, `,"bbox":`...)
		dst = g.BBox.AppendJSON(dst)
	}
	if g.ID.Exists() {
		dst = append(dst, `,"id":`...)
		dst = append(dst, g.ID.Raw...)
	}
	if g.Properties.Exists() {
		dst = append(dst, `,"properties":`...)
		dst = append(dst, g.Properties.Raw...)
	}
	dst = append(dst, '}')
	return dst
}

// JSON returns a json representation of the object
func (g Feature) JSON() string {
	return string(g.AppendJSON(nil))
}

// String returns a string representation of the object
func (g Feature) String() string {
	return g.JSON()
}

// Stats of the object
func (g Feature) Stats() Stats {
	gs := g.Geometry.Stats()
	return Stats{
		Weight:        int(unsafe.Sizeof(g)) + bboxWeight(g.BBox) + gs.Weight,
		PositionCount: bboxPositionCount(g.BBox) + gs.PositionCount,
	}
}

// Contains another object
func (g Feature) Contains(o Object) bool {
	if contains, certain := objectContainsRect(g.BBox, g, o); certain {
		return contains
	}
	return g.Geometry.Contains(o)
}

// Intersects another object
func (g Feature) Intersects(o Object) bool {
	if intersects, certain := objectIntersectsRect(g.BBox, g, o); certain {
		return intersects
	}
	return g.Geometry.Intersects(o)
}

// IntersectsPolyLine test if object intersect a polyline
func (g Feature) IntersectsPolyLine(line []Position) bool {
	if g.BBox != nil && g.BBox.Defined() {
		return g.BBox.Rect().IntersectsPolyLine(line)
	}
	return g.Geometry.IntersectsPolyLine(line)
}

// parseGeoJSONFeature will return a valid GeoJSON object.
func parseGeoJSONFeature(data string) (Object, error) {
	var g Feature
	rgeometry := gjson.Get(data, "geometry")
	if !rgeometry.Exists() {
		return nil, errGeometryMissing
	}
	var err error
	g.Geometry, err = Parse(rgeometry.Raw)
	if err != nil {
		return nil, err
	}
	g.BBox, err = parseBBox(data)
	if err != nil {
		return nil, err
	}
	g.ID = resultCopy(gjson.Get(data, "id"))
	g.Properties = resultCopy(gjson.Get(data, "properties"))
	return g, nil
}

func resultCopy(res gjson.Result) gjson.Result {
	if res.Exists() {
		if res.Type == gjson.String {
			res = gjson.Parse(string([]byte(res.Raw)))
		} else {
			res.Raw = string([]byte(res.Raw))
		}
	}
	return res
}
