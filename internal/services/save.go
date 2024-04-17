package services

import (
	"encoding/json"
	"os"

	"github.com/spf13/viper"
)

func Save(jsonFile *[]Item, newPlaylist Item) error {
	*jsonFile = append(*jsonFile, newPlaylist)

	updatedJson, err := json.MarshalIndent(jsonFile, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(viper.GetString(STORE_PATH), updatedJson, 0644)

	if err != nil {
		return err
	}

	return nil
}
