package song

import (
	"github.com/SimonVillalonIT/music-golang/cmd"
	"github.com/spf13/cobra"
)

var songCmd = &cobra.Command{
	Use:   "song",
	Short: "Handle all your songs",
}

func init() {
	cmd.RootCmd.AddCommand(songCmd)
}
