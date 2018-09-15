package geojson

import (
	"github.com/tidwall/tile38/pkg/geojson/geo"
	"github.com/tidwall/tile38/pkg/geojson/geohash"
	"github.com/tidwall/tile38/pkg/geojson/poly"
)

// Point is a geojson object with the type "Point"
type Point struct {
	Coordinates Position
	BBox        *BBox
	bboxDefined bool
}

func fillSimplePointOrPoint(coordinates Position, bbox *BBox, err error) (Object, error) {
	if coordinates.Z == 0 && bbox == nil {
		return fillSimplePoint(coordinates, bbox, err)
	}
	return fillPoint(coordinates, bbox, err)
}

func fillPoint(coordinates Position, bbox *BBox, err error) (Point, error) {
	bboxDefined := bbox != nil
	if !bboxDefined {
		cbbox := level1CalculatedBBox(coordinates, nil)
		bbox = &cbbox
	}
	return Point{
		Coordinates: coordinates,
		BBox:        bbox,
		bboxDefined: bboxDefined,
	}, err
}

// CalculatedBBox is exterior bbox containing the object.
func (g Point) CalculatedBBox() BBox {
	return level1CalculatedBBox(g.Coordinates, g.BBox)
}

// CalculatedPoint is a point representation of the object.
func (g Point) CalculatedPoint() Position {
	if g.BBox == nil {
		return g.Coordinates
	}
	return g.CalculatedBBox().center()
}

// Geohash converts the object to a geohash value.
func (g Point) Geohash(precision int) (string, error) {
	p := g.CalculatedPoint()
	return geohash.Encode(p.Y, p.X, precision)
}

// MarshalJSON allows the object to be encoded in json.Marshal calls.
func (g Point) MarshalJSON() ([]byte, error) {
	return g.appendJSON(nil), nil
}

func (g Point) appendJSON(json []byte) []byte {
	return appendLevel1JSON(json, "Point", g.Coordinates, g.BBox, g.bboxDefined)
}

// JSON is the json representation of the object. This might not be exactly the same as the original.
func (g Point) JSON() string {
	return string(g.appendJSON(nil))
}

// String returns a string representation of the object. This might be JSON or something else.
func (g Point) String() string {
	return g.JSON()
}

// PositionCount return the number of coordinates.
func (g Point) PositionCount() int {
	return level1PositionCount(g.Coordinates, g.BBox)
}

// Weight returns the in-memory size of the object.
func (g Point) Weight() int {
	return level1Weight(g.Coordinates, g.BBox)
}
func (g Point) bboxPtr() *BBox {
	return g.BBox
}
func (g Point) hasPositions() bool {
	return true
}

// WithinBBox detects if the object is fully contained inside a bbox.
func (g Point) WithinBBox(bbox BBox) bool {
	if g.bboxDefined {
		return rectBBox(g.CalculatedBBox()).InsideRect(rectBBox(bbox))
	}
	return poly.Point(g.Coordinates).InsideRect(rectBBox(bbox))
}

// IntersectsBBox detects if the object intersects a bbox.
func (g Point) IntersectsBBox(bbox BBox) bool {
	if g.bboxDefined {
		return rectBBox(g.CalculatedBBox()).IntersectsRect(rectBBox(bbox))
	}
	return poly.Point(g.Coordinates).InsideRect(rectBBox(bbox))
}

// Within detects if the object is fully contained inside another object.
func (g Point) Within(o Object) bool {
	return withinObjectShared(g, o,
		func(v Polygon) bool {
			return poly.Point(g.Coordinates).Inside(polyExteriorHoles(v.Coordinates))
		},
	)
}

// WithinCircle detects if the object is fully contained inside a circle.
func (g Point) WithinCircle(center Position, meters float64) bool {
	return geo.DistanceTo(g.Coordinates.Y, g.Coordinates.X, center.Y, center.X) < meters
}

// Intersects detects if the object intersects another object.
func (g Point) Intersects(o Object) bool {
	return intersectsObjectShared(g, o,
		func(v Polygon) bool {
			return poly.Point(g.Coordinates).Intersects(polyExteriorHoles(v.Coordinates))
		},
	)
}

// IntersectsCircle detects if the object intersects a circle.
func (g Point) IntersectsCircle(center Position, meters float64) bool {
	return geo.DistanceTo(g.Coordinates.Y, g.Coordinates.X, center.Y, center.X) <= meters
}

// Nearby detects if the object is nearby a position.
func (g Point) Nearby(center Position, meters float64) bool {
	return geo.DistanceTo(g.Coordinates.Y, g.Coordinates.X, center.Y, center.X) <= meters
}

// IsBBoxDefined returns true if the object has a defined bbox.
func (g Point) IsBBoxDefined() bool {
	return g.bboxDefined
}

// IsGeometry return true if the object is a geojson geometry object. false if it something else.
func (g Point) IsGeometry() bool {
	return true
}

// Clip returns the object obtained by clipping this object by a bbox.
func (g Point) Clipped(bbox BBox) Object {
	if g.IntersectsBBox(bbox) {
		return g
	}

	res, _ := fillMultiPoint([]Position{}, nil, nil)
	return res
}
