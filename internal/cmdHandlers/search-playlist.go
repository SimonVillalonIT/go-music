package cmdhandlers

import (
	"encoding/json"
	"os"

	"github.com/SimonVillalonIT/music-golang/internal/services"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SearchPlaylistHandler(cmd *cobra.Command, args []string) {
	search, limitInt, err := services.Prompt("Search for a playlist")

	if err != nil {
		cobra.CheckErr(err)
	}

	result, err := services.GetPlaylists(search, limitInt)

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

	if err != nil {
		cobra.CheckErr(err)
	}

	for _, idx := range idxs {
		if err := services.Save(&jsonFile, result[idx]); err != nil {
			cobra.CheckErr(err)
		}
	}
}
