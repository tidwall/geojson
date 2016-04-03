package geojson

import (
	"bytes"
	"strconv"

	"github.com/tidwall/tile38/geojson/geo"
	"github.com/tidwall/tile38/geojson/poly"
)

// BBox is a bounding box
type BBox struct {
	Min Position
	Max Position
}

// New2DBBox creates a new bounding box
func New2DBBox(minX, minY, maxX, maxY float64) BBox {
	return BBox{Min: Position{X: minX, Y: minY, Z: 0}, Max: Position{X: maxX, Y: maxY, Z: 0}}
}

func fillBBox(m map[string]interface{}) (*BBox, error) {
	var bbox *BBox
	var ok bool
	switch v := m["bbox"].(type) {
	default:
		return nil, errBBoxInvalidType
	case nil:
	case []interface{}:
		if !(len(v) == 4 || len(v) == 6) {
			return nil, errBBoxInvalidNumberOfValues
		}
		bbox = &BBox{}
		if bbox.Min.X, ok = v[0].(float64); !ok {
			return nil, errBBoxInvalidValue
		}
		if bbox.Min.Y, ok = v[1].(float64); !ok {
			return nil, errBBoxInvalidValue
		}
		i := 2
		if len(v) == 6 {
			if bbox.Min.Z, ok = v[2].(float64); !ok {
				return nil, errBBoxInvalidValue
			}
			i = 3
		} else {
			bbox.Min.Z = nilz
		}
		if bbox.Max.X, ok = v[i+0].(float64); !ok {
			return nil, errBBoxInvalidValue
		}
		if bbox.Max.Y, ok = v[i+1].(float64); !ok {
			return nil, errBBoxInvalidValue
		}
		if len(v) == 6 {
			if bbox.Max.Z, ok = v[i+2].(float64); !ok {
				return nil, errBBoxInvalidValue
			}
		} else {
			bbox.Max.Z = nilz
		}
	}
	return bbox, nil
}

func (b *BBox) isCordZDefined() bool {
	return b != nil && (b.Min.Z != nilz || b.Max.Z != nilz)
}

func (b *BBox) write(buf *bytes.Buffer) {
	if b == nil {
		return
	}
	hasZ := b.Min.Z != nilz && b.Max.Z != nilz
	buf.WriteString(`,"bbox":[`)
	buf.WriteString(strconv.FormatFloat(b.Min.X, 'f', -1, 64))
	buf.WriteByte(',')
	buf.WriteString(strconv.FormatFloat(b.Min.Y, 'f', -1, 64))
	if hasZ {
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatFloat(b.Min.Z, 'f', -1, 64))
	}
	buf.WriteByte(',')
	buf.WriteString(strconv.FormatFloat(b.Max.X, 'f', -1, 64))
	buf.WriteByte(',')
	buf.WriteString(strconv.FormatFloat(b.Max.Y, 'f', -1, 64))
	if hasZ {
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatFloat(b.Max.Z, 'f', -1, 64))
	}
	buf.WriteByte(']')
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

// ExternalJSON is the simple json representation of the bounding box used for external applications.
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
				Min: Position{X: x, Y: y, Z: 0},
				Max: Position{X: x + xsize, Y: y + ysize, Z: 0},
			})
		}
	}
	return bboxes
}

// BBoxesFromCenter calculates the bounding box surrounding a circle.
func BBoxesFromCenter(lat, lon, meters float64) (outer BBox) {
	outer.Max.Y, _ = geo.DestinationPoint(lat, lon, meters, 0)
	outer.Min.Y, _ = geo.DestinationPoint(lat, lon, meters, 180)
	_, outer.Min.X = geo.DestinationPoint(lat, lon, meters, 270)
	_, outer.Max.X = geo.DestinationPoint(lat, lon, meters, 90)
	return outer
}
