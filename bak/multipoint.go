package geojson

import "github.com/tidwall/gjson"

// MultiPoint GeoJSON type
type MultiPoint struct {
	Points []Point
	BBox   BBox
}

// BBoxDefined return true if there is a defined GeoJSON "bbox" member
func (g MultiPoint) BBoxDefined() bool {
	return g.BBox != nil && g.BBox.Defined()
}

// Rect returns the outer minimum bounding rectangle
func (g MultiPoint) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	var rect Rect
	for i := 0; i < len(g.Points); i++ {
		if i == 0 {
			rect = g.Points[i].Rect()
		} else {
			rect = rect.Union(g.Points[i].Rect())
		}
	}
	return rect
}

// Center returns the center position of the object
func (g MultiPoint) Center() Position {
	return g.Rect().Center()
}

// AppendJSON appends the GeoJSON reprensentation to dst
func (g MultiPoint) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"MultiPoint","coordinates":[`...)
	for i, g := range g.Points {
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
func (g MultiPoint) ForEachChild(iter func(child Object) bool) {
	for _, child := range g.Points {
		if !iter(child) {
			return
		}
	}
}

// Within is the inverse of contains
func (g MultiPoint) Within(other Object) bool {
	return other.Contains(g)
}

// Contains returns true if object contains other object
func (g MultiPoint) Contains(other Object) bool {
	return collectionContains(g, other)
}

// Intersects returns true if object intersects with other object
func (g MultiPoint) Intersects(other Object) bool {
	return collectionIntersects(g, other)
}

func loadJSONMultiPoint(data string) (Object, error) {
	var g MultiPoint
	var err error
	rcoords := gjson.Get(data, "coordinates")
	if !rcoords.Exists() {
		return nil, errCoordinatesMissing
	}
	if !rcoords.IsArray() {
		return nil, errCoordinatesInvalid
	}
	var coords Position
	var ex *Extra
	rcoords.ForEach(func(_, value gjson.Result) bool {
		coords, ex, err = loadJSONPointCoords("", value)
		if err != nil {
			return false
		}
		g.Points = append(g.Points, Point{Coordinates: coords, Extra: ex})
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
