package geojson

import (
	"github.com/tidwall/tile38/pkg/geojson/geohash"
)

// MultiPolygon is a geojson object with the type "MultiPolygon"
type MultiPolygon struct {
	Coordinates [][][]Position
	BBox        *BBox
	bboxDefined bool
	polygons    []Polygon
}

func fillMultiPolygon(coordinates [][][]Position, bbox *BBox, err error) (MultiPolygon, error) {
	polygons := make([]Polygon, len(coordinates))
	if err == nil {
		for i, ps := range coordinates {
			polygons[i], err = fillPolygon(ps, nil, nil)
			if err != nil {
				break
			}
		}
	}
	bboxDefined := bbox != nil
	if !bboxDefined {
		cbbox := calculatedBBox(polygons, nil)
		bbox = &cbbox
	}
	return MultiPolygon{
		Coordinates: coordinates,
		BBox:        bbox,
		bboxDefined: bboxDefined,
		polygons:    polygons,
	}, err
}

func calculatedBBox(polygons []Polygon, bbox *BBox) BBox {
	if bbox != nil {
		return *bbox
	}
	var cbbox BBox
	for i, p := range polygons {
		if i == 0 {
			cbbox = p.CalculatedBBox()
		} else {
			cbbox = cbbox.union(p.CalculatedBBox())
		}
	}
	return cbbox
}

// CalculatedBBox is exterior bbox containing the object.
func (g MultiPolygon) CalculatedBBox() BBox {
	return calculatedBBox(g.polygons, g.BBox)
}

// CalculatedPoint is a point representation of the object.
func (g MultiPolygon) CalculatedPoint() Position {
	return g.CalculatedBBox().center()
}

// Geohash converts the object to a geohash value.
func (g MultiPolygon) Geohash(precision int) (string, error) {
	p := g.CalculatedPoint()
	return geohash.Encode(p.Y, p.X, precision)
}

// PositionCount return the number of coordinates.
func (g MultiPolygon) PositionCount() int {
	return level4PositionCount(g.Coordinates, g.BBox)
}

// Weight returns the in-memory size of the object.
func (g MultiPolygon) Weight() int {
	return level4Weight(g.Coordinates, g.BBox)
}

// MarshalJSON allows the object to be encoded in json.Marshal calls.
func (g MultiPolygon) MarshalJSON() ([]byte, error) {
	return g.appendJSON(nil), nil
}

func (g MultiPolygon) appendJSON(json []byte) []byte {
	return appendLevel4JSON(json, "MultiPolygon", g.Coordinates, g.BBox, g.bboxDefined)
}

// JSON is the json representation of the object. This might not be exactly the same as the original.
func (g MultiPolygon) JSON() string {
	return string(g.appendJSON(nil))
}

// String returns a string representation of the object. This might be JSON or something else.
func (g MultiPolygon) String() string {
	return g.JSON()
}

func (g MultiPolygon) bboxPtr() *BBox {
	return g.BBox
}
func (g MultiPolygon) hasPositions() bool {
	if g.bboxDefined {
		return true
	}
	for _, c := range g.Coordinates {
		for _, c := range c {
			if len(c) > 0 {
				return true
			}
		}
	}
	return false
}

func (g MultiPolygon) getPolygon(index int) Polygon {
	if index < len(g.polygons) {
		return g.polygons[index]
	}
	return Polygon{Coordinates: g.Coordinates[index]}
}

// WithinBBox detects if the object is fully contained inside a bbox.
func (g MultiPolygon) WithinBBox(bbox BBox) bool {
	if g.bboxDefined {
		return rectBBox(g.CalculatedBBox()).InsideRect(rectBBox(bbox))
	}
	if len(g.Coordinates) == 0 {
		return false
	}
	for i := range g.Coordinates {
		if !g.getPolygon(i).WithinBBox(bbox) {
			return false
		}
	}
	return true
}

// IntersectsBBox detects if the object intersects a bbox.
func (g MultiPolygon) IntersectsBBox(bbox BBox) bool {
	if g.bboxDefined {
		return rectBBox(g.CalculatedBBox()).IntersectsRect(rectBBox(bbox))
	}
	for i := range g.Coordinates {
		if g.getPolygon(i).IntersectsBBox(bbox) {
			return true
		}
	}
	return false
}

// Within detects if the object is fully contained inside another object.
func (g MultiPolygon) Within(o Object) bool {
	return withinObjectShared(g, o,
		func(v Polygon) bool {
			if len(g.Coordinates) == 0 {
				return false
			}
			if !v.Within(g) {
				return false
			}
			return true
		},
	)
}

// WithinCircle detects if the object is fully contained inside a circle.
func (g MultiPolygon) WithinCircle(center Position, meters float64) bool {
	if len(g.polygons) == 0 {
		return false
	}
	for _, polygon := range g.polygons {
		if !polygon.WithinCircle(center, meters) {
			return false
		}
	}
	return true
}

// Intersects detects if the object intersects another object.
func (g MultiPolygon) Intersects(o Object) bool {
	return intersectsObjectShared(g, o,
		func(v Polygon) bool {
			if len(g.Coordinates) == 0 {
				return false
			}
			if v.Intersects(g) {
				return true
			}
			return false
		},
	)
}

// IntersectsCircle detects if the object intersects a circle.
func (g MultiPolygon) IntersectsCircle(center Position, meters float64) bool {
	for _, polygon := range g.polygons {
		if polygon.IntersectsCircle(center, meters) {
			return true
		}
	}
	return false
}

// Nearby detects if the object is nearby a position.
func (g MultiPolygon) Nearby(center Position, meters float64) bool {
	return nearbyObjectShared(g, center.X, center.Y, meters)
}

// IsBBoxDefined returns true if the object has a defined bbox.
func (g MultiPolygon) IsBBoxDefined() bool {
	return g.bboxDefined
}

// IsGeometry return true if the object is a geojson geometry object. false if it something else.
func (g MultiPolygon) IsGeometry() bool {
	return true
}

// Clip returns the object obtained by clipping this object by a bbox.
func (g MultiPolygon) Clipped(bbox BBox) Object {
	var new_coordinates [][][]Position

	for _, polygon := range g.polygons {
		clippedPolygon, _ := polygon.Clipped(bbox).(Polygon)
		new_coordinates = append(new_coordinates, clippedPolygon.Coordinates)
	}

	res, _ := fillMultiPolygon(new_coordinates, nil, nil)
	return res
}
