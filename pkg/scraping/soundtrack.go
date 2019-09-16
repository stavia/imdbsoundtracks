package scraping

// Artist defines the properties of a artist
type Artist struct {
	Name   string `json:"name"`
	Role   string `json:"role"`
	ImdbID string `json:"imdbID"`
	Image  string `json:"image"`
}

// Soundtrack defines the properties of a soundtrack
type Soundtrack struct {
	Name    string `json:"name"`
	Artists []Artist
}
