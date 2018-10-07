// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package geos

import "testing"

func TestCircleNewCircle(t *testing.T) {
	circle := NewCircle(P(-112, 33), 1000, 2)
	expect(t, circle.ContainsPoint(P(-112, 33)))
}
