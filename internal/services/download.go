package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/viper"
)

func Download(jsonFile *[]Item, item Item) error {
	found := false
	downloadPath := path.Join(viper.GetString(DOWNLOADS_FOLDER), item.Name)

	if _, err := os.Stat(downloadPath + ".mp3"); !os.IsNotExist(err) {
		return fmt.Errorf("Song already downloaded!")
	}

	for i, currentPlaylist := range *jsonFile {
		if item.Name == currentPlaylist.Name {
			if (*jsonFile)[i].Content != nil {
				(*jsonFile)[i].DownloadPath = downloadPath
				found = true
				break
			} else {
				(*jsonFile)[i].DownloadPath = downloadPath + ".mp3"
				found = true
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("item with title %q not found", item.Name)
	}

	updatedData, err := json.MarshalIndent(jsonFile, "", "    ")
	if err != nil {
		return err
	}

	if item.Content != nil {
		err := os.Mkdir(downloadPath, 0700)

		if err != nil {
			return err
		}

		err = downloadYouTubeVideo(item.URL, path.Join(downloadPath, "%(title)s.mp3"))
		if err != nil {
			log.Println(err)
			return err
		}

	} else {
		err = downloadYouTubeVideo(item.URL, downloadPath+".mp3")
		if err != nil {
			log.Println(err)
			return err
		}
	}

	if err := os.WriteFile(viper.GetString(STORE_PATH), updatedData, 0644); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func downloadYouTubeVideo(url string, outputPath string) error {
    cmd := exec.Command("mpv", "--ytdl-format=best","--no-terminal", "--no-video", "-o", outputPath, url)
	return cmd.Run()
}
