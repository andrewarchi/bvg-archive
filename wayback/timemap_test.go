package wayback

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGetTimeMap(t *testing.T) {
	timemap, err := GetTimeMap("https://www.bvg.de/de/Fahrinfo/Downloads/BVG-Liniennetz")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(timemap)
}

func TestGetPage(t *testing.T) {
	resp, err := GetPage("https://www.bvg.de/de/Fahrinfo/Downloads/BVG-Liniennetz", "20200229212100")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}
