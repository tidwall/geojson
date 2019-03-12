package geojson

import (
	"testing"

	"github.com/tidwall/geojson/geometry"
)

func TestCircleNew(t *testing.T) {
	expectJSON(t,
		`{"type":"Feature","geometry":{"type":"Point","coordinates":[-112,33]},"properties":{"type":"Circle","radius":"5000"}}`,
		`{"type":"Feature","geometry":{"type":"Point","coordinates":[-112,33]},"properties":{"type":"Circle","radius":5000,"radius_units":"m"}}`,
	)
	g, err := Parse(`{
		"type":"Feature",
		"geometry":{"type":"Point","coordinates":[-112.2693,33.5123]},
		"properties": {
		  "type": "Circle",
		  "radius": 1000
		 }
	  }`, nil)
	if err != nil {
		t.Fatal(err)
	}
	expect(t, g.Contains(PO(-112.26, 33.51)))

	circle := NewCircle(P(-112, 33), 123456.654321, 64)
	expectJSON(t, circle.JSON(), `{"type":"Feature","geometry":{"type":"Point","coordinates":[-112,33]},"properties":{"type":"Circle","radius":123456.654321,"radius_units":"m"}}`)

}

func TestCircleContains(t *testing.T) {
	g := NewCircle(P(-122.4412, 37.7335), 1000, 64)
	expect(t, g.Contains(PO(-122.4412, 37.7335)))
	expect(t, g.Contains(PO(-122.44121, 37.7335)))
	expect(t, g.Contains(
		MPO([]geometry.Point{
			P(-122.4408378, 37.7341129),
			P(-122.4408378, 37.733)})))
	expect(t, g.Contains(
		NewCircle(P(-122.44121, 37.7335), 500, 64)))
	expect(t, g.Contains(
		LO([]geometry.Point{
			P(-122.4408378, 37.7341129),
			P(-122.4408378, 37.733)})))
	expect(t, g.Contains(
		MLO([]*geometry.Line{
			L([]geometry.Point{
				P(-122.4408378, 37.7341129),
				P(-122.4408378, 37.733),
			}),
			L([]geometry.Point{
				P(-122.44, 37.733),
				P(-122.44, 37.7341129),
			})})))
	expect(t, g.Contains(
		PPO(
			[]geometry.Point{
				P(-122.4408378, 37.7341129),
				P(-122.4408378, 37.733),
				P(-122.44, 37.733),
				P(-122.44, 37.7341129),
				P(-122.4408378, 37.7341129),
			},
			[][]geometry.Point{})))

	// Does-not-contain
	expect(t, !g.Contains(PO(-122.265, 37.826)))
	expect(t, !g.Contains(
		NewCircle(P(-122.265, 37.826), 100, 64)))
	expect(t, !g.Contains(
		LO([]geometry.Point{
			P(-122.265, 37.826),
			P(-122.210, 37.860)})))
	expect(t, !g.Contains(
		MPO([]geometry.Point{
			P(-122.4408378, 37.7341129),
			P(-122.198181, 37.7490)})))
	expect(t, !g.Contains(
		MLO([]*geometry.Line{
			L([]geometry.Point{
				P(-122.265, 37.826),
				P(-122.265, 37.860),
			}),
			L([]geometry.Point{
				P(-122.44, 37.733),
				P(-122.44, 37.7341129),
			})})))
	expect(t, !g.Contains(PPO(
		[]geometry.Point{
			P(-122.265, 37.826),
			P(-122.265, 37.860),
			P(-122.210, 37.860),
			P(-122.210, 37.826),
			P(-122.265, 37.826),
		},
		[][]geometry.Point{})))
}

func TestCircleIntersects(t *testing.T) {
	g := NewCircle(P(-122.4412, 37.7335), 1000, 64)
	expect(t, g.Intersects(PO(-122.4412, 37.7335)))
	expect(t, g.Intersects(PO(-122.44121, 37.7335)))
	expect(t, g.Intersects(
		MPO([]geometry.Point{
			P(-122.4408378, 37.7341129),
			P(-122.4408378, 37.733)})))
	expect(t, g.Intersects(
		NewCircle(P(-122.44121, 37.7335), 500, 64)))
	expect(t, g.Intersects(
		LO([]geometry.Point{
			P(-122.4408378, 37.7341129),
			P(-122.4408378, 37.733)})))
	expect(t, g.Intersects(
		LO([]geometry.Point{
			P(-122.4408378, 37.7341129),
			P(-122.265, 37.826)})))
	expect(t, g.Intersects(
		MLO([]*geometry.Line{
			L([]geometry.Point{
				P(-122.4408378, 37.7341129),
				P(-122.265, 37.826),
			}),
			L([]geometry.Point{
				P(-122.44, 37.733),
				P(-122.44, 37.7341129),
			})})))
	expect(t, g.Intersects(
		PPO(
			[]geometry.Point{
				P(-122.4408378, 37.7341129),
				P(-122.265, 37.860),
				P(-122.210, 37.826),
				P(-122.44, 37.7341129),
				P(-122.4408378, 37.7341129),
			},
			[][]geometry.Point{})))
	expect(t, g.Intersects(
		MPO([]geometry.Point{
			P(-122.4408378, 37.7341129),
			P(-122.198181, 37.7490)})))
	expect(t, g.Intersects(
		MLO([]*geometry.Line{
			L([]geometry.Point{
				P(-122.265, 37.826),
				P(-122.265, 37.860),
			}),
			L([]geometry.Point{
				P(-122.44, 37.733),
				P(-122.44, 37.7341129),
			})})))

	// Does-not-intersect
	expect(t, !g.Intersects(PO(-122.265, 37.826)))
	expect(t, !g.Intersects(
		NewCircle(P(-122.265, 37.826), 100, 64)))
	expect(t, !g.Intersects(
		LO([]geometry.Point{
			P(-122.265, 37.826),
			P(-122.210, 37.860)})))
}

// This snippet tests 100M comparisons.
// On my box this takes 24.5s without haversine trick, and 13.7s with the trick.
//
//func TestCirclePerformance(t *testing.T) {
//	g := NewCircle(P(-122.4412, 37.7335), 1000, 64)
//	r := rand.New(rand.NewSource(42))
//	for i:= 0; i < 100000000; i++ {
//		g.Contains(PO(r.Float64()*360 - 180, r.Float64()*180 - 90))
//	}
//	expect(t, true)
//}

func TestPointCircle(t *testing.T) {
	p := NewPoint(geometry.Point{X: -0.8856761, Y: 52.7563759})
	c := NewCircle(geometry.Point{X: -0.8838196, Y: 52.7563395}, 200, 20)
	if !p.Within(c) {
		t.Fatal("expected true")
	}
}
