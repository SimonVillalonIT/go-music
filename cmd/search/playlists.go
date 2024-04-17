package search

import (
	cmdhandlers "github.com/SimonVillalonIT/music-golang/internal/cmdHandlers"
	"github.com/spf13/cobra"
)

var playlistsCmd = &cobra.Command{
	Use:   "playlists",
	Short: "Search for playlists",
	Run:   cmdhandlers.SearchPlaylistHandler,
}

func init() {
	searchCmd.AddCommand(playlistsCmd)
}
