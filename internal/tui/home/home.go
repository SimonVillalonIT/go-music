package tui

import (
	"github.com/DexterLB/mpvipc"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/SimonVillalonIT/go-music/internal/services"
	cmds "github.com/SimonVillalonIT/go-music/internal/tui/commands"
	"github.com/SimonVillalonIT/go-music/internal/tui/constants"
	fancyList "github.com/SimonVillalonIT/go-music/internal/tui/list"
	"github.com/SimonVillalonIT/go-music/internal/tui/player"
	"github.com/SimonVillalonIT/go-music/internal/tui/playlist"
	"github.com/SimonVillalonIT/go-music/internal/tui/search"
)

type Model struct {
	state           constants.SessionState
	width           int
	height          int
	list            fancyList.Model
	player          player.Model
	playlist        playlist.Model
	help            help.Model
	search          search.Model
	err             error
	jsonFile        []services.Item
	conn            *mpvipc.Connection
	isConnected     bool
	downloadLoading bool
	showHelp        bool
}

func New() tea.Model {
	jsonfile, items, _ := cmds.SearchSongs()

	list := fancyList.NewModel(items)
	player := player.NewModel()
	playlist := playlist.NewModel()
	search := search.NewModel()
	help := help.New()

	return Model{list: list, jsonFile: jsonfile, isConnected: false, player: player, playlist: playlist, search: search, downloadLoading: false, help: help, showHelp: true}
}

func (m Model) Init() tea.Cmd {
	return ChangeState((*uint)(&m.state))
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
		m.err = msg
	case cmds.DownloadSuccessMsg:
		m.updateData()
		m.list.StopSpinner()
	case cmds.SearchSuccessMsg:
		m.updateData()
	case cmds.QuitSearchMsg:
		m.state = constants.ListState
	case tea.KeyMsg:
		if key.Matches(msg, constants.Keymap.Help) {
			m.showHelp = !m.showHelp
		}
		if key.Matches(msg, constants.Keymap.Download) {
			if m.state == constants.ListState {
				m.list.StartSpinner()
				commands = append(commands, cmds.DownloadCmd(&m.jsonFile, m.list.SelectedItem().(services.Item)))
			}
		}
		if key.Matches(msg, constants.Keymap.Move) {
			if m.state == constants.ListState && m.playlist.IsLoaded() {
				m.state = constants.PlaylistState
			} else {
				m.state = constants.ListState
			}
			commands = append(commands, ChangeState((*uint)(&m.state)))
		}
		if key.Matches(msg, constants.Keymap.Quit) {
			if m.state != constants.SearchState {
				cmds.Kill()
				return m, tea.Quit
			}
		}
		if key.Matches(msg, constants.Keymap.Search) {
			m.state = constants.SearchState
		}
		if key.Matches(msg, constants.Keymap.Remove) {
			if m.isConnected {
				if m.state == constants.PlaylistState {
					commands = append(commands, cmds.RemoveCmd(m.conn, m.playlist.GetPosition()))
				}
				if m.state == constants.ListState {
					commands = append(commands, cmds.DeleteCmd(&m.jsonFile, m.list.SelectedItem().(services.Item)))
				}

			}
		}
		if key.Matches(msg, constants.Keymap.Play) {
			if m.isConnected {
				if m.state == constants.PlaylistState {
					commands = append(commands, cmds.PlayByPos(m.conn, m.playlist.GetPosition()))
				}
				if m.state == constants.ListState {
					commands = append(commands, cmds.PlayCmd(m.conn, m.list.SelectedItem().(services.Item)))
				}
				if m.state == constants.SearchState {
					if m.search.GetLength() > 0 {
						current := m.search.GetCurrent()
						commands = append(commands, cmds.SaveCmd(&m.jsonFile, current))
						m.state = constants.ListState
					} else {
						if m.search.GetChoice() != "" {
							commands = append(commands, cmds.SearchCmd(m.search.GetAnswer(), m.search.GetChoice()))
						}
					}
				}
			}
		}
		if key.Matches(msg, constants.Keymap.Pause) {
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

	if m.state == constants.SearchState {
		updatedSearch, searchCmd := m.search.Update(msg)
		m.search = updatedSearch.(search.Model)
		commands = append(commands, searchCmd)
	} else {
		updatedPlayer, cmd := m.player.Update(msg)
		m.player = updatedPlayer.(player.Model)

		if m.state == constants.ListState {
			updatedList, listCmd := m.list.Update(msg)
			m.list = updatedList.(fancyList.Model)
			commands = append(commands, listCmd)
		}
		updatedPlaylist, playlistCmd := m.playlist.Update(msg)
		m.playlist = updatedPlaylist.(playlist.Model)

		if m.state == constants.PlaylistState {
			commands = append(commands, playlistCmd)
		}

		commands = append(commands, cmd)
	}

	m.help, cmd = m.help.Update(msg)

	commands = append(commands, cmd)
	return m, tea.Batch(commands...)
}

func (m Model) View() string {
	list := m.list.View()
	playlist := m.playlist.View()

	if !m.isConnected {
		return "Loading..."
	}
	if m.conn.IsClosed() {
		return "Loading..."
	}

	if m.state == constants.SearchState {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, m.search.View(), m.help.FullHelpView(constants.SearchKeymap.FullHelp())))
	}

	if m.state == constants.ListState {
		list = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Render(list)
		playlist = lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Render(playlist)
	} else {
		playlist = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Render(playlist)
		list = lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Render(list)
	}

	var help string

	if m.state == constants.ListState {
		help = m.help.FullHelpView(constants.Keymap.FullHelp())
	}
	if m.state == constants.SearchState {
		help = m.help.FullHelpView(constants.SearchKeymap.FullHelp())
	}
	if m.state == constants.PlaylistState {
		help = m.help.FullHelpView(constants.PlaylistKeymap.FullHelp())
	}

	main := lipgloss.JoinHorizontal(lipgloss.Left, list, playlist)

	page := lipgloss.JoinVertical(lipgloss.Left, main, m.player.View())
	pageWithHelp := lipgloss.JoinVertical(lipgloss.Left, main, m.player.View(), help)

	if m.showHelp {
		return lipgloss.Place(m.width, m.height, 0, 0, pageWithHelp)
	}

	return lipgloss.Place(m.width, m.height, 0, 0, page)
}

func ChangeState(state *uint) tea.Cmd {
	return func() tea.Msg {
		return cmds.StateMsg(state)
	}
}

func (m *Model) SetState(newState uint) {
	m.state = constants.SessionState(newState)
}

func (m *Model) updateData() {
	jsonfile, items, _ := cmds.SearchSongs()

	m.jsonFile = jsonfile
	m.list.SetItems(items)
}
