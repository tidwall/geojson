package geom

import "testing"

func TestCircleNewCircle(t *testing.T) {
	circle := NewCircle(P(-112, 33), 1000, 2)
	expect(t, circle.ContainsPoint(P(-112, 33)))
}
