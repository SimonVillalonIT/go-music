package services

import (
	"fmt"
	"io"
	"net/http"

	"regexp"
	"sync"

	"github.com/google/uuid"
)

func getSongsDetails(videoURL string) (string, string, string) {
	resp, err := http.Get(videoURL)
	if err != nil {
		return "", "", ""
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", ""
	}

	bodyString := string(bodyBytes)

	videoTitle := GetFromRegex(bodyString, `<title>(.*?) - YouTube<\/title>`)

	videoOwner := GetFromRegex(bodyString, `"author":"([^"]+)"`)

	videoViews := GetFromRegex(bodyString, `"viewCount":{"simpleText":"([^"]+)"}`)

	return videoTitle, videoOwner, videoViews
}

func GetSongs(keyword string, limit int) ([]Item, error) {
	var result []Item
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
				title, owner, views := getSongsDetails(url)

				result = append(result, Item{ID: uuid.NewString(), URL: url, Name: title, Owner: owner, Views: views})
			}(videoURL)
			counter++
		}
		mutex.Unlock()
	}

	wg.Wait()

	return result, nil
}
