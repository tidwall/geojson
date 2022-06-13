package geojson

import "errors"

/*
Point               0001   1001   2001   3001
LineString          0002   1002   2002   3002
Polygon             0003   1003   2003   3003
MultiPoint          0004   1004   2004   3004
MultiLineString     0005   1005   2005   3005
MultiPolygon        0006   1006   2006   3006
GeometryCollection  0007   1007   2007   3007
*/

func trimPrefixSpace(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] <= ' ' {
			switch s[i] {
			case ' ', '\t', '\n', '\r':
				continue
			}
		}
		return s[i:]
	}
	return ""
}

func hasWKTPrefix(s string, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	for i := 0; i < len(s); i++ {
		sc := s[i]
		if sc >= 'A' && sc <= 'Z' {
			sc += 32
		}
		if sc != prefix[i] {
			return false
		}
	}
	return true
}

func parseWKT(data string, opts *ParseOptions) (Object, error) {
	// whitespace is already parsed
	switch {
	case hasWKTPrefix(data, "POINT"):
		return parseWKTPoint(data[5:], opts)
	case hasWKTPrefix(data, "LINESTRING"):
		return parseWKTLineString(data[10:], opts)
	case hasWKTPrefix(data, "POLYGON"):
		return parseWKTPolygon(data[7:], opts)
	case hasWKTPrefix(data, "MULTIPOINT"):
		return parseWKTMultiPoint(data[10:], opts)
	case hasWKTPrefix(data, "MULTILINESTRING"):
		return parseWKTMultiLineString(data[15:], opts)
	case hasWKTPrefix(data, "MULTIPOLYGON"):
		return parseWKTMultiPolygon(data[12:], opts)
	case hasWKTPrefix(data, "GEOMETRYCOLLECTION"):
		return parseWKTGeometryCollection(data[18:], opts)
	default:
		return nil, errDataInvalid
	}
}

func parseWKTPoint(data string, opts *ParseOptions) (Object, error) {

	return nil, errors.New("WKT POINT not supported")
}
func parseWKTLineString(data string, opts *ParseOptions) (Object, error) {
	return nil, errors.New("WKT LINESTRING not supported")
}
func parseWKTPolygon(data string, opts *ParseOptions) (Object, error) {
	return nil, errors.New("WKT POLYGON not supported")
}
func parseWKTMultiPoint(data string, opts *ParseOptions) (Object, error) {
	return nil, errors.New("WKT MULTIPOINT not supported")
}
func parseWKTMultiLineString(data string, opts *ParseOptions) (Object, error) {
	return nil, errors.New("WKT MULTILINESTRING not supported")
}
func parseWKTMultiPolygon(data string, opts *ParseOptions) (Object, error) {
	return nil, errors.New("WKT MULTIPOLYGON not supported")
}
func parseWKTGeometryCollection(data string, opts *ParseOptions) (Object, error) {
	return nil, errors.New("WKT GEOMETRYCOLLECTION not supported")
}
