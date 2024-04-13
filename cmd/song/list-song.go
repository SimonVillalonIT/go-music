package song

import (
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Handle all your songs",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	songCmd.AddCommand(listCommand)
}
