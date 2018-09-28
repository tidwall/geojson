package geom

// Line ...
type Line interface {
	Scan(iter func(seg Segment) bool)
	Search(rect Rect, iter func(seg Segment, index int) bool)
	Points() []Point
	Rect() Rect
	IsClosed() bool
}
