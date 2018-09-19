package geojson

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tidwall/gjson"
)

var (
	fmtErrTypeIsUnknown   = "type '%s' is unknown"
	errDataInvalid        = errors.New("invalid data")
	errTypeInvalid        = errors.New("type is invalid")
	errTypeMissing        = errors.New("type is missing")
	errCoordinatesInvalid = errors.New("coordinates is invalid")
	errCoordinatesMissing = errors.New("coordinates is missing")
	errGeometryInvalid    = errors.New("geometry is invalid")
	errGeometryMissing    = errors.New("geometry is missing")
	errFeaturesMissing    = errors.New("features is missing")
	errFeaturesInvalid    = errors.New("features is invalid")
	errGeometriesMissing  = errors.New("geometries is missing")
	errGeometriesInvalid  = errors.New("geometries is invalid")
	errBBoxInvalid        = errors.New("bbox is invalid")
)

// Parse an object. Only GeoJSON is currently supported.
func Parse(data string) (Object, error) {
	// look at the first byte
	for i := 0; ; i++ {
		if len(data) == 0 {
			return nil, errDataInvalid
		}
		switch data[0] {
		default:
			// well-known text is not supported yet
			return nil, errDataInvalid
		case 0, 1:
			if i > 0 {
				// 0x00 or 0x01 must be the first bytes
				return nil, errDataInvalid
			}
			// well-known binary is not supported yet
			return nil, errDataInvalid
		case '"':
			var s string
			err := json.Unmarshal([]byte(data), &s)
			if err != nil {
				return nil, err
			}
			return String(s), nil
		case ' ', '\t', '\n', '\r':
			// strip whitespace
			data = data[1:]
			continue
		case '{':
			return parseGeoJSON(data)
		}
	}
}

func parseGeoJSON(data string) (Object, error) {
	rtype := gjson.Get(data, "type")
	if !rtype.Exists() {
		return nil, errTypeMissing
	}
	if rtype.Type != gjson.String {
		return nil, errTypeInvalid
	}
	switch rtype.String() {
	default:
		return nil, fmt.Errorf(fmtErrTypeIsUnknown, rtype.String())
	case "Point":
		return parseGeoJSONPoint(data)
	case "LineString":
		return parseGeoJSONLineString(data)
	// case "Polygon":
	// 	return parseGeoJSONPolygon(data)
	case "Feature":
		return parseGeoJSONFeature(data)
	case "FeatureCollection":
		return parseGeoJSONFeatureCollection(data)
	case "GeometryCollection":
		return parseGeoJSONGeometryCollection(data)
	}
}

func parseCoords1(data string) (Position, *extra, error) {
	var coords Position
	var ex *extra
	rcoords := gjson.Get(data, "coordinates")
	if !rcoords.Exists() {
		return coords, nil, errCoordinatesMissing
	}
	if !rcoords.IsArray() {
		return coords, nil, errCoordinatesInvalid
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
		ex = new(extra)
		ex.z = true
		if count > 3 {
			ex.m = true
		}
		ex.coords = make([]float64, count-2)
		for i := 2; i < count; i++ {
			ex.coords[i-2] = nums[i]
		}
	}
	return coords, ex, nil
}

func parseCoords2(data string) ([]Position, *extra, error) {
	var err error
	var coords []Position
	var ex *extra
	var dims int
	rcoords := gjson.Get(data, "coordinates")
	if !rcoords.Exists() {
		return nil, nil, errCoordinatesMissing
	}
	if !rcoords.IsArray() {
		return nil, nil, errCoordinatesInvalid
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
				ex = new(extra)
				ex.z = true
				if count > 3 {
					ex.m = true
				}
				dims = ex.dims()
			}
		}
		if ex != nil {
			for i := 0; i < dims; i++ {
				ex.coords = append(ex.coords, nums[2+i])
			}
		}
		return true
	})
	if err != nil {
		return nil, nil, err
	}
	return coords, ex, err
}

func parseCoords3(data string) ([][]Position, *extra, error) {
	var err error
	var coords [][]Position
	var ex *extra
	var dims int
	value := gjson.Get(data, "coordinates")
	if !value.Exists() {
		return nil, nil, errCoordinatesMissing
	}
	if !value.IsArray() {
		return nil, nil, errCoordinatesInvalid
	}
	value.ForEach(func(key, value gjson.Result) bool {
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
					ex = new(extra)
					ex.z = true
					if count > 3 {
						ex.m = true
					}
					dims = ex.dims()
				}
			}
			if ex != nil {
				for i := 0; i < dims; i++ {
					ex.coords = append(ex.coords, nums[2+i])
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

func parseCoords4(data string) ([][][]Position, *extra, error) {
	var err error
	var coords [][][]Position
	var ex *extra
	var dims int
	value := gjson.Get(data, "coordinates")
	if !value.Exists() {
		return nil, nil, errCoordinatesMissing
	}
	if !value.IsArray() {
		return nil, nil, errCoordinatesInvalid
	}
	value.ForEach(func(key, value gjson.Result) bool {
		if !value.IsArray() {
			err = errCoordinatesInvalid
			return false
		}
		coords = append(coords, [][]Position{})
		ii := len(coords) - 1
		value.ForEach(func(key, value gjson.Result) bool {
			if !value.IsArray() {
				err = errCoordinatesInvalid
				return false
			}
			coords[ii] = append(coords[ii], []Position{})
			jj := len(coords[ii]) - 1
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
				coords[ii][jj] =
					append(coords[ii][jj], Position{X: nums[0], Y: nums[1]})
				if ex == nil {
					if count > 2 {
						ex = new(extra)
						ex.z = true
						if count > 3 {
							ex.m = true
						}
						dims = ex.dims()
					}
				}
				if ex != nil {
					for i := 0; i < dims; i++ {
						ex.coords = append(ex.coords, nums[2+i])
					}
				}
				return true
			})
			return err == nil
		})
		return err == nil
	})
	if err != nil {
		return nil, nil, err
	}
	return coords, ex, err
}
