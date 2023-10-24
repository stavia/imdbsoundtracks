package main

import (
	"flag"
	"net/http"

	"github.com/stavia/imdbsoundtracks/pkg/scraping"
)

func main() {
	arg1 := flag.String("imdbID", "", "imdbID of the film/tvSerie to get the soundtracks")

	flag.Parse()
	if *arg1 == "" {
		flag.PrintDefaults()
		return
	}

	client := &http.Client{}
	url := "https://www.imdb.com"
	scraper := scraping.NewScraper(client, url)
	soundtracks, _ := scraper.Soundtracks(*arg1)
	for _, soundtrack := range soundtracks {
		soundtrack.PrettyPrint()
	}
}
