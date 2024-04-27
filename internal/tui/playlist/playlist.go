package playlist

import (
	"fmt"
	"io"
	"log"
	"math"
	"strings"

	cmds "github.com/SimonVillalonIT/music-golang/internal/tui/commands"
	"github.com/SimonVillalonIT/music-golang/internal/tui/constants"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	ID           float64 `json:"id"`
	Filename     string  `json:"filename"`
	Playing      bool    `json:"playing"`
	Current      bool    `json:"current"`
	PlaylistPath string  `json:"playlist-path"`
}

func (i item) FilterValue() string { return i.Filename }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.Filename)

	fn := itemStyle.Render

	if i.Current && index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, ""))
		}
	} else if i.Current {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(strings.Join(s, ""))
		}
	} else if index == m.Index() {
		fn = func(s ...string) string {
			return "> " + strings.Join(s, "")
		}
	} else {
		fn = func(s ...string) string {
			return strings.Join(s, "")
		}
	}

	fmt.Fprint(w, fn(str))
}

type Model struct {
	list   list.Model
	choice string
	state  *uint
	Loaded bool
}

func NewModel() Model {
	l := list.New(make([]list.Item, 0), itemDelegate{}, 0, listHeight)
	l.Title = "Current playlist"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return Model{list: l, Loaded: false}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var commands []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(int(math.Round(float64(msg.Width) * 0.3)))
	case cmds.StateMsg:
		m.state = msg
		m.Loaded = true
	case cmds.PlaylistInfoMsg:
		items, err := getItems(msg)
		if err == nil {
			m.list.SetItems(items)
		} else {
			log.Print(err)
		}
	}
	if !m.Loaded {
		return m, nil
	}
	if *m.state == constants.PlaylistState {
		m.list, cmd = m.list.Update(msg)
		commands = append(commands, cmd)
	}

	return m, tea.Batch(commands...)
}

func (m Model) View() string {
	return m.list.View()
}

func getItems(data []interface{}) ([]list.Item, error) {
	var items []list.Item

	for _, m := range data {
		mMap, ok := m.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected type in array, expected map[string]interface{}")
		}

		var i item

		if i.Current, ok = mMap["current"].(bool); !ok {
			i.Current = false
		}

		if i.Filename, ok = mMap["filename"].(string); !ok {
			return nil, fmt.Errorf("filename field missing or not a string")
		}
		if i.ID, ok = mMap["id"].(float64); !ok {
			return nil, fmt.Errorf("id field missing or not an integer")
		}

		if i.Playing, ok = mMap["playing"].(bool); !ok {
			i.Current = false
		}

		if i.PlaylistPath, ok = mMap["playlist-path"].(string); !ok {
			return nil, fmt.Errorf("playlist-path field missing or not a string")
		}

		items = append(items, i)
	}

	return items, nil
}

func (m *Model) GetIndex() float64 {
	return m.list.SelectedItem().(item).ID - 1
}
