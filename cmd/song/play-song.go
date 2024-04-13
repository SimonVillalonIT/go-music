package song

import (
	cmdhandlers "github.com/SimonVillalonIT/music-golang/internal/cmdHandlers"
	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Handle all your songs",
	Run:   cmdhandlers.PlaySongHandler,
}

func init() {
	songCmd.AddCommand(playCmd)
}
