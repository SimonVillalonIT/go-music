package tui

import (
	"log"

	"github.com/DexterLB/mpvipc"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/SimonVillalonIT/music-golang/internal/services"
	cmds "github.com/SimonVillalonIT/music-golang/internal/tui/commands"
	"github.com/SimonVillalonIT/music-golang/internal/tui/constants"
	fancyList "github.com/SimonVillalonIT/music-golang/internal/tui/list"
	"github.com/SimonVillalonIT/music-golang/internal/tui/player"
)

type homeState uint

const (

)

type Model struct {
	width       int
	height      int
	list        fancyList.Model
	player      player.Model
	err         error
	jsonFile    []services.Item
	conn        *mpvipc.Connection
	isConnected bool
}

func New() tea.Model {
	jsonfile, items, err := cmds.SearchSongs()
	if err != nil {
		log.Print(err)
	}

	list := fancyList.NewModel(items)

	return Model{list: list, jsonFile: jsonfile, isConnected: false}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(m.width, m.height)
	case cmds.ConnMsg:
		if msg != nil {
			m.conn = msg
			m.isConnected = true
		}
	case cmds.ErrMsg:
		log.Print(msg)
		m.err = msg
	case tea.KeyMsg:
		if key.Matches(msg, constants.Keymap.Quit) {
			cmds.Kill()
			return m, tea.Quit
		}
		if key.Matches(msg, constants.Keymap.Enter) {
			if m.isConnected {
				commands = append(commands, cmds.PlayCmd(m.conn, m.list.SelectedItem().(services.Item)))
			}
		}
		if key.Matches(msg, constants.Keymap.Space) {
			if m.player.Playlist.Length > 0 {
				commands = append(commands, cmds.PauseCmd(m.conn))
			}
		}
		if key.Matches(msg, constants.Keymap.Next) {
			if m.player.Playlist.Length > 0 {
				commands = append(commands, cmds.NextCmd(m.conn))
			}
		}
		if key.Matches(msg, constants.Keymap.Prev) {
			if m.player.Playlist.Length > 0 {
				commands = append(commands, cmds.PrevCmd(m.conn))
			}
		}
		if key.Matches(msg, constants.Keymap.Minus) {
			if m.player.Playlist.Length > 0 {
				commands = append(commands, cmds.DecreaseCmd(m.conn))
			}
		}
		if key.Matches(msg, constants.Keymap.Plus) {
			if m.player.Playlist.Length > 0 {
				commands = append(commands, cmds.IncreaseCmd(m.conn))
			}
		}
	}
	updatedList, cmd := m.list.Update(msg)
	m.list = updatedList.(fancyList.Model)

	updatedPlayer, cmd := m.player.Update(msg)
	m.player = updatedPlayer.(player.Model)
	commands = append(commands, cmd)

	return m, tea.Batch(commands...)
}

func (m Model) View() string {
	list := m.list.View()
	if !m.isConnected {
		return "Loading..."
	}
	if m.conn.IsClosed() {
		return "Loading..."
	}
	return lipgloss.Place(m.width, m.height, 0, 0, lipgloss.JoinVertical(lipgloss.Left, list, m.player.View()))
}
