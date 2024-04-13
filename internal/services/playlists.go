package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"

	"github.com/spf13/viper"
)

type Playlist struct {
	Title        string
	URL          string
	Content      []string
	DownloadPath string
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

	err = os.WriteFile(viper.GetString(STORE_PATH), updatedJson, 0644)

	if err != nil {
		return err
	}

	return nil
}

func DownloadPlaylist(jsonFile *JsonFile, playlist Playlist) error {
	found := false
	playlistPath := path.Join(viper.GetString(DOWNLOADS_FOLDER), playlist.Title)
	err := os.Mkdir(playlistPath, 0700)

	if err != nil {
		return err
	}

	for i, currentPlaylist := range jsonFile.Playlists {
		if playlist.Title == currentPlaylist.Title {
			jsonFile.Playlists[i] = playlist
			jsonFile.Playlists[i].DownloadPath = playlistPath
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("playlist with title %q not found", playlist.Title)
	}

	updatedData, err := json.MarshalIndent(jsonFile, "", "    ")
	if err != nil {
		return err
	}

	command := exec.Command("yt-dlp",
		"--extract-audio",
		"--audio-format", "mp3",
		"--output", path.Join(playlistPath, "%(title)s.%(ext)s"),
		"--playlist-start", "1",
		"--playlist-end", "10",
		playlist.URL,
	)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// Run the command
	if err := command.Run(); err != nil {
		return err
	}
	// Write the updated JSON data back to the file
	if err := os.WriteFile(viper.GetString(STORE_PATH), updatedData, 0644); err != nil {
		return err
	}

	return nil
}
