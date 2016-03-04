package geohash

import (
	"fmt"
	"testing"
)

func fixed(f float64, d int) string {
	return fmt.Sprintf(fmt.Sprintf("%%0.%df", d), f)
}

func TestABC(t *testing.T) {
	lat, lon := 33.52345123, -115.512345123
	hash, err := Encode(lat, lon, 32)
	if err != nil {
		t.Fatal(err)
	}
	lat2, lon2, err := Decode(hash)
	if err != nil {
		t.Fatal(err)
	}
	if fixed(lat, 10) != fixed(lat2, 10) || fixed(lon, 10) != fixed(lon2, 10) {
		t.Fatalf("bad geohash %v,%v %v,%v", lat, lon, lat2, lon2)
	}
}
