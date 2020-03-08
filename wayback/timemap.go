package wayback

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TimeMap struct {
	URLKey     string
	Timestamp  string
	Original   string
	MIMEType   string
	StatusCode string
	Digest     string
	Redirect   string
	RobotFlags string
	Length     string
	Offset     string
	Filename   string
}

func GetTimeMap(url string) ([]TimeMap, error) {
	res, err := http.Get("http://web.archive.org/web/timemap/json/" + url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var lines [][]string
	if err := json.NewDecoder(res.Body).Decode(&lines); err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, errors.New("empty response")
	}
	timemap := make([]TimeMap, len(lines)-1)
	for i, line := range lines[1:] {
		if len(line) != 11 {
			return nil, fmt.Errorf("time map entry %d not length 11", i)
		}
		timemap[i] = TimeMap{line[0], line[1], line[2], line[3], line[4],
			line[5], line[6], line[7], line[8], line[9], line[10]}
	}
	return timemap, nil
}

func GetPage(url, timestamp string) ([]byte, error) {
	res, err := http.Get("http://web.archive.org/web/" + timestamp + "id_/" + url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}
