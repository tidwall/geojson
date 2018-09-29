package geom

// // Line ...
// type Line interface {
// 	Scan(iter func(seg Segment) bool)
// 	Search(rect Rect, iter func(seg Segment, index int) bool)
// 	Points() []Point
// 	Rect() Rect
// 	IsClosed() bool

// 	ContainsPoint(point Point) bool
// 	IntersectsPoint(point Point) bool
// }

// // NewLine ...
// func NewLine(points []Point, index int) Line {
// 	if index >= 0 && len(points) > index {
// 		return newLineIndexed(points)
// 	}
// 	return newLineSimple(points)
// }
