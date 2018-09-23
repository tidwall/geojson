package geojson

import (
	"github.com/tidwall/gjson"
)

// Point GeoJSON type
type Point struct {
	Coordinates Position
	BBox        BBox
	Extra       *Extra
}

// BBoxDefined return true if there is a defined GeoJSON "bbox" member
func (g Point) BBoxDefined() bool {
	return g.BBox != nil && g.BBox.Defined()
}

// Rect returns the outer minimum bounding rectangle
func (g Point) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	return Rect{Min: g.Coordinates, Max: g.Coordinates}
}

// Center returns the center position of the object
func (g Point) Center() Position {
	return g.Rect().Center()
}

// AppendJSON appends the GeoJSON reprensentation to dst
func (g Point) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Point","coordinates":`...)
	dst = appendJSONPosition(dst, g.Coordinates, g.Extra, 0)
	if g.BBox != nil && g.BBox.Defined() {
		dst = append(dst, `,"bbox":`...)
		dst = g.BBox.AppendJSON(dst)
	}
	dst = append(dst, '}')
	return dst
}

// ForEachChild iterates over child objects.
func (g Point) ForEachChild(func(child Object) bool) {}

// Within is the inverse of contains
func (g Point) Within(other Object) bool {
	return other.Contains(g)
}

// Contains returns true if object contains other object
func (g Point) Contains(other Object) bool {
	return objectContains(g, other)
}

// Intersects returns true if object intersects with other object
func (g Point) Intersects(other Object) bool {
	return objectIntersects(g, other)
}

func (g Point) primativeContains(other Object) bool {
	ppoint := polyPoint(g.Coordinates)
	switch other := other.(type) {
	case Position:
		return polyPoint(other).InsidePoint(ppoint)
	case Rect:
		return polyRect(other).InsidePoint(ppoint)
	case Point:
		return polyPoint(other.Coordinates).InsidePoint(ppoint)
	case LineString:
		return polyLine(other.Coordinates).InsidePoint(ppoint)
	case Polygon:
		return polyPolygon(other.Coordinates).InsidePoint(ppoint)
	}
	return false
}
func (g Point) primativeIntersects(other Object) bool {
	ppoint := polyPoint(g.Coordinates)
	switch other := other.(type) {
	case Position:
		return ppoint.IntersectsPoint(polyPoint(other))
	case Rect:
		return ppoint.IntersectsRect(polyRect(other))
	case Point:
		return ppoint.IntersectsPoint(polyPoint(other.Coordinates))
	case LineString:
		return ppoint.IntersectsLine(polyLine(other.Coordinates))
	case Polygon:
		return ppoint.IntersectsPolygon(polyPolygon(other.Coordinates))
	}
	return false
}

func loadJSONPoint(data string) (Object, error) {
	var g Point
	var err error
	g.Coordinates, g.Extra, err = loadJSONPointCoords(data, gjson.Result{})
	if err != nil {
		return nil, err
	}
	g.BBox, err = loadBBox(data)
	if err != nil {
		return nil, err
	}
	if g.BBox == nil && g.Extra == nil {
		return g.Coordinates, nil
	}
	return g, nil
}

func loadJSONPointCoords(data string, rcoords gjson.Result) (
	Position, *Extra, error,
) {
	var coords Position
	var ex *Extra
	if !rcoords.Exists() {
		rcoords = gjson.Get(data, "coordinates")
		if !rcoords.Exists() {
			return coords, nil, errCoordinatesMissing
		}
		if !rcoords.IsArray() {
			return coords, nil, errCoordinatesInvalid
		}
	}
	var err error
	var count int
	var nums [4]float64
	rcoords.ForEach(func(key, value gjson.Result) bool {
		if count == 4 {
			return false
		}
		if value.Type != gjson.Number {
			err = errCoordinatesInvalid
			return false
		}
		nums[count] = value.Float()
		count++
		return true
	})
	if err != nil {
		return coords, nil, err
	}
	if count < 2 {
		return coords, nil, errCoordinatesInvalid
	}
	coords = Position{X: nums[0], Y: nums[1]}
	if count > 2 {
		ex = new(Extra)
		if count > 3 {
			ex.Dims = DimsZM
		} else {
			ex.Dims = DimsZ
		}
		ex.Positions = make([]float64, count-2)
		for i := 2; i < count; i++ {
			ex.Positions[i-2] = nums[i]
		}
	}
	return coords, ex, nil
}
