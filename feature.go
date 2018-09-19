package geojson

import "github.com/tidwall/gjson"

type Feature struct {
	BBox       BBox
	Geometry   Object
	ID         gjson.Result
	Properties gjson.Result
}

func (g Feature) HasBBox() bool {
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
	if g.ID.Exists() {
		dst = append(dst, `,"id":`...)
		dst = append(dst, g.ID.Raw...)
	}
	if g.Properties.Exists() {
		dst = append(dst, `,"properties":`...)
		dst = append(dst, g.Properties.Raw...)
	}
	dst = append(dst, '}')
	return dst
}
func (g Feature) ForEach(iter func(child Object) bool) {
	iter(g.Geometry)
}
func (g Feature) Within(other Object) bool {
	panic("unsupported")
}
func (g Feature) Intersects(other Object) bool {
	panic("unsupported")
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
	g.ID = resultCopy(gjson.Get(data, "id"))
	g.Properties = resultCopy(gjson.Get(data, "properties"))
	if g.BBox == nil {
		g.BBox = bboxRect{g.Rect()}
	}
	return g, nil
}

func resultCopy(res gjson.Result) gjson.Result {
	if res.Exists() {
		if res.Type == gjson.String {
			res = gjson.Parse(string([]byte(res.Raw)))
		} else {
			res.Raw = string([]byte(res.Raw))
		}
	}
	return res
}
