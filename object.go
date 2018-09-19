package geojson

import (
	"errors"
	"fmt"

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

type Object interface {
	Rect() Rect
	Center() Position
	AppendJSON(dst []byte) []byte
	// Contains(other Object) bool
	// Overlaps(other Object) bool
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
