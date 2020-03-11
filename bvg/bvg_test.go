package bvg

import "testing"

func TestGetLineDownloads(t *testing.T) {
	downloads, err := GetLineDownloads("20200229212100")
	if err != nil {
		t.Fatal(err)
	}
	for _, download := range downloads {
		t.Error(download)
	}
}
