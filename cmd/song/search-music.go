package song

import (
	cmdhandlers "github.com/SimonVillalonIT/music-golang/internal/cmdHandlers"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a video",
	Run:   cmdhandlers.SearchSongHandler,
}

func init() {
	songCmd.AddCommand(searchCmd)
}
