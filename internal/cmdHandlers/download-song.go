package cmdhandlers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SimonVillalonIT/music-golang/internal/services"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DownloadSongHandler(cmd *cobra.Command, args []string) {
	rawFile, err := os.ReadFile(viper.GetString(services.STORE_PATH))

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

	err = services.DownloadSong(&jsonFile, jsonFile.Songs[idx])

	if err != nil {
		cobra.CheckErr(err)
	}

	fmt.Println("File Downloaded Succesfully")
}
