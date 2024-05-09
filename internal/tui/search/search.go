package search

import (
	"github.com/SimonVillalonIT/go-music/internal/services"
	cmds "github.com/SimonVillalonIT/go-music/internal/tui/commands"
	"github.com/SimonVillalonIT/go-music/internal/tui/constants"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	question Question
	list     list.Model
	choice   choiceModel
	width    int
	height   int
	styles   *Styles
	results  []services.Item
	error    string
}

func NewModel() Model {
	styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Focus()
	delegate := itemDelegate{}
	request := NewQuestion("Search here: ", "")
	list := list.New(make([]list.Item, 0), delegate, 80, 20)
	choice := choiceModel{}

	return Model{question: request, styles: styles, list: list, choice: choice}
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
		m.list.SetSize(msg.Width, msg.Height)
	case cmds.SearchResultMsg:
		m.results = msg
		m.list.SetItems(getItems(msg))
	case tea.KeyMsg:
		if key.Matches(msg, constants.Keymap.Quit) {
			commands = append(commands, cmds.QuitSearchCmd)
		}
		if msg.String() == "q" {
            m.question.Input.SetValue("q")
			return m, nil
		}
	}
	updatedChoice, cmd := m.choice.Update(msg)
	m.choice = updatedChoice.(choiceModel)

	commands = append(commands, cmd)

	if m.choice.choice != "" {
		m.question.Input, cmd = m.question.Input.Update(msg)

		commands = append(commands, cmd)

		updatedList, cmd := m.list.Update(msg)

		m.list = updatedList

		commands = append(commands, cmd)

	}

	return m, tea.Batch(commands...)
}

func (m Model) View() string {
	renderedError := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render(m.error)
	renderedInput := m.styles.InputField.Render(m.question.Input.View())

	if m.choice.choice == "" {
		return m.choice.View()
	}

	if len(m.results) > 0 {
		return m.list.View()
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(
		lipgloss.Center,
		m.question.Request,
		renderedError,
		renderedInput,
	))
}

func (m Model) GetChoice() string {
	return m.choice.choice
}

func (m *Model) ClearData() {
	m.results = []services.Item{}
	m.list.SetItems(make([]list.Item, 0))
	m.question.Answer = ""
	m.question.Request = ""
	m.ClearAnswer()
	m.choice.ClearChoice()
}
