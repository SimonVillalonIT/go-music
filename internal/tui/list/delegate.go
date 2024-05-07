package custom_list

import (
	cmds "github.com/SimonVillalonIT/go-music/internal/tui/commands"
	"github.com/SimonVillalonIT/go-music/internal/tui/constants"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func newItemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		switch msg := msg.(type) {
		case cmds.DownloadSuccessMsg:
			m.StopSpinner()
			return m.NewStatusMessage(statusMessageStyle("Downloaded successfully"))
		case cmds.DownloadErrorMsg:
			m.StopSpinner()
			return m.NewStatusMessage(errorMessageStyle(string(msg)))
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, constants.Keymap.Play):
				return m.NewStatusMessage(statusMessageStyle("Added to playlist"))
			case key.Matches(msg, constants.Keymap.Remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					constants.Keymap.Remove.SetEnabled(false)
				}
				return m.NewStatusMessage(statusMessageStyle("Removed from playlist"))
			}
		}

		return nil
	}

	return d
}
