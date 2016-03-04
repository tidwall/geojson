package geojson

import "github.com/tidwall/tile38/geojson/geohash"

// LineString is a geojson object with the type "LineString"
type LineString struct {
	Coordinates []Position
	BBox        *BBox
}

func fillLineString(coordinates []Position, bbox *BBox, b []byte, err error) (LineString, []byte, error) {
	if err == nil {
		if len(coordinates) < 2 {
			err = errLineStringInvalidCoordinates
		}
	}
	return LineString{
		Coordinates: coordinates,
		BBox:        bbox,
	}, b, err
}

// CalculatedBBox is exterior bbox containing the object.
func (g LineString) CalculatedBBox() BBox {
	return level2CalculatedBBox(g.Coordinates, g.BBox)
}

// CalculatedPoint is a point representation of the object.
func (g LineString) CalculatedPoint() Position {
	return g.CalculatedBBox().center()
}

// Geohash converts the object to a geohash value.
func (g LineString) Geohash(precision int) (string, error) {
	p := g.CalculatedPoint()
	return geohash.Encode(p.Y, p.X, precision)
}

// PositionCount return the number of coordinates.
func (g LineString) PositionCount() int {
	return level2PositionCount(g.Coordinates, g.BBox)
}

// Weight returns the in-memory size of the object.
func (g LineString) Weight() int {
	return level2Weight(g.Coordinates, g.BBox)
}

// MarshalJSON allows the object to be encoded in json.Marshal calls.
func (g LineString) MarshalJSON() ([]byte, error) {
	return []byte(g.JSON()), nil
}

// JSON is the json representation of the object. This might not be exactly the same as the original.
func (g LineString) JSON() string {
	return level2JSON("LineString", g.Coordinates, g.BBox)
}

// Bytes is the bytes representation of the object.
func (g LineString) Bytes() []byte {
	return level2Bytes(lineString, g.Coordinates, g.BBox)
}
func (g LineString) bboxPtr() *BBox {
	return g.BBox
}
func (g LineString) hasPositions() bool {
	return g.BBox != nil || len(g.Coordinates) > 0
}

// WithinBBox detects if the object is fully contained inside a bbox.
func (g LineString) WithinBBox(bbox BBox) bool {
	if g.BBox != nil {
		return rectBBox(g.CalculatedBBox()).InsideRect(rectBBox(bbox))
	}
	return polyPositions(g.Coordinates).InsideRect(rectBBox(bbox))
}

// IntersectsBBox detects if the object intersects a bbox.
func (g LineString) IntersectsBBox(bbox BBox) bool {
	if g.BBox != nil {
		return rectBBox(g.CalculatedBBox()).IntersectsRect(rectBBox(bbox))
	}
	return polyPositions(g.Coordinates).IntersectsRect(rectBBox(bbox))
}

// Within detects if the object is fully contained inside another object.
func (g LineString) Within(o Object) bool {
	return withinObjectShared(g, o,
		func(v Polygon) bool {
			return polyPositions(g.Coordinates).Inside(polyExteriorHoles(v.Coordinates))
		},
		func(v MultiPolygon) bool {
			for _, c := range v.Coordinates {
				if !polyPositions(g.Coordinates).Inside(polyExteriorHoles(c)) {
					return false
				}
			}
			return true
		},
	)
}

// Intersects detects if the object intersects another object.
func (g LineString) Intersects(o Object) bool {
	return intersectsObjectShared(g, o,
		func(v Polygon) bool {
			return polyPositions(g.Coordinates).Intersects(polyExteriorHoles(v.Coordinates))
		},
		func(v MultiPolygon) bool {
			for _, c := range v.Coordinates {
				if polyPositions(g.Coordinates).Intersects(polyExteriorHoles(c)) {
					return true
				}
			}
			return false
		},
	)
}

// Nearby detects if the object is nearby a position.
func (g LineString) Nearby(center Position, meters float64) bool {
	return nearbyObjectShared(g, center.X, center.Y, meters)
}

// IsBBoxDefined returns true if the object has a defined bbox.
func (g LineString) IsBBoxDefined() bool {
	return g.BBox != nil
}
