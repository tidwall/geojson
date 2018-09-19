package geojson

type Dims int

const (
	DimsZ  Dims = 1
	DimsZM Dims = 2
)

type Extra struct {
	Dims      Dims
	Positions []float64
}
