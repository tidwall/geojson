package geojson

import "testing"

func TestCirclePolygon(t *testing.T) {
	circle := CirclePolygon(-115, 33, 10000, 20)
	point := Point{Coordinates: Position{-115, 33, 0}}
	if !point.Intersects(circle) {
		t.Fatal("should intersect")
	}
	circle2 := CirclePolygon(-115, 33, 20000, 20)
	if !circle2.Intersects(circle) {
		t.Fatal("should intersect")
	}
	if !circle.Intersects(circle2) {
		t.Fatal("should intersect")
	}
	rect := Polygon{
		Coordinates: [][]Position{
			{
				{X: -120, Y: 20, Z: 0},
				{X: -120, Y: 40, Z: 0},
				{X: -100, Y: 40, Z: 0},
				{X: -100, Y: 40, Z: 0},
				{X: -120, Y: 20, Z: 0},
			},
		},
	}
	if !circle.Intersects(rect) {
		t.Fatal("should intersect")
	}
	if !rect.Intersects(circle) {
		t.Fatal("should intersect")
	}
	line := LineString{
		Coordinates: []Position{
			{X: -116, Y: 23, Z: 0},
			{X: -114, Y: 43, Z: 0},
		},
	}
	if !line.Intersects(circle) {
		t.Fatal("should intersect")
	}
}

func TestIssue281(t *testing.T) {
	p := testJSON(t, `{"type":"Polygon","coordinates":[[
		[-74.008283,40.718249],
		[-74.007339,40.713305],
		[-73.999013,40.714866],
		[-74.001760,40.720851],
		[-74.008283,40.718249]
	]]}`)

	// intersects polygon
	ls1 := testJSON(t, `{"type":"LineString","coordinates":[
		[-74.003648,40.717533],
		[-73.99575233459473,40.72046126415031],
		[-73.99721145629883,40.72338850378556]
	]}`)

	// outside polygon
	ls2 := testJSON(t, `{"type":"LineString","coordinates":[
		[-74.007682,40.722998],
		[-74.001932,40.728462],
		[-74.001846,40.723583]
	]}`)

	// inside polygon
	ls3 := testJSON(t, `{"type":"LineString","coordinates":[
          [-74.006910,40.717598],
          [-74.006137,40.715387],
          [-74.001331,40.715907]
        ]}`)

	if !ls1.Intersects(p) {
		t.Fatalf("expected true")
	}
	if !p.Intersects(ls1) {
		t.Fatalf("expected true")
	}
	if ls1.Within(p) {
		t.Fatalf("expected false")
	}
	if p.Within(ls1) {
		t.Fatalf("expected false")
	}
	if ls2.Intersects(p) {
		t.Fatalf("expected false")
	}
	if p.Intersects(ls2) {
		t.Fatalf("expected false")
	}
	if ls2.Within(p) {
		t.Fatalf("expected false")
	}
	if p.Within(ls2) {
		t.Fatalf("expected false")
	}
	if !ls3.Intersects(p) {
		t.Fatalf("expected true")
	}
	if !p.Intersects(ls3) {
		t.Fatalf("expected true")
	}
	if !ls3.Within(p) {
		t.Fatalf("expected true")
	}
	if p.Within(ls3) {
		t.Fatalf("expected false")
	}

}
