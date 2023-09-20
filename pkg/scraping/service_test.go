package scraping

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var update = flag.Bool("update", false, "Update golden files")

func TestGetGoQueryDoc(t *testing.T) {
	if *update {
		imdbID := "tt5626004"
		service := Service{}
		soundtracks := service.Soundtracks(imdbID)
		if len(soundtracks) != 10 {
			t.Errorf("Expected %v soundtracks, \n got %v", 10, len(soundtracks))
		}
	}
}

func TestGetSoundtrackAmericanMary(t *testing.T) {
	if *update {
		content := getHtml("tt1959332")
		os.WriteFile(filepath.Join("test-fixtures", "american_mary.html"), content, 0644)
	}
	file, _ := os.Open(filepath.Join("test-fixtures", "american_mary.html"))
	doc, _ := goquery.NewDocumentFromReader(file)
	service := Service{}
	soundtracks := service.GetSoundtracks(doc)
	jsonData, _ := json.Marshal(soundtracks)
	goldenData, err := ioutil.ReadFile(filepath.Join("test-fixtures", "american_mary.golden"))
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	if !bytes.Equal(jsonData, goldenData) {
		t.Errorf("JSON does not match .golden file")
		t.Errorf("Expected %v, \n got %v", string(goldenData), string(jsonData))
	}
}

func TestGetSoundtrackJoker(t *testing.T) {
	if *update {
		content := getHtml("tt7286456")
		os.WriteFile(filepath.Join("test-fixtures", "joker.html"), content, 0644)
	}
	file, _ := os.Open(filepath.Join("test-fixtures", "joker.html"))
	doc, _ := goquery.NewDocumentFromReader(file)
	service := Service{}
	soundtracks := service.GetSoundtracks(doc)
	jsonData, _ := json.Marshal(soundtracks)
	goldenData, err := ioutil.ReadFile(filepath.Join("test-fixtures", "joker.golden"))
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	if !bytes.Equal(jsonData, goldenData) {
		t.Errorf("JSON does not match .golden file")
		t.Errorf("Expected \n%v, got \n%v", string(goldenData), string(jsonData))
	}
}

func TestGetSoundtrackIrishConnection(t *testing.T) {
	if *update {
		content := getHtml("tt10164206")
		os.WriteFile(filepath.Join("test-fixtures", "irish_connection.html"), content, 0644)
	}
	file, _ := os.Open(filepath.Join("test-fixtures", "irish_connection.html"))
	doc, _ := goquery.NewDocumentFromReader(file)
	service := Service{}
	soundtracks := service.GetSoundtracks(doc)
	jsonData, _ := json.Marshal(soundtracks)
	goldenData, err := ioutil.ReadFile(filepath.Join("test-fixtures", "irish_connection.golden"))
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	if !bytes.Equal(jsonData, goldenData) {
		t.Errorf("JSON does not match .golden file")
		t.Errorf("Expected \n%v, got \n%v", string(goldenData), string(jsonData))
	}
}

func TestGetSoundtrackMarryMe(t *testing.T) {
	if *update {
		content := getHtml("tt10223460")
		os.WriteFile(filepath.Join("test-fixtures", "marry_me.html"), content, 0644)
	}
	file, _ := os.Open(filepath.Join("test-fixtures", "marry_me.html"))
	doc, _ := goquery.NewDocumentFromReader(file)
	service := Service{}
	soundtracks := service.GetSoundtracks(doc)
	jsonData, _ := json.Marshal(soundtracks)
	goldenData, err := ioutil.ReadFile(filepath.Join("test-fixtures", "marry_me.golden"))
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	if !bytes.Equal(jsonData, goldenData) {
		t.Errorf("JSON does not match .golden file")
		t.Errorf("Expected \n%v, got \n%v", string(goldenData), string(jsonData))
	}
}

func TestGetSoundtrackSeniorYear(t *testing.T) {
	if *update {
		content := getHtml("tt5315212")
		os.WriteFile(filepath.Join("test-fixtures", "senior_year.html"), content, 0644)
	}
	file, _ := os.Open(filepath.Join("test-fixtures", "senior_year.html"))
	doc, _ := goquery.NewDocumentFromReader(file)
	service := Service{}
	soundtracks := service.GetSoundtracks(doc)
	jsonData, _ := json.Marshal(soundtracks)
	goldenData, err := ioutil.ReadFile(filepath.Join("test-fixtures", "senior_year.golden"))
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	if !bytes.Equal(jsonData, goldenData) {
		t.Errorf("JSON does not match .golden file")
		t.Errorf("Expected \n%v, got \n%v", string(goldenData), string(jsonData))
	}
}

func TestGetSoundtrackArmageddonTime(t *testing.T) {
	if *update {
		content := getHtml("tt10343028")
		os.WriteFile(filepath.Join("test-fixtures", "armageddon_time.html"), content, 0644)
	}
	file, _ := os.Open(filepath.Join("test-fixtures", "armageddon_time.html"))
	doc, _ := goquery.NewDocumentFromReader(file)
	service := Service{}
	soundtracks := service.GetSoundtracks(doc)
	jsonData, _ := json.Marshal(soundtracks)
	goldenData, err := ioutil.ReadFile(filepath.Join("test-fixtures", "armageddon_time.golden"))
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	if !bytes.Equal(jsonData, goldenData) {
		t.Errorf("JSON does not match .golden file")
		t.Errorf("Expected \n%v, got \n%v", string(goldenData), string(jsonData))
	}
}

func TestGetSoundtrackAlienAbduction(t *testing.T) {
	if *update {
		content := getHtml("tt15213332")
		os.WriteFile(filepath.Join("test-fixtures", "alien_abduction.html"), content, 0644)
	}
	file, _ := os.Open(filepath.Join("test-fixtures", "alien_abduction.html"))
	doc, _ := goquery.NewDocumentFromReader(file)
	service := Service{}
	soundtracks := service.GetSoundtracks(doc)
	jsonData, _ := json.Marshal(soundtracks)
	goldenData, err := ioutil.ReadFile(filepath.Join("test-fixtures", "alien_abduction.golden"))
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	if !bytes.Equal(jsonData, goldenData) {
		t.Errorf("JSON does not match .golden file")
		t.Errorf("Expected \n%v, got \n%v", string(goldenData), string(jsonData))
	}
}

func TestGetSoundtrackSwitch(t *testing.T) {
	if *update {
		content := getHtml("tt7139110")
		os.WriteFile(filepath.Join("test-fixtures", "switch.html"), content, 0644)
	}
	file, _ := os.Open(filepath.Join("test-fixtures", "switch.html"))
	doc, _ := goquery.NewDocumentFromReader(file)
	service := Service{}
	soundtracks := service.GetSoundtracks(doc)
	jsonData, _ := json.Marshal(soundtracks)
	goldenData, err := ioutil.ReadFile(filepath.Join("test-fixtures", "switch.golden"))
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	if !bytes.Equal(jsonData, goldenData) {
		t.Errorf("JSON does not match .golden file")
		t.Errorf("Expected \n%v, got \n%v", string(goldenData), string(jsonData))
	}
}

func getHtml(imdbID string) []byte {
	url := fmt.Sprintf("https://www.imdb.com/title/%s/soundtrack", imdbID)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalln(fmt.Errorf("status code error: %d %s, url: %s", res.StatusCode, res.Status, url))
	}
	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return html
}
