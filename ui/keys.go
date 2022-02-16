package ui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up            key.Binding
	Down          key.Binding
	TogglePreview key.Binding
	OpenBrowser   key.Binding
	Help          key.Binding
	Quit          key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.TogglePreview, k.OpenBrowser},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	TogglePreview: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "open in Preview"),
	),
	OpenBrowser: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "open in Browser"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
