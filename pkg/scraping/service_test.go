package scraping

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var update = flag.Bool("update", false, "Update golden files")

func TestGetSoundtrack1(t *testing.T) {
	file, _ := os.Open(filepath.Join("test-fixtures", "soundtrack1.html"))
	doc, _ := goquery.NewDocumentFromReader(file)
	service := Service{}
	soundtracks := service.GetSoundtracks(doc)
	jsonData, _ := json.Marshal(soundtracks)
	goldenData, err := ioutil.ReadFile(filepath.Join("test-fixtures", "soundtrack1.golden"))
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	if !bytes.Equal(jsonData, goldenData) {
		t.Errorf("JSON does not match .golden file")
		t.Errorf("Expected %v, got %v", string(goldenData), string(jsonData))
	}
}

func TestGetSoundtrack2(t *testing.T) {
	file, _ := os.Open(filepath.Join("test-fixtures", "soundtrack2.html"))
	doc, _ := goquery.NewDocumentFromReader(file)
	service := Service{}
	soundtracks := service.GetSoundtracks(doc)
	jsonData, _ := json.Marshal(soundtracks)
	goldenData, err := ioutil.ReadFile(filepath.Join("test-fixtures", "soundtrack2.golden"))
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	if !bytes.Equal(jsonData, goldenData) {
		t.Errorf("JSON does not match .golden file")
		t.Errorf("Expected %v, got %v", string(goldenData), string(jsonData))
	}
}

func TestGetSoundtrackArtistWithAndToken(t *testing.T) {
	file, _ := os.Open(filepath.Join("test-fixtures", "soundtrack3.html"))
	doc, _ := goquery.NewDocumentFromReader(file)
	service := Service{}
	soundtracks := service.GetSoundtracks(doc)
	jsonData, _ := json.Marshal(soundtracks)
	goldenData, err := ioutil.ReadFile(filepath.Join("test-fixtures", "soundtrack3.golden"))
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	if !bytes.Equal(jsonData, goldenData) {
		t.Errorf("JSON does not match .golden file")
		t.Errorf("Expected \n%v, got \n%v", string(goldenData), string(jsonData))
	}
}

func TestRejectOrchestraArtists(t *testing.T) {
	file, _ := os.Open(filepath.Join("test-fixtures", "soundtrack4.html"))
	doc, _ := goquery.NewDocumentFromReader(file)
	service := Service{}
	soundtracks := service.GetSoundtracks(doc)
	jsonData, _ := json.Marshal(soundtracks)
	goldenData, err := ioutil.ReadFile(filepath.Join("test-fixtures", "soundtrack4.golden"))
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	if !bytes.Equal(jsonData, goldenData) {
		t.Errorf("JSON does not match .golden file")
		t.Errorf("Expected \n%v, got \n%v", string(goldenData), string(jsonData))
	}
}

func TestGetSoundtrack5(t *testing.T) {
	file, _ := os.Open(filepath.Join("test-fixtures", "soundtrack5.html"))
	doc, _ := goquery.NewDocumentFromReader(file)
	service := Service{}
	soundtracks := service.GetSoundtracks(doc)
	jsonData, _ := json.Marshal(soundtracks)
	goldenData, err := ioutil.ReadFile(filepath.Join("test-fixtures", "soundtrack5.golden"))
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}
	if !bytes.Equal(jsonData, goldenData) {
		t.Errorf("JSON does not match .golden file")
		t.Errorf("Expected \n%v, got \n%v", string(goldenData), string(jsonData))
	}
}
