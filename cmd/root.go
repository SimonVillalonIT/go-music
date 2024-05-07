package cmd

import (
	"os"

	"github.com/SimonVillalonIT/music-golang/internal/services"
	"github.com/SimonVillalonIT/music-golang/internal/tui"
	cmds "github.com/SimonVillalonIT/music-golang/internal/tui/commands"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "music-golang",
	Short: "The music manager built in Golang",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(tui.NewMainModel(tui.HomeView), tea.WithAltScreen())

		p.SetWindowTitle("Go-Music")

		go services.StartMpv()

		go cmds.SetupConnection(p)

		if _, err := p.Run(); err != nil {
			cobra.CheckErr(err)
		}
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(InitConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.music-golang.json)")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func InitConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigName(".music-golang")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		p := tea.NewProgram(tui.NewMainModel(tui.ConfigView), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			cobra.CheckErr(err)
		}
	}
}
