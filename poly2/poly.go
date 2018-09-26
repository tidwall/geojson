package poly

// // Point ...
// type Point struct {
// 	X, Y float64
// }

// // Ring ...
// type Ring []Point

// func ringIsConvex(ring Ring) bool {
// 	if len(ring) < 3 {
// 		return false
// 	}
// 	var convex := true
// 	var dir int
// 	var a, b, c Point
// 	for i := 0; i < len(ring); i++ {
// 		a = ring[i]
// 		if i == len(ring)-1 {
// 			b = ring[0]
// 			c = ring[1]
// 		} else if i == len(ring)-2 {
// 			b = ring[i+1]
// 			c = ring[0]
// 		} else {
// 			b = ring[i+1]
// 			c = ring[i+2]
// 		}
// 		dx1 := b.X - a.X
// 		dy1 := b.Y - a.Y
// 		dx2 := c.X - b.X
// 		dy2 := c.Y - b.Y
// 		zCrossProduct := dx1*dy2 - dy1*dx2
// 		if dir == 0 {
// 			if zCrossProduct < 0 {
// 				dir = -1
// 			} else if zCrossProduct > 0 {
// 				dir = 1
// 			}
// 		} else if zCrossProduct < 0 {
// 			if dir == 1 {
// 				convex = false
// 				//return false
// 			}
// 		} else if zCrossProduct > 0 {
// 			if dir == -1 {
// 				convex = false
// 				//return false
// 			}
// 		}
// 	}
// 	return convex
// }

// func pointInRing(point Point, ring Ring, exterior bool) bool {
// 	return false
// 	// if len(ring) < 3 {
// 	// 	return false
// 	// }
// 	// var lastZ float64
// 	// var a, b, c Point
// 	// in := false
// 	// for i := 0; i < len(ring); i++ {
// 	// 	a = ring[i]
// 	// 	if i == len(ring)-1 {
// 	// 		b = ring[0]
// 	// 		c = ring[1]
// 	// 	} else if i == len(ring)-2 {
// 	// 		b = ring[i+1]
// 	// 		c = ring[0]
// 	// 	} else {
// 	// 		b = ring[i+1]
// 	// 		c = ring[i+2]
// 	// 	}

// 	// 	dx1 := b.X - a.X
// 	// 	dy1 := b.Y - a.Y
// 	// 	dx2 := c.X - b.X
// 	// 	dy2 := c.Y - b.Y
// 	// 	zcrossproduct := dx1*dy2 - dy1*dx2
// 	// 	println(zcrossproduct)

// 	// 	//, b, c = , ring[(i+1)%len(ring)], ring[(i+2)%len(ring)]
// 	// 	_, _, _ = a, b, c
// 	// 	// res := raycast(point, a, b)
// 	// 	// if res.on {
// 	// 	// 	return exterior
// 	// 	// }
// 	// 	// if res.in {
// 	// 	// 	in = !in
// 	// 	// }
// 	// }
// 	// return in
// }

// // func segmentAngle(a, b Point) float64 {
// // 	return math.Atan2(b.Y-a.Y, b.X-a.X)
// // }
