package geojson

import "testing"

func TestCircle(t *testing.T) {
	g, err := Parse(`{
	"type":"Feature",
	"geometry":{"type":"Point","coordinates":[-112.2693,33.5123]},  
	"properties": { 
	  "type": "Circle",
	  "radius": 1000
	 }
  }`, nil)
	if err != nil {
		t.Fatal(err)
	}
	println(g.Contains(PO(-112.26, 33.51)))
}
