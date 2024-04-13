package playlist

import (
	"github.com/SimonVillalonIT/music-golang/cmd"
	"github.com/spf13/cobra"
)

var playlistCmd = &cobra.Command{
	Use:   "playlist",
	Short: "Handle all your playlists",
}

func init() {
	cmd.RootCmd.AddCommand(playlistCmd)
}
