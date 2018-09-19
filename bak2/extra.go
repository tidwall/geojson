package geojson

import (
	"unsafe"
)

type extra struct {
	z, m   bool
	coords []float64
}

func (ex *extra) dims() int {
	if ex != nil {
		if ex.z {
			if ex.m {
				return 2
			}
			return 1
		}
		if ex.m {
			return 1
		}
	}
	return 0
}

func (ex *extra) weight() int {
	if ex == nil {
		return 0
	}
	return int(unsafe.Sizeof(*ex)) + len(ex.coords)*8
}
