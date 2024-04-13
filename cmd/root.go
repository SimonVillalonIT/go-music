package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/SimonVillalonIT/music-golang/internal/services"
	"github.com/Songmu/prompter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "music-golang",
	Short: "The music manager built in Golang",
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.music-golang.json)")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".music-golang")
	}

	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Not Config file detected")
		generateConfig()
		return
	}

	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
}

func generateConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	downloadFolder := prompter.Prompt("Enter the folder that you want to use for downloading songs and playlists", "")
	_, err = os.Stat(downloadFolder)

	for err != nil {
		downloadFolder = prompter.Prompt("Enter the folder that you want to use for downloading songs and playlists", "")
		_, err = os.Stat(downloadFolder)
	}

	viper.Set(services.DOWNLOADS_FOLDER, downloadFolder)

	storePath := prompter.Prompt("Enter the location where you want the store to be", "")
	_, err = os.Stat(storePath)

	for err != nil {
		storePath = prompter.Prompt("Enter the location where you want the store to be", "")
		_, err = os.Stat(storePath)
	}

	storePath = path.Join(storePath, "store.json")

	file, err := os.Create(storePath)

	if err != nil {
		cobra.CheckErr(err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(services.JsonFile{}); err != nil {
		cobra.CheckErr(err)
	}

	viper.Set(services.STORE_PATH, storePath)

	configPath := prompter.Prompt("Enter the location for the config file", home)
	_, err = os.Stat(configPath)

	for err != nil {
		configPath = prompter.Prompt("Enter the location for the config file", home)
		_, err = os.Stat(storePath)
	}

	err = viper.WriteConfigAs(path.Join(configPath, ".music-golang.yaml"))
	viper.SetConfigFile(path.Join(configPath, ".music-golang.yaml"))
	cobra.CheckErr(err)
}
