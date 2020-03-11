package bvg

import "testing"

func TestGetLineDownloads(t *testing.T) {
	links, err := GetLineDownloads("20200229212100")
	if err != nil {
		t.Fatal(err)
	}
	for _, link := range links {
		t.Error(link.Render())
	}
}
