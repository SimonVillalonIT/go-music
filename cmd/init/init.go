package init

import (
	"encoding/json"
	"os"
	"path"

	"github.com/SimonVillalonIT/music-golang/cmd"
	"github.com/SimonVillalonIT/music-golang/internal/services"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fileDir string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Start your configuration",
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		jsonFile := services.JsonFile{}

		filePath := path.Join(fileDir, "store.json")

		file, err := os.Create(filePath)

		if err != nil {
			cobra.CheckErr(err)
		}

		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "    ")

		if err := encoder.Encode(jsonFile); err != nil {
			cobra.CheckErr(err)
		}

		viper.Set("store", filePath)

		viper.SafeWriteConfigAs(home + "/.music-golang.yaml")
	},
}

func init() {
	initCmd.Flags().StringVarP(&fileDir, "fileDir", "f", ".", "Directory to create store")
	initCmd.MarkFlagDirname("fileDir")
	initCmd.MarkFlagRequired("fileDir")
	cmd.RootCmd.AddCommand(initCmd)
}
