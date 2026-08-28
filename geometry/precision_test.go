package geometry

import (
	"testing"
)

func TestFloatComparisons(t *testing.T) {
	// Test basic float comparisons
	a := 1.0
	b := 1.0 + epsilon/2 // Should be considered equal
	c := 1.0 + epsilon*2 // Should be considered different

	if !FloatEqual(a, b) {
		t.Errorf("Expected %f and %f to be equal", a, b)
	}

	if FloatEqual(a, c) {
		t.Errorf("Expected %f and %f to not be equal", a, c)
	}

	if !FloatLess(a, c) {
		t.Errorf("Expected %f to be less than %f", a, c)
	}

	if !FloatGreater(c, a) {
		t.Errorf("Expected %f to be greater than %f", c, a)
	}
}

func TestPointComparisons(t *testing.T) {
	// Test point comparisons
	p1 := Point{X: 1.0, Y: 2.0}
	p2 := Point{X: 1.0 + epsilon/2, Y: 2.0 + epsilon/2} // Should be equal
	p3 := Point{X: 1.0 + epsilon*2, Y: 2.0}             // Should be different

	if !PointEqual(p1, p2) {
		t.Errorf("Expected points %v and %v to be equal", p1, p2)
	}

	if PointEqual(p1, p3) {
		t.Errorf("Expected points %v and %v to not be equal", p1, p3)
	}
}

func TestSegmentIntersection(t *testing.T) {
	// Test segment intersection with floating point precision
	seg1 := Segment{A: Point{X: 0, Y: 0}, B: Point{X: 1, Y: 1}}
	seg2 := Segment{A: Point{X: 0, Y: 1}, B: Point{X: 1, Y: 0}}

	// These segments should intersect at (0.5, 0.5)
	if !seg1.IntersectsSegment(seg2) {
		t.Error("Expected segments to intersect")
	}

	// Test with very close but not identical points
	seg3 := Segment{
		A: Point{X: 0 + epsilon/2, Y: 0 + epsilon/2},
		B: Point{X: 1 + epsilon/2, Y: 1 + epsilon/2},
	}
	seg4 := Segment{
		A: Point{X: 0 + epsilon/2, Y: 1 + epsilon/2},
		B: Point{X: 1 + epsilon/2, Y: 0 + epsilon/2},
	}

	if !seg3.IntersectsSegment(seg4) {
		t.Error("Expected segments with epsilon differences to still intersect")
	}
}

func TestIntersectsSegmentPrecision(t *testing.T) {
	// Test boundary conditions with epsilon precision

	// Test case 1: Segments that barely intersect
	seg1 := Segment{A: Point{X: 0, Y: 0}, B: Point{X: 1, Y: 0}}
	seg2 := Segment{A: Point{X: 0.5, Y: -epsilon / 2}, B: Point{X: 0.5, Y: epsilon / 2}}

	if !seg1.IntersectsSegment(seg2) {
		t.Error("Expected segments with epsilon-level intersection to intersect")
	}

	// Test case 2: Segments that barely don't intersect
	seg3 := Segment{A: Point{X: 0, Y: 0}, B: Point{X: 1, Y: 0}}
	seg4 := Segment{A: Point{X: 0.5, Y: epsilon * 2}, B: Point{X: 0.5, Y: epsilon * 3}}

	if seg3.IntersectsSegment(seg4) {
		t.Error("Expected segments with clear separation to not intersect")
	}

	// Test case 3: Collinear segments with epsilon differences
	seg5 := Segment{A: Point{X: 0, Y: 0}, B: Point{X: 1, Y: 1}}
	seg6 := Segment{A: Point{X: 0.5 + epsilon/2, Y: 0.5 + epsilon/2}, B: Point{X: 1.5 + epsilon/2, Y: 1.5 + epsilon/2}}

	if !seg5.IntersectsSegment(seg6) {
		t.Error("Expected overlapping collinear segments with epsilon differences to intersect")
	}

	// Test case 4: Parallel segments
	seg7 := Segment{A: Point{X: 0, Y: 0}, B: Point{X: 1, Y: 0}}
	seg8 := Segment{A: Point{X: 0, Y: epsilon * 2}, B: Point{X: 1, Y: epsilon * 2}}

	if seg7.IntersectsSegment(seg8) {
		t.Error("Expected parallel segments to not intersect")
	}
}

func TestRaycastPrecision(t *testing.T) {
	// Test raycast with floating point precision
	seg := Segment{A: Point{X: 0, Y: 0}, B: Point{X: 1, Y: 0}}

	// Point exactly on the segment
	p1 := Point{X: 0.5, Y: 0}
	result1 := seg.Raycast(p1)
	if !result1.On {
		t.Error("Expected point to be on segment")
	}

	// Point very close to the segment (within epsilon)
	p2 := Point{X: 0.5, Y: epsilon / 2}
	result2 := seg.Raycast(p2)
	_ = result2 // This should be considered as close enough depending on implementation

	// Point clearly off the segment
	p3 := Point{X: 0.5, Y: epsilon * 10}
	result3 := seg.Raycast(p3)
	if result3.On {
		t.Error("Expected point to not be on segment")
	}
}

func TestEpsilonChoiceValidation(t *testing.T) {
	// Test to validate the choice of epsilon = 1e-8 vs 1e-10

	// Simulate typical geometric computation errors
	baseValue := 1.0

	// Error from a typical floating point multiplication chain
	// e.g., result of several coordinate transformations
	computedValue1 := baseValue + 1e-9 // Small computational error
	computedValue2 := baseValue + 5e-9 // Larger computational error
	computedValue3 := baseValue + 1e-7 // Clear difference

	// With epsilon = 1e-8, these should be the expected results:

	// Very small error should be considered equal
	if !FloatEqual(baseValue, computedValue1) {
		t.Errorf("With epsilon=1e-8, expected %e and %e to be equal (computational noise)", baseValue, computedValue1)
	}

	// Moderate error should still be considered equal (accumulated precision loss)
	if !FloatEqual(baseValue, computedValue2) {
		t.Errorf("With epsilon=1e-8, expected %e and %e to be equal (accumulated error)", baseValue, computedValue2)
	}

	// Large error should be considered different
	if FloatEqual(baseValue, computedValue3) {
		t.Errorf("With epsilon=1e-8, expected %e and %e to be different", baseValue, computedValue3)
	}

	// Test with coordinate-scale values (typical for geographic data)
	lng1 := -122.4194155 // San Francisco longitude
	lng2 := lng1 + 1e-9  // Tiny GPS precision difference
	lng3 := lng1 + 1e-7  // Meaningful geographic difference (~1cm at equator)

	if !FloatEqual(lng1, lng2) {
		t.Errorf("GPS coordinates with noise should be considered equal: %f vs %f", lng1, lng2)
	}

	if FloatEqual(lng1, lng3) {
		t.Errorf("GPS coordinates with meaningful difference should not be equal: %f vs %f", lng1, lng3)
	}
}
