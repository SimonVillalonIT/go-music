package playlist

import (
	"github.com/DexterLB/mpvipc"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	Current      bool   `json:"current"`
	Filename     string `json:"filename"`
	ID           int    `json:"id"`
	Playing      bool   `json:"playing"`
	PlaylistPath string `json:"playlist-path"`
}

type Model struct {
	Conn  *mpvipc.Connection
	Items []item
}

func New(conn *mpvipc.Connection) Model {
	return Model{Conn: conn}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return ""
}

