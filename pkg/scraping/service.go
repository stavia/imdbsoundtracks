package scraping

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Scraper provides soundtrack scrapping operations
type Scraper interface {
	Soundtracks(imdbID string) (soundtracks []Soundtrack)
}

type Service struct {
}

// Soundtracks returns the soundtracks found for the given imdbID
func (s *Service) Soundtracks(imdbID string) (soundtracks []Soundtrack) {
	if !strings.Contains(imdbID, "tt") {
		imdbID = "tt" + imdbID
	}
	url := fmt.Sprintf("https://www.imdb.com/title/%s/soundtrack", imdbID)
	doc, err := s.GetGoqueryDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	return s.GetSoundtracks(doc)
}

// GetGoqueryDocument returns a Document that takes a string URL as argument
func (s *Service) GetGoqueryDocument(url string) (doc *goquery.Document, err error) {
	doc, err = goquery.NewDocument(url)
	return doc, err
}

// GetSoundtracks returns all soundtracks found for the given goquery Document
func (s *Service) GetSoundtracks(doc *goquery.Document) (soundtracks []Soundtrack) {
	doc.Find(".ipc-metadata-list").First().Find("li").Each(func(index int, selection *goquery.Selection) {
		soundtrack := s.GetSoundtrack(selection)
		if len(soundtrack.Artists) > 0 {
			soundtracks = append(soundtracks, soundtrack)
		}
	})

	return soundtracks
}

// GetSoundtrack extracts from the given text all the info of a soundtrack
func (s *Service) GetSoundtrack(doc *goquery.Selection) (soundtrack Soundtrack) {
	soundtrack.Name = getSoundtrackName(doc)
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
						artistFound = setArtistImdbIDFromGoquerySelection(artist)
						artistFound.Role = role
						soundtrack.Artists = appendArtist(soundtrack.Artists, artistFound)
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

func setArtistImdbIDFromGoquerySelection(artistSelection *goquery.Selection) (artist Artist) {
	artist.Name = artistSelection.Text()
	href, exist := artistSelection.Attr("href")
	if exist {
		artist.ImdbID = getArtistImdbID(href)
		artist.Image = getArtistImage(artist.ImdbID)
	}
	return artist
}

func getArtistImage(artistImdbID string) (urlImage string) {
	url := fmt.Sprintf("https://www.imdb.com/name/%s/", artistImdbID)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return urlImage
	}
	urlImage, _ = doc.Find("#name-poster").First().Attr("src")
	return urlImage
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

func getArtistImdbID(url string) string {
	re := regexp.MustCompile(`\/name\/(.*)\/`)
	matches := re.FindStringSubmatch(url)
	return matches[1]
}

func replaceAndByCommas(line string) string {
	return strings.Replace(line, ", and ", ", ", -1)
}
