package bvg

import (
	"crypto/sha512"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/andrewarchi/bvg-archive/dom"
	"github.com/andrewarchi/bvg-archive/wayback"
	"golang.org/x/net/html/atom"
)

type Download struct {
	URL   string
	Title string
	Date  time.Time
}

func GetNetworkMaps(timestamp string) ([]Download, error) {
	page, err := wayback.GetPage("https://www.bvg.de/de/Fahrinfo/Downloads/BVG-Liniennetz", timestamp)
	if err != nil {
		return nil, err
	}
	defer page.Body.Close()
	doc, err := dom.Parse(page.Body)
	if err != nil {
		return nil, err
	}
	download := doc.FindClass("article__body").FindClass("download")
	links := download.FindClass("link-list").FindTagAtomAll(atom.Li)
	downloads := make([]Download, len(links))
	for i, link := range links {
		url, _ := link.FindTagAtom(atom.A).LookupAttr("href")
		title := link.FindAttr("class", "link-list__text").TextContent()

		date, _ := link.FindTagAtom(atom.Img).LookupAttr("alt")
		date = strings.TrimPrefix(date, "Aktualisiert am: ")
		var t time.Time
		if date != "" {
			t, err = time.Parse("02.01.2006", date)
			if err != nil {
				return nil, err
			}
		}

		downloads[i] = Download{url, title, t}
	}
	return downloads, nil
}

func SavePDF(resp *http.Response, pathPrefix string) error {
	filename, err := getFilename(resp.Header)
	if err != nil {
		return err
	}
	path := pathPrefix + filename
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}
	lastModified, err := getLastModified(resp.Header)
	if err != nil {
		return err
	}
	if !lastModified.IsZero() {
		return os.Chtimes(path, lastModified, lastModified)
	}
	return nil
}

func getFilename(header http.Header) (string, error) {
	cd := header.Get("Content-Disposition")
	if cd == "" {
		return "", nil
	}
	_, params, err := mime.ParseMediaType(cd)
	if err != nil {
		return "", err
	}
	return params["filename"], nil
}

func getLastModified(header http.Header) (time.Time, error) {
	mod := header.Get("X-Archive-Orig-Last-Modified")
	if mod == "" {
		mod = header.Get("Last-Modified")
		if mod == "" {
			return time.Time{}, nil
		}
	}
	return time.Parse(time.RFC1123, mod)
}

func hash(r io.Reader) (sum [sha512.Size]byte, err error) {
	h := sha512.New()
	_, err = io.Copy(h, r)
	if err == nil {
		h.Sum(sum[:])
	}
	return
}
