package poly

func insideshpext(p Point, shape Ring, exterior bool) bool {
	// if len(shape) < 3 {
	// 	return false
	// }
	in := false
	for i := 0; i < len(shape); i++ {
		res := raycast(p, shape[i], shape[(i+1)%len(shape)])
		if res.on {
			return exterior
		}
		if res.in {
			in = !in
		}
	}
	return in
}
