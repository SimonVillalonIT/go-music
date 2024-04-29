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

type keymap struct {
	Download key.Binding
	Enter    key.Binding
	Space    key.Binding
    Search  key.Binding
	Rename   key.Binding
	Remove   key.Binding
	Back     key.Binding
	Quit     key.Binding
	Next     key.Binding
	Prev     key.Binding
	Minus    key.Binding
	Plus     key.Binding
	Move     key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Move: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "move"),
	),
	Download: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "download"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "exec"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	Remove: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "remove"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "select"),
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
}
