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
	"sync"

	"github.com/spf13/viper"
)

type Song struct {
	URL          string
	Title        string
	DownloadPath string
}

func getSongsDetails(videoURL string) string {
	resp, err := http.Get(videoURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	bodyString := string(bodyBytes)

	namePattern := `<title>(.*?) - YouTube<\/title>`
	regex := regexp.MustCompile(namePattern)
	match := regex.FindStringSubmatch(bodyString)

	var videoTitle string

	if len(match) >= 2 {
		videoTitle = match[1]
	}

	return videoTitle
}

func GetSongs(keyword string, limit int) ([]Song, error) {
	var result []Song
	search := "https://www.youtube.com/results?search_query=" + ClearStringForSearchQuery(keyword) + "&sp=EgIQAQ%253D%253D"
	resp, err := http.Get(search)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	bodyString := string(bodyBytes)

	videoIDPattern := `"videoId":"([A-Za-z0-9_-]{11})"`

	regex := regexp.MustCompile(videoIDPattern)
	videoIDs := regex.FindAllStringSubmatch(bodyString, -1)

	uniqueURLs := make(map[string]bool)

	var wg sync.WaitGroup
	var mutex sync.Mutex

	counter := 0

	for _, match := range videoIDs {
		if counter >= limit {
			break
		}

		videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", match[1])
		mutex.Lock()
		if _, ok := uniqueURLs[videoURL]; !ok {
			uniqueURLs[videoURL] = true
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				title := getSongsDetails(url)

				result = append(result, Song{URL: url, Title: title})
			}(videoURL)
			counter++
		}
		mutex.Unlock()
	}

	wg.Wait()

	return result, nil
}

func SaveSong(jsonFile JsonFile, newSong Song) error {
	jsonFile.Songs = append(jsonFile.Songs, newSong)

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

func DownloadSong(jsonFile *JsonFile, song Song) error {
	found := false
	songPath := path.Join(viper.GetString(DOWNLOADS_FOLDER), song.Title)

	for i, currentSong := range jsonFile.Songs {
		if song.Title == currentSong.Title {
			jsonFile.Songs[i] = song
			jsonFile.Songs[i].DownloadPath = songPath + ".mp3"
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("song with title %q not found", song.Title)
	}

	updatedData, err := json.MarshalIndent(jsonFile, "", "    ")
	if err != nil {
		return err
	}

	command := exec.Command("yt-dlp",
		"--extract-audio",
		"--audio-format", "mp3",
		"--output", songPath,
		"--playlist-start", "1",
		"--playlist-end", "10",
		song.URL,
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
