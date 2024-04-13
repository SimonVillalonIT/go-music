package song

import (
	cmdhandlers "github.com/SimonVillalonIT/music-golang/internal/cmdHandlers"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a song",
	Run:   cmdhandlers.DownloadSongHandler,
}

func init() {
	songCmd.AddCommand(downloadCmd)
}
