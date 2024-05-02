package cmds

import (
	"github.com/SimonVillalonIT/music-golang/internal/services"
	tea "github.com/charmbracelet/bubbletea"
)

func SearchCmd(query string, kind string) tea.Cmd {
	return func() tea.Msg {
		if kind == "Playlist" {
			songs, err := services.GetPlaylists(query, 10)
			if err != nil {
				return ErrMsg(err)
			}
			return SearchResultMsg(songs)

		}
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
		return SearchSuccessMsg(true)
	}
}

func QuitSearchCmd() tea.Msg {
	return QuitSearchMsg(true)
}
