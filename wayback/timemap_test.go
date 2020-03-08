package wayback

import "testing"

func TestGetTimeMap(t *testing.T) {
	timemap, err := GetTimeMap("https://www.bvg.de/de/Fahrinfo/Downloads/BVG-Liniennetz")
	if err != nil {
		t.Fatal(err)
	}
	t.Error(timemap)
}

func TestGetPage(t *testing.T) {
	body, err := GetPage("https://www.bvg.de/de/Fahrinfo/Downloads/BVG-Liniennetz", "20200229212100")
	if err != nil {
		t.Fatal(err)
	}
	t.Error(string(body))
}
