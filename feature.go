package geojson

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/tidwall/tile38/geojson/geohash"
)

// Feature is a geojson object with the type "Feature"
type Feature struct {
	Geometry   Object
	BBox       *BBox
	ID         interface{}
	Properties map[string]interface{}
}

func fillFeatureMap(m map[string]interface{}) (Feature, []byte, error) {
	var g Feature
	switch v := m["geometry"].(type) {
	default:
		return g, nil, errInvalidGeometryMember
	case nil:
		return g, nil, errGeometryMemberRequired
	case map[string]interface{}:
		var err error
		g.Geometry, err = objectMap(v, feat)
		if err != nil {
			return g, nil, err
		}
	}
	var err error
	g.BBox, err = fillBBox(m)
	if err != nil {
		return g, nil, err
	}
	switch v := m["properties"].(type) {
	default:
		return g, nil, errInvalidPropertiesMember
	case nil:
	case map[string]interface{}:
		g.Properties = v
	}
	g.ID = m["id"]
	return g, nil, err
}

func fillFeatureBytes(b []byte, bbox *BBox, isCordZ bool) (Feature, []byte, error) {
	var err error
	var g Feature
	g.BBox = bbox
	if len(b) < 4 {
		return g, nil, errNotEnoughData
	}
	l := int(binary.LittleEndian.Uint32(b))
	b = b[4:]
	if l > 0 {
		if len(b) < l {
			return g, nil, errNotEnoughData
		}
		var arr []interface{}
		err = json.Unmarshal(b[:l], &arr)
		if err != nil {
			return g, b, err
		}
		b = b[l:]

		if len(arr) > 0 {
			switch v := arr[0].(type) {
			default:
				return g, b, errInvalidData
			case nil:
			case map[string]interface{}:
				g.Properties = v
			}
		}
		if len(arr) > 1 {
			g.ID = arr[1]
		}
	}
	g.Geometry, b, err = objectBytes(b)
	return g, b, err
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
	if g.Properties != nil {
		b, _ := json.Marshal(g.Properties)
		res += len(b)
	}
	if g.ID != nil {
		b, _ := json.Marshal(g.ID)
		res += len(b) - 2
	}
	return res
}

// MarshalJSON allows the object to be encoded in json.Marshal calls.
func (g Feature) MarshalJSON() ([]byte, error) {
	return []byte(g.JSON()), nil
}

// JSON is the json representation of the object. This might not be exactly the same as the original.
func (g Feature) JSON() string {
	var buf bytes.Buffer
	buf.WriteString(`{"type":"Feature","geometry":`)
	buf.WriteString(g.Geometry.JSON())
	g.BBox.write(&buf)
	if g.Properties != nil {
		buf.WriteString(`,"properties":`)
		b, _ := json.Marshal(g.Properties)
		buf.Write(b)
	}
	if g.ID != nil {
		buf.WriteString(`,"id":`)
		b, _ := json.Marshal(g.ID)
		buf.Write(b)
	}
	buf.WriteByte('}')
	return buf.String()
}

// Bytes is the bytes representation of the object.
func (g Feature) Bytes() []byte {
	var buf bytes.Buffer
	isCordZ := g.BBox.isCordZDefined()
	writeHeader(&buf, feature, g.BBox, isCordZ)
	var pb []byte
	if g.Properties != nil || g.ID != nil {
		arr := []interface{}{g.Properties, g.ID}
		pb, _ = json.Marshal(arr)
	}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(len(pb)))
	buf.Write(b)
	buf.Write(pb)
	buf.Write(g.Geometry.Bytes())
	return buf.Bytes()
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
