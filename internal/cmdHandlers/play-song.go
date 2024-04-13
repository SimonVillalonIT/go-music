package cmdhandlers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/SimonVillalonIT/music-golang/internal/services"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func PlaySongHandler(cmd *cobra.Command, args []string) {
	rawFile, err := os.ReadFile(viper.GetString(services.STORE_PATH))

	if err != nil {
		cobra.CheckErr(err)
	}

	var jsonFile services.JsonFile

	err = json.Unmarshal(rawFile, &jsonFile)
	if err != nil {
		cobra.CheckErr(err)
	}

	songs := jsonFile.Songs

	idx, err := fuzzyfinder.Find(
		songs,
		func(i int) string {
			return songs[i].Title
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Song: %s (%s)\n",
				songs[i].Title,
				songs[i].URL,
			)
		}))

	if err != nil {
		cobra.CheckErr(err)
	}

	var songPath string

	if songs[idx].DownloadPath != "" {
		songPath = songs[idx].DownloadPath
		fmt.Println(songPath)
	} else {
		songPath = songs[idx].URL
	}

	command := exec.Command("mpv", "--no-video", songPath)

	// Set output to current console
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// Run the command
	if err := command.Run(); err != nil {
		cobra.CheckErr(err)
		return
	}

}
