package geojson

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/tidwall/geojson/geos"
)

func parseCollection(t *testing.T, data string, index bool) Collection {
	t.Helper()
	opts := *DefaultParseOptions
	if index {
		opts.IndexChildren = 1
	} else {
		opts.IndexChildren = 0
	}
	g, err := Parse(string(data), &opts)
	if err != nil {
		t.Fatal(err)
	}
	c, ok := g.(Collection)
	if !ok {
		t.Fatal("not searchable")
	}
	return c
}

func testCollectionSubset(
	t *testing.T, json string, subsetRect geos.Rect,
) []Collection {
	t.Helper()
	expectJSON(t, json, nil)
	start := time.Now()
	c := parseCollection(t, json, true)
	dur := time.Since(start)
	fmt.Printf("%d children in %s\n", len(c.Children()), dur)
	cs := []Collection{
		parseCollection(t, json, false),
		parseCollection(t, json, true),
	}
	var lastSubset string
	for i, c := range cs {
		var children []Object
		start := time.Now()
		c.Search(subsetRect, func(child Object) bool {
			children = append(children, child)
			return true
		})
		dur := time.Since(start)
		if i == 0 {
			fmt.Printf("Simple:  %v\n", dur)
		} else {
			fmt.Printf("Indexed: %v\n", dur)
		}
		var childrenJSONs []string
		for _, child := range children {
			childrenJSONs = append(childrenJSONs, string(child.AppendJSON(nil)))
		}
		sort.Strings(childrenJSONs)
		subset := `{"type":"GeometryCollection","geometries":[` +
			strings.Join(childrenJSONs, ",") + `]}`
		if i > 0 {
			if subset != lastSubset {
				t.Fatal("mismatch")
			}
		}
		lastSubset = subset
	}
	return cs

}

func TestCollectionBostonSubset(t *testing.T) {
	data, err := ioutil.ReadFile("test_files/boston_subset.geojson")
	expect(t, err == nil)
	cs := testCollectionSubset(t,
		string(data),
		R(-71.474046, 42.492479, -71.466321, 42.497415),
	)
	// expect(t, cs[0].(Object).Intersects(PO(-71.46723, 42.49432)))
	// expect(t, cs[1].(Object).Intersects(PO(-71.46723, 42.49432)))
	expect(t, !cs[0].(*FeatureCollection).Intersects(PO(-71.46713, 42.49431)))
	// expect(t, !cs[1].(Object).Intersects(PO(-71.46713, 42.49431)))

}

func TestCollectionUSASubset(t *testing.T) {
	data, err := ioutil.ReadFile("test_files/usa.geojson")
	expect(t, err == nil)
	testCollectionSubset(t,
		string(data),
		R(-90, -45, 90, 45),
	)
}
