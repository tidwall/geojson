package geojson

import (
	"github.com/tidwall/gjson"
)

type Point struct {
	Coordinates Position
	BBox        BBox
	Extra       *Extra
}

func (g Point) HasBBox() bool {
	return g.BBox != nil && g.BBox.Defined()
}

func (g Point) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	return Rect{Min: g.Coordinates, Max: g.Coordinates}
}
func (g Point) Center() Position {
	return g.Rect().Center()
}

func (g Point) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"Point","coordinates":`...)
	dst = appendJSONPosition(dst, g.Coordinates, g.Extra, 0)
	if g.BBox != nil && g.BBox.Defined() {
		dst = append(dst, `,"bbox":`...)
		dst = g.BBox.AppendJSON(dst)
	}
	dst = append(dst, '}')
	return dst
}
func (g Point) ForEach(func(child Object) bool) {}

func (g Point) Contains(other Object) bool {
	if g.HasBBox() {
		return g.Rect().Contains(other)
	}
	return g.Coordinates.Contains(other)
}
func (g Point) Intersects(other Object) bool {
	if g.HasBBox() {
		return g.Rect().Intersects(other)
	}
	return g.Coordinates.Intersects(other)
}

func loadJSONPoint(data string) (Object, error) {
	var g Point
	var err error
	g.Coordinates, g.Extra, err = loadJSONPointCoords(data, gjson.Result{})
	if err != nil {
		return nil, err
	}
	g.BBox, err = loadBBox(data)
	if err != nil {
		return nil, err
	}
	if g.BBox == nil && g.Extra == nil {
		return g.Coordinates, nil
	}
	return g, nil
}

func loadJSONPointCoords(data string, rcoords gjson.Result) (
	Position, *Extra, error,
) {
	var coords Position
	var ex *Extra
	if !rcoords.Exists() {
		rcoords = gjson.Get(data, "coordinates")
		if !rcoords.Exists() {
			return coords, nil, errCoordinatesMissing
		}
		if !rcoords.IsArray() {
			return coords, nil, errCoordinatesInvalid
		}
	}
	var err error
	var count int
	var nums [4]float64
	rcoords.ForEach(func(key, value gjson.Result) bool {
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
		return coords, nil, err
	}
	if count < 2 {
		return coords, nil, errCoordinatesInvalid
	}
	coords = Position{X: nums[0], Y: nums[1]}
	if count > 2 {
		ex = new(Extra)
		if count > 3 {
			ex.Dims = DimsZM
		} else {
			ex.Dims = DimsZ
		}
		ex.Positions = make([]float64, count-2)
		for i := 2; i < count; i++ {
			ex.Positions[i-2] = nums[i]
		}
	}
	return coords, ex, nil
}
