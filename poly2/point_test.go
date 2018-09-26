package poly

import "testing"

func TestPointContainsPoint(t *testing.T) {
	expect(t, P(10, 10).ContainsPoint(P(10, 10)))
	expect(t, !P(10, 10).ContainsPoint(P(11, 11)))
}
func TestPointContainsRect(t *testing.T) {
	expect(t, P(10, 10).ContainsRect(R(10, 10, 10, 10)))
	expect(t, !P(10, 10).ContainsRect(R(10, 10, 20, 20)))
}
func TestPointContainsLine(t *testing.T) {
	lineA := Line{P(10, 10), P(10, 10)}
	lineB := Line{P(0, 0), P(20, 20)}
	lineC := Line{P(10, 0), P(30, 20)}
	expect(t, !P(10, 10).ContainsLine(nil))
	expect(t, P(10, 10).ContainsLine(lineA))
	expect(t, !P(10, 10).ContainsLine(lineB))
	expect(t, !P(10, 10).ContainsLine(lineC))
}
func TestPointContainsRing(t *testing.T) {
	ringA := Ring{P(10, 10), P(10, 10), P(10, 10)}
	ringB := Ring{P(0, 0), P(20, 20), P(30, 0)}
	expect(t, !P(10, 10).ContainsRing(nil))
	expect(t, P(10, 10).ContainsRing(ringA))
	expect(t, !P(10, 10).ContainsRing(ringB))
}
func TestPointContainsPolygon(t *testing.T) {
	ringA := Ring{P(10, 10), P(10, 10), P(10, 10)}
	ringB := Ring{P(0, 0), P(20, 20), P(30, 0)}
	expect(t, !P(10, 10).ContainsPolygon(Polygon{}))
	expect(t, P(10, 10).ContainsPolygon(Polygon{Exterior: ringA}))
	expect(t, !P(10, 10).ContainsPolygon(Polygon{Exterior: ringB}))

}
func TestPointIntersectsPoint(t *testing.T) {
	expect(t, P(10, 10).IntersectsPoint(P(10, 10)))
	expect(t, !P(10, 10).IntersectsPoint(P(11, 11)))
}
func TestPointIntersectsRect(t *testing.T) {
	expect(t, P(10, 10).IntersectsRect(R(10, 10, 10, 10)))
	expect(t, P(10, 10).IntersectsRect(R(10, 10, 20, 20)))
	expect(t, P(20, 20).IntersectsRect(R(10, 10, 20, 20)))
	expect(t, P(15, 15).IntersectsRect(R(10, 10, 20, 20)))
	expect(t, !P(0, 0).IntersectsRect(R(10, 10, 20, 20)))
}
func TestPointIntersectsLine(t *testing.T) {
	lineA := Line{P(10, 10), P(10, 10), P(10, 10)}
	lineB := Line{P(0, 0), P(20, 20), P(40, 0)}
	expect(t, !P(10, 10).IntersectsLine(nil))
	expect(t, P(10, 10).IntersectsLine(lineA))
	expect(t, P(10, 10).IntersectsLine(lineB))
	expect(t, P(15, 15).IntersectsLine(lineB))
	expect(t, P(20, 20).IntersectsLine(lineB))
	expect(t, P(10, 10).IntersectsLine(lineB))
	expect(t, P(30, 10).IntersectsLine(lineB))
	expect(t, !P(20, 10).IntersectsLine(lineB))
}
func TestPointIntersectsRing(t *testing.T) {
	ringB := Ring{P(0, 0), P(20, 20), P(40, 0)}
	expect(t, !P(10, 10).IntersectsRing(nil))
	expect(t, P(10, 10).IntersectsRing(ringB))
	expect(t, P(10, 10).IntersectsRing(ringB))
	expect(t, P(20, 10).IntersectsRing(ringB))
	expect(t, !P(20, 21).IntersectsRing(ringB))
	expect(t, !P(20, -10).IntersectsRing(ringB))
}
func TestPointIntersectsPolygon(t *testing.T) {
	poly := Polygon{
		Ring{{10, 10}, {20, 10}, {20, 20}, {10, 20}, {10, 10}},
		[]Ring{
			Ring{{12, 12}, {14, 12}, {14, 14}, {12, 14}, {12, 12}},
			Ring{{16, 16}, {18, 16}, {18, 18}, {16, 18}, {16, 16}},
		},
	}
	expect(t, !P(10, 10).IntersectsPolygon(Polygon{}))
	expect(t, P(10, 10).IntersectsPolygon(poly))
	expect(t, P(20, 10).IntersectsPolygon(poly))
	expect(t, P(15, 10).IntersectsPolygon(poly))
	expect(t, P(15, 20).IntersectsPolygon(poly))
	expect(t, !P(13, 13).IntersectsPolygon(poly))
	expect(t, !P(17, 17).IntersectsPolygon(poly))
	expect(t, P(19, 17).IntersectsPolygon(poly))
}
