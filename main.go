package main

import (
	"fmt"
	"os"

	"github.com/SimonVillalonIT/music-golang/cmd"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	cmd.Execute()
}
