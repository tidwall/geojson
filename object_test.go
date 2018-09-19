package geojson

import "testing"

func TestObject(t *testing.T) {
	expectJSON(t, "", errDataInvalid)
	expectJSON(t, string([]byte{0, 1, 2, 3}), errDataInvalid)
	expectJSON(t, string([]byte{' ', 0}), errDataInvalid)
	expectJSON(t, `{}`, errTypeMissing)
	expectJSON(t, `{"}`, errDataInvalid)
	expectJSON(t, `{"type":null}`, errTypeInvalid)
	expectJSON(t, `{"type":"Square"}`, errTypeInvalid)
}
