package search

import (
	"fmt"
	"io"
	"strings"

	"github.com/SimonVillalonIT/go-music/internal/services"
	"github.com/charmbracelet/bubbles/list"
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

	NormalDesc = NormalTitle.Copy().
			Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	FilterMatch = lipgloss.NewStyle().Underline(true)

	SelectedTitle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
			Padding(0, 0, 0, 1)

	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
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

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 1 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {

	var (
		title, desc string
	)

	if i, ok := listItem.(services.Item); ok {
		title = i.Title()
		desc = i.Description()
	} else {
		return
	}

	if m.Width() <= 0 {
		return
	}

	// Prevent text from exceeding list width
	textwidth := uint(m.Width() - NormalTitle.GetPaddingLeft() - NormalTitle.GetPaddingRight())
	title = truncate.StringWithTail(title, textwidth, "...")
	var lines []string
	for i, line := range strings.Split(desc, "\n") {
		if i >= d.Height()-1 {
			break
		}
		lines = append(lines, truncate.StringWithTail(line, textwidth, "..."))
	}
	desc = strings.Join(lines, "\n")

	if index == m.Index() {
		title = selectedItemStyle.Render(title)
		desc = NormalDesc.Render(desc)
		fmt.Fprintf(w, "%s\n%s", title, desc)

		return
	} else {
		title = NormalTitle.Render(title)
		desc = NormalDesc.Render(desc)
		fmt.Fprintf(w, "%s\n%s", title, desc)

		return
	}
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
