package services

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const (
	DOWNLOADS_FOLDER = "downloadsFolder"
	STORE_PATH       = "storePath"
	Magenta          = "\033[35m"
	Reset            = "\033[0m"
)

type Item struct {
	Name         string
	URL          string
	Owner        string
	Views        string
	Content      []string
	DownloadPath string
}

func (i Item) Title() string       { return i.Name }
func (i Item) Description() string { return i.Owner + " -- " + i.Views }
func (i Item) FilterValue() string { return i.Name }

func ClearStringForSearchQuery(input string) string {
	cleaned := strings.ReplaceAll(input, " ", "+")
	encoded := url.QueryEscape(cleaned)

	return encoded
}

func GetFromRegex(bodyString, pattern string) string {
	regex := regexp.MustCompile(pattern)
	match := regex.FindStringSubmatch(bodyString)

	var result string

	if len(match) >= 2 {
		result = match[1]
	}

	return result
}

func CheckIfPathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func getItemPath(item Item) string {
	if item.DownloadPath != "" {
		return item.DownloadPath
	}
	return item.URL
}

func SecondsToHHMMSS(seconds float64) string {
	minutes := int(seconds / 60)
	seconds -= float64(minutes) * 60

	return fmt.Sprintf("%02d:%02d", minutes, int(seconds))
}
