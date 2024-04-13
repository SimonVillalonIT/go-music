package playlist

import (
	cmdhandlers "github.com/SimonVillalonIT/music-golang/internal/cmdHandlers"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchPlaylistCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a playlist",
	Run:   cmdhandlers.SearchPlaylistHandler,
}

func init() {
	playlistCmd.AddCommand(searchPlaylistCmd)
}
