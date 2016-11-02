package geojson

import (
	"bytes"

	"github.com/tidwall/gjson"
	"github.com/tidwall/tile38/geojson/geohash"
)

// Feature is a geojson object with the type "Feature"
type Feature struct {
	Geometry Object
	BBox     *BBox
	idprops  string // raw id and properties seperated by a '\0'
}

func fillFeatureMap(json string) (Feature, []byte, error) {
	var g Feature
	v := gjson.Get(json, "geometry")
	switch v.Type {
	default:
		return g, nil, errInvalidGeometryMember
	case gjson.Null:
		return g, nil, errGeometryMemberRequired
	case gjson.JSON:
		var err error
		g.Geometry, err = objectMap(v.Raw, feat)
		if err != nil {
			return g, nil, err
		}
	}
	var err error
	g.BBox, err = fillBBox(json)
	if err != nil {
		return g, nil, err
	}

	var propsExists bool
	props := gjson.Get(json, "properties")
	switch props.Type {
	default:
		return g, nil, errInvalidPropertiesMember
	case gjson.Null:
	case gjson.JSON:
		propsExists = true
	}
	id := gjson.Get(json, "id")
	if id.Exists() || propsExists {
		raw := make([]byte, len(id.Raw)+len(props.Raw)+1)
		copy(raw, id.Raw)
		copy(raw[len(id.Raw)+1:], props.Raw)
		g.idprops = string(raw)
	}
	return g, nil, err
}

// Geohash converts the object to a geohash value.
func (g Feature) Geohash(precision int) (string, error) {
	p := g.CalculatedPoint()
	return geohash.Encode(p.Y, p.X, precision)
}

// CalculatedPoint is a point representation of the object.
func (g Feature) CalculatedPoint() Position {
	return g.CalculatedBBox().center()
}

// CalculatedBBox is exterior bbox containing the object.
func (g Feature) CalculatedBBox() BBox {
	if g.BBox != nil {
		return *g.BBox
	}
	return g.Geometry.CalculatedBBox()
}

// PositionCount return the number of coordinates.
func (g Feature) PositionCount() int {
	res := g.Geometry.PositionCount()
	if g.BBox != nil {
		return 2 + res
	}
	return res
}

// Weight returns the in-memory size of the object.
func (g Feature) Weight() int {
	res := g.PositionCount() * sizeofPosition
	res += len(g.idprops)
	return res
}

// MarshalJSON allows the object to be encoded in json.Marshal calls.
func (g Feature) MarshalJSON() ([]byte, error) {
	return []byte(g.JSON()), nil
}

func (g Feature) getRaw() (id, props string) {
	for i := 0; i < len(g.idprops); i++ {
		if g.idprops[i] == 0 {
			return g.idprops[:i], g.idprops[i+1:]
		}
	}
	return "", ""
}

// JSON is the json representation of the object. This might not be exactly the same as the original.
func (g Feature) JSON() string {
	var buf bytes.Buffer
	buf.WriteString(`{"type":"Feature","geometry":`)
	buf.WriteString(g.Geometry.JSON())
	g.BBox.write(&buf)
	idRaw, propsRaw := g.getRaw()
	if propsRaw != "" {
		buf.WriteString(`,"properties":`)
		buf.WriteString(propsRaw)
	}
	if idRaw != "" {
		buf.WriteString(`,"id":`)
		buf.WriteString(idRaw)
	}
	buf.WriteByte('}')
	return buf.String()
}

// String returns a string representation of the object. This might be JSON or something else.
func (g Feature) String() string {
	return g.JSON()
}

// Bytes is the bytes representation of the object.
func (g Feature) Bytes() []byte {
	return []byte(g.JSON())
}
func (g Feature) bboxPtr() *BBox {
	return g.BBox
}
func (g Feature) hasPositions() bool {
	if g.BBox != nil {
		return true
	}
	return g.Geometry.hasPositions()
}

// WithinBBox detects if the object is fully contained inside a bbox.
func (g Feature) WithinBBox(bbox BBox) bool {
	if g.BBox != nil {
		return rectBBox(g.CalculatedBBox()).InsideRect(rectBBox(bbox))
	}
	return g.Geometry.WithinBBox(bbox)
}

// IntersectsBBox detects if the object intersects a bbox.
func (g Feature) IntersectsBBox(bbox BBox) bool {
	if g.BBox != nil {
		return rectBBox(g.CalculatedBBox()).IntersectsRect(rectBBox(bbox))
	}
	return g.Geometry.IntersectsBBox(bbox)
}

// Within detects if the object is fully contained inside another object.
func (g Feature) Within(o Object) bool {
	return withinObjectShared(g, o,
		func(v Polygon) bool {
			return g.Geometry.Within(o)
		},
		func(v MultiPolygon) bool {
			return g.Geometry.Within(o)
		},
	)
}

// Intersects detects if the object intersects another object.
func (g Feature) Intersects(o Object) bool {
	return intersectsObjectShared(g, o,
		func(v Polygon) bool {
			return g.Geometry.Intersects(o)
		},
		func(v MultiPolygon) bool {
			return g.Geometry.Intersects(o)
		},
	)
}

// Nearby detects if the object is nearby a position.
func (g Feature) Nearby(center Position, meters float64) bool {
	return nearbyObjectShared(g, center.X, center.Y, meters)
}

// IsBBoxDefined returns true if the object has a defined bbox.
func (g Feature) IsBBoxDefined() bool {
	return g.BBox != nil
}

// IsGeometry return true if the object is a geojson geometry object. false if it something else.
func (g Feature) IsGeometry() bool {
	return true
}
