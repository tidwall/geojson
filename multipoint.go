package geojson

import (
	"github.com/tidwall/tile38/pkg/geojson/geohash"
	"github.com/tidwall/tile38/pkg/geojson/poly"
)

// MultiPoint is a geojson object with the type "MultiPoint"
type MultiPoint struct {
	Coordinates []Position
	BBox        *BBox
	bboxDefined bool
}

func fillMultiPoint(coordinates []Position, bbox *BBox, err error) (MultiPoint, error) {
	bboxDefined := bbox != nil
	if !bboxDefined {
		cbbox := level2CalculatedBBox(coordinates, nil)
		bbox = &cbbox
	}
	return MultiPoint{
		Coordinates: coordinates,
		BBox:        bbox,
		bboxDefined: bboxDefined,
	}, err
}

// CalculatedBBox is exterior bbox containing the object.
func (g MultiPoint) CalculatedBBox() BBox {
	return level2CalculatedBBox(g.Coordinates, g.BBox)
}

// CalculatedPoint is a point representation of the object.
func (g MultiPoint) CalculatedPoint() Position {
	return g.CalculatedBBox().center()
}

// Geohash converts the object to a geohash value.
func (g MultiPoint) Geohash(precision int) (string, error) {
	p := g.CalculatedPoint()
	return geohash.Encode(p.Y, p.X, precision)
}

// PositionCount return the number of coordinates.
func (g MultiPoint) PositionCount() int {
	return level2PositionCount(g.Coordinates, g.BBox)
}

// Weight returns the in-memory size of the object.
func (g MultiPoint) Weight() int {
	return level2Weight(g.Coordinates, g.BBox)
}

// MarshalJSON allows the object to be encoded in json.Marshal calls.
func (g MultiPoint) MarshalJSON() ([]byte, error) {
	return g.appendJSON(nil), nil
}

func (g MultiPoint) appendJSON(json []byte) []byte {
	return appendLevel2JSON(json, "MultiPoint", g.Coordinates, g.BBox, g.bboxDefined)
}

// JSON is the json representation of the object. This might not be exactly the same as the original.
func (g MultiPoint) JSON() string {
	return string(g.appendJSON(nil))
}

// String returns a string representation of the object. This might be JSON or something else.
func (g MultiPoint) String() string {
	return g.JSON()
}

func (g MultiPoint) bboxPtr() *BBox {
	return g.BBox
}
func (g MultiPoint) hasPositions() bool {
	return g.bboxDefined || len(g.Coordinates) > 0
}

// WithinBBox detects if the object is fully contained inside a bbox.
func (g MultiPoint) WithinBBox(bbox BBox) bool {
	if g.bboxDefined {
		return rectBBox(g.CalculatedBBox()).InsideRect(rectBBox(bbox))
	}
	if len(g.Coordinates) == 0 {
		return false
	}
	for _, p := range g.Coordinates {
		if !poly.Point(p).InsideRect(rectBBox(bbox)) {
			return false
		}
	}
	return true
}

// IntersectsBBox detects if the object intersects a bbox.
func (g MultiPoint) IntersectsBBox(bbox BBox) bool {
	if g.bboxDefined {
		return rectBBox(g.CalculatedBBox()).IntersectsRect(rectBBox(bbox))
	}
	for _, p := range g.Coordinates {
		if poly.Point(p).InsideRect(rectBBox(bbox)) {
			return true
		}
	}
	return false
}

// Within detects if the object is fully contained inside another object.
func (g MultiPoint) Within(o Object) bool {
	return withinObjectShared(g, o,
		func(v Polygon) bool {
			if len(g.Coordinates) == 0 {
				return false
			}
			for _, p := range g.Coordinates {
				if !poly.Point(p).Inside(polyExteriorHoles(v.Coordinates)) {
					return false
				}
			}
			return true
		},
	)
}

// WithinCircle detects if the object is fully contained inside a circle.
func (g MultiPoint) WithinCircle(center Position, meters float64) bool {
	if len(g.Coordinates) == 0 {
		return false
	}
	for _, position := range g.Coordinates {
		if center.DistanceTo(position) >= meters {
			return false
		}
	}
	return true
}

// Intersects detects if the object intersects another object.
func (g MultiPoint) Intersects(o Object) bool {
	return intersectsObjectShared(g, o,
		func(v Polygon) bool {
			if len(g.Coordinates) == 0 {
				return false
			}
			for _, p := range g.Coordinates {
				if poly.Point(p).Intersects(polyExteriorHoles(v.Coordinates)) {
					return true
				}
			}
			return true
		},
	)
}

// IntersectsCircle detects if the object intersects a circle.
func (g MultiPoint) IntersectsCircle(center Position, meters float64) bool {
	for _, position := range g.Coordinates {
		if center.DistanceTo(position) <= meters {
			return true
		}
	}
	return false
}

// Nearby detects if the object is nearby a position.
func (g MultiPoint) Nearby(center Position, meters float64) bool {
	return nearbyObjectShared(g, center.X, center.Y, meters)
}

// IsBBoxDefined returns true if the object has a defined bbox.
func (g MultiPoint) IsBBoxDefined() bool {
	return g.bboxDefined
}

// IsGeometry return true if the object is a geojson geometry object. false if it something else.
func (g MultiPoint) IsGeometry() bool {
	return true
}

// Clip returns the object obtained by clipping this object by a bbox.
func (g MultiPoint) Clipped(bbox BBox) Object {
	var new_coordinates []Position

	for _, position := range g.Coordinates {
		if poly.Point(position).InsideRect(rectBBox(bbox)) {
			new_coordinates = append(new_coordinates, position)
		}
	}

	res, _ := fillMultiPoint(new_coordinates, nil, nil)

	return res
}
