package geom

import (
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	"github.com/tidwall/boxtree/d2"
	"github.com/tidwall/boxtree/res/tools"

	"github.com/tidwall/cities"
)

func TestRTree(t *testing.T) {
	// return
	// // N := 100000
	if true {
		for i := range cities.Cities {
			j := rand.Intn(i + 1)
			cities.Cities[i], cities.Cities[j] =
				cities.Cities[j], cities.Cities[i]
		}
	}

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

	var tr2 d2.BoxTree
	start = time.Now()
	for i, city := range cities.Cities {
		tr2.Insert(
			[]float64{city.Longitude, city.Latitude},
			[]float64{city.Longitude, city.Latitude},
			i,
		)
	}
	println(time.Since(start).String())
	ioutil.WriteFile("out2.svg", []byte(tools.SVG(&tr2)), 0666)

}
