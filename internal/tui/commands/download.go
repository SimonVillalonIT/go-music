package cmds

import (
	"time"

	"github.com/SimonVillalonIT/go-music/internal/services"
	tea "github.com/charmbracelet/bubbletea"
)

func DownloadCmd(jsonFile *[]services.Item, item services.Item) tea.Cmd {
	return func() tea.Msg {
		err := services.Download(jsonFile, item)
		if err != nil {
			return DownloadErrorMsg(err.Error())
		}
		return DownloadSuccessMsg(true)
	}
}

func ClearDownloadErrorCmd() tea.Cmd {
	return func() tea.Msg {
		return tea.Tick(1*time.Second, func(time.Time) tea.Msg {
			return ClearDownloadError(true)
		})
	}
}
