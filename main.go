package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/andrewarchi/bvg-archive/bvg"
	"github.com/andrewarchi/bvg-archive/wayback"
)

func main() {
	urls, titles, err := getNetworkMapURLs(false)
	if err != nil {
		exit(err)
	}
	for i, url := range urls {
		id := strings.TrimPrefix(url, "/de/index.php?section=downloads&cmd=58&download=")
		dir := filepath.Join("files", bvg.SanitizeFilename(id+" "+titles[i]))
		if err := bvg.SaveAllVersions("https://www.bvg.de"+url, dir); err != nil {
			exit(err)
		}
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
		timemap, err := wayback.GetTimeMap("https://www.bvg.de/de/Fahrinfo/Downloads/BVG-Liniennetz")
		if err != nil {
			return nil, nil, err
		}
		for _, entry := range timemap {
			downloads, err := bvg.GetNetworkMaps(entry.Timestamp)
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

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
