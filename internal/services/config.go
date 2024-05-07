package services

import (
	"encoding/json"
	"os"
	"path"

	"github.com/spf13/viper"
)

func GenerateConfig(downloadFolder, storePath string) error {
	home, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	if !path.IsAbs(downloadFolder) {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		downloadFolder = path.Join(dir, downloadFolder)
	}
	viper.Set(DOWNLOADS_FOLDER, downloadFolder)

	if !path.IsAbs(storePath) {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		storePath = path.Join(dir, storePath)
	}

	storePath = path.Join(storePath, "store.json")

	file, err := os.Create(storePath)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode([]Item{}); err != nil {
		return err
	}

	viper.Set(STORE_PATH, storePath)

	err = viper.WriteConfigAs(path.Join(home, ".go-music.yaml"))
	viper.SetConfigFile(path.Join(home, ".go-music.yaml"))

	return err
}
