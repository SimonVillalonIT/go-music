package tui

import (
	"fmt"
	"os"
	"path"

	"github.com/SimonVillalonIT/music-golang/internal/services"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

type Model struct {
	questions []Question
	index     int
	width     int
	height    int
	styles    *Styles
	done      bool
	error     string
}

type Question struct {
	Question string
	Answer   string
	Input    Input
}

func NewQuestion(question, def string) Question {
	return Question{Question: question, Input: NewInputField(def)}
}

func New() tea.Model {
	styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Focus()
	home, _ := os.UserHomeDir()
	questions := []Question{
		NewQuestion("Enter the path where you want to save your songs", path.Join(home, "Music")),
		NewQuestion("Enter the path where you want to save the store", home)}

	return &Model{questions: questions, styles: styles}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &m.questions[m.index]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.index == len(m.questions)-1 {
				m.done = true
			}
			if current.Input.Value() == "" {
				current.Answer = current.Input.Default()
			} else {
				current.Answer = current.Input.Value()
			}
			if services.CheckIfPathExists(current.Answer) {
				m.Next()
				m.error = ""
				return m, current.Input.Blur
			}
			m.error = "Invalid path:" + current.Answer
			return m, current.Input.Blur
		}
	}
	current.Input, cmd = current.Input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if len(m.questions) == 0 {
		return "No questions available"
	}
	current := m.questions[m.index]

	if m.done {
		if err := services.GenerateConfig(m.questions[0].Answer, m.questions[1].Answer); err != nil {
			m.error = fmt.Sprintf("Error: %s", err.Error())
			return m.error
		}
		return "Config saved successfully"
	}

	if m.width == 0 {
		return "loading..."
	}

	renderedError := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render(m.error)
	renderedInput := m.styles.InputField.Render(current.Input.View())

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(
		lipgloss.Center,
		current.Question,
		renderedError,
		renderedInput,
	))
}

func (m *Model) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
	} else {
		m.index = 0
	}
}
