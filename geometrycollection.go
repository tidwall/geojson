package geojson

import "github.com/tidwall/gjson"

// GeometryCollection is a GeoJSON GeometryCollection
type GeometryCollection struct {
	Geometries []Object
	BBox       BBox
}

// Rect returns a rectangle that contains the entire object
func (g GeometryCollection) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	return calcRectFromObjects(g.Geometries)
}

// Center is the center-most point of the object
func (g GeometryCollection) Center() Position {
	return g.Rect().Center()
}

// AppendJSON appends a json representation to destination
func (g GeometryCollection) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"GeometryCollection","geometries":[`...)
	for i := 0; i < len(g.Geometries); i++ {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = g.Geometries[i].AppendJSON(dst)
	}
	dst = append(dst, ']')
	if g.BBox != nil && g.BBox.Defined() {
		dst = append(dst, `,"bbox":`...)
		dst = g.BBox.AppendJSON(dst)
	}
	dst = append(dst, '}')
	return dst
}

// loadJSONGeometryCollection will return a valid GeoJSON object.
func loadJSONGeometryCollection(data string) (Object, error) {
	var g GeometryCollection
	rgeometries := gjson.Get(data, "geometries")
	if !rgeometries.Exists() {
		return nil, errGeometriesMissing
	}
	if !rgeometries.IsArray() {
		return nil, errGeometriesInvalid
	}
	var err error
	rgeometries.ForEach(func(key, value gjson.Result) bool {
		var f Object
		f, err = Load(value.Raw)
		if err != nil {
			return false
		}
		g.Geometries = append(g.Geometries, f)
		return true
	})
	if err != nil {
		return nil, err
	}
	g.BBox, err = loadBBox(data)
	if err != nil {
		return nil, err
	}
	if g.BBox == nil {
		g.BBox = bboxRect{g.Rect()}
	}
	return g, nil
}
