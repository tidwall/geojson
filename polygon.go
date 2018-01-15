package geojson

import "github.com/tidwall/tile38/geojson/geohash"

// Polygon is a geojson object with the type "Polygon"
type Polygon struct {
	Coordinates [][]Position
	BBox        *BBox
	bboxDefined bool
}

func fillPolygon(coordinates [][]Position, bbox *BBox, err error) (Polygon, error) {
	if err == nil {
		if len(coordinates) == 0 {
			err = errMustBeALinearRing
		}
	}
	if err == nil {
		for _, ps := range coordinates {
			if !isLinearRing(ps) {
				err = errMustBeALinearRing
				break
			}
		}
	}
	bboxDefined := bbox != nil
	if !bboxDefined {
		cbbox := level3CalculatedBBox(coordinates, nil, true)
		bbox = &cbbox
	}
	return Polygon{
		Coordinates: coordinates,
		BBox:        bbox,
		bboxDefined: bboxDefined,
	}, err
}

// CalculatedBBox is exterior bbox containing the object.
func (g Polygon) CalculatedBBox() BBox {
	return level3CalculatedBBox(g.Coordinates, g.BBox, true)
}

// CalculatedPoint is a point representation of the object.
func (g Polygon) CalculatedPoint() Position {
	return g.CalculatedBBox().center()
}

// Geohash converts the object to a geohash value.
func (g Polygon) Geohash(precision int) (string, error) {
	p := g.CalculatedPoint()
	return geohash.Encode(p.Y, p.X, precision)
}

// PositionCount return the number of coordinates.
func (g Polygon) PositionCount() int {
	return level3PositionCount(g.Coordinates, g.BBox)
}

// Weight returns the in-memory size of the object.
func (g Polygon) Weight() int {
	return level3Weight(g.Coordinates, g.BBox)
}

// MarshalJSON allows the object to be encoded in json.Marshal calls.
func (g Polygon) MarshalJSON() ([]byte, error) {
	return g.appendJSON(nil), nil
}

func (g Polygon) appendJSON(json []byte) []byte {
	return appendLevel3JSON(json, "Polygon", g.Coordinates, g.BBox, g.bboxDefined)
}

// JSON is the json representation of the object. This might not be exactly the same as the original.
func (g Polygon) JSON() string {
	return string(g.appendJSON(nil))
}

// String returns a string representation of the object. This might be JSON or something else.
func (g Polygon) String() string {
	return g.JSON()
}

func (g Polygon) bboxPtr() *BBox {
	return g.BBox
}
func (g Polygon) hasPositions() bool {
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
func (g Polygon) WithinBBox(bbox BBox) bool {
	if g.bboxDefined {
		return rectBBox(g.CalculatedBBox()).InsideRect(rectBBox(bbox))
	}
	if len(g.Coordinates) == 0 {
		return false
	}
	rbbox := rectBBox(bbox)
	ext, holes := polyExteriorHoles(g.Coordinates)
	if len(holes) > 0 {
		if rbbox.Max == rbbox.Min {
			return rbbox.Min.Inside(ext, holes)
		}
		return rbbox.Inside(ext, holes)
	}
	return ext.InsideRect(rectBBox(bbox))
}

// IntersectsBBox detects if the object intersects a bbox.
func (g Polygon) IntersectsBBox(bbox BBox) bool {
	if g.bboxDefined {
		return rectBBox(g.CalculatedBBox()).IntersectsRect(rectBBox(bbox))
	}
	if len(g.Coordinates) == 0 {
		return false
	}
	rbbox := rectBBox(bbox)
	ext, holes := polyExteriorHoles(g.Coordinates)
	if len(holes) > 0 {
		if rbbox.Max == rbbox.Min {
			return rbbox.Min.Intersects(ext, holes)
		}
		return rbbox.Intersects(ext, holes)
	}
	return ext.IntersectsRect(rectBBox(bbox))
}

// Within detects if the object is fully contained inside another object.
func (g Polygon) Within(o Object) bool {
	return withinObjectShared(g, o,
		func(v Polygon) bool {
			if len(g.Coordinates) == 0 {
				return false
			}
			return polyPositions(g.Coordinates[0]).Inside(polyExteriorHoles(v.Coordinates))
		},
	)
}

// Intersects detects if the object intersects another object.
func (g Polygon) Intersects(o Object) bool {
	return intersectsObjectShared(g, o,
		func(v Polygon) bool {
			if len(g.Coordinates) == 0 {
				return false
			}
			return polyPositions(g.Coordinates[0]).Intersects(polyExteriorHoles(v.Coordinates))
		},
	)
}

// Nearby detects if the object is nearby a position.
func (g Polygon) Nearby(center Position, meters float64) bool {
	return nearbyObjectShared(g, center.X, center.Y, meters)
}

// IsBBoxDefined returns true if the object has a defined bbox.
func (g Polygon) IsBBoxDefined() bool {
	return g.bboxDefined
}

// IsGeometry return true if the object is a geojson geometry object. false if it something else.
func (g Polygon) IsGeometry() bool {
	return true
}
