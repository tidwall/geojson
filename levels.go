package geojson

import (
	"bytes"
	"encoding/binary"
)

////////////////////////////////
// level 1
////////////////////////////////

func fillLevel1Map(m map[string]interface{}) (
	coordinates Position, bbox *BBox, bytesOut []byte, err error,
) {
	switch v := m["coordinates"].(type) {
	default:
		err = errInvalidCoordinates
		return
	case nil:
		err = errCoordinatesRequired
		return
	case []interface{}:
		coordinates, err = fillPosition(v)
	}
	if err == nil {
		bbox, err = fillBBox(m)
	}
	return
}

func fillLevel1Bytes(b []byte, bbox *BBox, isCordZ bool) (
	coordinates Position, bboxOut *BBox, bytesOut []byte, err error,
) {
	bboxOut = bbox
	coordinates, bytesOut, err = fillPositionBytes(b, isCordZ)
	return
}

func level1CalculatedBBox(coordinates Position, bbox *BBox) BBox {
	if bbox != nil {
		return *bbox
	}
	return BBox{
		Min: coordinates,
		Max: coordinates,
	}
}

func level1PositionCount(coordinates Position, bbox *BBox) int {
	if bbox != nil {
		return 3
	}
	return 1
}

func level1Weight(coordinates Position, bbox *BBox) int {
	return level1PositionCount(coordinates, bbox) * sizeofPosition
}

func level1JSON(name string, coordinates Position, bbox *BBox) string {
	isCordZ := level1IsCoordZDefined(coordinates, bbox)
	var buf bytes.Buffer
	buf.WriteString(`{"type":"`)
	buf.WriteString(name)
	buf.WriteString(`","coordinates":[`)
	coordinates.writeJSON(&buf, isCordZ)
	buf.WriteByte(']')
	bbox.write(&buf)
	buf.WriteByte('}')
	return buf.String()
}

func level1IsCoordZDefined(coordinates Position, bbox *BBox) bool {
	if bbox.isCordZDefined() {
		return true
	}
	return coordinates.Z != nilz
}

func level1Bytes(objType byte, coordinates Position, bbox *BBox) []byte {
	var buf bytes.Buffer
	isCordZ := level1IsCoordZDefined(coordinates, bbox)
	writeHeader(&buf, objType, bbox, isCordZ)
	coordinates.writeBytes(&buf, isCordZ)
	return buf.Bytes()
}

////////////////////////////////
// level 2
////////////////////////////////

func fillLevel2Map(m map[string]interface{}) (
	coordinates []Position, bbox *BBox, bytesOut []byte, err error,
) {
	switch v := m["coordinates"].(type) {
	default:
		err = errInvalidCoordinates
		return
	case nil:
		err = errCoordinatesRequired
		return
	case []interface{}:
		coordinates = make([]Position, len(v))
		for i, v := range v {
			v, ok := v.([]interface{})
			if !ok {
				err = errCoordinatesMustBeArray
				return
			}
			var p Position
			p, err = fillPosition(v)
			if err != nil {
				return
			}
			coordinates[i] = p
		}
	}
	bbox, err = fillBBox(m)
	return
}

func fillLevel2Bytes(b []byte, bbox *BBox, isCordZ bool) (
	coordinates []Position, bboxOut *BBox, bytesOut []byte, err error,
) {
	bboxOut = bbox
	if len(b) < 4 {
		err = errNotEnoughData
		return
	}
	coordinates = make([]Position, int(binary.LittleEndian.Uint32(b)))
	b = b[4:]
	for i := 0; i < len(coordinates); i++ {
		coordinates[i], b, err = fillPositionBytes(b, isCordZ)
		if err != nil {
			return
		}
	}
	bytesOut = b
	return
}

func level2CalculatedBBox(coordinates []Position, bbox *BBox) BBox {
	if bbox != nil {
		return *bbox
	}
	_, bbox2 := positionBBox(0, BBox{}, coordinates)
	return bbox2
}

func level2PositionCount(coordinates []Position, bbox *BBox) int {
	if bbox != nil {
		return 2 + len(coordinates)
	}
	return len(coordinates)
}

func level2Weight(coordinates []Position, bbox *BBox) int {
	return level2PositionCount(coordinates, bbox) * sizeofPosition
}

func level2JSON(name string, coordinates []Position, bbox *BBox) string {
	isCordZ := level2IsCoordZDefined(coordinates, bbox)
	var buf bytes.Buffer
	buf.WriteString(`{"type":"`)
	buf.WriteString(name)
	buf.WriteString(`","coordinates":[`)
	for i, p := range coordinates {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte('[')
		p.writeJSON(&buf, isCordZ)
		buf.WriteByte(']')
	}
	buf.WriteByte(']')
	bbox.write(&buf)
	buf.WriteByte('}')
	return buf.String()
}

func level2IsCoordZDefined(coordinates []Position, bbox *BBox) bool {
	if bbox.isCordZDefined() {
		return true
	}
	for _, p := range coordinates {
		if p.Z != nilz {
			return true
		}
	}
	return false
}

func level2Bytes(objType byte, coordinates []Position, bbox *BBox) []byte {
	var buf bytes.Buffer
	isCordZ := level2IsCoordZDefined(coordinates, bbox)
	writeHeader(&buf, objType, bbox, isCordZ)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(len(coordinates)))
	buf.Write(b)
	for _, p := range coordinates {
		p.writeBytes(&buf, isCordZ)
	}
	return buf.Bytes()
}

////////////////////////////////
// level 3
////////////////////////////////

func fillLevel3Map(m map[string]interface{}) (
	coordinates [][]Position, bbox *BBox, bytesOut []byte, err error,
) {
	switch v := m["coordinates"].(type) {
	default:
		err = errInvalidCoordinates
		return
	case nil:
		err = errCoordinatesRequired
		return
	case []interface{}:
		coordinates = make([][]Position, len(v))
		for i, v := range v {
			v, ok := v.([]interface{})
			if !ok {
				err = errInvalidCoordinatesValue
				return
			}
			ps := make([]Position, len(v))
			for i, v := range v {
				v, ok := v.([]interface{})
				if !ok {
					err = errInvalidCoordinatesValue
					return
				}
				var p Position
				p, err = fillPosition(v)
				if err != nil {
					return
				}
				ps[i] = p
			}
			coordinates[i] = ps
		}
	}
	bbox, err = fillBBox(m)
	return
}

func fillLevel3Bytes(b []byte, bbox *BBox, isCordZ bool) (
	coordinates [][]Position, bboxOut *BBox, bytesOut []byte, err error,
) {
	bboxOut = bbox
	if len(b) < 4 {
		err = errNotEnoughData
		return
	}
	coordinates = make([][]Position, int(binary.LittleEndian.Uint32(b)))
	b = b[4:]
	for i := 0; i < len(coordinates); i++ {
		if len(b) < 4 {
			err = errNotEnoughData
			return
		}
		ps := make([]Position, int(binary.LittleEndian.Uint32(b)))
		b = b[4:]
		for j := 0; j < len(ps); j++ {
			ps[j], b, err = fillPositionBytes(b, isCordZ)
			if err != nil {
				return
			}
		}
		coordinates[i] = ps
	}
	bytesOut = b
	return
}

func level3CalculatedBBox(coordinates [][]Position, bbox *BBox, isPolygon bool) BBox {
	if bbox != nil {
		return *bbox
	}
	var bbox2 BBox
	var i = 0
	for _, ps := range coordinates {
		i, bbox2 = positionBBox(i, bbox2, ps)
		if isPolygon {
			break // only the exterior ring should be calculated for a polygon
		}
	}
	return bbox2
}

func level3Weight(coordinates [][]Position, bbox *BBox) int {
	return level3PositionCount(coordinates, bbox) * sizeofPosition
}

func level3PositionCount(coordinates [][]Position, bbox *BBox) int {
	var res int
	for _, p := range coordinates {
		res += len(p)
	}
	if bbox != nil {
		return 2 + res
	}
	return res
}

func level3JSON(name string, coordinates [][]Position, bbox *BBox) string {
	isCordZ := level3IsCoordZDefined(coordinates, bbox)
	var buf bytes.Buffer
	buf.WriteString(`{"type":"`)
	buf.WriteString(name)
	buf.WriteString(`","coordinates":[`)
	for i, p := range coordinates {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte('[')
		for i, p := range p {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteByte('[')
			p.writeJSON(&buf, isCordZ)
			buf.WriteByte(']')
		}
		buf.WriteByte(']')
	}
	buf.WriteByte(']')
	bbox.write(&buf)
	buf.WriteByte('}')
	return buf.String()
}

func level3IsCoordZDefined(coordinates [][]Position, bbox *BBox) bool {
	if bbox.isCordZDefined() {
		return true
	}
	for _, p := range coordinates {
		for _, p := range p {
			if p.Z != nilz {
				return true
			}
		}
	}
	return false
}

func level3Bytes(objType byte, coordinates [][]Position, bbox *BBox) []byte {
	var buf bytes.Buffer
	isCordZ := level3IsCoordZDefined(coordinates, bbox)
	writeHeader(&buf, objType, bbox, isCordZ)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(len(coordinates)))
	buf.Write(b)
	for _, p := range coordinates {
		binary.LittleEndian.PutUint32(b, uint32(len(p)))
		buf.Write(b)
		for _, p := range p {
			p.writeBytes(&buf, isCordZ)
		}
	}
	return buf.Bytes()
}

////////////////////////////////
// level 4
////////////////////////////////

func fillLevel4Map(m map[string]interface{}) (
	coordinates [][][]Position, bbox *BBox, bytesOut []byte, err error,
) {
	switch v := m["coordinates"].(type) {
	default:
		err = errInvalidCoordinates
		return
	case nil:
		err = errCoordinatesRequired
		return
	case []interface{}:
		coordinates = make([][][]Position, len(v))
		for i, v := range v {
			v, ok := v.([]interface{})
			if !ok {
				err = errInvalidCoordinatesValue
				return
			}
			ps := make([][]Position, len(v))
			for i, v := range v {
				v, ok := v.([]interface{})
				if !ok {
					err = errInvalidCoordinatesValue
					return
				}
				pss := make([]Position, len(v))
				for i, v := range v {
					v, ok := v.([]interface{})
					if !ok {
						err = errInvalidCoordinatesValue
						return
					}
					var p Position
					p, err = fillPosition(v)
					if err != nil {
						return
					}
					pss[i] = p
				}
				ps[i] = pss
			}
			coordinates[i] = ps
		}
	}
	bbox, err = fillBBox(m)
	return
}

func fillLevel4Bytes(b []byte, bbox *BBox, isCordZ bool) (
	coordinates [][][]Position, bboxOut *BBox, bytesOut []byte, err error,
) {
	bboxOut = bbox
	if len(b) < 4 {
		err = errNotEnoughData
		return
	}
	coordinates = make([][][]Position, int(binary.LittleEndian.Uint32(b)))
	b = b[4:]
	for i := 0; i < len(coordinates); i++ {
		if len(b) < 4 {
			err = errNotEnoughData
			return
		}
		ps := make([][]Position, int(binary.LittleEndian.Uint32(b)))
		b = b[4:]
		for i := 0; i < len(ps); i++ {
			if len(b) < 4 {
				err = errNotEnoughData
				return
			}
			pss := make([]Position, int(binary.LittleEndian.Uint32(b)))
			b = b[4:]
			for i := 0; i < len(pss); i++ {
				pss[i], b, err = fillPositionBytes(b, isCordZ)
				if err != nil {
					return
				}
			}
			ps[i] = pss
		}
		coordinates[i] = ps
	}
	bytesOut = b
	return
}

func level4CalculatedBBox(coordinates [][][]Position, bbox *BBox) BBox {
	if bbox != nil {
		return *bbox
	}
	var bbox2 BBox
	var i = 0
	for _, ps := range coordinates {
		for _, ps := range ps {
			i, bbox2 = positionBBox(i, bbox2, ps)
		}
	}
	return bbox2
}

func level4Weight(coordinates [][][]Position, bbox *BBox) int {
	return level4PositionCount(coordinates, bbox) * sizeofPosition
}

func level4PositionCount(coordinates [][][]Position, bbox *BBox) int {
	var res int
	for _, p := range coordinates {
		for _, p := range p {
			res += len(p)
		}
	}
	if bbox != nil {
		return 2 + res
	}
	return res
}

func level4JSON(name string, coordinates [][][]Position, bbox *BBox) string {
	isCordZ := level4IsCoordZDefined(coordinates, bbox)
	var buf bytes.Buffer
	buf.WriteString(`{"type":"`)
	buf.WriteString(name)
	buf.WriteString(`","coordinates":[`)
	for i, p := range coordinates {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte('[')
		for i, p := range p {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteByte('[')
			for i, p := range p {
				if i > 0 {
					buf.WriteByte(',')
				}
				buf.WriteByte('[')
				p.writeJSON(&buf, isCordZ)
				buf.WriteByte(']')
			}
			buf.WriteByte(']')
		}
		buf.WriteByte(']')
	}
	buf.WriteByte(']')
	bbox.write(&buf)
	buf.WriteByte('}')
	return buf.String()
}

func level4IsCoordZDefined(coordinates [][][]Position, bbox *BBox) bool {
	if bbox.isCordZDefined() {
		return true
	}
	for _, p := range coordinates {
		for _, p := range p {
			for _, p := range p {
				if p.Z != nilz {
					return true
				}
			}
		}
	}
	return false
}

func level4Bytes(objType byte, coordinates [][][]Position, bbox *BBox) []byte {
	var buf bytes.Buffer
	isCordZ := level4IsCoordZDefined(coordinates, bbox)
	writeHeader(&buf, objType, bbox, isCordZ)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(len(coordinates)))
	buf.Write(b)
	for _, p := range coordinates {
		binary.LittleEndian.PutUint32(b, uint32(len(p)))
		buf.Write(b)
		for _, p := range p {
			binary.LittleEndian.PutUint32(b, uint32(len(p)))
			buf.Write(b)
			for _, p := range p {
				p.writeBytes(&buf, isCordZ)
			}
		}
	}
	return buf.Bytes()
}
