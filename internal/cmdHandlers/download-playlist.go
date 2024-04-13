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

func DownloadPlaylistHandler(cmd *cobra.Command, args []string) {
	rawFile, err := os.ReadFile(viper.GetString(services.STORE_PATH))

	if err != nil {
		cobra.CheckErr(err)
	}

	var jsonFile services.JsonFile

	err = json.Unmarshal(rawFile, &jsonFile)

	playlists := jsonFile.Playlists

	idx, err := fuzzyfinder.Find(
		playlists,
		func(i int) string {
			return playlists[i].Title
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Song: %s (%s)\n",
				playlists[i].Title,
				playlists[i].URL,
			)
		}))

	if err != nil {
		cobra.CheckErr(err)
	}

	err = services.DownloadPlaylist(&jsonFile, jsonFile.Playlists[idx])

	if err != nil {
		cobra.CheckErr(err)
	}
}
