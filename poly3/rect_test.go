package poly

import "testing"

func TestRectContainsRect(t *testing.T) {
	big := Rect{Point{0, 0}, Point{10, 10}}
	small := Rect{Point{4, 4}, Point{6, 6}}
	expect(t, big.ContainsRect(big))
	expect(t, big.ContainsRect(small))
	expect(t, big.ContainsRect(small))
	expect(t, !small.ContainsRect(big))
	expect(t, !big.ContainsRect(big.move(-1, -1)))
	expect(t, !big.ContainsRect(big.move(-1, 0)))
	expect(t, !big.ContainsRect(big.move(-1, 1)))
	expect(t, !big.ContainsRect(big.move(0, -1)))
	expect(t, !big.ContainsRect(big.move(0, 1)))
	expect(t, !big.ContainsRect(big.move(1, -1)))
	expect(t, !big.ContainsRect(big.move(1, 0)))
	expect(t, !big.ContainsRect(big.move(1, 1)))
}
func TestRectIntersectsRect(t *testing.T) {
	big := Rect{Point{0, 0}, Point{10, 10}}
	small := Rect{Point{4, 4}, Point{6, 6}}
	expect(t, big.IntersectsRect(big))
	expect(t, big.IntersectsRect(small))
	expect(t, big.IntersectsRect(small))
	expect(t, small.IntersectsRect(big))
	expect(t, big.IntersectsRect(big.move(-1, -1)))
	expect(t, big.IntersectsRect(big.move(-1, 0)))
	expect(t, big.IntersectsRect(big.move(-1, 1)))
	expect(t, big.IntersectsRect(big.move(0, -1)))
	expect(t, big.IntersectsRect(big.move(0, 1)))
	expect(t, big.IntersectsRect(big.move(1, -1)))
	expect(t, big.IntersectsRect(big.move(1, 0)))
	expect(t, big.IntersectsRect(big.move(1, 1)))
	expect(t, big.IntersectsRect(big.move(-10, 0)))
	expect(t, big.IntersectsRect(big.move(10, 0)))
	expect(t, big.IntersectsRect(big.move(-10, -10)))
	expect(t, big.IntersectsRect(big.move(10, 10)))
	expect(t, big.IntersectsRect(big.move(0, -10)))
	expect(t, big.IntersectsRect(big.move(0, 10)))
	expect(t, !big.IntersectsRect(big.move(-11, 0)))
	expect(t, !big.IntersectsRect(big.move(11, 0)))
	expect(t, !big.IntersectsRect(big.move(-11, -11)))
	expect(t, !big.IntersectsRect(big.move(11, 11)))
	expect(t, !big.IntersectsRect(big.move(0, -11)))
	expect(t, !big.IntersectsRect(big.move(0, 11)))
}
