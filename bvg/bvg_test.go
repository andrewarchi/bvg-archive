package bvg

import (
	"testing"

	"github.com/andrewarchi/internet-archive/dom"
)

func TestGetLineDownloads(t *testing.T) {
	links, err := GetLineDownloads("20200229212100")
	if err != nil {
		t.Fatal(err)
	}
	t.Error(dom.RenderNode(links))
}
