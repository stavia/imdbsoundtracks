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
	doc.Find("#soundtracks_content").Find(".list").Find("div").Each(func(index int, selection *goquery.Selection) {
		soundtrack := s.GetSoundtrack(selection)
		if len(soundtrack.Artists) > 0 {
			soundtracks = append(soundtracks, soundtrack)
		}
	})
	return soundtracks
}

// GetSoundtrack extracts from the given text all the info of a soundtrack
func (s *Service) GetSoundtrack(doc *goquery.Selection) (soundtrack Soundtrack) {
	html, _ := doc.Html()
	splits := strings.Split(html, "\n")
	soundtrack.Name = getSoundtrackName(splits[0])
	byRegexp := regexp.MustCompile(`(.*)by(.*)`)
	for i := 1; i < len(splits); i++ {
		matches := byRegexp.FindStringSubmatch(splits[i])
		rolesFound := standardRole(splits[i])
		if len(rolesFound) == 0 {
			continue
		}
		for _, role := range rolesFound {
			if len(matches) != 3 {
				continue
			}
			line := strings.Split(replaceAndByCommas(matches[2]), ",")
			for _, artist := range line {
				ampersandFound := strings.Contains(artist, "&amp;")
				var artistFound Artist
				artistDoc, _ := goquery.NewDocumentFromReader(strings.NewReader((artist)))
				artistNodes := artistDoc.Find("a")
				if ampersandFound {
					// check if there is any artist name that contains ampersand
					appendAmpersandArtist(artistNodes, artistDoc, role, &soundtrack)
				}
				if artistNodes.Length() > 0 {
					artistNodes.Each(func(index int, artist *goquery.Selection) {
						artistFound = setArtistImdbIDFromGoquerySelection(artist)
						artistFound.Role = role
						soundtrack.Artists = appendArtist(soundtrack.Artists, artistFound)
					})
				} else {
					artistFound = getArtistFromText(artistDoc.Text(), role)
					if artistFound.Role != "" {
						soundtrack.Artists = appendArtist(soundtrack.Artists, artistFound)
					}
				}
			}
		}
	}
	return soundtrack
}

func getSoundtrackName(text string) string {
	document, _ := goquery.NewDocumentFromReader(strings.NewReader((text)))
	return strings.TrimSpace(document.Text())
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

func appendAmpersandArtist(artistNodes *goquery.Selection, artistDoc *goquery.Document, role string, soundtrack *Soundtrack) {
	line := artistDoc.Text()
	artistNodes.Each(func(index int, artist *goquery.Selection) {
		if !strings.Contains(artist.Text(), "&") {
			line = strings.TrimSpace(strings.Replace(line, artist.Text(), "", -1))
		}
	})
	splits := strings.Split(line, "&")
	for _, split := range splits {
		if split == "" {
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
		//artist.Image = GetArtistImage(artist.ImdbID)
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
	for _, existArtist := range artists {
		if existArtist == artist {
			return artists
		}
	}
	return append(artists, artist)
}

func standardRole(role string) (result []string) {
	role = strings.ToLower(strings.TrimSpace(role))
	if strings.Contains(role, "music by") {
		result = append(result, "composer")
	}
	if strings.Contains(role, "performed") {
		result = append(result, "performer")
	}
	if strings.Contains(role, "written") ||
		strings.Contains(role, "lyrics") ||
		strings.Contains(role, "composed") {
		result = append(result, "writer")
	}
	if strings.HasPrefix(role, "by") {
		result = append(result, "composer")
	}
	if strings.Contains(role, "produced by") {
		result = append(result, "producer")
	}
	if strings.Contains(role, "arranged by") {
		result = append(result, "music arranger")
	}
	return result
}

func getImdbID(url string) string {
	re := regexp.MustCompile(`tt[0-9]{7}`)
	matches := re.FindStringSubmatch(url)
	return matches[0]
}

func getArtistImdbID(url string) string {
	re := regexp.MustCompile(`\/name\/(.*)\/`)
	matches := re.FindStringSubmatch(url)
	return matches[1]
}

func replaceAndByCommas(line string) string {
	var re = regexp.MustCompile(`(?m)^([^<]*)<\w+.*/\w+>([^<]*)$`)
	for i, match := range re.FindStringSubmatch(line) {
		if i != 0 {
			var rightText = strings.Replace(match, " and ", ",", -1)
			line = strings.Replace(line, match, rightText, -1)
		}
	}
	return line
}
