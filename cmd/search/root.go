package search

import (
	"github.com/SimonVillalonIT/music-golang/cmd"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for songs",
}

func init() {
	cmd.RootCmd.AddCommand(searchCmd)
}
