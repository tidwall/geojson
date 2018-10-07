package geojson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/tidwall/pretty"
)

func expectJSON(t testing.TB, data string, exp error) Object {
	if t != nil {
		t.Helper()
	}
	obj, err := Parse(data)
	if err != exp {
		if t == nil {
			panic(fmt.Sprintf("expected '%v', got '%v'", exp, err))
		} else {
			t.Fatalf("expected '%v', got '%v'", exp, err)
		}
	}
	return obj
}

func cleanJSON(data string) string {
	var v interface{}
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		panic(err)
	}
	dst, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	opts := *pretty.DefaultOptions
	opts.Width = 99999999
	return string(pretty.PrettyOptions(dst, &opts))
}

func testGeoJSONFile(t testing.TB, path string) Object {
	t.Helper()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	obj, err := Parse(string(data))
	if err != nil {
		t.Fatal(err)
	}
	orgJSON := cleanJSON(string(data))
	newJSON := cleanJSON(string(obj.AppendJSON(nil)))
	if orgJSON != newJSON {
		var ln int
		var col int
		for i := 0; i < len(orgJSON) && i < len(newJSON); i++ {
			if orgJSON[i] != newJSON[i] {
				break
			}
			if orgJSON[i] == '\n' {
				ln++
				col = 0
			} else {
				col++
			}
		}
		tpath1 := "/tmp/org.json"
		tpath2 := "/tmp/new.json"
		ioutil.WriteFile(tpath1, []byte(orgJSON), 0666)
		ioutil.WriteFile(tpath2, []byte(newJSON), 0666)
		t.Fatalf("%v (ln: %d, col: %d)\nfile://%s\nfile://%s",
			filepath.Base(path), ln, col, tpath1, tpath2)
	}
	return obj
}

func expect(t *testing.T, v bool) {
	t.Helper()
	if !v {
		t.Fatal("invalid expectation")
	}
}

func TestGeoJSON(t *testing.T) {
	fis, err := ioutil.ReadDir("test_files")
	if err != nil {
		t.Fatal(err)
	}
	for _, fi := range fis {
		testGeoJSONFile(t, filepath.Join("test_files", fi.Name()))
	}
}

func BenchmarkFeature(t *testing.B) {
	var r Object = R(0, 0, 20, 20)
	var p Object = P(10, 10)
	p = expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]}}`, nil)
	for i := 0; i < t.N; i++ {
		if !r.Intersects(p) {
			t.Fatal("bad")
		}
	}
}
func BenchmarkPosition(t *testing.B) {
	var r Object = R(0, 0, 20, 20)
	var p Object = P(10, 10)
	for i := 0; i < t.N; i++ {
		if !r.Intersects(p) {
			t.Fatal("bad")
		}
	}
}
