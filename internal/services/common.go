package services

import (
	"net/url"
	"strings"
)

type JsonFile struct {
	Songs     []Song
	Playlists []Playlist
}

func ClearStringForSearchQuery(input string) string {
	// Replace spaces with '+' signs
	cleaned := strings.ReplaceAll(input, " ", "+")

	// Encode the string so that special characters are properly handled
	encoded := url.QueryEscape(cleaned)

	return encoded
}
