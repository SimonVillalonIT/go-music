package search

import (
	"fmt"
	"io"
	"strings"

	"github.com/SimonVillalonIT/music-golang/internal/services"
	cmds "github.com/SimonVillalonIT/music-golang/internal/tui/commands"
	"github.com/SimonVillalonIT/music-golang/internal/tui/constants"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

var (
	NormalTitle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
			Padding(0, 0, 0, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
)

type Input interface {
	Value() string
	Default() string
	Blur() tea.Msg
	Update(tea.Msg) (Input, tea.Cmd)
	View() string
}

type InputField struct {
	textinput textinput.Model
	def       string
}

func NewInputField(def string) *InputField {
	ti := textinput.New()
	ti.Placeholder = "Write here... Default(" + def + ")"
	ti.Focus()
	return &InputField{textinput: ti, def: def}
}

func (i *InputField) Value() string {
	return i.textinput.Value()
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
	question Question
	list     list.Model
	width    int
	height   int
	styles   *Styles
	results  []services.Item
	error    string
}

type Question struct {
	Request string
	Answer  string
	Input   Input
}

func NewQuestion(request, def string) Question {
	return Question{Request: request, Input: NewInputField(def)}
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 1 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(services.Item)
	if !ok {
		return
	}
	textwidth := uint(m.Width() - NormalTitle.GetPaddingLeft() - NormalTitle.GetPaddingRight())

	str := fmt.Sprintf("%s", i.Title())

	fn := itemStyle.Render

	if index == m.Index() {
		fn = func(s ...string) string {
			toRender := ("> " + strings.Join(s, ""))
			toRender = truncate.StringWithTail(toRender, textwidth, "...")
			return itemStyle.Render(toRender)
		}
	} else {
		fn = func(s ...string) string {
			toRender := (strings.Join(s, ""))
			toRender = truncate.StringWithTail(toRender, textwidth, "...")
			return itemStyle.Render(toRender)
		}
	}
	fmt.Fprint(w, fn(str))
}
func NewModel() Model {
	styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Focus()
	delegate := itemDelegate{}
	request := NewQuestion("Search for a song: ", "")
	list := list.New(make([]list.Item, 0), delegate, 80, 20)

	return Model{question: request, styles: styles, list: list}
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
        if key.Matches(msg, constants.Keymap.Quit){
            commands = append(commands, cmds.QuitSearchCmd)
        }
	}
	m.question.Input, cmd = m.question.Input.Update(msg)

	commands = append(commands, cmd)

	updatedList, cmd := m.list.Update(msg)

	m.list = updatedList

	commands = append(commands, cmd)

	return m, tea.Batch(commands...)
}

func (m Model) View() string {
	renderedError := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render(m.error)
	renderedInput := m.styles.InputField.Render(m.question.Input.View())

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

func (m Model) Answer() string {
	return m.question.Input.Value()
}

func getItems(items []services.Item) []list.Item {
	var result []list.Item

	for _, item := range items {
		result = append(result, list.Item(item))
	}

	return result
}

func (m *Model) GetLength() int {
	return len(m.list.Items())
}

func (m *Model) GetCurrent() services.Item {
	return m.list.SelectedItem().(services.Item)
}
