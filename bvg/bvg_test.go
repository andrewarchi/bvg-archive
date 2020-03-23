package bvg

import "testing"

func TestGetNetworkMaps(t *testing.T) {
	maps, err := GetNetworkMaps("")
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range maps {
		t.Error(m)
	}
}
