package geojson

import "github.com/tidwall/gjson"

type MultiPolygon struct {
	Polygons []Polygon
	BBox     BBox
}

func (g MultiPolygon) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	var rect Rect
	for i := 0; i < len(g.Polygons); i++ {
		if i == 0 {
			rect = g.Polygons[i].Rect()
		} else {
			rect = rect.Union(g.Polygons[i].Rect())
		}
	}
	return rect
}

func (g MultiPolygon) Center() Position {
	return g.Rect().Center()
}

func (g MultiPolygon) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"MultiPolygon","coordinates":[`...)
	for i, g := range g.Polygons {
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

func loadJSONMultiPolygon(data string) (Object, error) {
	var g MultiPolygon
	var err error
	rcoords := gjson.Get(data, "coordinates")
	if !rcoords.Exists() {
		return nil, errCoordinatesMissing
	}
	if !rcoords.IsArray() {
		return nil, errCoordinatesInvalid
	}
	var coords [][]Position
	var ex *Extra
	rcoords.ForEach(func(_, value gjson.Result) bool {
		coords, ex, err = loadJSONPolygonCoords("", value)
		if err != nil {
			return false
		}
		g.Polygons = append(g.Polygons,
			Polygon{Coordinates: coords, Extra: ex},
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
