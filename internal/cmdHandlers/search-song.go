package cmdhandlers

import (
	"encoding/json"
	"os"

	"github.com/SimonVillalonIT/go-music/internal/services"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SearchSongHandler(cmd *cobra.Command, args []string) {
	search, limitInt, err := services.Prompt("Search for a song")

	if err != nil {
		cobra.CheckErr(err)
	}

	result, err := services.GetSongs(search, limitInt)

	if err != nil {
		cobra.CheckErr(err)
	}

	idxs, err := services.DisplayFuzzyFind(result)

	if err != nil {
		cobra.CheckErr(err)
	}

	rawFile, err := os.ReadFile(viper.GetString(services.STORE_PATH))

	if err != nil {
		cobra.CheckErr(err)
	}

	var jsonFile []services.Item

	err = json.Unmarshal(rawFile, &jsonFile)
	for _, idx := range idxs {
		if err = services.Save(&jsonFile, result[idx]); err != nil {
			cobra.CheckErr(err)
		}
	}
}
