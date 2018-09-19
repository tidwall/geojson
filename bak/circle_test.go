package geojson

import (
	"testing"
)

func TestSegmentIntersectsCirclePointsInside(t *testing.T) {
	// either start or end is within the circle
	start := Position{X: -122.4408378, Y: 37.7341129, Z: 0}
	end := Position{X: -122.4408378, Y: 37.733, Z: 0}
	center := Position{X: -122.4409, Y: 37.734, Z: 0}
	meters := 30.0
	if !SegmentIntersectsCircle(start, end, center, meters) {
		t.Fatal("!")
	}
	center = Position{X: -122.4409, Y: 37.733, Z: 0}
	if !SegmentIntersectsCircle(start, end, center, meters) {
		t.Fatal("!")
	}
}

func TestSegmentIntersectsCirclePointsOutside(t *testing.T) {
	// neither start nor end are within the circle, but the segment intersects it
	start := Position{X: -122.4408378, Y: 37.7341129, Z: 0}
	end := Position{X: -122.4408378, Y: 37.733, Z: 0}
	center := Position{X: -122.4412, Y: 37.7335, Z: 0}
	meters := 70.0
	if !SegmentIntersectsCircle(start, end, center, meters) {
		t.Fatal("!")
	}
}

func TestSegmentIntersectsCircleLineButNotSegment(t *testing.T) {
	// the line of the segment intersects the circle, but the segment does not
	start := Position{X: -122.4408378, Y: 37.7341129, Z: 0}
	end := Position{X: -122.4408378, Y: 37.733, Z: 0}
	center := Position{X: -122.4412, Y: 37.737, Z: 0}
	meters := 70.0
	if SegmentIntersectsCircle(start, end, center, meters) {
		t.Fatal("!")
	}
}
