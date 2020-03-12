package bvg

import (
	"strings"

	"github.com/andrewarchi/bvg-archive/dom"
	"github.com/andrewarchi/bvg-archive/wayback"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Download struct {
	Title string
	URL   string
	Date  string
}

func GetLineDownloads(timestamp string) ([]Download, error) {
	page, err := wayback.GetPage("https://www.bvg.de/de/Fahrinfo/Downloads/BVG-Liniennetz", timestamp)
	if err != nil {
		return nil, err
	}
	defer page.Body.Close()
	doc, err := html.Parse(page.Body)
	if err != nil {
		return nil, err
	}
	download := (*dom.Node)(doc).FindClass("article__body").FindClass("download")
	links := download.FindClass("link-list").FindTagAtomAll(atom.Li)
	downloads := make([]Download, len(links))
	for i, link := range links {
		title := link.FindAttr("class", "link-list__text").TextContent()
		url, _ := link.FindTagAtom(atom.A).LookupAttr("href")
		date, _ := link.FindTagAtom(atom.Img).LookupAttr("alt")
		date = strings.TrimPrefix(date, "Aktualisiert am: ")
		downloads[i] = Download{title, url, date}
	}
	return downloads, nil
}
