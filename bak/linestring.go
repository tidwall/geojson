package geojson

import "github.com/tidwall/gjson"

// LineString GeoJSON type
type LineString struct {
	Coordinates []Position
	BBox        BBox
	Extra       *Extra
}

// BBoxDefined return true if there is a defined GeoJSON "bbox" member
func (g LineString) BBoxDefined() bool {
	return g.BBox != nil && g.BBox.Defined()
}

// Rect returns the outer minimum bounding rectangle
func (g LineString) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	var rect Rect
	for i := 0; i < len(g.Coordinates); i++ {
		if i == 0 {
			rect.Min = g.Coordinates[i]
			rect.Max = g.Coordinates[i]
		} else {
			rect = rect.Expand(g.Coordinates[i])
		}
	}
	return rect
}

// Center returns the center position of the object
func (g LineString) Center() Position {
	return g.Rect().Center()
}

// AppendJSON appends the GeoJSON reprensentation to dst
func (g LineString) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"LineString","coordinates":[`...)
	for i, p := range g.Coordinates {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = appendJSONPosition(dst, p, g.Extra, i)
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
func (g LineString) ForEachChild(func(child Object) bool) {}

// Within is the inverse of contains
func (g LineString) Within(other Object) bool {
	return other.Contains(g)
}

// Contains returns true if object contains other object
func (g LineString) Contains(other Object) bool {
	return objectContains(g, other)
}

// Intersects returns true if object intersects with other object
func (g LineString) Intersects(other Object) bool {
	return objectIntersects(g, other)
}

func (g LineString) primativeContains(other Object) bool {
	pline := polyLine(g.Coordinates)
	switch other := other.(type) {
	case Position:
		return polyPoint(other).InsideLine(pline)
	case Rect:
		return polyRect(other).InsideLine(pline)
	case Point:
		return polyPoint(other.Coordinates).InsideLine(pline)
	case LineString:
		return polyLine(other.Coordinates).InsideLine(pline)
	case Polygon:
		return polyPolygon(other.Coordinates).InsideLine(pline)
	}
	return false
}
func (g LineString) primativeIntersects(other Object) bool {
	pline := polyLine(g.Coordinates)
	switch other := other.(type) {
	case Position:
		return pline.IntersectsPoint(polyPoint(other))
	case Rect:
		return pline.IntersectsRect(polyRect(other))
	case Point:
		return pline.IntersectsPoint(polyPoint(other.Coordinates))
	case LineString:
		return pline.IntersectsLine(polyLine(other.Coordinates))
	case Polygon:
		return pline.IntersectsPolygon(polyPolygon(other.Coordinates))
	}
	return false
}

func loadJSONLineString(data string) (Object, error) {
	var g LineString
	var err error
	g.Coordinates, g.Extra, err = loadJSONLineStringCoords(data, gjson.Result{})
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

func loadJSONLineStringCoords(data string, rcoords gjson.Result) (
	[]Position, *Extra, error,
) {
	var err error
	var coords []Position
	var ex *Extra
	var dims int
	if !rcoords.Exists() {
		rcoords = gjson.Get(data, "coordinates")
		if !rcoords.Exists() {
			return nil, nil, errCoordinatesMissing
		}
		if !rcoords.IsArray() {
			return nil, nil, errCoordinatesInvalid
		}
	}
	rcoords.ForEach(func(key, value gjson.Result) bool {
		if !value.IsArray() {
			err = errCoordinatesInvalid
			return false
		}
		var count int
		var nums [4]float64
		value.ForEach(func(key, value gjson.Result) bool {
			if count == 4 {
				return false
			}
			if value.Type != gjson.Number {
				err = errCoordinatesInvalid
				return false
			}
			nums[count] = value.Float()
			count++
			return true
		})
		if err != nil {
			return false
		}
		if count < 2 {
			err = errCoordinatesInvalid
			return false
		}
		coords = append(coords, Position{X: nums[0], Y: nums[1]})
		if ex == nil {
			if count > 2 {
				ex = new(Extra)
				if count > 3 {
					ex.Dims = DimsZM
				} else {
					ex.Dims = DimsZ
				}
				dims = int(ex.Dims)
			}
		}
		if ex != nil {
			for i := 0; i < dims; i++ {
				ex.Positions = append(ex.Positions, nums[2+i])
			}
		}
		return true
	})
	if err != nil {
		return nil, nil, err
	}
	return coords, ex, err
}
