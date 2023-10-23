package scraping

import (
	"encoding/json"
	"fmt"
)

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

// PrettyPrint provides pretty-print for a soundtrack
func (s *Soundtrack) PrettyPrint() {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err == nil {
		fmt.Println(string(bytes))
	}
}
