package geojson

import (
	"github.com/tidwall/tile38/pkg/geojson/geohash"
	"github.com/tidwall/tile38/pkg/geojson/poly"
)

// MultiLineString is a geojson object with the type "MultiLineString"
type MultiLineString struct {
	Coordinates [][]Position
	BBox        *BBox
	bboxDefined bool
}

func fillMultiLineString(coordinates [][]Position, bbox *BBox, err error) (MultiLineString, error) {
	if err == nil {
		for _, coordinates := range coordinates {
			if len(coordinates) < 2 {
				err = errLineStringInvalidCoordinates
				break
			}
		}
	}
	bboxDefined := bbox != nil
	if !bboxDefined {
		cbbox := level3CalculatedBBox(coordinates, nil, false)
		bbox = &cbbox
	}
	return MultiLineString{
		Coordinates: coordinates,
		BBox:        bbox,
		bboxDefined: bboxDefined,
	}, err
}

func (g MultiLineString) getLineString(index int) LineString {
	return LineString{Coordinates: g.Coordinates[index]}
}

// CalculatedBBox is exterior bbox containing the object.
func (g MultiLineString) CalculatedBBox() BBox {
	return level3CalculatedBBox(g.Coordinates, g.BBox, false)
}

// CalculatedPoint is a point representation of the object.
func (g MultiLineString) CalculatedPoint() Position {
	return g.CalculatedBBox().center()
}

// Geohash converts the object to a geohash value.
func (g MultiLineString) Geohash(precision int) (string, error) {
	p := g.CalculatedPoint()
	return geohash.Encode(p.Y, p.X, precision)
}

// PositionCount return the number of coordinates.
func (g MultiLineString) PositionCount() int {
	return level3PositionCount(g.Coordinates, g.BBox)
}

// Weight returns the in-memory size of the object.
func (g MultiLineString) Weight() int {
	return level3Weight(g.Coordinates, g.BBox)
}

// MarshalJSON allows the object to be encoded in json.Marshal calls.
func (g MultiLineString) MarshalJSON() ([]byte, error) {
	return g.appendJSON(nil), nil
}

func (g MultiLineString) appendJSON(json []byte) []byte {
	return appendLevel3JSON(json, "MultiLineString", g.Coordinates, g.BBox, g.bboxDefined)
}

// JSON is the json representation of the object. This might not be exactly the same as the original.
func (g MultiLineString) JSON() string {
	return string(g.appendJSON(nil))
}

// String returns a string representation of the object. This might be JSON or something else.
func (g MultiLineString) String() string {
	return g.JSON()
}

func (g MultiLineString) bboxPtr() *BBox {
	return g.BBox
}
func (g MultiLineString) hasPositions() bool {
	if g.bboxDefined {
		return true
	}
	for _, c := range g.Coordinates {
		if len(c) > 0 {
			return true
		}
	}
	return false
}

// WithinBBox detects if the object is fully contained inside a bbox.
func (g MultiLineString) WithinBBox(bbox BBox) bool {
	if g.bboxDefined {
		return rectBBox(g.CalculatedBBox()).InsideRect(rectBBox(bbox))
	}
	if len(g.Coordinates) == 0 {
		return false
	}
	for _, ls := range g.Coordinates {
		if len(ls) == 0 {
			return false
		}
		for _, p := range ls {
			if !poly.Point(p).InsideRect(rectBBox(bbox)) {
				return false
			}
		}
	}
	return true
}

// IntersectsBBox detects if the object intersects a bbox.
func (g MultiLineString) IntersectsBBox(bbox BBox) bool {
	if g.bboxDefined {
		return rectBBox(g.CalculatedBBox()).IntersectsRect(rectBBox(bbox))
	}
	for _, ls := range g.Coordinates {
		if polyPositions(ls).IntersectsRect(rectBBox(bbox)) {
			return true
		}
	}
	return false
}

// Within detects if the object is fully contained inside another object.
func (g MultiLineString) Within(o Object) bool {
	return withinObjectShared(g, o,
		func(v Polygon) bool {
			if len(g.Coordinates) == 0 {
				return false
			}
			for _, ls := range g.Coordinates {
				if !polyPositions(ls).Inside(polyExteriorHoles(v.Coordinates)) {
					return false
				}
			}
			return true
		},
	)
}

// WithinCircle detects if the object is fully contained inside a circle.
func (g MultiLineString) WithinCircle(center Position, meters float64) bool {
	if len(g.Coordinates) == 0 {
		return false
	}
	for _, ls := range g.Coordinates {
		if len(ls) == 0 {
			return false
		}
		for _, position := range ls {
			if center.DistanceTo(position) >= meters {
				return false
			}
		}
	}
	return true
}

// Intersects detects if the object intersects another object.
func (g MultiLineString) Intersects(o Object) bool {
	return intersectsObjectShared(g, o,
		func(v Polygon) bool {
			if len(g.Coordinates) == 0 {
				return false
			}
			for _, ls := range g.Coordinates {
				if polyPositions(ls).Intersects(polyExteriorHoles(v.Coordinates)) {
					return true
				}
			}
			return false
		},
	)
}

// IntersectsCircle detects if the object intersects a circle.
func (g MultiLineString) IntersectsCircle(center Position, meters float64) bool {
	for _, ls := range g.Coordinates {
		for i := 0; i < len(ls) - 1 ; i++ {
			if SegmentIntersectsCircle(ls[i], ls[i + 1], center, meters) {
				return true
			}
		}
	}
	return false
}

// Nearby detects if the object is nearby a position.
func (g MultiLineString) Nearby(center Position, meters float64) bool {
	return nearbyObjectShared(g, center.X, center.Y, meters)
}

// IsBBoxDefined returns true if the object has a defined bbox.
func (g MultiLineString) IsBBoxDefined() bool {
	return g.bboxDefined
}

// IsGeometry return true if the object is a geojson geometry object. false if it something else.
func (g MultiLineString) IsGeometry() bool {
	return true
}

// Clip returns the object obtained by clipping this object by a bbox.
func (g MultiLineString) Clipped(bbox BBox) Object {
	var new_coordinates [][]Position

	for ix := range g.Coordinates {
		clippedMultiLineString, _ := g.getLineString(ix).Clipped(bbox).(MultiLineString)
		for _, ls := range clippedMultiLineString.Coordinates {
			new_coordinates = append(new_coordinates, ls)
		}
	}

	res, _ := fillMultiLineString(new_coordinates, nil, nil)
	return res
}
