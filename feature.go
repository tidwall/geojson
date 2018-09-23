package geojson

import "github.com/tidwall/gjson"

type Feature struct {
	BBox       BBox
	Geometry   Object
	ID         string
	Properties string
}

func (g Feature) BBoxDefined() bool {
	return g.BBox != nil && g.BBox.Defined()
}

func (g Feature) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	return g.Geometry.Rect()
}
func (g Feature) Center() Position {
	return g.Rect().Center()
}

func (g Feature) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Feature","geometry":`...)
	dst = g.Geometry.AppendJSON(dst)
	if g.BBox != nil && g.BBox.Defined() {
		dst = append(dst, `,"bbox":`...)
		dst = g.BBox.AppendJSON(dst)
	}
	if g.ID != "" {
		dst = append(dst, `,"id":`...)
		dst = append(dst, g.ID...)
	}
	if g.Properties != "" {
		dst = append(dst, `,"properties":`...)
		dst = append(dst, g.Properties...)
	}
	dst = append(dst, '}')
	return dst
}

func (g Feature) ForEach(iter func(child Object) bool) {
	iter(g.Geometry)
}

func (g Feature) Contains(other Object) bool {
	return collectionContains(g, other, true)
}

func (g Feature) Intersects(other Object) bool {
	return collectionIntersects(g, other, true)
}

// loadJSONFeature will return a valid GeoJSON object.
func loadJSONFeature(data string) (Object, error) {
	var g Feature
	rgeometry := gjson.Get(data, "geometry")
	if !rgeometry.Exists() {
		return nil, errGeometryMissing
	}
	var err error
	g.Geometry, err = Load(rgeometry.Raw)
	if err != nil {
		return nil, err
	}
	g.BBox, err = loadBBox(data)
	if err != nil {
		return nil, err
	}
	id := gjson.Get(data, "id").Raw
	properties := gjson.Get(data, "properties").Raw
	if len(id) > 0 || len(properties) > 0 {
		combined := id + " " + properties
		g.ID = combined[:len(id)]
		g.Properties = combined[len(id)+1:]
	}
	if g.BBox == nil {
		switch geometry := g.Geometry.(type) {
		default:
			g.BBox = bboxRect{g.Rect()}
		case Position, Rect:
		case Point:
			g.BBox = geometry.BBox
		}
	}
	return g, nil
}

func resultCopy(res gjson.Result) string {
	if len(res.Raw) > 0 {
		return string([]byte(res.Raw))
	}
	return ""
}
