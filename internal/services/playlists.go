package services

import (
	"io"
	"net/http"
	"regexp"
)

func getPlaylistDetails(url string) (string, string, []string) {
	var content []string
	resp, err := http.Get(url)
	if err != nil {
		return "", "", content
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", content
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

	ownerPattern := `\"ownerText\":\{\"runs\":\[\{\"text\":\"([^\"]+)\"`

	// Compile the regex pattern
	ownerRegex := regexp.MustCompile(ownerPattern)

	// Find the matches
	ownerMatches := ownerRegex.FindStringSubmatch(bodyString)

	var videoOwner string

	if len(ownerMatches) >= 2 {
		videoOwner = ownerMatches[1]
	}

	return playlistName, videoOwner, content
}

func GetPlaylists(keyword string, limit int) ([]Item, error) {
	var results []Item
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
			title, owner, content := getPlaylistDetails(url)
			newPlaylist := Item{Name: title, URL: url, Owner: owner, Content: content}
			results = append(results, newPlaylist)
		}
		counter++
	}
	return results, err
}
