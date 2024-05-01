package tui

import (
	config "github.com/SimonVillalonIT/music-golang/internal/tui/config"
	help "github.com/SimonVillalonIT/music-golang/internal/tui/config"
	home "github.com/SimonVillalonIT/music-golang/internal/tui/home"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	HomeView   sessionState = iota
	ConfigView sessionState = 2
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
	case HomeView:
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
	case HomeView:
		return m.home.View()
	}
	return "No view specified"
}
