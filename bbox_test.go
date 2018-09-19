package geojson

import "testing"

func TestBBox(t *testing.T) {
	_, err := loadBBox(`{"bbox":[]}`)
	if err != errBBoxInvalid {
		t.Fatalf("expected '%v', got '%v'", errBBoxInvalid, err)
	}
	_, err = loadBBox(`{"bbox":null}`)
	if err != errBBoxInvalid {
		t.Fatalf("expected '%v', got '%v'", errBBoxInvalid, err)
	}
	_, err = loadBBox(`{"bbox":[0,1,2,3,4,5,6,7,8]}`)
	if err != nil {
		t.Fatalf("expected '%v', got '%v'", nil, err)
	}
	_, err = loadBBox(`{"bbox":["sadf"]}`)
	if err != errBBoxInvalid {
		t.Fatalf("expected '%v', got '%v'", errBBoxInvalid, err)
	}
	bbox, err := loadBBox(`{"bbox":[1,2,3,4]}`)
	if err != nil {
		t.Fatalf("expected '%v', got '%v'", nil, err)
	}
	if bboxWeight(bbox) != 32 {
		t.Fatalf("expected '%v', got '%v'", 32, bboxWeight(bbox))
	}
	if bboxWeight(nil) != 0 {
		t.Fatalf("expected '%v', got '%v'", 0, bboxWeight(nil))
	}
	if bboxPositionCount(bbox) != 2 {
		t.Fatalf("expected '%v', got '%v'", 2, bboxPositionCount(bbox))
	}
	if bboxPositionCount(nil) != 0 {
		t.Fatalf("expected '%v', got '%v'", 0, bboxPositionCount(nil))
	}
}
