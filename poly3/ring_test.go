package poly

import (
	"testing"
)

func TestRingContainsRing(t *testing.T) {
	other := Ring{{2, 2}, {8, 2}, {8, 8}, {2, 8}, {2, 2}}
	expect(t, rectangle.ContainsRing(rectangle))
	expect(t, rectangle.ContainsRing(other))
	expect(t, !other.ContainsRing(rectangle))
	tbasic := func(delta float64, ex bool) {
		expect(t, ex == rectangle.ContainsRing(other.move(-delta, -delta)))
		expect(t, ex == rectangle.ContainsRing(other.move(-delta, 0)))
		expect(t, ex == rectangle.ContainsRing(other.move(-delta, delta)))
		expect(t, ex == rectangle.ContainsRing(other.move(0, -delta)))
		expect(t, ex == rectangle.ContainsRing(other.move(0, delta)))
		expect(t, ex == rectangle.ContainsRing(other.move(delta, -delta)))
		expect(t, ex == rectangle.ContainsRing(other.move(delta, 0)))
		expect(t, ex == rectangle.ContainsRing(other.move(delta, delta)))
	}
	tbasic(0, true)
	tbasic(1, true)
	tbasic(2, true)
	tbasic(3, false)
	tbasic(4, false)
	tbasic(5, false)
}

func TestRingContainsRingConvex(t *testing.T) {
	small := Ring{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}
	expect(t, !concave1.ContainsRing(small))
	expect(t, concave1.ContainsRing(small.move(1, 0)))
	expect(t, concave1.ContainsRing(small.move(2, 0)))
	expect(t, !concave1.ContainsRing(small.move(-1, 0)))
	expect(t, !concave1.ContainsRing(small.move(-2, 0)))
	big := Ring{{3, 3}, {7, 3}, {7, 7}, {3, 7}, {3, 3}}
	for x := -4.0; x <= 4; x++ {
		for y := -4.0; y <= 4; y++ {
			expect(t, !bowtie.ContainsRing(big.move(x, y)))
		}
	}
	expect(t, bowtie.ContainsRing(small))
	expect(t, bowtie.ContainsRing(small.move(-1, 0)))
	expect(t, bowtie.ContainsRing(small.move(+1, 0)))
	expect(t, !bowtie.ContainsRing(small.move(0, -1)))
	expect(t, !bowtie.ContainsRing(small.move(0, +1)))
}

func TestRingIntersectsRing(t *testing.T) {
	center := Ring{{2, 2}, {8, 2}, {8, 8}, {2, 8}, {2, 2}}
	expect(t, rectangle.IntersectsRing(rectangle))
	expect(t, rectangle.IntersectsRing(center))
	tbasic := func(delta float64, ex bool) {
		expect(t, ex == rectangle.IntersectsRing(center.move(-delta, -delta)))
		expect(t, ex == rectangle.IntersectsRing(center.move(-delta, 0)))
		expect(t, ex == rectangle.IntersectsRing(center.move(-delta, delta)))
		expect(t, ex == rectangle.IntersectsRing(center.move(0, -delta)))
		expect(t, ex == rectangle.IntersectsRing(center.move(0, delta)))
		expect(t, ex == rectangle.IntersectsRing(center.move(delta, -delta)))
		expect(t, ex == rectangle.IntersectsRing(center.move(delta, 0)))
		expect(t, ex == rectangle.IntersectsRing(center.move(delta, delta)))
		expect(t, ex == center.move(-delta, -delta).IntersectsRing(rectangle))
		expect(t, ex == center.move(-delta, 0).IntersectsRing(rectangle))
		expect(t, ex == center.move(-delta, delta).IntersectsRing(rectangle))
		expect(t, ex == center.move(0, -delta).IntersectsRing(rectangle))
		expect(t, ex == center.move(0, delta).IntersectsRing(rectangle))
		expect(t, ex == center.move(delta, -delta).IntersectsRing(rectangle))
		expect(t, ex == center.move(delta, 0).IntersectsRing(rectangle))
		expect(t, ex == center.move(delta, delta).IntersectsRing(rectangle))
	}
	for i := 0.0; i < 8; i++ {
		tbasic(i, true)
	}
	tbasic(9, false)
	tbasic(10, false)
	big := Ring{{3, 3}, {7, 3}, {7, 7}, {3, 7}, {3, 3}}
	expect(t, bowtie.IntersectsRing(big))
	expect(t, !bowtie.IntersectsRing(big.move(8, 0)))
	expect(t, !bowtie.IntersectsRing(big.move(-8, 0)))
	expect(t, !bowtie.IntersectsRing(big.move(0, -5)))
	expect(t, !bowtie.IntersectsRing(big.move(0, 5)))
	expect(t, bowtie.IntersectsRing(big.move(0, -4)))
	expect(t, bowtie.IntersectsRing(big.move(0, 4)))
	expect(t, big.IntersectsRing(bowtie))
	expect(t, !big.move(8, 0).IntersectsRing(bowtie))
	expect(t, !big.move(-8, 0).IntersectsRing(bowtie))
	expect(t, !big.move(0, -5).IntersectsRing(bowtie))
	expect(t, !big.move(0, 5).IntersectsRing(bowtie))
	expect(t, big.move(0, -4).IntersectsRing(bowtie))
	expect(t, big.move(0, 4).IntersectsRing(bowtie))
}

func algoRingIntersectsRing(a, b Ring, allowOn bool) bool {
	// test if any points from A are within B
	for _, point := range a {
		if pointInRing(point, b, allowOn) {
			return true
		}
	}
	// test if any points from B are within A
	for _, point := range b {
		if pointInRing(point, a, allowOn) {
			return true
		}
	}

	return false
}

func TestRingIntersectsRingExterior(t *testing.T) {
	small := Ring{{4, 4}, {6, 4}, {6, 6}, {4, 6}, {4, 4}}

	expect(t, algoRingIntersectsRing(small, small, true))
}
