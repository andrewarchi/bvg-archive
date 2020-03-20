package main

import (
	"fmt"
	"os"
	"time"

	"github.com/andrewarchi/bvg-archive/bvg"
	"github.com/andrewarchi/bvg-archive/wayback"
)

type downloadInfo struct {
	Title   string
	URL     string
	Version time.Time
	Capture string
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
		for _, download := range downloads {
			fmt.Printf("%s\t%s\n", download.Date, download.URL)
			captures[download.URL] = append(captures[download.URL], downloadInfo{
				Title:   download.Title,
				URL:     download.URL,
				Version: download.Date,
				Capture: entry.Timestamp,
			})
		}
	}
	for url := range captures {
		if err := bvg.SaveAllVersions("https://www.bvg.de"+url, "files/"); err != nil {
			exit(err)
		}
	}
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
