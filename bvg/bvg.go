package bvg

import (
	"crypto/sha512"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/andrewarchi/bvg-archive/wayback"
)

type NetworkMap struct {
	URL   string
	Title string
	Date  time.Time
}

func GetNetworkMaps(timestamp string) ([]NetworkMap, error) {
	url := "https://www.bvg.de/de/Fahrinfo/Downloads/BVG-Liniennetz"
	if timestamp != "" {
		url = wayback.PageURL(url, timestamp)
	}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}
	links := doc.Find(".article__body .download .link-list li.link-list__item")
	maps := make([]NetworkMap, links.Length())
	links.Each(func(i int, s *goquery.Selection) {
		url, _ := s.Find("a").Attr("href")
		title := s.Find(".link-list__text").First().Text()
		date, _ := s.Find("img").Attr("alt")
		date = strings.TrimPrefix(date, "Aktualisiert am: ")

		var t time.Time
		if date != "" {
			t, _ = time.Parse("02.01.2006", date)
		}

		maps[i] = NetworkMap{url, title, t}
	})
	return maps, nil
}

type LineInfo struct {
	LongName  string
	ShortName string
	PDFURL    string
	ImageURL  string
}

func GetLineInfo(timestamp string) ([]LineInfo, error) {
	url := "https://www.bvg.de/de/Fahrinfo/Haltestelleinfo"
	if timestamp != "" {
		url = wayback.PageURL(url, timestamp)
	}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}
	var info []LineInfo
	doc.Find(".tab-list__body tr").Each(func(i int, s *goquery.Selection) {
		icon := s.Find(".tab-list__icon .icon-t").First()
		if icon.Length() != 0 {
			long := icon.Find(".visuallyhidden").Text()
			short := strings.TrimSpace(icon.Nodes[0].NextSibling.Data)
			pdf, _ := s.Find(".tab-list__text a").Attr("href")
			image, _ := s.Find(".tab-list__text ~ .tab-list__text a").Attr("href")
			info = append(info, LineInfo{long, short, pdf, image})
		}
	})
	return info, nil
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
	if filename == "" {
		filename = SanitizeFilename(path.Base(resp.Request.URL.Path))
	}

	path := pathPrefix + filename
	if stat, err := os.Stat(path); err == nil && stat.Size() > 0 {
		fmt.Println("Skipped")
		return nil // skip existing
	}
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
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	savedTimes, savedLive, err := getSavedTimes(dir, 2*time.Hour)
	if err != nil {
		return err
	}

	if !savedLive {
		t := time.Now().Format(wayback.TimestampFormat)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		if err := SaveFile(resp, filepath.Join(dir, t+"live_")); err != nil {
			return err
		}
	}

	timemap, err := wayback.GetTimeMap(url)
	if err != nil {
		return err
	}
	for i := len(timemap) - 1; i >= 0; i-- {
		timestamp := timemap[i].Timestamp
		if _, ok := savedTimes[timestamp]; ok {
			continue
		}
		resp, err := wayback.GetPage(url, timestamp)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if err := SaveFile(resp, filepath.Join(dir, timestamp+"_")); err != nil {
			return err
		}
	}
	return nil
}

func getSavedTimes(dir string, d time.Duration) (map[string]struct{}, bool, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, false, err
	}
	times := make(map[string]struct{})
	hasLive := false
	for _, file := range files {
		if !file.IsDir() && file.Size() > 0 {
			n := file.Name()
			l := len(wayback.TimestampFormat)
			if len(n) < l {
				continue
			}
			timestamp := n[:l]
			t, err := time.Parse(wayback.TimestampFormat, timestamp)
			if err != nil {
				continue
			}
			if len(n) >= l+4 && n[l:l+4] == "live" {
				if time.Since(t) <= d {
					hasLive = true
				}
			} else {
				times[timestamp] = struct{}{}
			}
		}
	}
	return times, hasLive, nil
}

var illegal = regexp.MustCompile(`[\0-\x1f:?"*/\\<>|]`)

func SanitizeFilename(filename string) string {
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

func getRetrieved(header http.Header) (time.Time, error) {
	return getDate(header, "X-Archive-Orig-Date", "Memento-Datetime", "Date")
}

func getLastModified(header http.Header) (time.Time, error) {
	return getDate(header, "X-Archive-Orig-Last-Modified", "Last-Modified")
}

func getDate(header http.Header, keys ...string) (time.Time, error) {
	for _, key := range keys {
		if h := header.Get(key); h != "" {
			if t, err := time.Parse(time.RFC1123, h); err == nil {
				return t, nil
			}
			return time.Parse(time.RFC850, h)
		}
	}
	return time.Time{}, nil
}

func hash(r io.Reader) (sum [sha512.Size]byte, err error) {
	h := sha512.New()
	_, err = io.Copy(h, r)
	if err == nil {
		h.Sum(sum[:])
	}
	return
}
