package geojson

import "github.com/tidwall/gjson"

// MultiLineString GeoJSON type
type MultiLineString struct {
	LineStrings []LineString
	BBox        BBox
}

// BBoxDefined return true if there is a defined GeoJSON "bbox" member
func (g MultiLineString) BBoxDefined() bool {
	return g.BBox != nil && g.BBox.Defined()
}

// Rect returns the outer minimum bounding rectangle
func (g MultiLineString) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	var rect Rect
	for i := 0; i < len(g.LineStrings); i++ {
		if i == 0 {
			rect = g.LineStrings[i].Rect()
		} else {
			rect = rect.Union(g.LineStrings[i].Rect())
		}
	}
	return rect
}

// Center returns the center position of the object
func (g MultiLineString) Center() Position {
	return g.Rect().Center()
}

// AppendJSON appends the GeoJSON reprensentation to dst
func (g MultiLineString) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"MultiLineString","coordinates":[`...)
	for i, g := range g.LineStrings {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = append(dst,
			gjson.GetBytes(g.AppendJSON(nil), "coordinates").String()...)
	}
	dst = append(dst, ']')
	if g.BBox != nil && g.BBox.Defined() {
		dst = append(dst, `,"bbox":`...)
		dst = g.BBox.AppendJSON(dst)
	}
	dst = append(dst, '}')
	return dst
}

// ForEachChild iterates over child objects.
func (g MultiLineString) ForEachChild(iter func(child Object) bool) {
	for _, child := range g.LineStrings {
		if !iter(child) {
			return
		}
	}
}

// Contains returns true if object contains other object
func (g MultiLineString) Contains(other Object) bool {
	return collectionContains(g, other)
}

// Intersects returns true if object intersects with other object
func (g MultiLineString) Intersects(other Object) bool {
	return collectionIntersects(g, other)
}

func loadJSONMultiLineString(data string) (Object, error) {
	var g MultiLineString
	var err error
	rcoords := gjson.Get(data, "coordinates")
	if !rcoords.Exists() {
		return nil, errCoordinatesMissing
	}
	if !rcoords.IsArray() {
		return nil, errCoordinatesInvalid
	}
	var coords []Position
	var ex *Extra
	rcoords.ForEach(func(_, value gjson.Result) bool {
		coords, ex, err = loadJSONLineStringCoords("", value)
		if err != nil {
			return false
		}
		g.LineStrings = append(g.LineStrings,
			LineString{Coordinates: coords, Extra: ex},
		)
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
