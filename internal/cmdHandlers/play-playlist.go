package cmdhandlers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/SimonVillalonIT/music-golang/internal/services"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func PlayPlaylistHandler(cmd *cobra.Command, args []string) {
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
			return fmt.Sprintf("Playlist: %s\nContent:\n%s",
				playlists[i].Title,
				strings.Join(playlists[i].Content, "\n"),
			)
		}))

	if err != nil {
		cobra.CheckErr(err)
	}

	var playlistPath string

	if playlists[idx].DownloadPath != "" {
		playlistPath = playlists[idx].DownloadPath
	} else {
		playlistPath = playlists[idx].URL
	}

	command := exec.Command("mpv", "--no-video", playlistPath)

	// Set output to current console
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// Run the command
	if err := command.Run(); err != nil {
		cobra.CheckErr(err)
		return
	}

}
