package geojson

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/tidwall/geojson/geo"
	"github.com/tidwall/geojson/geometry"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
)

var (
	fmtErrTypeIsUnknown         = "type '%s' is unknown"
	errDataInvalid              = errors.New("invalid data")
	errTypeInvalid              = errors.New("invalid type")
	errTypeMissing              = errors.New("missing type")
	errCoordinatesInvalid       = errors.New("invalid coordinates")
	errCoordinatesMissing       = errors.New("missing coordinates")
	errGeometryInvalid          = errors.New("invalid geometry")
	errGeometryMissing          = errors.New("missing geometry")
	errFeaturesMissing          = errors.New("missing features")
	errFeaturesInvalid          = errors.New("invalid features")
	errGeometriesMissing        = errors.New("missing geometries")
	errGeometriesInvalid        = errors.New("invalid geometries")
	errCircleRadiusUnitsInvalid = errors.New("invalid circle radius units")
)

// Object is a GeoJSON type
type Object interface {
	Empty() bool
	Valid() bool
	Rect() geometry.Rect
	Center() geometry.Point
	Contains(other Object) bool
	Within(other Object) bool
	Intersects(other Object) bool
	AppendJSON(dst []byte) []byte
	JSON() string
	String() string
	Distance(obj Object) float64
	NumPoints() int
	ForEach(iter func(geom Object) bool) bool
	Spatial() Spatial
	MarshalJSON() ([]byte, error)
	AppendBinary(dst []byte) []byte
	Binary() []byte
}

var _ = []Object{
	&Point{}, &LineString{}, &Polygon{}, &Feature{},
	&MultiPoint{}, &MultiLineString{}, &MultiPolygon{},
	&GeometryCollection{}, &FeatureCollection{},
	&Rect{}, &Circle{}, &SimplePoint{},
}

// Collection is a searchable collection type.
type Collection interface {
	Children() []Object
	Indexed() bool
	Search(rect geometry.Rect, iter func(child Object) bool)
}

var _ = []Collection{
	&MultiPoint{}, &MultiLineString{}, &MultiPolygon{},
	&FeatureCollection{}, &GeometryCollection{},
}

type extra struct {
	dims   byte      // number of extra coordinate values, 1 or 2
	values []float64 // extra coordinate values
	// valid json object that includes extra members such as
	// "bbox", "id", "properties", and foreign members
	members string
}

// ParseOptions ...
type ParseOptions struct {
	// IndexChildren option will cause the object to index their children
	// objects when the number of children is greater than or equal to the
	// provided value. Setting this value to 0 will disable indexing.
	// The default is 64.
	IndexChildren int
	// IndexGeometry option will cause the object to index it's geometry
	// when the number of points in it's base polygon or linestring is greater
	// that or equal to the provided value. Setting this value to 0 will
	// disable indexing.
	// The default is 64.
	IndexGeometry int
	// IndexGeometryKind is the kind of index implementation.
	// Default is QuadTreeCompressed
	IndexGeometryKind geometry.IndexKind
	// RequireValid option cause parse to fail when a geojson object is invalid.
	RequireValid bool
	// AllowSimplePoints options will force to parse to return the SimplePoint
	// type when a geojson point only consists of an 2D x/y coord and no extra
	// json members.
	AllowSimplePoints bool
	// DisableCircleType disables the special Circle syntax that is unique to
	// only Tile38.
	DisableCircleType bool
	// AllowRects options will force to parse and return the Rect type when a
	// geojson polygon only consists of a perfect rectangle, where there are
	// exactly 5 points with the first point being the min x/y and the
	// following point winding counter clockwise creating a closed rectangle.
	AllowRects bool
}

// DefaultParseOptions ...
var DefaultParseOptions = &ParseOptions{
	IndexChildren:     64,
	IndexGeometry:     64,
	IndexGeometryKind: geometry.QuadTree,
	RequireValid:      false,
	AllowSimplePoints: false,
	DisableCircleType: false,
	AllowRects:        false,
}

// Parse a GeoJSON object
func Parse(data string, opts *ParseOptions) (Object, error) {
	if opts == nil {
		// opts should never be nil
		opts = DefaultParseOptions
	}
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
			return parseJSON(data, opts)
		}
	}
}

func toGeometryOpts(opts *ParseOptions) geometry.IndexOptions {
	var gopts geometry.IndexOptions
	if opts == nil {
		gopts = *geometry.DefaultIndexOptions
	} else {
		gopts.Kind = opts.IndexGeometryKind
		gopts.MinPoints = opts.IndexGeometry
	}
	return gopts
}

type parseKeys struct {
	rCoordinates gjson.Result
	rGeometries  gjson.Result
	rGeometry    gjson.Result
	rFeatures    gjson.Result
	members      string // a valid payload with all extra members
}

func parseJSON(data string, opts *ParseOptions) (Object, error) {
	if !gjson.Valid(data) {
		return nil, errDataInvalid
	}
	var keys parseKeys
	var fmembers []byte
	var rType gjson.Result
	gjson.Parse(data).ForEach(func(key, val gjson.Result) bool {
		switch key.String() {
		case "type":
			rType = val
		case "coordinates":
			keys.rCoordinates = val
		case "geometries":
			keys.rGeometries = val
		case "geometry":
			keys.rGeometry = val
		case "features":
			keys.rFeatures = val
		default:
			if len(fmembers) == 0 {
				fmembers = append(fmembers, '{')
			} else {
				fmembers = append(fmembers, ',')
			}
			fmembers = append(fmembers, key.Raw...)
			fmembers = append(fmembers, ':')
			fmembers = append(fmembers, val.Raw...)
		}
		return true
	})
	if len(fmembers) > 0 {
		fmembers = append(fmembers, '}')
		fmembers = pretty.UglyInPlace(fmembers)
		keys.members = string(fmembers)
	}
	if !rType.Exists() {
		return nil, errTypeMissing
	}
	if rType.Type != gjson.String {
		return nil, errTypeInvalid
	}
	switch rType.String() {
	default:
		return nil, fmt.Errorf(fmtErrTypeIsUnknown, rType.String())
	case "Point":
		return parseJSONPoint(&keys, opts)
	case "LineString":
		return parseJSONLineString(&keys, opts)
	case "Polygon":
		return parseJSONPolygon(&keys, opts)
	case "Feature":
		return parseJSONFeature(&keys, opts)
	case "MultiPoint":
		return parseJSONMultiPoint(&keys, opts)
	case "MultiLineString":
		return parseJSONMultiLineString(&keys, opts)
	case "MultiPolygon":
		return parseJSONMultiPolygon(&keys, opts)
	case "GeometryCollection":
		return parseJSONGeometryCollection(&keys, opts)
	case "FeatureCollection":
		return parseJSONFeatureCollection(&keys, opts)
	}
}

func parseBBoxAndExtras(ex **extra, keys *parseKeys, opts *ParseOptions) error {
	if keys.members == "" {
		return nil
	}
	if *ex == nil {
		*ex = new(extra)
	}
	(*ex).members = keys.members
	return nil
}

func appendJSONPoint(dst []byte, point geometry.Point, ex *extra, idx int) []byte {
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

func (ex *extra) appendJSONExtra(dst []byte, propertiesRequired bool) []byte {
	if ex != nil && ex.members != "" {
		dst = append(dst, ',')
		dst = append(dst, ex.members[1:len(ex.members)-1]...)
		if propertiesRequired {
			if !gjson.Get(ex.members, "properties").Exists() {
				dst = append(dst, `,"properties":{}`...)
			}
		}
	} else if propertiesRequired {
		dst = append(dst, `,"properties":{}`...)
	}

	return dst
}

func (ex *extra) appendBinary(dst []byte) []byte {
	if ex == nil {
		return append(dst, 0)
	}
	dst = append(dst, 1)
	dst = append(dst, ex.dims)
	dst = appendUvarint(dst, uint64(len(ex.values)))
	for _, x := range ex.values {
		dst = appendFloat64(dst, x)
	}
	dst = appendUvarint(dst, uint64(len(ex.members)))
	dst = append(dst, ex.members...)
	return dst
}

func parseBinaryPoints(src []byte) ([]geometry.Point, int) {
	mark := len(src)
	nvals, n := binary.Uvarint(src)
	if n <= 0 {
		return nil, 0
	}
	src = src[n:]
	if uint64(len(src)) < nvals*16 {
		return nil, 0
	}
	points := make([]geometry.Point, nvals)
	for i := 0; i < len(points); i++ {
		points[i] = geometry.Point{
			X: math.Float64frombits(binary.LittleEndian.Uint64(src[i*16:])),
			Y: math.Float64frombits(binary.LittleEndian.Uint64(src[i*16+8:])),
		}
	}
	src = src[nvals*16:]
	return points, mark - len(src)
}

func parseBinaryFloat64s(src []byte) ([]float64, int) {
	mark := len(src)
	nvals, n := binary.Uvarint(src)
	if n <= 0 {
		return nil, 0
	}
	src = src[n:]
	if uint64(len(src)) < nvals*8 {
		return nil, 0
	}
	values := make([]float64, nvals)
	for i := 0; i < len(values); i++ {
		x := math.Float64frombits(binary.LittleEndian.Uint64(src[i*8:]))
		values[i] = x
	}
	src = src[nvals*8:]
	return values, mark - len(src)
}

func parseBinaryExtra(src []byte) (*extra, int) {
	mark := len(src)
	if len(src) == 0 {
		return nil, 0
	}
	if src[0] == 0 {
		return nil, 1
	}
	if src[0] != 1 {
		return nil, 0
	}
	src = src[1:]

	ex := &extra{}
	if len(src) == 0 {
		return nil, 0
	}
	ex.dims = src[0]
	src = src[1:]

	var n int
	ex.values, n = parseBinaryFloat64s(src)
	if n <= 0 {
		return nil, 0
	}
	src = src[n:]

	mlen, n := binary.Uvarint(src)
	if n <= 0 {
		return nil, 0
	}
	src = src[n:]

	if uint64(len(src)) < mlen {
		return nil, 0
	}
	ex.members = string(src[:mlen])
	src = src[mlen:]

	return ex, mark - len(src)
}

func appendUvarint(dst []byte, x uint64) []byte {
	var buf [10]byte
	n := binary.PutUvarint(buf[:], x)
	return append(dst, buf[:n]...)
}

func appendFloat64(dst []byte, x float64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(x))
	return append(dst, buf[:]...)
}

func appendBinaryPoint(dst []byte, point geometry.Point) []byte {
	dst = appendFloat64(dst, point.X)
	return appendFloat64(dst, point.Y)
}

func parseBinaryPoint(src []byte) geometry.Point {
	// The size has already been checked by caller
	return geometry.Point{
		X: math.Float64frombits(binary.LittleEndian.Uint64(src[0:8])),
		Y: math.Float64frombits(binary.LittleEndian.Uint64(src[8:16])),
	}
}

func appendJSONSeries(
	dst []byte, series geometry.Series, ex *extra, pidx int,
) (ndst []byte, npidx int) {
	dst = append(dst, '[')
	nPoints := series.NumPoints()
	for i := 0; i < nPoints; i++ {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = appendJSONPoint(dst, series.PointAt(i), ex, pidx)
		pidx++
	}
	dst = append(dst, ']')
	return dst, pidx
}

func appendBinarySeries(dst []byte, series geometry.Series) (ndst []byte) {
	nPoints := series.NumPoints()
	dst = appendUvarint(dst, uint64(nPoints))
	for i := 0; i < nPoints; i++ {
		dst = appendBinaryPoint(dst, series.PointAt(i))
	}
	return dst
}

func appendBinaryPoints(dst []byte, points []geometry.Point) (ndst []byte) {
	dst = appendUvarint(dst, uint64(len(points)))
	for i := 0; i < len(points); i++ {
		dst = appendBinaryPoint(dst, points[i])
	}
	return dst
}

func unionRects(a, b geometry.Rect) geometry.Rect {
	if b.Min.X < a.Min.X {
		a.Min.X = b.Min.X
	}
	if b.Max.X > a.Max.X {
		a.Max.X = b.Max.X
	}
	if b.Min.Y < a.Min.Y {
		a.Min.Y = b.Min.Y
	}
	if b.Max.Y > a.Max.Y {
		a.Max.Y = b.Max.Y
	}
	return a
}

func geoDistancePoints(a, b geometry.Point) float64 {
	return geo.DistanceTo(a.Y, a.X, b.Y, b.X)
}

// These types do not necessarily align with WKT/WKB type integer codes.
// New values may be added in the future, but do not change the older ones
// as that might break compatibilty with dependent applications.
const (
	binPoint              byte = 1
	binLineString         byte = 2
	binPolygon            byte = 3
	binMultiPoint         byte = 4
	binMultiLineString    byte = 5
	binMultiPolygon       byte = 6
	binGeometryCollection byte = 7
	binFeature            byte = 128
	binFeatureCollection  byte = 129
	binRect               byte = 130
	binSimplePoint        byte = 131
	binCircle             byte = 132
)

// ParseBinary from bytes that where generated from the Object.Binary()
// method. Reutrns nil if there was a problem parsing.
//
// Only the fields relating to geometry indexing are used. The others are
// ignored.
func ParseBinary(src []byte, opts *ParseOptions) (Object, int) {
	if opts == nil {
		opts = DefaultParseOptions
	}
	mark := len(src)
	if len(src) == 0 || src[0] != ':' {
		return nil, 0
	}
	src = src[1:]
	if len(src) == 0 {
		return nil, 0
	}
	kind := src[0]
	src = src[1:]
	var obj Object
	var n int
	switch kind {
	case binPoint:
		obj, n = parseBinaryPointObject(src, opts)
	case binLineString:
		obj, n = parseBinaryLineStringObject(src, opts)
	case binPolygon:
		obj, n = parseBinaryPolygonObject(src, opts)
	case binMultiPoint:
		obj, n = parseBinaryMultiPointObject(src, opts)
	case binMultiLineString:
		obj, n = parseBinaryMultiLineStringObject(src, opts)
	case binMultiPolygon:
		obj, n = parseBinaryMultiPolygonObject(src, opts)
	case binGeometryCollection:
		obj, n = parseBinaryGeometryCollectionObject(src, opts)
	case binFeature:
		obj, n = parseBinaryFeatureObject(src, opts)
	case binFeatureCollection:
		obj, n = parseBinaryFeatureCollectionObject(src, opts)
	case binRect:
		obj, n = parseBinaryRectObject(src, opts)
	case binSimplePoint:
		obj, n = parseBinarySimplePointObject(src, opts)
	case binCircle:
		obj, n = parseBinaryCircleObject(src, opts)
	default:
		return nil, 0
	}
	if n <= 0 {
		return nil, 0
	}
	src = src[n:]
	return obj, mark - len(src)
}
