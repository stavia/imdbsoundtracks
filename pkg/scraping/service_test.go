package scraping

import (
	"bytes"
	_ "bytes"
	"encoding/json"
	_ "encoding/json"
	"flag"
	"fmt"
	_ "fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/PuerkitoBio/goquery"
)

var update = flag.Bool("update", false, "Update golden files")
var e2e = flag.Bool("e2e", false, "Execute e2e tests")

func TestGetLastDanceSoundtrack(t *testing.T) {
	if *e2e {
		imdbID := "tt8420184"
		client := &http.Client{}
		service := NewScraper(client, "https://www.imdb.com")
		soundtracks, _ := service.Soundtracks(imdbID)
		expected := 32
		if len(soundtracks) != expected {
			t.Errorf("Expected %v soundtracks, \n got %v", expected, len(soundtracks))
		}
	}
}

func TestGetSoundtrackAmericanMary(t *testing.T) {
	imdbID := "tt1959332"
	if *update {
		content := getHtml(imdbID)
		os.WriteFile(filepath.Join("test-fixtures", "american_mary.html"), content, 0644)
	}
	expected, _ := ioutil.ReadFile(filepath.Join("test-fixtures", "american_mary.html"))
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(expected))
	}))
	defer svr.Close()
	client := &http.Client{}
	service := NewScraper(client, svr.URL)
	soundtracks, _ := service.Soundtracks(imdbID)
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
	imdbID := "tt7286456"
	if *update {
		content := getHtml(imdbID)
		os.WriteFile(filepath.Join("test-fixtures", "joker.html"), content, 0644)
	}
	expected, _ := ioutil.ReadFile(filepath.Join("test-fixtures", "joker.html"))
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(expected))
	}))
	defer svr.Close()
	client := &http.Client{}
	service := NewScraper(client, svr.URL)
	soundtracks, _ := service.Soundtracks(imdbID)
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
	imdbID := "tt10164206"
	if *update {
		content := getHtml("tt10164206")
		os.WriteFile(filepath.Join("test-fixtures", "irish_connection.html"), content, 0644)
	}
	expected, _ := ioutil.ReadFile(filepath.Join("test-fixtures", "irish_connection.html"))
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(expected))
	}))
	defer svr.Close()
	client := &http.Client{}
	service := NewScraper(client, svr.URL)
	soundtracks, _ := service.Soundtracks(imdbID)
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
	imdbID := "tt10223460"
	if *update {
		content := getHtml(imdbID)
		os.WriteFile(filepath.Join("test-fixtures", "marry_me.html"), content, 0644)
	}
	expected, _ := ioutil.ReadFile(filepath.Join("test-fixtures", "marry_me.html"))
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(expected))
	}))
	defer svr.Close()
	client := &http.Client{}
	service := NewScraper(client, svr.URL)
	soundtracks, _ := service.Soundtracks(imdbID)
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
	imdbID := "tt5315212"
	if *update {
		content := getHtml(imdbID)
		os.WriteFile(filepath.Join("test-fixtures", "senior_year.html"), content, 0644)
	}
	expected, _ := ioutil.ReadFile(filepath.Join("test-fixtures", "senior_year.html"))
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(expected))
	}))
	defer svr.Close()
	client := &http.Client{}
	service := NewScraper(client, svr.URL)
	soundtracks, _ := service.Soundtracks(imdbID)
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
	imdbID := "tt10343028"
	if *update {
		content := getHtml(imdbID)
		os.WriteFile(filepath.Join("test-fixtures", "armageddon_time.html"), content, 0644)
	}
	expected, _ := ioutil.ReadFile(filepath.Join("test-fixtures", "armageddon_time.html"))
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(expected))
	}))
	defer svr.Close()
	client := &http.Client{}
	service := NewScraper(client, svr.URL)
	soundtracks, _ := service.Soundtracks(imdbID)
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
	imdbID := "tt15213332"
	if *update {
		content := getHtml(imdbID)
		os.WriteFile(filepath.Join("test-fixtures", "alien_abduction.html"), content, 0644)
	}
	expected, _ := ioutil.ReadFile(filepath.Join("test-fixtures", "alien_abduction.html"))
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(expected))
	}))
	defer svr.Close()
	client := &http.Client{}
	service := NewScraper(client, svr.URL)
	soundtracks, _ := service.Soundtracks(imdbID)
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
	imdbID := "tt7139110"
	if *update {
		content := getHtml(imdbID)
		os.WriteFile(filepath.Join("test-fixtures", "switch.html"), content, 0644)
	}
	expected, _ := ioutil.ReadFile(filepath.Join("test-fixtures", "switch.html"))
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(expected))
	}))
	defer svr.Close()
	client := &http.Client{}
	service := NewScraper(client, svr.URL)
	soundtracks, _ := service.Soundtracks(imdbID)
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
