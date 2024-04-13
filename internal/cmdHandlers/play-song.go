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
	rawFile, err := os.ReadFile(viper.GetString("store"))

	if err != nil {
		cobra.CheckErr(err)
	}

	var jsonFile services.JsonFile

	err = json.Unmarshal(rawFile, &jsonFile)

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

	command := exec.Command("mpv", "--no-video", songs[idx].URL)

	// Set output to current console
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// Run the command
	if err := command.Run(); err != nil {
		cobra.CheckErr(err)
		return
	}

}
