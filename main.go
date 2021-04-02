package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/andrewarchi/bvg-archive/bvg"
	"github.com/andrewarchi/urlhero/ia"
)

func main() {
	// Usage: bvg-archive [dir]

	dir := "files"
	if len(os.Args) >= 2 {
		dir = os.Args[1]
	}
	urls, titles, err := getNetworkMapURLs(false)
	if err != nil {
		exit(err)
	}
	for i, url := range urls {
		id := strings.TrimPrefix(url, "/de/index.php?section=downloads&cmd=58&download=")
		out := filepath.Join(dir, bvg.SanitizeFilename(id+" "+titles[i]))
		if err := bvg.SaveAllVersions("https://www.bvg.de"+url, out); err != nil {
			exit(err)
		}
	}
	if err := saveLineInfo(dir + "/lines/"); err != nil {
		exit(err)
	}
}

func getNetworkMapURLs(archived bool) ([]string, []string, error) {
	urlMap := make(map[string]string)
	downloads, err := bvg.GetNetworkMaps("")
	if err != nil {
		return nil, nil, err
	}
	for _, download := range downloads {
		urlMap[download.URL] = download.Title
	}

	if archived {
		timemap, err := ia.GetTimemap("https://www.bvg.de/de/Fahrinfo/Downloads/BVG-Liniennetz", &ia.TimemapOptions{
			Fields: []string{"timestamp"},
			Limit:  100000,
		})
		if err != nil {
			return nil, nil, err
		}
		for _, entry := range timemap {
			downloads, err := bvg.GetNetworkMaps(entry[0])
			if err != nil {
				return nil, nil, err
			}
			for _, download := range downloads {
				if _, ok := urlMap[download.URL]; !ok {
					urlMap[download.URL] = download.Title
				}
			}
		}
	}

	urls := make([]string, 0, len(urlMap))
	for url := range urlMap {
		urls = append(urls, url)
	}
	sort.Strings(urls)
	titles := make([]string, len(urlMap))
	for i, url := range urls {
		titles[i] = urlMap[url]
	}
	return urls, titles, nil
}

func saveLineInfo(dir string) error {
	info, err := bvg.GetLineInfo("")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	for _, line := range info {
		url := "https://www.bvg.de" + line.PDFURL
		fmt.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		if err := bvg.SaveFile(resp, dir); err != nil {
			return err
		}
	}
	return nil
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
