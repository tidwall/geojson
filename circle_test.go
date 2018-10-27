package geojson

import (
	"testing"

	"github.com/tidwall/geojson/geometry"
)

func TestCircle(t *testing.T) {
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
}

func TestCircle_Contains(t *testing.T) {
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

func TestCircle_Intersects(t *testing.T) {
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
