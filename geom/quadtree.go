package geom

type quadTreeQuad struct {
	data interface{}
}

type quadTree struct {
	root quadTreeQuad
	min  [2]float64
	max  [2]float64
}
