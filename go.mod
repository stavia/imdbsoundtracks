module imdbsoundtracks

go 1.20

require github.com/stavia/imdbsoundtracks/pkg/scraping v0.0.0-20231013142406-2f65fdadfaca
replace github.com/stavia/imdbsoundtracks/pkg/scraping => ./pkg/scraping

require (
	github.com/PuerkitoBio/goquery v1.8.1 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	golang.org/x/net v0.17.0 // indirect
)
