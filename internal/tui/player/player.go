package player

import (
	"github.com/DexterLB/mpvipc"
	cmds "github.com/SimonVillalonIT/music-golang/internal/tui/commands"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	title    string
	author   string
	duration string
	width    int
	heigth   int
	Conn     *mpvipc.Connection
}

func NewModel(conn *mpvipc.Connection) Model {
	return Model{Conn: conn}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd
	switch msg := msg.(type) {
	case cmds.EventMsg:
		m.title = string(msg)
		commands = append(commands, cmds.ObserveTrack(m.Conn))
	case cmds.RetryMsg:
		commands = append(commands, cmds.ObserveTrack(m.Conn))
	case tea.WindowSizeMsg:
		m.heigth = msg.Height
		m.width = msg.Width
	}
	return m, tea.Batch(commands...)
}

func (m Model) View() string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Width(m.width).Render(lipgloss.JoinHorizontal(lipgloss.Left, "Playing: ", m.title))
}
