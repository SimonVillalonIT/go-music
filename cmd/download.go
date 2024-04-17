package cmd

import (
	cmdhandlers "github.com/SimonVillalonIT/music-golang/internal/cmdHandlers"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download any stored song or playlist",
	Run:   cmdhandlers.DownloadHandler,
}

func init() {
	RootCmd.AddCommand(downloadCmd)
}
