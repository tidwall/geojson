package ring

// Ring ...
type Ring interface {
	Scan(iter func(seg Segment) bool)
	Search(rect Rect, iter func(seg Segment, index int) bool)
	Points() []Point
	Rect() Rect
	Convex() bool
	IntersectsSegment(seg Segment, allowOnEdge bool) bool
	ContainsPoint(point Point, allowOnEdge bool) bool
	IntersectsRing(ring Ring) bool
	//ContainsRing(ring Ring) bool
}

// NewRing ...
func NewRing(points []Point, indexed bool) Ring {
	if indexed {
		return newTreeRing(points)
	}
	return newSimpleRing(points)
}
