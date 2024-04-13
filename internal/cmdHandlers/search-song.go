package cmdhandlers

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/SimonVillalonIT/music-golang/internal/services"
	"github.com/Songmu/prompter"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SearchSongHandler(cmd *cobra.Command, args []string) {
	search := prompter.Prompt("Search for a song", "")

	intPattern := `^\d+$`

	regex := regexp.MustCompile(intPattern)
	limit := prompter.Prompt("Enter a limit of result:", "10")

	for !regex.MatchString(limit) {
		fmt.Println("The limit must be a valid number!")
		limit = prompter.Prompt("Enter a limit of result:", "10")
	}

	limitInt, err := strconv.Atoi(limit)

	cobra.CheckErr(err)

	result, err := services.GetSongs(search, limitInt)

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
			return fmt.Sprintf("Song: %s (%s)\n",
				result[i].Title,
				result[i].URL,
			)
		}))

	if err != nil {
		cobra.CheckErr(err)
	}

	rawFile, err := os.ReadFile(viper.GetString("store"))

	if err != nil {
		cobra.CheckErr(err)
	}

	var jsonFile services.JsonFile

	err = json.Unmarshal(rawFile, &jsonFile)

	services.SaveSong(jsonFile, result[idx])
}
