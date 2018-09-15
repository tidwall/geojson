package geojson

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/tile38/pkg/geojson/geohash"
)

// FeatureCollection is a geojson object with the type "FeatureCollection"
type FeatureCollection struct {
	Features    []Object
	BBox        *BBox
	bboxDefined bool
}

func fillFeatureCollectionMap(json string) (FeatureCollection, error) {
	var g FeatureCollection
	res := gjson.Get(json, "features")
	switch res.Type {
	default:
		return g, errInvalidFeaturesMember
	case gjson.Null:
		return g, errFeaturesMemberRequired
	case gjson.JSON:
		if !resIsArray(res) {
			return g, errInvalidFeaturesMember
		}
		v := res.Array()
		g.Features = make([]Object, len(v))
		for i, res := range v {
			if res.Type != gjson.JSON {
				return g, errInvalidFeature
			}
			o, err := objectMap(res.Raw, fcoll)
			if err != nil {
				return g, err
			}
			g.Features[i] = o
		}
	}
	var err error
	g.BBox, err = fillBBox(json)
	if err != nil {
		return g, err
	}
	g.bboxDefined = g.BBox != nil
	if !g.bboxDefined {
		cbbox := g.CalculatedBBox()
		g.BBox = &cbbox
	}
	return g, err
}

// Geohash converts the object to a geohash value.
func (g FeatureCollection) Geohash(precision int) (string, error) {
	p := g.CalculatedPoint()
	return geohash.Encode(p.Y, p.X, precision)
}

// CalculatedPoint is a point representation of the object.
func (g FeatureCollection) CalculatedPoint() Position {
	return g.CalculatedBBox().center()
}

// CalculatedBBox is exterior bbox containing the object.
func (g FeatureCollection) CalculatedBBox() BBox {
	if g.BBox != nil {
		return *g.BBox
	}
	var bbox BBox
	for i, g := range g.Features {
		if i == 0 {
			bbox = g.CalculatedBBox()
		} else {
			bbox = bbox.union(g.CalculatedBBox())
		}
	}
	return bbox
}

// PositionCount return the number of coordinates.
func (g FeatureCollection) PositionCount() int {
	var res int
	for _, g := range g.Features {
		res += g.PositionCount()
	}
	if g.BBox != nil {
		return 2 + res
	}
	return res
}

// Weight returns the in-memory size of the object.
func (g FeatureCollection) Weight() int {
	var res int
	for _, g := range g.Features {
		res += g.Weight()
	}
	return res
}

// MarshalJSON allows the object to be encoded in json.Marshal calls.
func (g FeatureCollection) MarshalJSON() ([]byte, error) {
	return g.appendJSON(nil), nil
}

func (g FeatureCollection) appendJSON(json []byte) []byte {
	json = append(json, `{"type":"FeatureCollection","features":[`...)
	for i, g := range g.Features {
		if i != 0 {
			json = append(json, ',')
		}
		json = append(json, g.JSON()...)
	}
	json = append(json, ']')
	if g.bboxDefined {
		json = appendBBoxJSON(json, g.BBox)
	}
	return append(json, '}')
}

// JSON is the json representation of the object. This might not be exactly the same as the original.
func (g FeatureCollection) JSON() string {
	return string(g.appendJSON(nil))
}

// String returns a string representation of the object. This might be JSON or something else.
func (g FeatureCollection) String() string {
	return g.JSON()
}

// Bytes is the bytes representation of the object.
func (g FeatureCollection) Bytes() []byte {
	return []byte(g.JSON())
}
func (g FeatureCollection) bboxPtr() *BBox {
	return g.BBox
}
func (g FeatureCollection) hasPositions() bool {
	if g.BBox != nil {
		return true
	}
	for _, g := range g.Features {
		if g.hasPositions() {
			return true
		}
	}
	return false
}

// WithinBBox detects if the object is fully contained inside a bbox.
func (g FeatureCollection) WithinBBox(bbox BBox) bool {
	if g.bboxDefined {
		return rectBBox(g.CalculatedBBox()).InsideRect(rectBBox(bbox))
	}
	if len(g.Features) == 0 {
		return false
	}
	for _, g := range g.Features {
		if !g.WithinBBox(bbox) {
			return false
		}
	}
	return true
}

// IntersectsBBox detects if the object intersects a bbox.
func (g FeatureCollection) IntersectsBBox(bbox BBox) bool {
	if g.bboxDefined {
		return rectBBox(g.CalculatedBBox()).IntersectsRect(rectBBox(bbox))
	}
	for _, g := range g.Features {
		if g.IntersectsBBox(bbox) {
			return true
		}
	}
	return false
}

// Within detects if the object is fully contained inside another object.
func (g FeatureCollection) Within(o Object) bool {
	return withinObjectShared(g, o,
		func(v Polygon) bool {
			if len(g.Features) == 0 {
				return false
			}
			for _, f := range g.Features {
				if !f.Within(o) {
					return false
				}
			}
			return true
		},
	)
}

// WithinCircle detects if the object is fully contained inside a circle.
func (g FeatureCollection) WithinCircle(center Position, meters float64) bool {
	if len(g.Features) == 0 {
		return false
	}
	for _, feature := range g.Features {
		if !feature.WithinCircle(center, meters) {
			return false
		}
	}
	return true
}

// Intersects detects if the object intersects another object.
func (g FeatureCollection) Intersects(o Object) bool {
	return intersectsObjectShared(g, o,
		func(v Polygon) bool {
			if len(g.Features) == 0 {
				return false
			}
			for _, f := range g.Features {

				if f.Intersects(o) {
					return true
				}
			}
			return false
		},
	)
}

// IntersectsCircle detects if the object intersects a circle.
func (g FeatureCollection) IntersectsCircle(center Position, meters float64) bool {
	for _, feature := range g.Features {
		if feature.IntersectsCircle(center, meters) {
			return true
		}
	}
	return false
}

// Nearby detects if the object is nearby a position.
func (g FeatureCollection) Nearby(center Position, meters float64) bool {
	return nearbyObjectShared(g, center.X, center.Y, meters)
}

// IsBBoxDefined returns true if the object has a defined bbox.
func (g FeatureCollection) IsBBoxDefined() bool {
	return g.bboxDefined
}

// IsGeometry return true if the object is a geojson geometry object. false if it something else.
func (g FeatureCollection) IsGeometry() bool {
	return true
}

// Clip returns the object obtained by clipping this object by a bbox.
func (g FeatureCollection) Clipped(bbox BBox) Object {
	var new_features []Object
	for _, feature := range g.Features {
		new_features = append(new_features, feature.Clipped(bbox))
	}

	fc := FeatureCollection{Features: new_features}
	cbbox := fc.CalculatedBBox()
	fc.BBox = &cbbox
	return fc
}
