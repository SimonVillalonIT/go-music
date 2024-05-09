package playlist

import (
	"fmt"
	"io"
	"log"
	"math"
	"path/filepath"
	"strings"

	"github.com/SimonVillalonIT/go-music/internal/services"
	cmds "github.com/SimonVillalonIT/go-music/internal/tui/commands"
	"github.com/SimonVillalonIT/go-music/internal/tui/constants"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
	"github.com/spf13/viper"
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

type item struct {
	ID           float64 `json:"id"`
	Filename     string  `json:"filename"`
	Playing      bool    `json:"playing"`
	Current      bool    `json:"current"`
	PlaylistPath string  `json:"playlist-path"`
	Position     int
}

func (i item) FilterValue() string { return i.Filename }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 1 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	textwidth := uint(m.Width() - NormalTitle.GetPaddingLeft() - NormalTitle.GetPaddingRight())

	str := fmt.Sprintf("%s", i.Filename)

	fn := itemStyle.Render

	if i.Current && index == m.Index() {
		fn = func(s ...string) string {
			toRender := ("> " + strings.Join(s, ""))
			toRender = truncate.StringWithTail(toRender, textwidth, "...")
			return itemStyle.Render(selectedItemStyle.Render(toRender))
		}
	} else if i.Current {
		fn = func(s ...string) string {
			toRender := (strings.Join(s, ""))
			toRender = truncate.StringWithTail(toRender, textwidth, "...")
			return itemStyle.Render(selectedItemStyle.Render(toRender))
		}
	} else if index == m.Index() {
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

type Model struct {
	width  int
	height int
	list   list.Model
	choice string
	state  *uint
	loaded bool
}

func NewModel() Model {
	l := list.New(make([]list.Item, 0), itemDelegate{}, 0, 0)
	l.Title = "Current playlist"
	l.Styles.Title = titleStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	return Model{list: l, loaded: false}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var commands []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = int(math.Round(float64(msg.Width) * 0.3))
		m.height = int(math.Round(float64(msg.Height) * 0.8))
		m.list.SetSize(m.width, int(math.Round(float64(msg.Height)*0.8)))
	case cmds.StateMsg:
		m.state = msg
	case cmds.PlaylistInfoMsg:
		items, err := getItems(msg)
		if err == nil {
			m.list.SetItems(items)
			m.loaded = true
		} else {
			log.Print(err)
		}
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, nil
		}
	}
	if !m.loaded {
		return m, nil
	}
	if *m.state == constants.PlaylistState {
		m.list, cmd = m.list.Update(msg)
		commands = append(commands, cmd)
	}

	return m, tea.Batch(commands...)
}

func (m Model) View() string {

	if len(m.list.Items()) < 1 {
		return lipgloss.NewStyle().Width(m.width).Height(m.height).Render("")
	}

	return lipgloss.NewStyle().Width(m.width).Height(m.height).Render(m.list.View())
}

func getItems(data []interface{}) ([]list.Item, error) {
	var items []list.Item

	for pos, m := range data {
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

		i.Filename = filepath.Base(viper.GetString(services.DOWNLOADS_FOLDER) + "/" + i.Filename)

		i.Filename = strings.TrimSuffix(i.Filename, ".mp3")

		if i.ID, ok = mMap["id"].(float64); !ok {
			return nil, fmt.Errorf("id field missing or not an integer")
		}

		if i.Playing, ok = mMap["playing"].(bool); !ok {
			i.Current = false
		}

		if i.PlaylistPath, ok = mMap["playlist-path"].(string); !ok {
			return nil, fmt.Errorf("playlist-path field missing or not a string")
		}

		i.Position = pos

		items = append(items, i)
	}

	return items, nil
}

func (m *Model) GetPosition() int {
	return m.list.SelectedItem().(item).Position
}

func (m *Model) IsLoaded() bool {
	return m.loaded && len(m.list.Items()) > 0
}
