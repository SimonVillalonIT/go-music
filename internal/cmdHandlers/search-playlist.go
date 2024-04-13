package cmdhandlers

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/SimonVillalonIT/music-golang/internal/services"
	"github.com/Songmu/prompter"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SearchPlaylistHandler(cmd *cobra.Command, args []string) {
	search := prompter.Prompt("Search for a playlist", "")

	intPattern := `^\d+$`

	regex := regexp.MustCompile(intPattern)
	limit := prompter.Prompt("Enter a limit of result:", "10")

	for !regex.MatchString(limit) {
		fmt.Println("The limit must be a valid number!")
		limit = prompter.Prompt("Enter a limit of result:", "10")
	}

	limitInt, err := strconv.Atoi(limit)

	cobra.CheckErr(err)

	result, err := services.GetPlaylists(search, limitInt)

	if err != nil {
		cobra.CheckErr(err)
	}

	idx, err := fuzzyfinder.Find(
		result,
		func(i int) string {
			return result[i].Title
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Playlist: %s\nContent:\n%s",
				result[i].Title,
				strings.Join(result[i].Content, "\n"),
			)
		}))

	if err != nil {
		cobra.CheckErr(err)
	}

	rawFile, err := os.ReadFile(viper.GetString(services.STORE_PATH))

	if err != nil {
		cobra.CheckErr(err)
	}

	var jsonFile services.JsonFile

	err = json.Unmarshal(rawFile, &jsonFile)

	services.SavePlaylist(jsonFile, result[idx])
}
