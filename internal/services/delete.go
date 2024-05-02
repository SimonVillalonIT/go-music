package services

import (
	"encoding/json"
	"os"

	"github.com/spf13/viper"
)

func Delete(jsonFile *[]Item, item Item) error {
	if item.DownloadPath != "" {
		var err error
		info, err := os.Stat(item.DownloadPath)
		if err != nil {
			return err
		}

		if info.IsDir() {
			err = os.RemoveAll(item.DownloadPath)
		} else {
			err = os.Remove(item.DownloadPath)
		}

		if err != nil {
			return err
		}
	}

	var newJsonFile []Item
	for _, i := range *jsonFile {
		if i.Name != item.Name {
			newJsonFile = append(newJsonFile, i)
		}
	}

	updatedJsonFile, err := json.MarshalIndent(newJsonFile, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(viper.GetString(STORE_PATH), updatedJsonFile, 0644)

	if err != nil {
		return err
	}

	return nil
}
