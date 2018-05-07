package geojson

// String is a not a geojson object, but just a string
type String string

func (s String) bboxPtr() *BBox {
	return nil
}
func (s String) hasPositions() bool {
	return false
}

// WithinBBox detects if the object is fully contained inside a bbox.
func (s String) WithinBBox(bbox BBox) bool {
	return false
}

// IntersectsBBox detects if the object intersects a bbox.
func (s String) IntersectsBBox(bbox BBox) bool {
	return false
}

// Within detects if the object is fully contained inside another object.
func (s String) Within(o Object) bool {
	return false
}

// WithinCircle detects if the object is fully contained inside a circle.
func (s String) WithinCircle(center Position, meters float64) bool {
	return false
}

// Intersects detects if the object intersects another object.
func (s String) Intersects(o Object) bool {
	return false
}
func (s String) IntersectsCircle(center Position, meters float64) bool {
	return false
}

// Nearby detects if the object is nearby a position.
func (s String) Nearby(center Position, meters float64) bool {
	return false
}

// CalculatedBBox is exterior bbox containing the object.
func (s String) CalculatedBBox() BBox {
	return BBox{}
}

// CalculatedPoint is a point representation of the object.
func (s String) CalculatedPoint() Position {
	return Position{}
}

func (s String) appendJSON(json []byte) []byte {
	b, _ := s.MarshalJSON()
	return append(json, b...)
}

// JSON is the json representation of the object. This might not be exactly the same as the original.
func (s String) JSON() string {
	return string(s.appendJSON(nil))
}

// String returns a string representation of the object. This might be JSON or something else.
func (s String) String() string {
	return string(s)
}

// IsGeometry return true if the object is a geojson geometry object. false if it something else.
func (s String) IsGeometry() bool {
	return false
}

// Clip returns the object obtained by clipping this object by a bbox.
func (s String) Clipped(bbox BBox) Object {
	return s
}

// Bytes is the bytes representation of the object.
func (s String) Bytes() []byte {
	return []byte(s.String())
}

// PositionCount return the number of coordinates.
func (s String) PositionCount() int {
	return 0
}

// Weight returns the in-memory size of the object.
func (s String) Weight() int {
	return len(s)
}

// MarshalJSON allows the object to be encoded in json.Marshal calls.
func (s String) MarshalJSON() ([]byte, error) {
	return jsonMarshalString(string(s)), nil
}

// Geohash converts the object to a geohash value.
func (s String) Geohash(precision int) (string, error) {
	return "", nil
}

// IsBBoxDefined returns true if the object has a defined bbox.
func (s String) IsBBoxDefined() bool {
	return false
}
