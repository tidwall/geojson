package geojson

import "testing"

func P(x, y float64) Position {
	return Position{X: x, Y: y}
}

func TestPosition(t *testing.T) {
	json := string(P(1, 2).AppendJSON(nil))
	exp := `{"type":"Point","coordinates":[1,2]}`
	if json != exp {
		t.Fatalf("expected '%v', got '%v'", exp, json)
	}
	if P(1, 2) != P(1, 2).Center() {
		t.Fatalf("expected '%v', got '%v'", P(1, 2), P(1, 2).Center())
	}
}
