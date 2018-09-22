package geojson

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/tidwall/geojson/poly"
	"github.com/tidwall/gjson"
)

var (
	fmtErrTypeIsUnknown   = "type '%s' is unknown"
	errDataInvalid        = errors.New("invalid data")
	errTypeInvalid        = errors.New("invalid type")
	errTypeMissing        = errors.New("missing type")
	errCoordinatesInvalid = errors.New("invalid coordinates")
	errCoordinatesMissing = errors.New("missing coordinates")
	errGeometryInvalid    = errors.New("invalid geometry")
	errGeometryMissing    = errors.New("missing geometry")
	errFeaturesMissing    = errors.New("missing features")
	errFeaturesInvalid    = errors.New("invalid features")
	errGeometriesMissing  = errors.New("missing geometries")
	errGeometriesInvalid  = errors.New("invalid geometries")
	errBBoxInvalid        = errors.New("invalid bbox")
	errMustBeALinearRing  = errors.New("invalid polygon")
)

// Object is a geo object
type Object interface {
	// BBoxDefined return true if there is a defined GeoJSON "bbox" member
	BBoxDefined() bool
	// Rect returns the outer minimum bounding rectangle
	Rect() Rect
	// Center returns the center position of the object
	Center() Position
	// AppendJSON appends the GeoJSON reprensentation to dst
	AppendJSON(dst []byte) []byte
	// ForEach iterates over child objects. Used for GeoJSON types:
	//   MultiPoint, MultiLineString, MultiPolygon, Feature, FeatureCollection,
	//   and GeometryCollection
	ForEach(func(child Object) bool)
	// Contains returns true if object contains other object
	Contains(other Object) bool
	// Contains returns true if object intersects with other object
	Intersects(other Object) bool
}

var _ = []Object{
	Position{}, Rect{},
	Point{}, LineString{}, Polygon{},
	MultiPoint{}, MultiLineString{}, MultiPolygon{},
	GeometryCollection{},
	Feature{}, FeatureCollection{},
}

func Load(data string) (Object, error) {
	// look at the first byte
	for i := 0; ; i++ {
		if len(data) == 0 {
			return nil, errDataInvalid
		}
		switch data[0] {
		default:
			// well-known text is not supported yet
			return nil, errDataInvalid
		case 0, 1:
			if i > 0 {
				// 0x00 or 0x01 must be the first bytes
				return nil, errDataInvalid
			}
			// well-known binary is not supported yet
			return nil, errDataInvalid
		case ' ', '\t', '\n', '\r':
			// strip whitespace
			data = data[1:]
			continue
		case '{':
			return loadJSON(data)
		}
	}
}

func loadJSON(data string) (Object, error) {
	if !gjson.Valid(data) {
		return nil, errDataInvalid
	}
	rtype := gjson.Get(data, "type")
	if !rtype.Exists() {
		return nil, errTypeMissing
	}
	if rtype.Type != gjson.String {
		return nil, errTypeInvalid
	}
	switch rtype.String() {
	default:
		return nil, fmt.Errorf(fmtErrTypeIsUnknown, rtype.String())
	case "Point":
		return loadJSONPoint(data)
	case "LineString":
		return loadJSONLineString(data)
	case "Polygon":
		return loadJSONPolygon(data)
	case "MultiPoint":
		return loadJSONMultiPoint(data)
	case "MultiLineString":
		return loadJSONMultiLineString(data)
	case "MultiPolygon":
		return loadJSONMultiPolygon(data)
	case "GeometryCollection":
		return loadJSONGeometryCollection(data)
	case "Feature":
		return loadJSONFeature(data)
	case "FeatureCollection":
		return loadJSONFeatureCollection(data)
	}
}

func polyPoint(posn Position) poly.Point {
	return *(*poly.Point)(unsafe.Pointer(&posn))
}
func polyRect(rect Rect) poly.Rect {
	return *(*poly.Rect)(unsafe.Pointer(&rect))
}
func polyLine(line []Position) poly.Line {
	return *(*poly.Line)(unsafe.Pointer(&line))
}
func polyPolygon(polygon [][]Position) poly.Polygon {
	var newPoly poly.Polygon
	if len(polygon) > 0 {
		newPoly.Exterior = *(*poly.Ring)(unsafe.Pointer(&polygon[0]))
		if len(polygon) > 1 {
			newPoly.Holes = (*(*[]poly.Ring)(unsafe.Pointer(&polygon)))[1:]
		}
	}
	return newPoly
}

func collectionObjectContains(g, other Object) bool {
	if g.BBoxDefined() {
		return g.Rect().Contains(other)
	}
	var contains bool
	g.ForEach(func(child Object) bool {
		if child.Contains(other) {
			contains = true
			return false
		}
		return true
	})
	return contains
}

func collectionObjectIntersects(g, other Object) bool {
	if g.BBoxDefined() {
		return g.Rect().Intersects(other)
	}
	var intersects bool
	g.ForEach(func(child Object) bool {
		if child.Intersects(other) {
			intersects = true
			return false
		}
		return true
	})
	return intersects
}
