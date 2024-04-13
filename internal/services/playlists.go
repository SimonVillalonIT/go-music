package services

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/spf13/viper"
)

type Playlist struct {
	Title   string
	URL     string
	Content []string
}

func getPlaylistDetails(url string) (string, []string) {
	var content []string
	resp, err := http.Get(url)
	if err != nil {
		return "", content
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", content
	}

	bodyString := string(bodyBytes)

	nameRegex := regexp.MustCompile(`<title>(.*?) - YouTube<\/title>`)

	nameMatch := nameRegex.FindStringSubmatch(bodyString)

	contentRegex := regexp.MustCompile(`"title":{"runs":\[\{"text":"([^"]+)"\}`)

	contentMatches := contentRegex.FindAllStringSubmatch(bodyString, -1)

	for _, match := range contentMatches {
		title := match[1]
		if title == "Descripción" {
			break
		}
		content = append(content, title)
	}

	// Extract playlist name and author
	var playlistName string
	if len(nameMatch) >= 2 {
		playlistName = nameMatch[1]
	}

	return playlistName, content
}

func GetPlaylists(keyword string, limit int) ([]Playlist, error) {
	var results []Playlist
	keyword = ClearStringForSearchQuery(keyword)
	resp, err := http.Get("https://www.youtube.com/results?search_query=" + keyword + "&sp=EgIQAw%253D%253D")
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return results, err
	}

	bodyString := string(bodyBytes)

	// Regular expression pattern to match playlist IDs
	playlistIDPattern := `\/playlist\?list=([A-Za-z0-9_-]+)`

	regex := regexp.MustCompile(playlistIDPattern)
	playlistIDs := regex.FindAllStringSubmatch(bodyString, -1)

	// Create a map to store unique playlist IDs
	uniqueIDs := make(map[string]bool)

	counter := 0

	for _, match := range playlistIDs {
		if counter >= limit {
			break
		}

		playlistID := match[1]
		// Check if the ID is already in the map
		if _, ok := uniqueIDs[playlistID]; !ok {
			// If not, add it to the map and construct the playlist URL
			uniqueIDs[playlistID] = true
			url := "https://www.youtube.com/playlist?list=" + playlistID
			// Fetch playlist details
			title, content := getPlaylistDetails(url)
			newPlaylist := Playlist{Title: title, URL: url, Content: content}
			results = append(results, newPlaylist)
		}
		counter++
	}
	return results, err
}

func SavePlaylist(jsonFile JsonFile, newPlaylist Playlist) error {
	jsonFile.Playlists = append(jsonFile.Playlists, newPlaylist)

	updatedJson, err := json.MarshalIndent(jsonFile, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(viper.GetString("store"), updatedJson, 0644)

	if err != nil {
		return err
	}

	return nil
}
