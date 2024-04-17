package cmdhandlers

import (
	"encoding/json"
	"os"

	"github.com/SimonVillalonIT/music-golang/internal/services"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DownloadHandler(cmd *cobra.Command, args []string) {
	rawFile, err := os.ReadFile(viper.GetString(services.STORE_PATH))

	if err != nil {
		cobra.CheckErr(err)
	}

	var jsonFile []services.Item

	err = json.Unmarshal(rawFile, &jsonFile)

	if err != nil {
		cobra.CheckErr(err)
	}

	idxs, err := services.DisplayFuzzyFind(jsonFile)

	if err != nil {
		cobra.CheckErr(err)
	}

	for _, idx := range idxs {
		if err = services.Download(&jsonFile, jsonFile[idx]); err != nil {
			cobra.CheckErr(err)
		}
	}

}
