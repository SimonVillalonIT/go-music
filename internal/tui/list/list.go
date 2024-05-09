package custom_list

import (
	"math"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	list   list.Model
	width  int
	height int
}

func NewModel(items []list.Item) Model {

	delegate := newItemDelegate()
	musicList := list.New(items, delegate, 0, 0)
	musicList.Title = "Songs & Playlists"
	musicList.Styles.Title = titleStyle
	musicList.SetShowHelp(false)
	musicList.SetSpinner(spinner.Dot)
	return Model{
		list: musicList,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = int(math.Round(float64(msg.Width) * 0.65))
		m.height = int(math.Round(float64(msg.Height) * 0.8))
		m.SetSize(m.width, m.height)
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, nil
		}
		if m.list.FilterState() == list.Filtering {
			break
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.NewStyle().Width(m.width).Height(m.height).Render((m.list.View()))
}

func (m Model) SelectedItem() list.Item {
	return m.list.SelectedItem()
}

func (m *Model) SetSize(w, h int) {
	m.list.SetSize(w, h)
}

func (m *Model) SetItems(items []list.Item) {
	m.list.SetItems(items)
}

func (m *Model) StartSpinner() {
	m.list.StartSpinner()
}

func (m *Model) StopSpinner() {
	m.list.StopSpinner()
}
