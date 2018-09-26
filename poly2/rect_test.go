package poly

import "testing"

func TestRectContainsPoint(t *testing.T) {
	rect := R(10, 10, 20, 20)
	expect(t, rect.ContainsPoint(P(10, 10)))
	expect(t, rect.ContainsPoint(P(20, 20)))
	expect(t, rect.ContainsPoint(P(15, 15)))
	expect(t, !rect.ContainsPoint(P(25, 15)))
	expect(t, !rect.ContainsPoint(P(5, 15)))
	expect(t, !rect.ContainsPoint(P(15, 5)))
	expect(t, !rect.ContainsPoint(P(15, 25)))
}
func TestRectContainsRect(t *testing.T) {
	rect := R(10, 10, 20, 20)
	expect(t, rect.ContainsRect(rect))
	expect(t, rect.ContainsRect(R(20, 20, 20, 20)))
	expect(t, rect.ContainsRect(R(15, 15, 16, 16)))
	expect(t, !rect.ContainsRect(R(12, 12, 25, 15)))
	expect(t, !rect.ContainsRect(R(5, 15, 18, 18)))
	expect(t, !rect.ContainsRect(R(15, 5, 18, 18)))
	expect(t, !rect.ContainsRect(R(12, 12, 15, 25)))
}
func TestRectContainsLine(t *testing.T) {
	rect := R(10, 10, 20, 20)
	expect(t, !rect.ContainsLine(Line{}))
	expect(t, rect.ContainsLine(Line{P(10, 10)}))
	expect(t, rect.ContainsLine(Line{P(15, 15)}))
	expect(t, rect.ContainsLine(Line{P(15, 15), P(20, 20)}))
	expect(t, !rect.ContainsLine(Line{P(15, 15), P(20, 20), P(30, 30)}))
}
func TestRectContainsRing(t *testing.T) {
	rect := R(10, 10, 20, 20)
	expect(t, !rect.ContainsRing(Ring{}))
	expect(t, rect.ContainsRing(Ring{P(10, 10)}))
	expect(t, rect.ContainsRing(Ring{P(15, 15)}))
	expect(t, rect.ContainsRing(Ring{P(15, 15), P(20, 20)}))
	expect(t, !rect.ContainsRing(Ring{P(15, 15), P(20, 20), P(30, 30)}))
}
func TestRectContainsPolygon(t *testing.T) {
	rect := R(10, 10, 20, 20)
	expect(t, !rect.ContainsPolygon(Polygon{Exterior: Ring{}}))
	expect(t, rect.ContainsPolygon(Polygon{Exterior: Ring{P(10, 10)}}))
	expect(t, rect.ContainsPolygon(Polygon{Exterior: Ring{P(15, 15)}}))
	expect(t, rect.ContainsPolygon(Polygon{Exterior: Ring{P(15, 15), P(20, 20)}}))
	expect(t, !rect.ContainsPolygon(Polygon{Exterior: Ring{P(15, 15), P(20, 20), P(30, 30)}}))
}
func TestRectIntersectsPoint(t *testing.T) {
	rect := R(10, 10, 20, 20)
	expect(t, rect.IntersectsPoint(P(10, 10)))
	expect(t, rect.IntersectsPoint(P(20, 20)))
	expect(t, rect.IntersectsPoint(P(15, 15)))
	expect(t, !rect.IntersectsPoint(P(25, 15)))
	expect(t, !rect.IntersectsPoint(P(5, 15)))
	expect(t, !rect.IntersectsPoint(P(15, 5)))
	expect(t, !rect.IntersectsPoint(P(15, 25)))
}
func TestRectIntersectsRect(t *testing.T) {
	rect := R(10, 10, 20, 20)
	expect(t, rect.IntersectsRect(rect))
	expect(t, rect.IntersectsRect(R(20, 20, 20, 20)))
	expect(t, rect.IntersectsRect(R(15, 15, 16, 16)))
	expect(t, rect.IntersectsRect(R(12, 12, 25, 15)))
	expect(t, rect.IntersectsRect(R(5, 15, 18, 18)))
	expect(t, rect.IntersectsRect(R(15, 5, 18, 18)))
	expect(t, rect.IntersectsRect(R(12, 12, 15, 25)))
	expect(t, !rect.IntersectsRect(R(5, 5, 8, 40)))
	expect(t, !rect.IntersectsRect(R(25, 5, 28, 40)))
	expect(t, !rect.IntersectsRect(R(5, 5, 40, 8)))
	expect(t, !rect.IntersectsRect(R(5, 35, 40, 38)))
}
func TestRectIntersectsLine(t *testing.T) {
	rect := R(10, 10, 20, 20)
	expect(t, !rect.IntersectsLine(Line{}))
	expect(t, rect.IntersectsLine(Line{P(15, 15), P(18, 18)}))
	expect(t, rect.IntersectsLine(Line{P(5, 5), P(28, 28)}))
	expect(t, !rect.IntersectsLine(Line{P(5, 55), P(28, 28)}))
	expect(t, rect.IntersectsLine(Line{P(5, 55), P(28, 28), P(0, 10)}))
}
func TestRectIntersectsRing(t *testing.T) {
	rect := R(10, 10, 20, 20)
	expect(t, !rect.IntersectsRing(Ring{}))
	expect(t, rect.IntersectsRing(Ring{P(5, 5), P(18, 18)}))
	// expect(t, rect.IntersectsRing(Ring{P(5, 5), P(28, 28)}))
	// expect(t, !rect.IntersectsRing(Ring{P(5, 55), P(28, 28)}))
	// expect(t, rect.IntersectsRing(Ring{P(5, 55), P(28, 28), P(0, 10)}))
}
