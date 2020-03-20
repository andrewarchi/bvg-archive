package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/andrewarchi/bvg-archive/bvg"
	"github.com/andrewarchi/bvg-archive/wayback"
)

var illegalPattern = regexp.MustCompile("[/\\?&]")

type downloadInfo struct {
	Title   string
	URL     string
	Version time.Time
	Capture time.Time
}

func main() {
	timemap, err := wayback.GetTimeMap("https://www.bvg.de/de/Fahrinfo/Downloads/BVG-Liniennetz")
	if err != nil {
		exit(err)
	}
	captures := make(map[string][]downloadInfo)
	for _, entry := range timemap {
		fmt.Printf("%s:\n", entry.Timestamp)
		downloads, err := bvg.GetNetworkMaps(entry.Timestamp)
		if err != nil {
			exit(err)
		}
		capture, err := time.Parse("20060102150405", entry.Timestamp)
		if err != nil {
			exit(err)
		}
		for _, download := range downloads {
			fmt.Printf("%s\t%s\n", download.Date, download.URL)
			captures[download.URL] = append(captures[download.URL], downloadInfo{
				Title:   download.Title,
				URL:     download.URL,
				Version: download.Date,
				Capture: capture,
			})
		}
	}
	for url, info := range captures {
		fmt.Println(url)
		for _, capture := range info {
			fmt.Printf("%v\t%v\t%v\n", capture.Capture, capture.Version, capture.Title)
		}
		fullURL := "https://www.bvg.de" + url
		timemap, err := wayback.GetTimeMap(fullURL)
		if err != nil {
			exit(err)
		}
		dir := "files/" + illegalPattern.ReplaceAllString(url, "_")
		if err := os.MkdirAll(dir, 0o777); err != nil {
			exit(err)
		}
		for _, entry := range timemap {
			resp, err := wayback.GetPage(fullURL, entry.Timestamp)
			if err != nil {
				exit(err)
			}
			defer resp.Body.Close()
			if err := bvg.SavePDF(resp, fmt.Sprintf("%s/%s_", dir, entry.Timestamp)); err != nil {
				exit(err)
			}
		}
	}
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
