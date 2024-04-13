package playlist

import (
	cmdhandlers "github.com/SimonVillalonIT/music-golang/internal/cmdHandlers"
	"github.com/spf13/cobra"
)

var playPlaylistCmd = &cobra.Command{
	Use:   "play",
	Short: "Play from the playlist that you stored",
	Run:   cmdhandlers.PlayPlaylistHandler,
}

func init() {
	playlistCmd.AddCommand(playPlaylistCmd)
}
