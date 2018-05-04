package geojson

import (
	"github.com/tidwall/tile38/pkg/geojson/geo"
)

func SegmentIntersectsCircle(start, end, center Position, meters float64) bool {
	// These are faster checks.  If they succeed there's no need do complicate things.
	if center.DistanceTo(start) <= meters {
		return true
	}
	if center.DistanceTo(end) <= meters {
		return true
	}

	// Distance between start and end
	l := geo.DistanceTo(start.Y, start.X, end.Y, end.X)

	// Unit direction vector
	dx := (end.X - start.X) / l
	dy := (end.Y - start.Y) / l

	// Point of the line closest to the center
	t := dx * (center.X - start.X) + dy * (center.Y - start.Y)
	px := t * dx + start.X
	py := t * dy + start.Y
	if px < start.X || px > end.X || py < start.Y || py > end.Y {
		// closest point is outside the segment
		return false
	}

	// Distance from the closest point to the center
	return geo.DistanceTo(center.Y, center.X, py, px) <= meters
}
