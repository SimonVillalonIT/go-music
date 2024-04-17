package search

import (
	cmdhandlers "github.com/SimonVillalonIT/music-golang/internal/cmdHandlers"
	"github.com/spf13/cobra"
)

var songCmd = &cobra.Command{
	Use:   "songs",
	Short: "Search for songs",
	Run:   cmdhandlers.SearchSongHandler,
}

func init() {
	searchCmd.AddCommand(songCmd)
}
