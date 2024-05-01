package constants

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* CONSTANTS */
type SessionState uint

const (
	ListState = iota
	PlaylistState
	SearchState
)

var (
	// P the current tea program
	P *tea.Program
	// WindowSize store the size of the terminal window

	WindowSize tea.WindowSizeMsg
)

/* STYLING */

// DocStyle styling for viewports
var DocStyle = lipgloss.NewStyle().Margin(0, 2)

// HelpStyle styling for help context menu
var HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

// ErrStyle provides styling for error messages
var ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#bd534b")).Render

// AlertStyle provides styling for alert messages
var AlertStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Render

type KeymapType struct {
	Play     key.Binding
	Pause    key.Binding
	Move     key.Binding
	Search   key.Binding
	Download key.Binding
	Remove   key.Binding
	Next     key.Binding
	Prev     key.Binding
	Minus    key.Binding
	Plus     key.Binding
	Quit     key.Binding
	Help     key.Binding
}

var Keymap = KeymapType{
	Move: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "move"),
	),
	Download: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "download"),
	),
	Play: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "play"),
	),
	Remove: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "remove"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q", "esc"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
	Pause: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "pause"),
	),
	Prev: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "prev"),
	),
	Next: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next"),
	),
	Minus: key.NewBinding(
		key.WithKeys("-"),
		key.WithHelp("-", "decrease"),
	),
	Plus: key.NewBinding(
		key.WithKeys("+"),
		key.WithHelp("+", "increase"),
	),
	Search: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "search"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

func (k KeymapType) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

func (k KeymapType) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Play, k.Pause},
		{k.Next, k.Prev},
		{k.Download, k.Search}, // second column
		{k.Plus, k.Minus},
		{k.Move, k.Remove},
		{k.Quit, k.Help},
	}
}
