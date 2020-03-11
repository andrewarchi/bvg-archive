package bvg

import (
	"bytes"

	"github.com/andrewarchi/internet-archive/dom"
	"github.com/andrewarchi/internet-archive/wayback"
	"golang.org/x/net/html"
)

func GetLineDownloads(timestamp string) ([]*dom.Node, error) {
	page, err := wayback.GetPage("https://www.bvg.de/de/Fahrinfo/Downloads/BVG-Liniennetz", timestamp)
	if err != nil {
		return nil, err
	}
	doc, err := html.Parse(bytes.NewReader(page))
	if err != nil {
		return nil, err
	}
	return (*dom.Node)(doc).FindClass("article__body").FindClass("download").FindClass("link-list").FindTagAll("li"), nil
}
