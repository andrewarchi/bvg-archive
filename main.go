package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
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
			page, err := wayback.GetPage(fullURL, entry.Timestamp)
			if err != nil {
				exit(err)
			}
			defer page.Body.Close()
			filename, err := getFilename(page.Header)
			if err != nil {
				exit(err)
			}
			path := fmt.Sprintf("%s/%s_%s", dir, entry.Timestamp, filename)
			file, err := os.Create(path)
			if err != nil {
				exit(err)
			}
			defer file.Close()
			if _, err := io.Copy(file, page.Body); err != nil {
				exit(err)
			}
			lastModified, err := getLastModified(page.Header)
			if err != nil {
				exit(err)
			}
			if err := os.Chtimes(path, lastModified, lastModified); err != nil {
				exit(err)
			}
		}
	}
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
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
	}
	return time.Parse(time.RFC1123, mod)
}
