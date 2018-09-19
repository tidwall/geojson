package geojson

import (
	"unsafe"
)

// LineString GeoJSON Object
type LineString struct {
	Coordinates []Position
	BBox        BBox
	extra       *extra
}

// Rect returns a rectangle that contains the entire object
func (g LineString) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	return calculateRect2(g.Coordinates)
}

// Center is the center-most point of the object
func (g LineString) Center() Position {
	return g.Rect().Center()
}

// AppendJSON appends a json representation to destination
func (g LineString) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"LineString","coordinates":[`...)
	for i := 0; i < len(g.Coordinates); i++ {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = g.Coordinates[i].appendJSON(dst, g.extra, i)
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
func (g LineString) JSON() string {
	return string(g.AppendJSON(nil))
}

// String returns a string representation of the object
func (g LineString) String() string {
	return g.JSON()
}

// Stats of the object
func (g LineString) Stats() Stats {
	return Stats{
		Weight: int(unsafe.Sizeof(g)) + len(g.Coordinates)*16 +
			bboxWeight(g.BBox) + g.extra.weight(),
		PositionCount: len(g.Coordinates) + bboxPositionCount(g.BBox),
	}
}

// Contains another object
func (g LineString) Contains(o Object) bool {
	// if contains, certain := objectContainsRect(g.BBox, g, o); certain {
	// 	return contains
	// }
	// oRect := o.Rect()
	// if oRect.Min == oRect.Max {
	// 	return SimplePoint{oRect.Min}.IntersectsPolyLine(g.Coordinates)
	// }
	// it's not possible for an object with area to be withing a polyline
	return false
}

// Intersects another object
func (g LineString) Intersects(o Object) bool {
	if intersects, certain := objectIntersectsRect(g.BBox, g, o); certain {
		return intersects
	}
	return o.IntersectsPolyLine(g.Coordinates)
}

// IntersectsPolyLine test if object intersect a polyline
func (g LineString) IntersectsPolyLine(line []Position) bool {
	if g.BBox != nil && g.BBox.Defined() {
		return g.BBox.Rect().IntersectsPolyLine(line)
	}
	return polyLine(g.Coordinates).LineStringIntersectsLineString(
		polyLine(line),
	)
}

// parseGeoJSONLineString will return a valid GeoJSON object.
func parseGeoJSONLineString(data string) (Object, error) {
	var g LineString
	var err error
	g.Coordinates, g.extra, err = parseCoords2(data)
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
