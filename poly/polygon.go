package poly

// Polygon is a closed shape that is of an exterior ring and interior holes.
type Polygon struct {
	Exterior Ring
	Holes    []Ring
}
