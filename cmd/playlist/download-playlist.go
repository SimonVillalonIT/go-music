package playlist

import (
	cmdhandlers "github.com/SimonVillalonIT/music-golang/internal/cmdHandlers"
	"github.com/spf13/cobra"
)

var downloadPlaylist = &cobra.Command{
	Use:   "download",
	Short: "Play from the playlist that you stored",
	Run:   cmdhandlers.DownloadPlaylistHandler,
}

func init() {
	playlistCmd.AddCommand(downloadPlaylist)
}
