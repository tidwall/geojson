package geojson

// func (a Point) contains(b Object) bool {
// 	switch b := b.(type) {
// 	default:
// 		return false
// 	case SimplePoint:
// 		return a.containsSimplePoint(b)
// 	case Point:
// 		return a.containsPoint(b)
// 	case MultiPoint:
// 		return a.containsMultiPoint(b)
// 	case LineString:
// 		return a.containsLineString(b)
// 	case MultiLineString:
// 		return a.containsMultiLineString(b)
// 	case Polygon:
// 		return a.containsPolygon(b)
// 	case MultiPolygon:
// 		return a.containsPolygon(b)
// 	}
// }

// func withinSimplePoint(a Object, b SimplePoint) bool {
// 	return a.WithinBBox(b.CalculatedBBox())
// }
// func withinPoint(a Object, b Point) bool {
// 	return a.WithinBBox(b.CalculatedBBox())
// }
// func withinMultiPoint(a Object, b MultiPoint) bool {
// 	for i := range b.Coordinates {
// 		if a.Within(Point{Coordinates: b.Coordinates[i]}) {
// 			return true
// 		}
// 	}
// 	return false
// }
// func withinLineString(a Object, b LineString) bool {
// 	if len(b.Coordinates) == 0 {
// 		return false
// 	}
// 	switch a := a.(type) {
// 	default:
// 		return false
// 	case SimplePoint:
// 		return poly.Point(Position{X: a.X, Y: a.Y, Z: 0}).
// 			IntersectsLineString(
// 				polyPositions(b.Coordinates),
// 			)
// 	case Point:
// 		return poly.Point(a.Coordinates).IntersectsLineString(
// 			polyPositions(b.Coordinates),
// 		)
// 	case MultiPoint:
// 		if len(b.Coordinates) == 0 {
// 			return false
// 		}
// 		for _, p := range b.Coordinates {
// 			if !poly.Point(p).IntersectsLineString(
// 				polyPositions(b.Coordinates),
// 			) {
// 				return false
// 			}
// 		}
// 		return true
// 	}
// }

// func withinObjectShared2(a Object, b Object) bool {
// 	bbp := b.bboxPtr()
// 	if bbp != nil {
// 		if !a.WithinBBox(*bbp) {
// 			return false
// 		}
// 		if b.IsBBoxDefined() {
// 			return true
// 		}
// 	}
// 	switch b := b.(type) {
// 	default:
// 		return false
// 	case SimplePoint:
// 		return withinSimplePoint(a, b)
// 	case Point:
// 		return withinPoint(a, b)
// 	case MultiPoint:
// 		return withinMultiPoint(a, b)
// 	case LineString:
// 		return withinLineString(a, b)
// 	case MultiLineString:
// 		for i := range b.Coordinates {
// 			if a.Within(b.getLineString(i)) {
// 				return true
// 			}
// 		}
// 		return false
// 	case Polygon:
// 		if len(b.Coordinates) == 0 {
// 			return false
// 		}

// 		return pin(b)
// 	case MultiPolygon:
// 		for i := range b.Coordinates {
// 			if pin(b.getPolygon(i)) {
// 				return true
// 			}
// 		}
// 		return false
// 	case Feature:
// 		return a.Within(b.Geometry)
// 	case FeatureCollection:
// 		if len(b.Features) == 0 {
// 			return false
// 		}
// 		for _, f := range b.Features {
// 			if !a.Within(f) {
// 				return false
// 			}
// 		}
// 		return true
// 	case GeometryCollection:
// 		if len(b.Geometries) == 0 {
// 			return false
// 		}
// 		for _, f := range b.Geometries {
// 			if !a.Within(f) {
// 				return false
// 			}
// 		}
// 		return true
// 	}
// }
