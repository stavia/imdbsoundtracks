package scraping

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var update = flag.Bool("update", false, "Update golden files")

func TestGetSoundtrack1(t *testing.T) {
	imdbID := "tt1959332"
	goldenPath := fmt.Sprintf("./test-fixtures/%s.html", imdbID)
	if *update {
		url := fmt.Sprintf("https://www.imdb.com/title/%s/soundtrack", imdbID)
		resp, _ := http.Get(url)
		body, _ := ioutil.ReadAll(resp.Body)
		ioutil.WriteFile(goldenPath, []byte(body), 0644)
	}
	file, _ := os.Open(goldenPath)
	doc, _ := goquery.NewDocumentFromReader(file)
	service := Service{}
	soundtracks := service.GetSoundtracks(doc)
	jsonData, _ := json.Marshal(soundtracks)
	fmt.Println(string(jsonData))
}
