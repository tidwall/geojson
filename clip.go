package geojson

// Cohen-Sutherland Line Clipping
// https://www.cs.helsinki.fi/group/goa/viewing/leikkaus/lineClip.html
func ClipSegment(start, end Position, bbox BBox) (resStart, resEnd Position, rejected bool) {
	startCode := getCode(bbox, start)
	endCode := getCode(bbox, end)

	if (startCode | endCode) == 0 {
		// trivially accept
		resStart, resEnd = start, end
	} else if (startCode & endCode) != 0 {
		// trivially reject
		rejected = true
	} else if startCode != 0 {
		// start is outside. get new start.
		newStart := intersect(bbox, startCode, start, end)
		resStart, resEnd, rejected = ClipSegment(newStart, end, bbox)
	} else {
		// end is outside. get new end.
		newEnd := intersect(bbox, endCode, start, end)
		resStart, resEnd, rejected = ClipSegment(start, newEnd, bbox)
	}

	return
}

// Sutherland-Hodgman Polygon Clipping
// https://www.cs.helsinki.fi/group/goa/viewing/leikkaus/intro2.html
func ClipRing(ring[] Position, bbox BBox) (resRing []Position) {

	if len(ring) < 4 {
		// under 4 elements this is not a polygon ring!
		return
	}

	var edge uint8
	var inside, prevInside bool
	var prev Position

	for edge = 1; edge <= 8; edge *=2 {
		prev = ring[len(ring) - 2]
		prevInside = (getCode(bbox, prev) & edge) == 0

		for _, p := range ring {

			inside = (getCode(bbox, p) & edge) == 0

			if prevInside && inside {
				// Staying inside
				resRing = append(resRing, p)
			} else if prevInside && !inside {
				// Leaving
				resRing = append(resRing, intersect(bbox, edge, prev, p))
			} else if !prevInside && inside {
				// Entering
				resRing = append(resRing, intersect(bbox, edge, prev, p))
				resRing = append(resRing, p)
			} else  {
				// Staying outside
			}

			prev, prevInside = p, inside
		}

		if resRing[0] != resRing[len(resRing)-1] {
			resRing = append(resRing, resRing[0])
		}
		ring, resRing = resRing, []Position{}
	}

	resRing = ring
	return
}


func getCode(bbox BBox, point Position) (code uint8) {
	code = 0

	if point.X < bbox.Min.X {
		code |= 1  // left
	} else if point.X > bbox.Max.X {
		code |= 2  // right
	}

	if point.Y < bbox.Min.Y {
		code |= 4  // bottom
	} else if point.Y > bbox.Max.Y {
		code |= 8  // top
	}

	return
}


func intersect(bbox BBox, code uint8, start, end Position) (new Position) {
	if (code & 8) != 0 {  // top
		new = Position{
			X: start.X + (end.X - start.X) * (bbox.Max.Y - start.Y) / (end.Y - start.Y),
			Y: bbox.Max.Y,
		}
	} else if (code & 4) != 0 {  // bottom
		new = Position{
			X: start.X + (end.X - start.X) * (bbox.Min.Y - start.Y) / (end.Y - start.Y),
			Y: bbox.Min.Y,
		}
	} else if (code & 2) != 0 {  //right
		new = Position{
			X: bbox.Max.X,
			Y: start.Y + (end.Y - start.Y) * (bbox.Max.X - start.X) / (end.X - start.X),
		}
	} else if (code & 1) != 0 {  // left
		new = Position{
			X: bbox.Min.X,
			Y: start.Y + (end.Y - start.Y) * (bbox.Min.X - start.X) / (end.X - start.X),
		}
	} else {  // should not call intersect with the zero code
	}

	return
}
