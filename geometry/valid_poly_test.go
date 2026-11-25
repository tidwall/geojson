package geometry

import "testing"

// Regression test for issue #25: Valid() must reject polygons whose rings
// are not simple, not closed, or have too few points.
func TestPolyValidGeometry(t *testing.T) {
	pt := func(x, y float64) Point { return Point{X: x, Y: y} }
	tests := []struct {
		name  string
		poly  *Poly
		valid bool
	}{
		{
			name: "simple square is valid",
			poly: NewPoly([]Point{
				pt(0, 0), pt(0, 1), pt(1, 1), pt(1, 0), pt(0, 0),
			}, nil, nil),
			valid: true,
		},
		{
			name: "self-intersecting bowtie is invalid",
			poly: NewPoly([]Point{
				pt(0, 0), pt(2, 2), pt(0, 2), pt(2, 0), pt(0, 0),
			}, nil, nil),
			valid: false,
		},
		{
			name: "unclosed ring is invalid",
			poly: NewPoly([]Point{
				pt(0, 0), pt(0, 1), pt(1, 1), pt(1, 0),
			}, nil, nil),
			valid: false,
		},
		{
			name: "degenerate two-point ring is invalid",
			poly: NewPoly([]Point{
				pt(0, 0), pt(1, 1), pt(0, 0),
			}, nil, nil),
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.poly.Valid(); got != tt.valid {
				t.Fatalf("Valid() = %v, want %v", got, tt.valid)
			}
		})
	}
}
