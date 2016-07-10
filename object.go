package geojson

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/tidwall/tile38/geojson/poly"
)

const (
	point              = 0
	multiPoint         = 1
	lineString         = 2
	multiLineString    = 3
	polygon            = 4
	multiPolygon       = 5
	geometryCollection = 6
	feature            = 7
	featureCollection  = 8
)

var (
	errNotEnoughData = errors.New("not enough data")
	errTooMuchData   = errors.New("too much data")
	errInvalidData   = errors.New("invalid data")
)

var ( // json errors
	fmtErrTypeIsUnknown              = "The type '%s' is unknown"
	errInvalidTypeMember             = errors.New("Type member is invalid. Expecting a string")
	errInvalidCoordinates            = errors.New("Coordinates member is invalid. Expecting an array")
	errCoordinatesRequired           = errors.New("Coordinates member is required.")
	errInvalidGeometries             = errors.New("Geometries member is invalid. Expecting an array")
	errGeometriesRequired            = errors.New("Geometries member is required.")
	errInvalidGeometryMember         = errors.New("Geometry member is invalid. Expecting an object")
	errGeometryMemberRequired        = errors.New("Geometry member is required.")
	errInvalidFeaturesMember         = errors.New("Features member is invalid. Expecting an array")
	errFeaturesMemberRequired        = errors.New("Features member is required")
	errInvalidFeature                = errors.New("Invalid feature in collection")
	errInvalidPropertiesMember       = errors.New("Properties member in invalid. Expecting an array")
	errInvalidCoordinatesValue       = errors.New("Coordinates member has an invalid value")
	errLineStringInvalidCoordinates  = errors.New("Coordinates must be an array of two or more positions")
	errInvalidNumberOfPositionValues = errors.New("Position must have two or more numbers")
	errInvalidPositionValue          = errors.New("Position has an invalid value")
	errCoordinatesMustBeArray        = errors.New("Coordinates member must be an array of positions")
	errMustBeALinearRing             = errors.New("Polygon must have at least 4 positions and the first and last position must be the same")
	errBBoxInvalidType               = errors.New("BBox member is an invalid. Expecting an array")
	errBBoxInvalidNumberOfValues     = errors.New("BBox member requires exactly 4 or 6 values")
	errBBoxInvalidValue              = errors.New("BBox has an invalid value")
	errInvalidGeometry               = errors.New("Invalid geometry in collection")
)

const nilz = 0

// Object is a geojson object
type Object interface {
	bboxPtr() *BBox
	hasPositions() bool
	// WithinBBox detects if the object is fully contained inside a bbox.
	WithinBBox(bbox BBox) bool
	// IntersectsBBox detects if the object intersects a bbox.
	IntersectsBBox(bbox BBox) bool
	// Within detects if the object is fully contained inside another object.
	Within(o Object) bool
	// Intersects detects if the object intersects another object.
	Intersects(o Object) bool
	// Nearby detects if the object is nearby a position.
	Nearby(center Position, meters float64) bool
	// CalculatedBBox is exterior bbox containing the object.
	CalculatedBBox() BBox
	// CalculatedPoint is a point representation of the object.
	CalculatedPoint() Position
	// JSON is the json representation of the object. This might not be exactly the same as the original.
	JSON() string
	// String returns a string represenation of the object. This may be JSON or something else.
	String() string
	// Bytes is the bytes representation of the object.
	Bytes() []byte
	// PositionCount return the number of coordinates.
	PositionCount() int
	// Weight returns the in-memory size of the object.
	Weight() int
	// MarshalJSON allows the object to be encoded in json.Marshal calls.
	MarshalJSON() ([]byte, error)
	// Geohash converts the object to a geohash value.
	Geohash(precision int) (string, error)
	// IsBBoxDefined returns true if the object has a defined bbox.
	IsBBoxDefined() bool
	// IsGeometry return true if the object is a geojson geometry object. false if it something else.
	IsGeometry() bool
}

func writeHeader(buf *bytes.Buffer, objType byte, bbox *BBox, isCordZ bool) {
	header := objType
	if bbox != nil {
		header |= 1 << 4
		if bbox.Min.Z != nilz || bbox.Max.Z != nilz {
			header |= 1 << 5
		}
	}
	if isCordZ {
		header |= 1 << 6
	}
	buf.WriteByte(header)
	if bbox != nil {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, math.Float64bits(bbox.Min.X))
		buf.Write(b)
		binary.LittleEndian.PutUint64(b, math.Float64bits(bbox.Min.Y))
		buf.Write(b)
		if bbox.Min.Z != nilz || bbox.Max.Z != nilz {
			binary.LittleEndian.PutUint64(b, math.Float64bits(bbox.Min.Z))
			buf.Write(b)
		}
		binary.LittleEndian.PutUint64(b, math.Float64bits(bbox.Max.X))
		buf.Write(b)
		binary.LittleEndian.PutUint64(b, math.Float64bits(bbox.Max.Y))
		buf.Write(b)
		if bbox.Min.Z != nilz || bbox.Max.Z != nilz {
			binary.LittleEndian.PutUint64(b, math.Float64bits(bbox.Max.Z))
			buf.Write(b)
		}
	}
}

func positionBBox(i int, bbox BBox, ps []Position) (int, BBox) {
	for _, p := range ps {
		if i == 0 {
			bbox.Min = p
			bbox.Max = p
		} else {
			if p.X < bbox.Min.X {
				bbox.Min.X = p.X
			}
			if p.Y < bbox.Min.Y {
				bbox.Min.Y = p.Y
			}
			if p.X > bbox.Max.X {
				bbox.Max.X = p.X
			}
			if p.Y > bbox.Max.Y {
				bbox.Max.Y = p.Y
			}
		}
		i++
	}
	return i, bbox
}

func isLinearRing(ps []Position) bool {
	return len(ps) >= 4 && ps[0] == ps[len(ps)-1]
}

// ObjectJSON parses geojson and returns an Object
func ObjectJSON(s string) (Object, error) {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, err
	}
	return objectMap(m, root)
}

var (
	root  = 0 // accept all types
	gcoll = 1 // accept only geometries
	feat  = 2 // accept only geometries
	fcoll = 3 // accept only features
)

func objectMap(m map[string]interface{}, from int) (Object, error) {
	var err error
	typ, ok := m["type"].(string)
	if !ok {
		return nil, errInvalidTypeMember
	}

	if from != root {
		ok = false
		switch from {
		case gcoll, feat:
			switch typ {
			case "Point", "MultiPoint", "LineString", "MultiLineString", "Polygon", "MultiPolygon", "GeometryCollection":
				ok = true
			}
		case fcoll:
			switch typ {
			case "Feature":
				ok = true
			}
		}
		if !ok {
			return nil, fmt.Errorf(fmtErrTypeIsUnknown, typ)
		}
	}

	var o Object
	switch typ {
	default:
		return nil, fmt.Errorf(fmtErrTypeIsUnknown, typ)
	case "Point":
		o, _, err = fillSimplePointOrPoint(fillLevel1Map(m))
	case "MultiPoint":
		o, _, err = fillMultiPoint(fillLevel2Map(m))
	case "LineString":
		o, _, err = fillLineString(fillLevel2Map(m))
	case "MultiLineString":
		o, _, err = fillMultiLineString(fillLevel3Map(m))
	case "Polygon":
		o, _, err = fillPolygon(fillLevel3Map(m))
	case "MultiPolygon":
		o, _, err = fillMultiPolygon(fillLevel4Map(m))
	case "GeometryCollection":
		o, _, err = fillGeometryCollectionMap(m)
	case "Feature":
		o, _, err = fillFeatureMap(m)
	case "FeatureCollection":
		o, _, err = fillFeatureCollectionMap(m)
	}
	return o, err
}

// ObjectBytes parses geojson bytes and returns an Object
func ObjectBytes(b []byte) (Object, error) {
	var o Object
	var err error
	o, b, err = objectBytes(b)
	if err != nil {
		return nil, err
	}
	if len(b) > 0 {
		return nil, errTooMuchData
	}
	return o, nil
}

// ObjectAuto parses geojson bytes or json and returns an Object
func ObjectAuto(b []byte) (Object, error) {
	if len(b) == 0 {
		return nil, errNotEnoughData
	}
	// Check both routes. Take an educated guess at which to try first.
	var o Object
	var err error
	switch b[0] {
	default:
		o, err = ObjectBytes(b)
		if err != nil {
			o, err = ObjectJSON(string(b))
		}
	case '{', ' ', '\r', '\n':
		o, err = ObjectJSON(string(b))
		if err != nil {
			o, err = ObjectBytes(b)
		}
	}
	return o, err
}

func objectBytes(b []byte) (Object, []byte, error) {
	if len(b) == 0 {
		return nil, b, errNotEnoughData
	}
	var objType = b[0] & 0xF
	var hasBBox = (b[0]>>4)&1 == 1
	var isBBoxZ = (b[0]>>5)&1 == 1
	var isCordZ = (b[0]>>6)&1 == 1
	var bbox *BBox
	b = b[1:] // strip header
	if hasBBox {
		bbox = &BBox{}
		if len(b) < 8 {
			return nil, b, errNotEnoughData
		}
		bbox.Min.X = math.Float64frombits(binary.LittleEndian.Uint64(b))
		b = b[8:]
		if len(b) < 8 {
			return nil, b, errNotEnoughData
		}
		bbox.Min.Y = math.Float64frombits(binary.LittleEndian.Uint64(b))
		b = b[8:]
		if isBBoxZ {
			if len(b) < 8 {
				return nil, b, errNotEnoughData
			}
			bbox.Min.Z = math.Float64frombits(binary.LittleEndian.Uint64(b))
			b = b[8:]
		} else {
			bbox.Min.Z = nilz
		}
		bbox.Max.X = math.Float64frombits(binary.LittleEndian.Uint64(b))
		b = b[8:]
		if len(b) < 8 {
			return nil, b, errNotEnoughData
		}
		bbox.Max.Y = math.Float64frombits(binary.LittleEndian.Uint64(b))
		b = b[8:]
		if isBBoxZ {
			if len(b) < 8 {
				return nil, b, errNotEnoughData
			}
			bbox.Max.Z = math.Float64frombits(binary.LittleEndian.Uint64(b))
			b = b[8:]
		} else {
			bbox.Max.Z = nilz
		}
	}
	var err error
	var o Object
	switch objType {
	default:
		return nil, b, errors.New("invalid type")
	case point:
		o, b, err = fillSimplePointOrPoint(fillLevel1Bytes(b, bbox, isCordZ))
	case multiPoint:
		o, b, err = fillMultiPoint(fillLevel2Bytes(b, bbox, isCordZ))
	case lineString:
		o, b, err = fillLineString(fillLevel2Bytes(b, bbox, isCordZ))
	case multiLineString:
		o, b, err = fillMultiLineString(fillLevel3Bytes(b, bbox, isCordZ))
	case polygon:
		o, b, err = fillPolygon(fillLevel3Bytes(b, bbox, isCordZ))
	case multiPolygon:
		o, b, err = fillMultiPolygon(fillLevel4Bytes(b, bbox, isCordZ))
	case geometryCollection:
		o, b, err = fillGeometryCollectionBytes(b, bbox, isCordZ)
	case feature:
		o, b, err = fillFeatureBytes(b, bbox, isCordZ)
	case featureCollection:
		o, b, err = fillFeatureCollectionBytes(b, bbox, isCordZ)
	}
	if err != nil {
		return nil, b, err
	}
	return o, b, nil
}

func withinObjectShared(g Object, o Object, pin func(v Polygon) bool, mpin func(v MultiPolygon) bool) bool {
	bbp := o.bboxPtr()
	if bbp != nil {
		return g.WithinBBox(*bbp)
	}
	switch v := o.(type) {
	default:
		return false
	case Polygon:
		if len(v.Coordinates) == 0 {
			return false
		}
		return pin(v)
	case MultiPolygon:
		if len(v.Coordinates) == 0 {
			return false
		}
		return mpin(v)
	case Feature:
		return g.Within(v.Geometry)
	case FeatureCollection:
		if len(v.Features) == 0 {
			return false
		}
		for _, f := range v.Features {
			if !g.Within(f) {
				return false
			}
		}
		return true
	case GeometryCollection:
		if len(v.Geometries) == 0 {
			return false
		}
		for _, f := range v.Geometries {
			if !g.Within(f) {
				return false
			}
		}
		return true
	}
}

func intersectsObjectShared(g Object, o Object, pin func(v Polygon) bool, mpin func(v MultiPolygon) bool) bool {
	bbp := o.bboxPtr()
	if bbp != nil {
		return g.IntersectsBBox(*bbp)
	}
	switch v := o.(type) {
	default:
		return false
	case Polygon:
		if len(v.Coordinates) == 0 {
			return false
		}
		return pin(v)
	case MultiPolygon:
		if len(v.Coordinates) == 0 {
			return false
		}
		return mpin(v)
	case Feature:
		return g.Intersects(v.Geometry)
	case FeatureCollection:
		if len(v.Features) == 0 {
			return false
		}
		for _, f := range v.Features {
			if g.Intersects(f) {
				return true
			}
		}
		return false
	case GeometryCollection:
		if len(v.Geometries) == 0 {
			return false
		}
		for _, f := range v.Geometries {
			if g.Intersects(f) {
				return true
			}
		}
		return false
	}
}

// CirclePolygon returns a Polygon around the radius.
func CirclePolygon(x, y, meters float64, steps int) Polygon {
	if steps < 3 {
		steps = 3
	}
	p := Polygon{
		Coordinates: [][]Position{make([]Position, steps+1)},
	}
	center := Position{X: x, Y: y, Z: 0}
	step := 360.0 / float64(steps)
	i := 0
	for deg := float64(0); deg < 360; deg += step {
		c := Position(poly.Point(center.Destination(meters, deg)))
		p.Coordinates[0][i] = c
		i++
	}
	p.Coordinates[0][i] = p.Coordinates[0][0]
	return p
}

// The object's calculated bounding box must intersect the radius of the circle to pass.
func nearbyObjectShared(g Object, x, y float64, meters float64) bool {
	if !g.hasPositions() {
		return false
	}
	center := Position{X: x, Y: y, Z: 0}
	bbox := g.CalculatedBBox()
	if bbox.Min.X == bbox.Max.X && bbox.Min.Y == bbox.Max.Y {
		// just a point, return is point is inside of the circle
		return center.DistanceTo(bbox.Min) <= meters
	}
	circlePoly := CirclePolygon(x, y, meters, 12)
	return g.Intersects(circlePoly)
}
