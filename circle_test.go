package geojson

import "testing"

func TestCircle(t *testing.T) {
	expectJSON(t,
		`{"type":"Feature","geometry":{"type":"Point","coordinates":[-112,33]},"properties":{"type":"Circle","radius":"5000"}}`,
		`{"type":"Feature","geometry":{"type":"Point","coordinates":[-112,33]},"properties":{"type":"Circle","radius":5000,"radius_units":"m"}}`,
	)
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
	expect(t, g.Contains(PO(-112.26, 33.51)))
}
