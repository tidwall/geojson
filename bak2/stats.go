package geojson

// Stats of an object
type Stats struct {
	// Weight is the estimated memory size of the object in bytes
	Weight int
	// PositionCount is the number of point in the object
	PositionCount int
}
