package main

import (
	"flag"
)

func main() {
	arg1 := flag.String("imdbID", "", "imdbID of the film/tvSerie to get the soundtracks")

	flag.Parse()
	if *arg1 == "" {
		flag.PrintDefaults()
		return
	}

}
