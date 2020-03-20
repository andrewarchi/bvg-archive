package bvg

import (
	"crypto/sha512"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/andrewarchi/bvg-archive/dom"
	"github.com/andrewarchi/bvg-archive/wayback"
	"golang.org/x/net/html/atom"
)

var illegal = regexp.MustCompile(`[\0-\x1f:?"*/\\<>|]`)

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

func SaveFile(resp *http.Response, pathPrefix string) error {
	filename, err := getFilename(resp.Header)
	if err != nil {
		return err
	}
	lastModified, err := getLastModified(resp.Header)
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
	if !lastModified.IsZero() {
		return os.Chtimes(path, lastModified, lastModified)
	}
	return nil
}

func SaveAllVersions(url, dir string) error {
	dir = filepath.Join(dir, sanitizeFilename(url))
	if err := os.MkdirAll(dir, 0o777); err != nil {
		return err
	}

	t := time.Now().Format("20060102150405")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if err := SaveFile(resp, filepath.Join(dir, t+"live_")); err != nil {
		return err
	}

	timemap, err := wayback.GetTimeMap(url)
	if err != nil {
		return err
	}
	for _, entry := range timemap {
		resp, err := wayback.GetPage(url, entry.Timestamp)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if err := SaveFile(resp, filepath.Join(dir, entry.Timestamp+"_")); err != nil {
			return err
		}
	}
	return nil
}

func sanitizeFilename(filename string) string {
	filename = strings.TrimPrefix(filename, "https://")
	filename = strings.TrimPrefix(filename, "http://")
	filename = strings.TrimPrefix(filename, "www.")
	filename = illegal.ReplaceAllString(filename, "_")
	return filename
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
	if t, err := time.Parse(time.RFC1123, mod); err == nil {
		return t, nil
	}
	return time.Parse(time.RFC850, mod)
}

func hash(r io.Reader) (sum [sha512.Size]byte, err error) {
	h := sha512.New()
	_, err = io.Copy(h, r)
	if err == nil {
		h.Sum(sum[:])
	}
	return
}
