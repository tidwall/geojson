package geojson

import "testing"

func R(minX, minY, maxX, maxY float64) Rect {
	return Rect{
		Min: Position{X: minX, Y: minY},
		Max: Position{X: maxX, Y: maxY},
	}
}

func P(x, y float64) Position {
	return Position{X: x, Y: y}
}

func BenchmarkPosition(t *testing.B) {
	rr := R(0, 0, 20, 20)
	pp := P(10, 10)
	var r Object = &rr
	var p Object = &pp
	// p = &Point{Coordinates: P(10, 10)}
	for i := 0; i < t.N; i++ {
		if !r.Contains(p) {
			t.Fatal("bad")
		}
	}
}
