package main

import (
	"flag"

	"github.com/stavia/imdbsoundtracks/pkg/scraping"
)

func main() {
	arg1 := flag.String("imdbID", "", "imdbID of the film/tvSerie to get the soundtracks")

	flag.Parse()
	if *arg1 == "" {
		flag.PrintDefaults()
		return
	}

	scraper := new(scraping.Service)
	soundtracks := scraper.Soundtracks(*arg1)
	for _, soundtrack := range soundtracks {
		soundtrack.PrettyPrint()
	}
}
