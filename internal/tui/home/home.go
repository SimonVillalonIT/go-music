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
	"github.com/SimonVillalonIT/music-golang/internal/tui/download"
	fancyList "github.com/SimonVillalonIT/music-golang/internal/tui/list"
	"github.com/SimonVillalonIT/music-golang/internal/tui/player"
	"github.com/SimonVillalonIT/music-golang/internal/tui/playlist"
)

type Model struct {
	state       constants.SessionState
	width       int
	height      int
	list        fancyList.Model
	player      player.Model
	playlist    playlist.Model
	search      search.Model
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
	player := player.NewModel()
	playlist := playlist.NewModel()
	search := search.New()

	return Model{list: list, jsonFile: jsonfile, isConnected: false, player: player, playlist: playlist, search: search}
}

func (m Model) Init() tea.Cmd {
	return changeState((*uint)(&m.state))
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
	case tea.KeyMsg:
		if key.Matches(msg, constants.Keymap.Move) {
			if m.state == constants.ListState && m.playlist.IsLoaded() {
				m.state = constants.PlaylistState
			} else {
				m.state = constants.ListState
			}
			commands = append(commands, changeState((*uint)(&m.state)))
		}
		if key.Matches(msg, constants.Keymap.Quit) {
			if m.state != constants.SearchState {
				cmds.Kill()
				return m, tea.Quit
			} else {
				m.state = constants.ListState
			}
		}
		if key.Matches(msg, constants.Keymap.Search) {
			m.state = constants.SearchState
		}
		if key.Matches(msg, constants.Keymap.Remove) && m.state == constants.PlaylistState {
			if m.isConnected {
				commands = append(commands, cmds.RemoveCmd(m.conn, m.playlist.GetPosition()))
			}
		}
		if key.Matches(msg, constants.Keymap.Enter) {
			if m.isConnected {
				if m.state == constants.ListState {
					commands = append(commands, cmds.PlayCmd(m.conn, m.list.SelectedItem().(services.Item)))
				}
				if m.state == constants.SearchState {
					if m.search.GetLength() > 0 {
						current := m.search.GetCurrent()
                        commands = append(commands, cmds.SaveCmd(&m.jsonFile, current))
					} else {
						commands = append(commands, cmds.SearchCmd(m.search.Answer()))
					}
				}
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
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.search.View())
	}

	if m.state == constants.ListState {
		list = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Render(list)
		playlist = lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Render(playlist)
	} else {
		playlist = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Render(playlist)
		list = lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Render(list)
	}

	main := lipgloss.JoinHorizontal(lipgloss.Left, list, playlist)
	return lipgloss.Place(m.width, m.height, 0, 0, lipgloss.JoinVertical(lipgloss.Left, main, m.player.View()))
}

func changeState(state *uint) tea.Cmd {
	return func() tea.Msg {
		return cmds.StateMsg(state)
	}
}
