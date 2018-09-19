package geojson

import (
	"encoding/json"
	"fmt"
)

func P(x, y float64) Position {
	return Position{X: x, Y: y}
}
func R(minX, minY, maxX, maxY float64) Rect {
	return Rect{Min: P(minX, minY), Max: P(maxX, maxY)}
}

func expect(exp, val interface{}) {
	if fmt.Sprintf("%v", exp) != fmt.Sprintf("%v", val) {
		panic(fmt.Sprintf("expected '%v', got '%v'", exp, val))
	}
}

func cleanJSON(data string) string {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		panic(err)
	}
	out, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func mustParseGeoJSON(data string) Object {
	obj, err := Parse(data)
	if err != nil {
		panic(err)
	}
	if cleanJSON(obj.JSON()) != cleanJSON(data) {
		panic(fmt.Sprintf("\nexp '%v'\ngot '%v'",
			cleanJSON(data), cleanJSON(obj.JSON())))
	}
	return obj
}
