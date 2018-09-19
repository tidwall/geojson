package geojson

import (
	"unsafe"
)

// Point is a point
type Point struct {
	Coordinates Position
	BBox        BBox
	extra       *extra
}

// Rect returns a rectangle that contains the entire object
func (g Point) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	return Rect{Min: g.Coordinates, Max: g.Coordinates}
}

// Center is the center-most point of the object
func (g Point) Center() Position {
	if g.BBox != nil {
		return g.BBox.Rect().Center()
	}
	return g.Coordinates
}

// AppendJSON appends a json representation to destination
func (g Point) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Point","coordinates":`...)
	dst = g.Coordinates.appendJSON(dst, g.extra, 0)
	if g.BBox != nil && g.BBox.Defined() {
		dst = append(dst, `,"bbox":`...)
		dst = g.BBox.AppendJSON(dst)
	}
	dst = append(dst, '}')
	return dst
}

// JSON returns a json representation of the object
func (g Point) JSON() string {
	return string(g.AppendJSON(nil))
}

// String returns a string representation of the object
func (g Point) String() string {
	return g.JSON()
}

// Stats of the object
func (g Point) Stats() Stats {
	return Stats{
		Weight: int(unsafe.Sizeof(g)) + bboxWeight(g.BBox) +
			g.extra.weight(),
		PositionCount: 1 + bboxPositionCount(g.BBox),
	}
}

// Contains another object
func (g Point) Contains(o Object) bool {
	return g.Rect().ContainsRect(o.Rect())
}

// Intersects another object
func (g Point) Intersects(o Object) bool {
	return g.Rect().IntersectsRect(o.Rect())
}

// IntersectsPolyLine test if object intersect a polyline
func (g Point) IntersectsPolyLine(line []Position) bool {
	if g.BBox != nil && g.BBox.Defined() {
		return g.BBox.Rect().IntersectsPolyLine(line)
	}
	return g.Coordinates.IntersectsPolyLine(line)
}

// parseGeoJSONPoint will return a valid GeoJSON object.
func parseGeoJSONPoint(data string) (Object, error) {
	var g Point
	var err error
	g.Coordinates, g.extra, err = parseCoords1(data)
	if err != nil {
		return nil, err
	}
	g.BBox, err = parseBBox(data)
	if err != nil {
		return nil, err
	}
	// if g.extra == nil && g.BBox == nil {
	// 	return SimplePoint{Coordinates: g.Coordinates}, nil
	// }
	return g, nil
}
