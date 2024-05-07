package search

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var choices = []string{"Playlist", "Song"}

type choiceModel struct {
	cursor int
	choice string
}

func (m choiceModel) Init() tea.Cmd {
	return nil
}

func (m choiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.choice = choices[m.cursor]

		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	}

	return m, nil
}

func (m choiceModel) View() string {
	s := strings.Builder{}
	s.WriteString("What do you want to search for?\n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}

	return s.String()
}

func (m *choiceModel) ClearChoice() {
	m.choice = ""
}
