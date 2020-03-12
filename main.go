package main

import (
	"fmt"
	"os"
	"time"

	"github.com/andrewarchi/internet-archive/bvg"
	"github.com/andrewarchi/internet-archive/wayback"
)

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
	for _, entry := range timemap[:1] {
		fmt.Printf("%s:\n", entry.Timestamp)
		downloads, err := bvg.GetLineDownloads(entry.Timestamp)
		if err != nil {
			exit(err)
		}
		capture, err := time.Parse("20060102150405", entry.Timestamp)
		if err != nil {
			exit(err)
		}
		for _, download := range downloads {
			version, err := time.Parse("02.01.2006", download.Date)
			if err != nil {
				exit(err)
			}
			fmt.Printf("%s\t%s\n", download.Date, download.URL)
			captures[download.URL] = append(captures[download.URL], downloadInfo{
				Title:   download.Title,
				URL:     download.URL,
				Version: version,
				Capture: capture,
			})
		}
	}
	for url, info := range captures {
		fmt.Println(url)
		for _, capture := range info {
			fmt.Printf("%v\t%v\t%v\n", capture.Capture, capture.Version, capture.Title)
		}
	}
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
