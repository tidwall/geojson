package geojson

import "github.com/tidwall/gjson"

type MultiLineString struct {
	LineStrings []LineString
	BBox        BBox
}

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

func (g MultiLineString) Center() Position {
	return g.Rect().Center()
}

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
