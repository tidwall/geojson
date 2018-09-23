package geojson

import "testing"

func R(minX, minY, maxX, maxY float64) Rect {
	return Rect{
		Min: Position{X: minX, Y: minY},
		Max: Position{X: maxX, Y: maxY},
	}
}

func TestRect(t *testing.T) {
	bbox, err := loadBBox(`{"bbox":[1,2,3,4]}`)
	if err != nil {
		t.Fatal(err)
	}
	rect := bbox.Rect()
	if rect != R(1, 2, 3, 4) {
		t.Fatalf("expected '%v', got '%v'", R(1, 2, 3, 4), rect)
	}
	if rect.Center() != P(2, 3) {
		t.Fatalf("expected '%v', got '%v'", P(2, 3), rect.Center())
	}
	json := string(rect.AppendJSON(nil))
	exp := `{"type":"Polygon","coordinates":[[[1,2],[3,2],[3,4],[1,4],[1,2]]]}`
	if json != exp {
		t.Fatalf("expected '%v', got '%v'", exp, json)
	}
	json = string(R(1, 2, 1, 2).AppendJSON(nil))
	exp = `{"type":"Point","coordinates":[1,2]}`
	if json != exp {
		t.Fatalf("expected '%v', got '%v'", exp, json)
	}

	bbox = bboxRect{rect: rect}
	json = string(bbox.AppendJSON(nil))
	exp = `{"type":"Polygon","coordinates":[[[1,2],[3,2],[3,4],[1,4],[1,2]]]}`
	if json != exp {
		t.Fatalf("expected '%v', got '%v'", exp, json)
	}
	if !R(10, 10, 20, 20).ContainsRect(R(10, 10, 20, 20)) {
		t.Fatal("expected true")
	}
	if R(11, 10, 20, 20).ContainsRect(R(10, 10, 20, 20)) {
		t.Fatal("expected false")
	}
	if R(10, 11, 20, 20).ContainsRect(R(10, 10, 20, 20)) {
		t.Fatal("expected false")
	}
	if R(10, 10, 19, 20).ContainsRect(R(10, 10, 20, 20)) {
		t.Fatal("expected false")
	}
	if R(10, 10, 20, 19).ContainsRect(R(10, 10, 20, 20)) {
		t.Fatal("expected false")
	}
	if !R(10, 10, 20, 20).IntersectsRect(R(10, 10, 20, 20)) {
		t.Fatal("expected true")
	}
	if !R(0, 0, 20, 20).IntersectsRect(R(10, 10, 20, 20)) {
		t.Fatal("expected true")
	}
	if !R(0, 0, 10, 10).IntersectsRect(R(10, 10, 20, 20)) {
		t.Fatal("expected true")
	}
	if R(0, 0, 9, 9).IntersectsRect(R(10, 10, 20, 20)) {
		t.Fatal("expected true")
	}
	if R(20, 21, 29, 29).IntersectsRect(R(10, 10, 20, 20)) {
		t.Fatal("expected true")
	}
}

func TestRectPoly(t *testing.T) {
	r := R(10, 10, 20, 20)
	expect(t, !r.BBoxDefined())
	r.ForEachChild(func(Object) bool { panic("should not be reached") })
	expect(t, r.ContainsPosition(P(15, 15)))
	expect(t, r.ContainsPosition(P(10, 10)))
	expect(t, r.ContainsPosition(P(20, 20)))
	expect(t, !r.ContainsPosition(P(21, 20)))
	expect(t, !r.ContainsPosition(P(20, 21)))
	expect(t, !r.ContainsPosition(P(10, 9)))
	expect(t, !r.ContainsPosition(P(9, 10)))
	expect(t, r.Contains(R(10, 10, 20, 20)))
	expect(t, r.Contains(R(11, 11, 19, 19)))
	expect(t, !r.Contains(R(9, 10, 20, 20)))
	expect(t, r.Contains(P(15, 15)))
	expect(t, !r.Contains(P(21, 20)))
	expect(t, r.Contains(
		expectJSON(t, `{"type":"Point","coordinates":[15,15,10]}`, nil),
	))
	expect(t, r.Contains(
		expectJSON(t, `{"type":"Point","coordinates":[5,5],"bbox":[11,11,19,19]}`, nil),
	))
	expect(t, r.Intersects(R(10, 10, 20, 20)))
	expect(t, r.Intersects(R(11, 11, 19, 19)))
	expect(t, r.Intersects(R(9, 10, 20, 20)))
	expect(t, !r.Intersects(R(2, 10, 9, 20)))
	expect(t, r.Intersects(P(15, 15)))
	expect(t, !r.Intersects(P(21, 20)))
	expect(t, r.Intersects(
		expectJSON(t, `{"type":"Point","coordinates":[15,15,10]}`, nil),
	))
	expect(t, !r.Intersects(
		expectJSON(t, `{"type":"Point","coordinates":[5,15,10]}`, nil),
	))
	expect(t, r.Intersects(
		expectJSON(t, `{"type":"Point","coordinates":[5,5],"bbox":[11,11,19,19]}`, nil),
	))
	expect(t, r.Intersects(
		expectJSON(t, `{"type":"LineString","coordinates":[[5,5],[25,25]]}`, nil),
	))
	expect(t, !r.Intersects(
		expectJSON(t, `{"type":"LineString","coordinates":[
			[9,9],[9,21],[21,21],[21,9]
		]}`, nil),
	))
	expect(t, !r.Intersects(
		expectJSON(t, `{"type":"Polygon","coordinates":[
			[[9,9],[9,21],[21,21],[21,9],[9,9]],
			[[9.5,9.5],[9.5,20.5],[20.5,20.5],[20.5,9.5],[9.5,9.5]]
		]}`, nil),
	))

	expect(t, r.Intersects(
		expectJSON(t, `{"type":"Feature","geometry":
			{"type":"Point","coordinates":[15,15,10]}
		}`, nil),
	))
	expect(t, !r.Intersects(
		expectJSON(t, `{"type":"Feature","geometry":
			{"type":"Polygon","coordinates":[
				[[9,9],[9,21],[21,21],[21,9],[9,9]],
				[[9.5,9.5],[9.5,20.5],[20.5,20.5],[20.5,9.5],[9.5,9.5]]
			]}
		}`, nil),
	))

	expect(t, R(9, 10, 20, 20).Contains(
		expectJSON(t, `{"type":"Feature","geometry":
			{"type":"Polygon","coordinates":[
				[[10,10],[10,20],[20,20],[20,10],[10,10]]
			]}
		}`, nil),
	))

}

func TestRectAux(t *testing.T) {
	expect(t, R(10, 10, 20, 20).Contains(
		expectJSON(t, `{"type":"Feature","geometry":
			{"type":"Polygon","coordinates":[
				[[10,10],[10,20],[20,20],[20,10],[10,10]]
			]}
		}`, nil),
	))
	expect(t, R(0, 0, 20, 20).primativeContains(
		expectJSON(t, `{"type":"LineString","coordinates":[[5,5],[8,8]]}`, nil),
	))
	expect(t, !R(10, 10, 20, 20).primativeContains(
		expectJSON(t, `{"type":"LineString","coordinates":[[5,5],[8,8]]}`, nil),
	))
	expect(t, !R(10, 10, 20, 20).primativeContains(nil))
	expect(t, !R(10, 10, 20, 20).primativeIntersects(nil))
}
