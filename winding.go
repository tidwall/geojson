package geojson

import "github.com/tidwall/geojson/geometry"

// ringOrientation returns +1 for counterclockwise, -1 for clockwise, and 0 for
// degenerate rings (zero signed area or fewer than three distinct points).
// GeoJSON uses geographic coordinates with latitude (Y) increasing northward,
// so the standard Cartesian shoelace sign convention matches RFC 7946 §3.1.6.
// This is a planar computation: rings that cross the antimeridian (±180°
// longitude) may be classified with the wrong sign.
func ringOrientation(pts []geometry.Point) int {
	n := len(pts)
	if n < 3 {
		return 0
	}
	var area float64
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += pts[i].X*pts[j].Y - pts[j].X*pts[i].Y
	}
	if area > 0 {
		return 1
	}
	if area < 0 {
		return -1
	}
	return 0
}

func reverseRing(pts []geometry.Point) {
	for i, j := 0, len(pts)-1; i < j; i, j = i+1, j-1 {
		pts[i], pts[j] = pts[j], pts[i]
	}
}

// reverseExtraRange reverses the per-point extra values for one ring. The
// ring occupies len(ring) consecutive points starting at pointOffset inside
// ex.values, with each point holding dims floats.
func reverseExtraRange(ex *extra, pointOffset, nPoints int) {
	if ex == nil || ex.dims == 0 || nPoints < 2 {
		return
	}
	dims := int(ex.dims)
	base := pointOffset * dims
	end := base + nPoints*dims
	if base < 0 || end > len(ex.values) {
		return
	}
	for i, j := 0, nPoints-1; i < j; i, j = i+1, j-1 {
		for k := 0; k < dims; k++ {
			ex.values[base+i*dims+k], ex.values[base+j*dims+k] =
				ex.values[base+j*dims+k], ex.values[base+i*dims+k]
		}
	}
}

// applyRFC7946Winding enforces RFC 7946 §3.1.6 ring winding on a single
// polygon's coords (exterior at index 0, holes after). It honours the
// FixWinding and RequireWinding parse options.
func applyRFC7946Winding(
	coords [][]geometry.Point, ex *extra, opts *ParseOptions,
) error {
	if opts == nil || (!opts.FixWinding && !opts.RequireWinding) {
		return nil
	}
	offset := 0
	for i, ring := range coords {
		orient := ringOrientation(ring)
		wantCCW := i == 0
		if orient != 0 {
			compliant := (wantCCW && orient == 1) || (!wantCCW && orient == -1)
			if !compliant {
				if opts.RequireWinding {
					return errCoordinatesInvalid
				}
				if opts.FixWinding {
					reverseRing(ring)
					reverseExtraRange(ex, offset, len(ring))
				}
			}
		}
		offset += len(ring)
	}
	return nil
}
