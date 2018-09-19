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
