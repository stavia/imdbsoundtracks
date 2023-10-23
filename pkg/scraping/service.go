package scraping

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Scraper provides soundtrack scrapping operations
type Scraper interface {
	Soundtracks(imdbID string) (soundtracks []Soundtrack, err error)
}

type ScraperHttpClient struct {
	Client *http.Client
	Url    string
}

func NewScraper(Client *http.Client, Url string) Scraper {
	return &ScraperHttpClient{Client, Url}
}

// https://golang.cafe/blog/golang-httptest-example.html
// https://speedscale.com/blog/testing-golang-with-httptest/
// https://quii.gitbook.io/learn-go-with-tests/build-an-application/http-server

// Soundtracks returns the soundtracks found for the given imdbID
func (s *ScraperHttpClient) Soundtracks(imdbID string) (soundtracks []Soundtrack, err error) {
	if !strings.Contains(imdbID, "tt") {
		imdbID = "tt" + imdbID
	}
	soundtracks, err = s.getSoundtracks(imdbID)
	if err != nil {
		return soundtracks, err
	}
	if len(soundtracks) == 0 {
		soundtracks, err = s.getEpisodeSoundtracks(imdbID)
	}

	return soundtracks, err
}

func (s *ScraperHttpClient) getGoqueryDocument(url string) (doc *goquery.Document, err error) {
	//client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	//req, err := s.Client.NewRequest("GET", url, nil)
	if err != nil {
		return doc, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	//res, err := client.Do(req)
	res, err := s.Client.Do(req)
	if err != nil {
		return doc, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		errorMessage := fmt.Sprintf("status code error: %d %s, url: %s", res.StatusCode, res.Status, url)
		err = errors.New(errorMessage)
		return doc, err
	}
	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return doc, err
	}
	return doc, err
}

// getSoundtracks returns all soundtracks found for the given imdbID
func (s *ScraperHttpClient) getSoundtracks(imdbID string) (soundtracks []Soundtrack, err error) {
	//url := fmt.Sprintf("https://www.imdb.com/title/%s/soundtrack", s.Url, imdbID)
	url := fmt.Sprintf("%s/title/%s/soundtrack", s.Url, imdbID)
	doc, err := s.getGoqueryDocument(url)
	if err != nil {
		return soundtracks, err
	}
	doc.Find(".ipc-metadata-list").First().Find("li").Each(func(index int, selection *goquery.Selection) {
		soundtrack := s.getSoundtrack(selection)
		if len(soundtrack.Artists) > 0 {
			soundtracks = append(soundtracks, soundtrack)
		}
	})

	return soundtracks, err
}

func (s *ScraperHttpClient) getEpisodeSoundtracks(imdbID string) (soundtracks []Soundtrack, err error) {
	season := 1
	numberEpisodesWithoutSoundtracks := 0
	for {
		url := fmt.Sprintf("%s/title/%s/episodes?season=%d", s.Url, imdbID, season)
		doc, err := s.getGoqueryDocument(url)
		if err != nil {
			return soundtracks, err
		}
		if numberEpisodesWithoutSoundtracks > 4 {
			break
		}
		var soundtracksFound []Soundtrack
		doc.Find(".episode-item-wrapper").EachWithBreak(func(index int, selection *goquery.Selection) bool {
			href, exists := selection.Find("a").First().Attr("href")
			if !exists {
				numberEpisodesWithoutSoundtracks++
				return false
			}
			soundtracksFound, err = s.getSoundtracks(getImdbID(href))
			if err != nil {
				numberEpisodesWithoutSoundtracks++
				return false
			}

			if len(soundtracksFound) > 0 {
				soundtracks = append(soundtracks, soundtracksFound...)
			} else {
				numberEpisodesWithoutSoundtracks++
			}

			if numberEpisodesWithoutSoundtracks > 4 {
				return false
			}
			return true
		})
		if len(soundtracksFound) == 0 {
			numberEpisodesWithoutSoundtracks++
		}
		season++
	}
	return soundtracks, err
}

func (s *ScraperHttpClient) getSoundtrack(doc *goquery.Selection) (soundtrack Soundtrack) {
	soundtrack.Name = getSoundtrackName(doc)
	if soundtrack.Name == "" {
		return soundtrack
	}
	byRegexp := regexp.MustCompile(`(.*)\sby\s(.*)`)
	doc.Find(".ipc-html-content-inner-div").Each(func(index int, selection *goquery.Selection) {
		itemText := selection.Text()
		itemHtml, _ := selection.Html()
		rolesFound := standardRole(itemText)
		if len(rolesFound) == 0 {
			return
		}
		matches := byRegexp.FindStringSubmatch(itemHtml)
		for _, role := range rolesFound {
			if len(matches) != 3 {
				continue
			}
			line := strings.Split(replaceAndByCommas(matches[2]), ",")
			for _, artist := range line {
				artist = cleanText(artist)
				if skipLine(artist) {
					continue
				}
				ampersandFound := strings.Contains(artist, "&amp;")
				var artistFound Artist
				artistDoc, _ := goquery.NewDocumentFromReader(strings.NewReader((artist)))
				artistNodes := artistDoc.Find("a")
				if artistNodes.Length() > 0 {
					if ampersandFound {
						// check if there is any artist name that contains ampersand
						appendAmpersandArtist(artistNodes, artistDoc, role, &soundtrack)
					}
					artistNodes.Each(func(index int, artist *goquery.Selection) {
						artistFound, err := s.setArtistImdbIDFromGoquerySelection(artist)
						if err == nil {
							artistFound.Role = role
							soundtrack.Artists = appendArtist(soundtrack.Artists, artistFound)
						}
					})
				} else {
					if strings.Contains(artistDoc.Text(), "&") {
						// Skip artist. Impossible to know the name of the artist
						// Example (tt10164206):
						// Performed by Paddy Nash & The Happy Enchiladas
						continue
					}
					artistFound = getArtistFromText(artistDoc.Text(), role)
					if artistFound.Role != "" {
						soundtrack.Artists = appendArtist(soundtrack.Artists, artistFound)
					}
				}
			}
		}
	})
	return soundtrack
}

func (s *ScraperHttpClient) getArtistImage(artistImdbID string) (urlImage string) {
	url := fmt.Sprintf("%s/name/%s/", s.Url, artistImdbID)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return urlImage
	}
	urlImage, _ = doc.Find("#name-poster").First().Attr("src")
	return urlImage
}

func (s *ScraperHttpClient) setArtistImdbIDFromGoquerySelection(artistSelection *goquery.Selection) (artist Artist, err error) {
	artist.Name = artistSelection.Text()
	href, exist := artistSelection.Attr("href")
	if exist {
		artist.ImdbID = getArtistImdbID(href)
		if artist.ImdbID == "" {
			return artist, errors.New("Artist not found")
		}
		artist.Image = s.getArtistImage(artist.ImdbID)
	}
	return artist, err
}

func getSoundtrackName(doc *goquery.Selection) string {
	return doc.Find(".ipc-metadata-list-item__label").Text()
}

func getArtistFromText(text string, role string) (artist Artist) {
	artistRegexp := regexp.MustCompile(`(.*) \(.*\)`)
	artistName := strings.TrimSpace(artistRegexp.ReplaceAllString(text, `$1`))
	if artistName == "chorus" {
		return artist
	}
	artist.Name = artistName
	artist.Role = role
	return artist
}

func skipLine(text string) bool {
	text = strings.TrimSpace(text)
	return text == "" || strings.Contains(text, "(uncredited)")
}

func skipArtist(artist Artist) bool {
	return strings.Contains(artist.Name, "Orchestra") || strings.Contains(artist.Name, "orchestra")
}

func cleanArtistName(artist *Artist) {
	var re = regexp.MustCompile(`(?mi)and his orchestra`)
	artist.Name = strings.TrimSpace(re.ReplaceAllString(artist.Name, ""))
}

func cleanText(text string) string {
	text = strings.TrimSpace(text)
	regex := `\[link: .*\]`
	r := regexp.MustCompile(regex)
	text = r.ReplaceAllString(text, "")
	regex = `<br.>|\:\s`
	r = regexp.MustCompile(regex)
	text = r.ReplaceAllString(text, "")
	regex = `\(.*?\)`
	r = regexp.MustCompile(regex)
	text = r.ReplaceAllString(text, "")

	return strings.TrimSpace(text)
}

func appendAmpersandArtist(artistNodes *goquery.Selection, artistDoc *goquery.Document, role string, soundtrack *Soundtrack) {
	line := artistDoc.Text()
	artistNodes.Each(func(index int, artist *goquery.Selection) {
		if !strings.Contains(artist.Text(), "&") {
			line = strings.TrimSpace(strings.Replace(line, artist.Text(), "", -1))
		}
	})
	splits := strings.Split(line, "&")
	for _, split := range splits {
		split = strings.TrimSpace(split)
		// feat. problem -> tt10223460
		if split == "" || split == "feat." || split == "/" {
			continue
		}
		// problem -> tt5315212
		if strings.Contains(split, "(as") {
			continue
		}
		artistFound := getArtistFromText(split, role)
		if artistFound.Role != "" {
			soundtrack.Artists = appendArtist(soundtrack.Artists, artistFound)
		}
	}
}

func appendArtist(artists []Artist, artist Artist) []Artist {
	cleanArtistName(&artist)
	if skipArtist(artist) {
		return artists
	}
	for _, existArtist := range artists {
		if existArtist == artist {
			return artists
		}
	}
	return append(artists, artist)
}

func standardRole(role string) (result []string) {
	role = strings.ToLower(strings.TrimSpace(role))
	if strings.Contains(role, "by arrangement") {
		return result
	}
	if strings.Contains(role, "performed") {
		result = append(result, "performer")
	}
	if strings.Contains(role, "written") ||
		strings.Contains(role, "lyrics") {
		result = append(result, "writer")
	}
	if strings.Contains(role, "produced by") {
		result = append(result, "producer")
	}
	if strings.Contains(role, "arranged by") {
		result = append(result, "music arranger")
	}
	if strings.Contains(role, "music by") ||
		strings.Contains(role, "composed by") {
		result = append(result, "composer")
	}
	if len(result) == 0 && strings.HasPrefix(role, "by") {
		result = append(result, "composer")
	}
	return result
}

func getArtistImdbID(url string) (artistImdbID string) {
	re := regexp.MustCompile(`\/name\/(.*)\/`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		artistImdbID = matches[1]
	}
	return artistImdbID
}

func replaceAndByCommas(line string) string {
	return strings.Replace(line, ", and ", ", ", -1)
}

func getImdbID(url string) string {
	re := regexp.MustCompile(`tt[0-9]{7,}`)
	matches := re.FindStringSubmatch(url)
	return matches[0]
}
