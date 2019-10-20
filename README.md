# imdbsoundtracks @
imdbsoundtracks is a Go library for retrieving soundtracks from IMDb.

## Usage ##

Download and install the package
```console
go get -u github.com/stavia/imdbsoundtracks/pkg/scraping
```

Import the package
```go
import "github.com/stavia/imdbsoundtracks/pkg/scraping"
```

Construct a new scraping service, then get all soundtracks for the IMDb id tt7286456 ([Joker 2019](https://www.imdb.com/title/tt7286456/))
```go
scraper := new(scraping.Service)
soundtracks := scraper.Soundtracks("tt7286456")
```

Loop through and print all soundtracks
```go
for _, soundtrack := range soundtracks {
    soundtrack.PrettyPrint()
}
```