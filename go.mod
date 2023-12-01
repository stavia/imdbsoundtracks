module imdbsoundtracks

go 1.20

require github.com/stavia/imdbsoundtracks/pkg/scraping v0.0.0-20231024065410-1d4583b39255

replace github.com/stavia/imdbsoundtracks/pkg/scraping => ./pkg/scraping

require (
	github.com/PuerkitoBio/goquery v1.8.1 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	golang.org/x/net v0.19.0 // indirect
)
