package tui

import (
	config "github.com/SimonVillalonIT/go-music/internal/tui/config"
	help "github.com/SimonVillalonIT/go-music/internal/tui/config"
	home "github.com/SimonVillalonIT/go-music/internal/tui/home"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	HomeView = iota
	ConfigView
	SearchView
)

type MainModel struct {
	state  sessionState
	home   tea.Model
	config tea.Model
	help   tea.Model
	p      *tea.Program
}

func NewMainModel(state sessionState) *MainModel {
	return &MainModel{help: help.New(), home: home.New(), config: config.New(), state: state}
}

func (m MainModel) Init() tea.Cmd {
	switch m.state {
	case ConfigView:
		return m.config.Init()
	default:
		return m.home.Init()
	}
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m.state {
	case ConfigView:
		newConfig, newCmd := m.config.Update(msg)
		configModel, ok := newConfig.(config.Model)
		if !ok {
			panic("could not perform assertion on config model")
		}
		m.config = configModel
		cmd = newCmd
	default:
		newHome, newCmd := m.home.Update(msg)
		homeModel, ok := newHome.(home.Model)
		if !ok {
			panic("could not perform assertion on home model")
		}
		m.home = homeModel
		cmd = newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	switch m.state {
	case ConfigView:
		return m.config.View()
	default:
		return m.home.View()
	}
}
