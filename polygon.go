package geojson

import "github.com/tidwall/gjson"

type Polygon struct {
	Coordinates [][]Position
	BBox        BBox
	Extra       *Extra
}

func (g Polygon) HasBBox() bool {
	return g.BBox != nil && g.BBox.Defined()
}

func (g Polygon) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	var rect Rect
	if len(g.Coordinates) > 0 {
		for i := 0; i < len(g.Coordinates[0]); i++ {
			if i == 0 {
				rect.Min = g.Coordinates[0][i]
				rect.Max = g.Coordinates[0][i]
			} else {
				rect = rect.Expand(g.Coordinates[0][i])
			}
		}
	}
	return rect
}

func (g Polygon) Center() Position {
	return g.Rect().Center()
}

func (g Polygon) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Polygon","coordinates":[`...)
	var n int
	for i, p := range g.Coordinates {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = append(dst, '[')
		for i, p := range p {
			if i > 0 {
				dst = append(dst, ',')
			}
			dst = appendJSONPosition(dst, p, g.Extra, n)
			n++
		}
		dst = append(dst, ']')
	}
	dst = append(dst, ']')
	if g.BBox != nil && g.BBox.Defined() {
		dst = append(dst, `,"bbox":`...)
		dst = g.BBox.AppendJSON(dst)
	}
	dst = append(dst, '}')
	return dst
}
func (g Polygon) ForEach(func(child Object) bool) {}

func (g Polygon) Within(other Object) bool {
	panic("unsupported")
}
func (g Polygon) Intersects(other Object) bool {
	panic("unsupported")
}

func loadJSONPolygon(data string) (Object, error) {
	var g Polygon
	var err error
	g.Coordinates, g.Extra, err = loadJSONPolygonCoords(data, gjson.Result{})
	if err != nil {
		return nil, err
	}
	if len(g.Coordinates) == 0 {
		return nil, errMustBeALinearRing
	}
	for _, p := range g.Coordinates {
		if len(p) < 4 || p[0] != p[len(p)-1] {
			return nil, errMustBeALinearRing
		}
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

func loadJSONPolygonCoords(data string, rcoords gjson.Result) (
	[][]Position, *Extra, error,
) {
	var err error
	var coords [][]Position
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
		coords = append(coords, []Position{})
		ii := len(coords) - 1
		value.ForEach(func(key, value gjson.Result) bool {
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
			coords[ii] = append(coords[ii], Position{X: nums[0], Y: nums[1]})
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
		return err == nil
	})
	if err != nil {
		return nil, nil, err
	}
	return coords, ex, err
}
