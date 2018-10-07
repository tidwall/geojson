package geojson

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/tidwall/geojson/geos"
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

// Object ...
type Object interface {
	Empty() bool
	Rect() geos.Rect
	Center() geos.Point
	AppendJSON(dst []byte) []byte
	Contains(other Object) bool
	Within(other Object) bool
	Intersects(other Object) bool

	withinRect(rect geos.Rect) bool
	withinPoint(point geos.Point) bool
	withinLine(line *geos.Line) bool
	withinPoly(poly *geos.Poly) bool
	intersectsRect(rect geos.Rect) bool
	intersectsPoint(point geos.Point) bool
	intersectsLine(line *geos.Line) bool
	intersectsPoly(poly *geos.Poly) bool
}

var _ = []Object{
	&Point{}, &LineString{}, &Polygon{}, &Feature{},
	&MultiPoint{},
}

type extra struct {
	bbox      *geos.Rect
	bboxExtra []float64
	dims      int
	values    []float64
}

// Parse a GeoJSON object
func Parse(data string) (Object, error) {
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
			return parseJSON(data)
		}
	}
}

func parseJSON(data string) (Object, error) {
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
		return parseJSONPoint(data)
		// case "LineString":
		// 	return loadJSONLineString(data)
		// case "Polygon":
		// 	return loadJSONPolygon(data)
		// case "MultiPoint":
		// 	return loadJSONMultiPoint(data)
		// case "MultiLineString":
		// 	return loadJSONMultiLineString(data)
		// case "MultiPolygon":
		// 	return loadJSONMultiPolygon(data)
		// case "GeometryCollection":
		// 	return loadJSONGeometryCollection(data)
		// case "Feature":
		// 	return loadJSONFeature(data)
		// case "FeatureCollection":
		// 	return loadJSONFeatureCollection(data)
	}
}

func parseBBox(data string) (bbox *geos.Rect, bboxExtra []float64, err error) {
	rbbox := gjson.Get(data, "bbox")
	if !rbbox.Exists() {
		return nil, nil, nil
	}
	if !rbbox.IsArray() {
		return nil, nil, errBBoxInvalid
	}
	var count int
	var nums [8]float64
	rbbox.ForEach(func(key, value gjson.Result) bool {
		if count == 8 {
			return false
		}
		if value.Type != gjson.Number {
			err = errBBoxInvalid
			return false
		}
		nums[count] = value.Float()
		count++
		return true
	})
	if err != nil {
		return nil, nil, err
	}
	if count < 4 || count%2 == 1 {
		return nil, nil, errBBoxInvalid
	}
	var rect geos.Rect
	rect.Min.X = nums[0]
	rect.Min.Y = nums[1]
	rect.Max.X = nums[count/2]
	rect.Max.Y = nums[count/2+1]
	if count == 4 {
		return &rect, nil, nil
	}
	if count == 6 {
		bboxExtra = []float64{
			nums[2],
			nums[count/2+2],
		}
		return &rect, bboxExtra, nil
	}
	bboxExtra = []float64{
		nums[2],
		nums[3],
		nums[count/2+2],
		nums[count/2+3],
	}
	return &rect, bboxExtra, nil
}

func appendJSONPoint(dst []byte, point geos.Point, ex *extra, idx int) []byte {
	dst = append(dst, '[')
	dst = strconv.AppendFloat(dst, point.X, 'f', -1, 64)
	dst = append(dst, ',')
	dst = strconv.AppendFloat(dst, point.Y, 'f', -1, 64)
	if ex != nil {
		dims := int(ex.dims)
		for i := 0; i < dims; i++ {
			dst = append(dst, ',')
			dst = strconv.AppendFloat(
				dst, ex.values[idx*dims+i], 'f', -1, 64,
			)
		}
	}
	dst = append(dst, ']')
	return dst
}

func (ex *extra) appendJSONBBox(dst []byte) []byte {
	if ex.bbox != nil {
		dst = append(dst, `,"bbox":[`...)
		dst = strconv.AppendFloat(dst, ex.bbox.Min.X, 'f', -1, 64)
		dst = append(dst, ',')
		dst = strconv.AppendFloat(dst, ex.bbox.Min.Y, 'f', -1, 64)
		if len(ex.bboxExtra) == 2 {
			dst = append(dst, ',')
			dst = strconv.AppendFloat(dst, ex.bboxExtra[0], 'f', -1, 64)
		} else if len(ex.bboxExtra) == 4 {
			dst = append(dst, ',')
			dst = strconv.AppendFloat(dst, ex.bboxExtra[0], 'f', -1, 64)
			dst = append(dst, ',')
			dst = strconv.AppendFloat(dst, ex.bboxExtra[1], 'f', -1, 64)
		}
		dst = append(dst, ',')
		dst = strconv.AppendFloat(dst, ex.bbox.Max.X, 'f', -1, 64)
		dst = append(dst, ',')
		dst = strconv.AppendFloat(dst, ex.bbox.Max.Y, 'f', -1, 64)
		if len(ex.bboxExtra) == 2 {
			dst = append(dst, ',')
			dst = strconv.AppendFloat(dst, ex.bboxExtra[1], 'f', -1, 64)
		} else if len(ex.bboxExtra) == 4 {
			dst = append(dst, ',')
			dst = strconv.AppendFloat(dst, ex.bboxExtra[3], 'f', -1, 64)
			dst = append(dst, ',')
			dst = strconv.AppendFloat(dst, ex.bboxExtra[4], 'f', -1, 64)
		}
		dst = append(dst, ']')
	}
	return dst
}
