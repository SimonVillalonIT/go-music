package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/spf13/viper"
)

type Song struct {
	URL   string
	Title string
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

	err = os.WriteFile(viper.GetString("store"), updatedJson, 0644)

	if err != nil {
		return err
	}

	return nil
}
