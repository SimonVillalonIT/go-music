package services

import (
	"net/url"
	"strings"
)

const (
	DOWNLOADS_FOLDER = "downloadsFolder"
	STORE_PATH       = "storePath"
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
