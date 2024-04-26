package services

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path"
)

func Download(jsonFile *[]Item, item Item) error {
	found := false
	downloadPath := path.Join(viper.GetString(DOWNLOADS_FOLDER), item.Name)

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

	var command *exec.Cmd

	if item.Content != nil {
		err := os.Mkdir(downloadPath, 0700)

		if err != nil {
			return err
		}

		command = exec.Command("yt-dlp",
			"--extract-audio",
			"--audio-format", "mp3",
			"--output", path.Join(downloadPath, "%(title)s.%(ext)s"),
			item.URL,
		)
	} else {
		command = exec.Command("yt-dlp",
			"--extract-audio",
			"--audio-format", "mp3",
			"--output", downloadPath+".mp3",
			item.URL,
		)
	}

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return err
	}
	if err := os.WriteFile(viper.GetString(STORE_PATH), updatedData, 0644); err != nil {
		return err
	}

	return nil
}