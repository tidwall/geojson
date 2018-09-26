package poly

func pointInRing(point Point, ring Ring, allowOnEdge bool) bool {
	in := false
	var a, b Point
	for i := 0; i < len(ring); i++ {
		a = ring[i]
		if i == len(ring)-1 {
			b = ring[0]
		} else {
			b = ring[i+1]
		}
		res := raycast(point, a, b)
		if res.on {
			return allowOnEdge
		}
		if res.in {
			in = !in
		}
	}
	return in
}
