package player

import (
	"fmt"

	"github.com/SimonVillalonIT/music-golang/internal/services"
	cmds "github.com/SimonVillalonIT/music-golang/internal/tui/commands"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Playlist cmds.Playlist
	width    int
	heigth   int
}

func NewModel() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.heigth = msg.Height
		m.width = msg.Width
	case cmds.TrackNameMsg:
		m.Playlist.CurrentTrack.Name = string(msg)
	case cmds.TrackCurrentFrameMsg:
		m.Playlist.CurrentTrack.CurrentFrame = float64(msg)
	case cmds.TrackDurationMsg:
		m.Playlist.CurrentTrack.Duration = float64(msg)

	case cmds.PlaylistPositionMsg:
		m.Playlist.Position = float64(msg)
	case cmds.PlaylistLengthMsg:
		m.Playlist.Length = float64(msg)
	}
	return m, tea.Batch(commands...)
}

func (m Model) View() string {
	if m.Playlist.Length < 1 {
		return ""
	}
	data := fmt.Sprintf("Playing: %s   %s/%s | %.0f/%.0f", m.Playlist.CurrentTrack.Name, services.SecondsToHHMMSS(m.Playlist.CurrentTrack.CurrentFrame), services.SecondsToHHMMSS(m.Playlist.CurrentTrack.Duration), m.Playlist.Position, m.Playlist.Length)
	return lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Width(m.width).Render(lipgloss.JoinHorizontal(lipgloss.Left), data)
}
