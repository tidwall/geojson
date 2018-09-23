package geojson

import "testing"

func TestPoint(t *testing.T) {
	p := expectJSON(t, `{"type":"Point","coordinates":[1,2,3]}`, nil)
	if p.Center() != P(1, 2) {
		t.Fatalf("expected '%v', got '%v'", P(1, 2), p.Center())
	}
	expectJSON(t, `{"type":"Point","coordinates":[1,null]}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Point","coordinates":[1,2],"bbox":null}`, errBBoxInvalid)
	expectJSON(t, `{"type":"Point"}`, errCoordinatesMissing)
	expectJSON(t, `{"type":"Point","coordinates":null}`, errCoordinatesInvalid)
	expectJSON(t, `{"type":"Point","coordinates":[1,2,3,4,5]}`, nil)
	expectJSON(t, `{"type":"Point","coordinates":[1]}`, errCoordinatesInvalid)
}

func TestPointPoly(t *testing.T) {
	p := expectJSON(t, `{"type":"Point","coordinates":[15,15,0]}`, nil)
	expect(t, !p.BBoxDefined())
	p.ForEachChild(func(Object) bool { panic("should not be reached") })
	expect(t, p.Contains(P(15, 15)))
	expect(t, p.Contains(R(15, 15, 15, 15)))
	expect(t, !p.Contains(R(10, 10, 15, 15)))
	expect(t, !p.Contains(P(10, 10)))
	expect(t, p.Intersects(P(15, 15)))
	expect(t, p.Intersects(R(10, 10, 20, 20)))
	expect(t, p.Intersects(
		expectJSON(t, `{"type":"Point","coordinates":[15,15,10]}`, nil),
	))
	expect(t, !p.Intersects(
		expectJSON(t, `{"type":"Point","coordinates":[9,15,10]}`, nil),
	))
	expect(t, p.Intersects(
		expectJSON(t, `{"type":"Point","coordinates":[9,15,10],"bbox":[10,10,20,20]}`, nil),
	))
	expect(t, p.Intersects(
		expectJSON(t, `{"type":"LineString","coordinates":[
			[10,10],[20,20]
		]}`, nil),
	))
	expect(t, !p.Intersects(
		expectJSON(t, `{"type":"LineString","coordinates":[
			[9,10],[20,20]
		]}`, nil),
	))
	expect(t, p.Intersects(
		expectJSON(t, `{"type":"Polygon","coordinates":[
			[[9,9],[9,21],[21,21],[21,9],[9,9]]
		]}`, nil),
	))

	expect(t, !p.Intersects(
		expectJSON(t, `{"type":"Polygon","coordinates":[
			[[9,9],[9,21],[21,21],[21,9],[9,9]],
			[[9.5,9.5],[9.5,20.5],[20.5,20.5],[20.5,9.5],[9.5,9.5]]
		]}`, nil),
	))
	expect(t, p.Intersects(
		expectJSON(t, `{"type":"Feature","geometry":
			{"type":"Point","coordinates":[15,15,10]}
		}`, nil),
	))
	expect(t, p.Intersects(
		expectJSON(t, `{"type":"Feature","geometry":
			{"type":"Polygon","coordinates":[
				[[9,9],[9,21],[21,21],[21,9],[9,9]]
			]}
		}`, nil),
	))
	expect(t, !p.Intersects(
		expectJSON(t, `{"type":"Feature","geometry":
			{"type":"Polygon","coordinates":[
				[[9,9],[9,21],[21,21],[21,9],[9,9]],
				[[9.5,9.5],[9.5,20.5],[20.5,20.5],[20.5,9.5],[9.5,9.5]]
			]}
		}`, nil),
	))
	expect(t, !expectJSON(t,
		`{"type":"Point","coordinates":[15,15],"bbox":[10,10,15,15]}`, nil,
	).Contains(P(7, 7)))
	expect(t, expectJSON(t,
		`{"type":"Point","coordinates":[15,15],"bbox":[10,10,15,15]}`, nil,
	).Contains(P(12, 12)))

	expect(t, !expectJSON(t,
		`{"type":"Point","coordinates":[15,15],"bbox":[10,10,15,15]}`, nil,
	).Intersects(P(7, 7)))
	expect(t, expectJSON(t,
		`{"type":"Point","coordinates":[15,15],"bbox":[10,10,15,15]}`, nil,
	).Intersects(P(12, 12)))
}

func TestPointAux(t *testing.T) {
	expect(t, (Point{Coordinates: P(10, 10)}).Contains(Point{Coordinates: P(10, 10)}))
	expect(t, !(Point{Coordinates: P(10, 10)}).Contains(Point{Coordinates: P(11, 10)}))
	expect(t, (Point{Coordinates: P(10, 10)}).Contains(
		LineString{Coordinates: []Position{P(10, 10)}},
	))
	expect(t, (Point{Coordinates: P(10, 10)}).Contains(
		Polygon{Coordinates: [][]Position{{P(10, 10)}}},
	))
	expect(t, !(Point{Coordinates: P(10, 10)}).primativeContains(nil))
	expect(t, !(Point{Coordinates: P(10, 10)}).primativeIntersects(nil))
}
