package search

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Question struct {
	Request string
	Answer  string
	Input   Input
}

func NewQuestion(request, def string) Question {
	return Question{Request: request, Input: NewInputField()}
}

type Input interface {
	Value() string
	Default() string
	Blur() tea.Msg
	Update(tea.Msg) (Input, tea.Cmd)
	View() string
	SetValue(string)
}

type InputField struct {
	textinput textinput.Model
	def       string
}

func NewInputField() *InputField {
	ti := textinput.New()
	ti.Placeholder = "Write here..."
	ti.Focus()
	return &InputField{textinput: ti}
}

func (i *InputField) Value() string {
	return i.textinput.Value()
}

func (i *InputField) SetValue(str string) {
	i.textinput.SetValue(i.textinput.Value() + str)
    i.textinput.CursorEnd()
}

func (i *InputField) Default() string {
	return i.def
}

func (i *InputField) Blur() tea.Msg {
	return i.textinput.Blur
}

func (i *InputField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	i.textinput, cmd = i.textinput.Update(msg)
	return i, cmd
}

func (i *InputField) View() string {
	return i.textinput.View()
}

func (m Model) GetAnswer() string {
	return m.question.Input.Value()
}

func (m *Model) ClearAnswer() {
	m.question.Input = NewInputField()
}
