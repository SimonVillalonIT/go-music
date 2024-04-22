package tui

import (
	"log"

	"github.com/DexterLB/mpvipc"
	"github.com/SimonVillalonIT/music-golang/internal/services"
	cmds "github.com/SimonVillalonIT/music-golang/internal/tui/commands"
	"github.com/SimonVillalonIT/music-golang/internal/tui/constants"
	custom_list "github.com/SimonVillalonIT/music-golang/internal/tui/list"
	fancyList "github.com/SimonVillalonIT/music-golang/internal/tui/list"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

//
// import (
// 	"bytes"
// 	"encoding/json"
// 	"log"
// 	"os"
// 	"os/exec"
// 	"slices"
// 	"strings"
//
// 	"github.com/DexterLB/mpvipc"
// 	"github.com/SimonVillalonIT/music-golang/internal/services"
// 	"github.com/SimonVillalonIT/music-golang/internal/tui/constants"
// 	"github.com/charmbracelet/bubbles/key"
// 	"github.com/charmbracelet/bubbles/list"
// 	"github.com/charmbracelet/lipgloss"
// 	"github.com/spf13/viper"
// )
//
// const (
// 	normal   = 0
// 	play     = 1
// 	download = 2
// 	search   = 3
// )
//
// var docStyle = lipgloss.NewStyle().Margin(1, 2)
//
// type Model struct {
// 	err       error
// 	help      tea.Model
// 	list      list.Model
// 	mode      uint
// 	jsonFile  []services.Item
// 	selected  []string
// 	conn      mpvipc.Connection
// 	connected bool
// }
//
// func New() tea.Model {
// 	jsonfile, items, err := searchsongs()
// 	if err != nil {
// 		log.print(err)
// 	}
//
// 	listModel := list.New(items, list.NewDefaultDelegate(), 0, 0)
// 	listModel.Title = "Your saved songs & playlists"
// 	listModel.AdditionalShortHelpKeys = func() []key.Binding {
// 		return []key.Binding{
// 			constants.Keymap.Enter,
// 			constants.Keymap.Download,
// 			constants.Keymap.Delete,
// 			constants.Keymap.Space,
// 			constants.Keymap.Back,
// 		}
// 	}
// 	return Model{jsonFile: jsonFile, list: listModel, mode: 0}
// }
//
// type mpvMsg bool
//
//
// type playMsg mpvipc.Connection
//
// type errMsg error

// type connMsg mpvipc.Connection
// type connErr error
//
//
//
// func startMpv() tea.Msg {
// 	if !isMpvRunning() {
// 		err := exec.Command("mpv", "--no-video", "--idle", "--really-quiet", "--input-ipc-server=/tmp/mpvsocket").Run()
//
// 		if err != nil {
// 			log.Printf("Failed to start mpv: %v", err)
//
// 			return errMsg(err)
// 		}
// 	}
//
// 	return mpvMsg(true)
// }
//
// func connCmd() tea.Msg {
// 	conn := mpvipc.NewConnection("/tmp/mpvsocket")
//
// 	err := conn.Open()
//
// 	if conn.IsClosed() {
// 		return connErr(err)
// 	}
//
// 	return connMsg(*conn)
// }
//
// func playCmd(conn mpvipc.Connection, items ...services.Item) tea.Cmd {
// 	return func() tea.Msg {
// 		err := services.Play(conn, items...)
//
// 		if err != nil {
// 			return errMsg(err)
// 		}
// 		return playMsg(conn)
// 	}
// }
//
// func kill() {
// 	exec.Command("killall", "mpv").Run()
// }
//
// func (m Model) Init() tea.Cmd {
// 	return startMpv
// }
//
// func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var cmd tea.Cmd
// 	var cmds []tea.Cmd
// 	switch msg := msg.(type) {
// 	case connErr:
// 		cmds = append(cmds, connCmd)
// 	case connMsg:
// 		m.conn = mpvipc.Connection(msg)
// 	case mpvMsg:
// 		if msg == true {
// 			cmds = append(cmds, connCmd)
// 		}
// 	case errMsg:
// 		log.Print(msg)
// 		m.err = msg
//
// 	case tea.WindowSizeMsg:
// 		h, v := docStyle.GetFrameSize()
// 		m.list.SetSize(msg.Width-h, msg.Height-v)
//
// 	case tea.KeyMsg:
// 		if key.Matches(msg, constants.Keymap.Quit) {
// 			kill()
// 			cmds = append(cmds, tea.Quit)
// 		}
// 		if key.Matches(msg, constants.Keymap.Enter) {
// 			if m.conn.IsClosed() {
// 				cmds = append(cmds, connCmd)
// 			} else {
// 				if m.mode == normal {
// 					cmds = append(cmds, playCmd(m.conn, m.getSelectedItems()...))
// 				}
// 				if m.mode == download {
// 				}
// 			}
// 		}
// 		if key.Matches(msg, constants.Keymap.Download) {
// 			if m.mode == download {
// 				m.mode = normal
// 			}
// 			if m.mode == normal {
// 				m.mode = download
// 			}
// 		}
// 		if key.Matches(msg, constants.Keymap.Space) {
// 			m.selected = m.selectItem()
// 		}
// 	}
//
// 	m.list, cmd = m.list.Update(msg)
// 	cmds = append(cmds, cmd)
// 	return m, tea.Batch(cmds...)
// }
//
// func (m Model) View() string {
// 	if m.err != nil {
// 		return m.err.Error()
// 	}
//
// 	mainList := docStyle.Render(m.list.View())
//
// 	selectedList := m.renderSelectedList()
//
// 	return lipgloss.JoinHorizontal(lipgloss.Top, mainList, selectedList)
// }
//
// func searchSongs() ([]services.Item, []list.Item, error) {
// 	rawFile, err := os.ReadFile(viper.GetString(services.STORE_PATH))
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	var jsonFile []services.Item
//
// 	err = json.Unmarshal(rawFile, &jsonFile)
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	items := make([]list.Item, len(jsonFile))
//
// 	for i, song := range jsonFile {
// 		items[i] = song
// 	}
//
// 	return jsonFile, items, nil
// }
//
// func (m Model) renderSelectedList() string {
// 	style := list.DefaultStyles()
// 	var selectedItems []string
//
// 	for _, item := range m.getSelectedItems() {
// 		selectedItems = append(selectedItems, item.Name)
// 	}
//
// 	selectedList := strings.Join(selectedItems, "\n")
//
// 	return lipgloss.NewStyle().Render(lipgloss.JoinVertical(lipgloss.Top, style.Title.Render("Selected Items"), selectedList))
// }
//
// func (m Model) selectItem() []string {
// 	if !slices.Contains(m.selected, m.list.SelectedItem().FilterValue()) {
// 		result := append(m.selected, m.list.SelectedItem().FilterValue())
//
// 		return result
// 	}
// 	result := slices.DeleteFunc(m.selected, func(e string) bool {
// 		return e == m.list.SelectedItem().FilterValue()
// 	})
// 	return result
// }
//
// func (m Model) getSelectedItems() []services.Item {
// 	var selectedItems []services.Item
// 	selectedMap := make(map[string]services.Item)
//
// 	for _, item := range m.jsonFile {
// 		if id := item.FilterValue(); id != "" { // Check for empty ID
// 			selectedMap[id] = item
// 		}
// 	}
//
// 	for _, id := range m.selected {
// 		if item, ok := selectedMap[id]; ok {
// 			selectedItems = append(selectedItems, item)
// 		}
// 	}
//
// 	return selectedItems
// }

type Model struct {
	list     fancyList.Model
	err      error
	help     tea.Model
	mode     uint
	jsonFile []services.Item
	conn     mpvipc.Connection
}

func New() tea.Model {
	jsonfile, items, err := cmds.SearchSongs()
	if err != nil {
		log.Print(err)
	}

	list := fancyList.NewModel(items)
	return Model{list: list, jsonFile: jsonfile}
}

func (m Model) Init() tea.Cmd {
	return cmds.StartMpv
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case cmds.ConnErr:
		commands = append(commands, cmds.ConnCmd)
	case cmds.ConnMsg:
		m.conn = mpvipc.Connection(msg)
	case cmds.MpvMsg:
		if msg == true {
			commands = append(commands, cmds.ConnCmd)
		}
	case cmds.ErrMsg:
		log.Print(msg)
		m.err = msg

	case tea.KeyMsg:
		if key.Matches(msg, constants.Keymap.Enter) {
			commands = append(commands, cmds.PlayCmd(m.conn, m.list.SelectedItem().(services.Item)))
		}

	}

	updatedList, cmd := m.list.Update(msg)
	m.list = updatedList.(custom_list.Model)
	commands = append(commands, cmd)

	return m, tea.Batch(commands...)
}
func (m Model) View() string {
	return m.list.View()
}
