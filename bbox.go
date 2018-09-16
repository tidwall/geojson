package geojson

import (
	"math"
	"strconv"

	"github.com/tidwall/geojson/poly"
	"github.com/tidwall/gjson"
)

// BBox is a bounding box
type BBox struct {
	Min Position
	Max Position
}

// New2DBBox creates a new bounding box
func New2DBBox(minX, minY, maxX, maxY float64) BBox {
	return BBox{
		Min: Position{X: minX, Y: minY},
		Max: Position{X: maxX, Y: maxY},
	}
}

func fillBBox(json string) (*BBox, error) {
	var bbox *BBox
	res := gjson.Get(json, "bbox")
	switch res.Type {
	default:
		return nil, errBBoxInvalidType
	case gjson.Null:
	case gjson.JSON:
		v := res.Array()
		if !(len(v) == 4 || len(v) == 6) {
			return nil, errBBoxInvalidNumberOfValues
		}
		bbox = &BBox{}
		for i := 0; i < len(v); i++ {
			if v[i].Type != gjson.Number {
				return nil, errBBoxInvalidValue
			}
		}
		bbox.Min.X = v[0].Float()
		bbox.Min.Y = v[1].Float()
		i := 2
		if len(v) == 6 {
			bbox.Min.Z = v[2].Float()
			i = 3
		} else {
			bbox.Min.Z = nilz
		}
		bbox.Max.X = v[i+0].Float()
		bbox.Max.Y = v[i+1].Float()
		if len(v) == 6 {
			bbox.Max.Z = v[i+2].Float()
			i = 3
		} else {
			bbox.Max.Z = nilz
		}
	}
	return bbox, nil
}

func (b *BBox) isCordZDefined() bool {
	return b != nil && (b.Min.Z != nilz || b.Max.Z != nilz)
}

func appendBBoxJSON(json []byte, b *BBox) []byte {
	if b == nil {
		return json
	}
	hasZ := b.Min.Z != nilz && b.Max.Z != nilz
	json = append(json, `,"bbox":[`...)
	json = strconv.AppendFloat(json, b.Min.X, 'f', -1, 64)
	json = append(json, ',')
	json = strconv.AppendFloat(json, b.Min.Y, 'f', -1, 64)
	if hasZ {
		json = append(json, ',')
		json = strconv.AppendFloat(json, b.Min.Z, 'f', -1, 64)
	}
	json = append(json, ',')
	json = strconv.AppendFloat(json, b.Max.X, 'f', -1, 64)
	json = append(json, ',')
	json = strconv.AppendFloat(json, b.Max.Y, 'f', -1, 64)
	if hasZ {
		json = append(json, ',')
		json = strconv.AppendFloat(json, b.Max.Z, 'f', -1, 64)
	}
	json = append(json, ']')
	return json
}

func (b BBox) center() Position {
	return Position{
		(b.Max.X-b.Min.X)/2 + b.Min.X,
		(b.Max.Y-b.Min.Y)/2 + b.Min.Y,
		0,
	}
}

func (b BBox) union(bbox BBox) BBox {
	if bbox.Min.X < b.Min.X {
		b.Min.X = bbox.Min.X
	}
	if bbox.Min.Y < b.Min.Y {
		b.Min.Y = bbox.Min.Y
	}
	if bbox.Max.X > b.Max.X {
		b.Max.X = bbox.Max.X
	}
	if bbox.Max.Y > b.Max.Y {
		b.Max.Y = bbox.Max.Y
	}
	return b
}

func (b BBox) exterior() []Position {
	return []Position{
		{b.Min.X, b.Min.Y, 0},
		{b.Min.X, b.Max.Y, 0},
		{b.Max.X, b.Max.Y, 0},
		{b.Max.X, b.Min.Y, 0},
		{b.Min.X, b.Min.Y, 0},
	}
}

func rectBBox(bbox BBox) poly.Rect {
	return poly.Rect{
		Min: poly.Point{X: bbox.Min.X, Y: bbox.Min.Y, Z: 0},
		Max: poly.Point{X: bbox.Max.X, Y: bbox.Max.Y, Z: 0},
	}
}

// ExternalJSON is the simple json representation of the bounding box used for
// external applications.
func (b BBox) ExternalJSON() string {
	sw, ne := b.Min, b.Max
	sw.Z, ne.Z = 0, 0
	return `{"sw":` + sw.ExternalJSON() + `,"ne":` + ne.ExternalJSON() + `}`
}

// Sparse returns back an evenly distributed number of sub bboxs.
func (b BBox) Sparse(amount byte) []BBox {
	if amount == 0 {
		return []BBox{b}
	}
	var bboxes []BBox
	split := 1 << amount
	var xsize, ysize float64
	if b.Max.X < b.Min.X {
		// crosses the prime meridian
		xsize = (b.Min.X - b.Max.X) / float64(split)
	} else {
		xsize = (b.Max.X - b.Min.X) / float64(split)
	}
	if b.Max.Y < b.Min.Y {
		// crosses the equator
		ysize = (b.Min.Y - b.Max.Y) / float64(split)
	} else {
		ysize = (b.Max.Y - b.Min.Y) / float64(split)
	}

	for y := b.Min.Y; y < b.Max.Y; y += ysize {
		for x := b.Min.X; x < b.Max.X; x += xsize {
			bboxes = append(bboxes, BBox{
				Min: Position{X: x, Y: y, Z: b.Min.Z},
				Max: Position{X: x + xsize, Y: y + ysize, Z: b.Max.Z},
			})
		}
	}
	return bboxes
}

// BBoxesFromCenter calculates the bounding box surrounding a circle.
func BBoxesFromCenter(lat, lon, meters float64) (outer BBox) {

	outer.Min.Y, outer.Min.X, outer.Max.Y, outer.Max.X =
		BoundsFromCenter(lat, lon, meters)
	if outer.Min.X == outer.Max.X {
		switch outer.Min.X {
		case -180:
			outer.Max.X = 180
		case 180:
			outer.Min.X = -180
		}
	}

	return outer
}

// BoundsFromCenter calculates the bounding box surrounding a circle.
func BoundsFromCenter(lat, lon, meters float64) (
	latMin, lonMin, latMax, lonMax float64,
) {

	// see http://janmatuschek.de/LatitudeLongitudeBoundingCoordinates#Latitude
	lat = lat * radians
	lon = lon * radians

	r := meters / earthRadius // angular radius

	latMin = lat - r
	latMax = lat + r

	latT := math.Asin(math.Sin(lat) / math.Cos(r))
	lonΔ := math.Acos((math.Cos(r) - math.Sin(latT)*math.Sin(lat)) /
		(math.Cos(latT) * math.Cos(lat)))

	lonMin = lon - lonΔ
	lonMax = lon + lonΔ

	// Adjust for north poll
	if latMax > math.Pi/2 {
		lonMin = -math.Pi
		latMax = math.Pi / 2
		lonMax = math.Pi
	}

	// Adjust for south poll
	if latMin < -math.Pi/2 {
		latMin = -math.Pi / 2
		lonMin = -math.Pi
		lonMax = math.Pi
	}

	// Adjust for wraparound. Remove this if the commented-out condition below
	// this block is added.
	if lonMin < -math.Pi || lonMax > math.Pi {
		lonMin = -math.Pi
		lonMax = math.Pi
	}

	// normalise to -180..+180°
	lonMin = math.Mod(lonMin+3*math.Pi, 2*math.Pi) - math.Pi
	lonMax = math.Mod(lonMax+3*math.Pi, 2*math.Pi) - math.Pi

	return latMin * degrees, lonMin * degrees,
		latMax * degrees, lonMax * degrees
}
