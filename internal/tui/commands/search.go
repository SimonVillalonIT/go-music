package cmds

import (
	"github.com/SimonVillalonIT/music-golang/internal/services"
	tea "github.com/charmbracelet/bubbletea"
)

func SearchCmd(query string) tea.Cmd {
	return func() tea.Msg {
		songs, err := services.GetSongs(query, 10)
		if err != nil {
			return ErrMsg(err)
		}
		return SearchResultMsg(songs)
	}
}

func SaveCmd(jsonFile *[]services.Item, item services.Item) tea.Cmd {
	return func() tea.Msg {
		err := services.Save(jsonFile, item)
		if err != nil {
			return ErrMsg(err)
		}
		return nil
	}
}
