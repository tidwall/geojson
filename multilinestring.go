package geojson

import (
	"github.com/tidwall/tile38/geojson/geohash"
)

// MultiLineString is a geojson object with the type "MultiLineString"
type MultiLineString struct {
	Coordinates [][]Position
	BBox        *BBox
	bboxDefined bool
	linestrings []LineString
}

func fillMultiLineString(coordinates [][]Position, bbox *BBox, err error) (MultiLineString, error) {
	linestrings := make([]LineString, len(coordinates))
	if err == nil {
		for i, ps := range coordinates {
			linestrings[i], err = fillLineString(ps, nil, nil)
			if err != nil {
				break
			}
		}
	}
	bboxDefined := bbox != nil
	if !bboxDefined {
		cbbox := mlCalculatedBBox(linestrings, nil)
		bbox = &cbbox
	}
	return MultiLineString{
		Coordinates: coordinates,
		BBox:        bbox,
		bboxDefined: bboxDefined,
		linestrings: linestrings,
	}, err
}

func mlCalculatedBBox(linestrings []LineString, bbox *BBox) BBox {
	if bbox != nil {
		return *bbox
	}
	var cbbox BBox
	for i, g := range linestrings {
		if i == 0 {
			cbbox = g.CalculatedBBox()
		} else {
			cbbox = cbbox.union(g.CalculatedBBox())
		}
	}
	return cbbox
}

// CalculatedBBox is exterior bbox containing the object.
func (g MultiLineString) CalculatedBBox() BBox {
	return mlCalculatedBBox(g.linestrings, g.BBox)
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

func (g MultiLineString) getLineString(index int) LineString {
	if index < len(g.linestrings) {
		return g.linestrings[index]
	}
	return LineString{Coordinates: g.Coordinates[index]}
}

// WithinBBox detects if the object is fully contained inside a bbox.
func (g MultiLineString) WithinBBox(bbox BBox) bool {
	if g.bboxDefined {
		return rectBBox(g.CalculatedBBox()).InsideRect(rectBBox(bbox))
	}
	if len(g.Coordinates) == 0 {
		return false
	}
	for i := range g.Coordinates {
		if !g.getLineString(i).WithinBBox(bbox) {
			return false
		}
	}
	return true
}

// IntersectsBBox detects if the object intersects a bbox.
func (g MultiLineString) IntersectsBBox(bbox BBox) bool {
	if g.bboxDefined {
		return rectBBox(g.CalculatedBBox()).IntersectsRect(rectBBox(bbox))
	}
	for i := range g.Coordinates {
		if g.getLineString(i).IntersectsBBox(bbox) {
			return true
		}
	}
	return false
}

// Within detects if the object is fully contained inside another object.
func (g MultiLineString) Within(o Object) bool {
	return withinObjectShared(g, o)
}

// Intersects detects if the object intersects another object.
func (g MultiLineString) Intersects(o Object) bool {
	return intersectsObjectShared(g, o)
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
