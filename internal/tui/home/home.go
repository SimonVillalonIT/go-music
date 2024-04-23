package tui

import (
	"log"

	"github.com/DexterLB/mpvipc"
	"github.com/SimonVillalonIT/music-golang/internal/services"
	cmds "github.com/SimonVillalonIT/music-golang/internal/tui/commands"
	"github.com/SimonVillalonIT/music-golang/internal/tui/constants"
	custom_list "github.com/SimonVillalonIT/music-golang/internal/tui/list"
	fancyList "github.com/SimonVillalonIT/music-golang/internal/tui/list"
	"github.com/SimonVillalonIT/music-golang/internal/tui/player"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	list       fancyList.Model
	player     player.Model
	err        error
	help       tea.Model
	mode       uint
	jsonFile   []services.Item
	conn       mpvipc.Connection
	width      int
	height     int
	playerShow bool
}

func New() tea.Model {
	jsonfile, items, err := cmds.SearchSongs()
	if err != nil {
		log.Print(err)
	}

	list := fancyList.NewModel(items)

	return Model{list: list, jsonFile: jsonfile, playerShow: false}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(cmds.StartMpv, cmds.ConnCmd)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(m.width, m.height)
	case cmds.ConnErr:
		commands = append(commands, cmds.ConnCmd)
	case cmds.ConnMsg:
		m.conn = mpvipc.Connection(msg)
		m.player = player.NewModel(&m.conn)
		commands = append(commands, cmds.ObserveTrack(&m.conn))
		m.playerShow = true
	case cmds.MpvMsg:
		if msg == true {
			commands = append(commands, cmds.ConnCmd)
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
			if m.conn.IsClosed() {
				commands = append(commands, cmds.ConnCmd)
			} else {
				commands = append(commands, cmds.PlayCmd(m.conn, m.list.SelectedItem().(services.Item)))
			}
		}
	}

	updatedList, cmd := m.list.Update(msg)
	m.list = updatedList.(custom_list.Model)
	commands = append(commands, cmd)

	if m.playerShow {
		updatedPlayer, cmd := m.player.Update(msg)
		m.player = updatedPlayer.(player.Model)
		commands = append(commands, cmd)
	}

	return m, tea.Batch(commands...)
}

func (m Model) View() string {
	var view string
	list := m.list.View()
	if m.conn.IsClosed() {
		return list
	} else {
		view = lipgloss.Place(m.width, m.height, 0, 0, lipgloss.JoinVertical(lipgloss.Left, list, m.player.View()))
	}
	return view
}
