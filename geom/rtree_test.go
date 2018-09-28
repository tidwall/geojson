package geom

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/tidwall/cities"
)

func TestRTree(t *testing.T) {
	// N := 100000
	var rects []rTreeRect
	for _, city := range cities.Cities {
		rects = append(rects, rTreeRect{
			min: [2]float64{city.Longitude, city.Latitude},
			max: [2]float64{city.Longitude, city.Latitude},
		})
	}
	var tr rTree
	start := time.Now()
	tr.load(rects)
	println(time.Since(start).String())
	ioutil.WriteFile("out.svg", []byte(tr.svg()), 0666)
}
