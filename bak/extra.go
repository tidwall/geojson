package geojson

// Dims are extra ZM dimensions
type Dims int

const (
	// DimsZ defines extra Z coordinate values
	DimsZ Dims = 1
	// DimsZM defines extra ZM coordinate values
	DimsZM Dims = 2
)

// Extra is extra ZM coordinate values
type Extra struct {
	Dims      Dims
	Positions []float64
}
