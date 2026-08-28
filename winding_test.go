package geojson

import (
	"testing"

	"github.com/tidwall/geojson/geometry"
)

// Exterior is given clockwise; a CCW fix should reverse the coords.
func TestPolygonFixWindingExterior(t *testing.T) {
	// CW square: (0,0)->(0,10)->(10,10)->(10,0)->(0,0)
	in := `{"type":"Polygon","coordinates":[[[0,0],[0,10],[10,10],[10,0],[0,0]]]}`
	want := `{"type":"Polygon","coordinates":[[[0,0],[10,0],[10,10],[0,10],[0,0]]]}`
	expectJSONOpts(t, in, want, &ParseOptions{FixWinding: true})
}

// Holes given CCW must be flipped to CW; exterior already CCW stays put.
func TestPolygonFixWindingHole(t *testing.T) {
	in := `{"type":"Polygon","coordinates":[` +
		`[[0,0],[10,0],[10,10],[0,10],[0,0]],` +
		`[[2,2],[8,2],[8,8],[2,8],[2,2]]]}`
	want := `{"type":"Polygon","coordinates":[` +
		`[[0,0],[10,0],[10,10],[0,10],[0,0]],` +
		`[[2,2],[2,8],[8,8],[8,2],[2,2]]]}`
	expectJSONOpts(t, in, want, &ParseOptions{FixWinding: true})
}

// Already-correct rings must not be modified by FixWinding.
func TestPolygonFixWindingNoop(t *testing.T) {
	in := `{"type":"Polygon","coordinates":[` +
		`[[0,0],[10,0],[10,10],[0,10],[0,0]],` +
		`[[2,2],[2,8],[8,8],[8,2],[2,2]]]}`
	expectJSONOpts(t, in, in, &ParseOptions{FixWinding: true})
}

// RequireWinding must reject polygons whose exterior is CW.
func TestPolygonRequireWindingRejectsExterior(t *testing.T) {
	in := `{"type":"Polygon","coordinates":[[[0,0],[0,10],[10,10],[10,0],[0,0]]]}`
	expectJSONOpts(t, in, errCoordinatesInvalid, &ParseOptions{RequireWinding: true})
}

// RequireWinding must reject polygons whose hole is CCW.
func TestPolygonRequireWindingRejectsHole(t *testing.T) {
	in := `{"type":"Polygon","coordinates":[` +
		`[[0,0],[10,0],[10,10],[0,10],[0,0]],` +
		`[[2,2],[8,2],[8,8],[2,8],[2,2]]]}`
	expectJSONOpts(t, in, errCoordinatesInvalid, &ParseOptions{RequireWinding: true})
}

// RequireWinding accepts correctly wound polygons unchanged.
func TestPolygonRequireWindingAccepts(t *testing.T) {
	in := `{"type":"Polygon","coordinates":[` +
		`[[0,0],[10,0],[10,10],[0,10],[0,0]],` +
		`[[2,2],[2,8],[8,8],[8,2],[2,2]]]}`
	expectJSONOpts(t, in, in, &ParseOptions{RequireWinding: true})
}

// Fixing winding must carry Z coordinates along with the points.
func TestPolygonFixWindingWithZ(t *testing.T) {
	in := `{"type":"Polygon","coordinates":[` +
		`[[0,0,1],[0,10,2],[10,10,3],[10,0,4],[0,0,1]]]}`
	want := `{"type":"Polygon","coordinates":[` +
		`[[0,0,1],[10,0,4],[10,10,3],[0,10,2],[0,0,1]]]}`
	expectJSONOpts(t, in, want, &ParseOptions{FixWinding: true})
}

// MultiPolygon: mix a CW polygon and a correct polygon; only the bad one flips.
func TestMultiPolygonFixWinding(t *testing.T) {
	in := `{"type":"MultiPolygon","coordinates":[` +
		`[[[0,0],[0,10],[10,10],[10,0],[0,0]]],` +
		`[[[20,0],[30,0],[30,10],[20,10],[20,0]]]]}`
	want := `{"type":"MultiPolygon","coordinates":[` +
		`[[[0,0],[10,0],[10,10],[0,10],[0,0]]],` +
		`[[[20,0],[30,0],[30,10],[20,10],[20,0]]]]}`
	expectJSONOpts(t, in, want, &ParseOptions{FixWinding: true})
}

func TestMultiPolygonRequireWindingRejects(t *testing.T) {
	in := `{"type":"MultiPolygon","coordinates":[` +
		`[[[0,0],[10,0],[10,10],[0,10],[0,0]]],` +
		`[[[20,0],[20,10],[30,10],[30,0],[20,0]]]]}`
	expectJSONOpts(t, in, errCoordinatesInvalid, &ParseOptions{RequireWinding: true})
}

// Degenerate rings (zero signed area) carry no orientation and must be
// accepted by RequireWinding rather than rejected.
func TestPolygonRequireWindingAcceptsDegenerate(t *testing.T) {
	in := `{"type":"Polygon","coordinates":[[[0,0],[5,5],[10,10],[0,0]]]}`
	expectJSONOpts(t, in, in, &ParseOptions{RequireWinding: true})
}

// Polygon with a hole plus z-coords: exercises pointOffset advancing across
// the exterior before reversing the hole's slice of extra.values.
func TestPolygonFixWindingHoleWithZ(t *testing.T) {
	in := `{"type":"Polygon","coordinates":[` +
		`[[0,0,1],[10,0,2],[10,10,3],[0,10,4],[0,0,1]],` +
		`[[2,2,5],[8,2,6],[8,8,7],[2,8,8],[2,2,5]]]}`
	want := `{"type":"Polygon","coordinates":[` +
		`[[0,0,1],[10,0,2],[10,10,3],[0,10,4],[0,0,1]],` +
		`[[2,2,5],[2,8,8],[8,8,7],[8,2,6],[2,2,5]]]}`
	expectJSONOpts(t, in, want, &ParseOptions{FixWinding: true})
}

// MultiPolygon where each sub-polygon carries its own z-values: confirms the
// per-polygon offset resets between children (each child has its own extra).
func TestMultiPolygonFixWindingPerChildZ(t *testing.T) {
	in := `{"type":"MultiPolygon","coordinates":[` +
		`[[[0,0,1],[0,10,2],[10,10,3],[10,0,4],[0,0,1]]],` +
		`[[[20,0,5],[30,0,6],[30,10,7],[20,10,8],[20,0,5]]]]}`
	want := `{"type":"MultiPolygon","coordinates":[` +
		`[[[0,0,1],[10,0,4],[10,10,3],[0,10,2],[0,0,1]]],` +
		`[[[20,0,5],[30,0,6],[30,10,7],[20,10,8],[20,0,5]]]]}`
	expectJSONOpts(t, in, want, &ParseOptions{FixWinding: true})
}

func TestRingOrientation(t *testing.T) {
	ccw := []geometry.Point{{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 10}, {X: 0, Y: 10}, {X: 0, Y: 0}}
	cw := []geometry.Point{{X: 0, Y: 0}, {X: 0, Y: 10}, {X: 10, Y: 10}, {X: 10, Y: 0}, {X: 0, Y: 0}}
	degen := []geometry.Point{{X: 0, Y: 0}, {X: 5, Y: 5}, {X: 10, Y: 10}, {X: 0, Y: 0}}

	if got := ringOrientation(ccw); got != 1 {
		t.Fatalf("ccw: expected 1, got %d", got)
	}
	if got := ringOrientation(cw); got != -1 {
		t.Fatalf("cw: expected -1, got %d", got)
	}
	if got := ringOrientation(degen); got != 0 {
		t.Fatalf("degenerate: expected 0, got %d", got)
	}
}
