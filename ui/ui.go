package ui

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	ready           bool
	Owner           string
	Name            string
	width           int
	contentViewport contentViewport
	previewViewport previewViewport
	sections        []Section
	data            []Discussion
	cursor          cursor
	help            help.Model
	keyMap          keyMap
}

type contentViewport struct {
	width  int
	height int
}

type previewViewport struct {
	width  int
	height int
}

type cursor struct {
	currDiscID    int
	currSectionID int
}

type initMsg struct {
	data     []Discussion
	sections []Section
}

func NewModel(owner, name string) Model {
	// helperModel := help.NewModel()
	// style := lipgloss.NewStyle().Foreground(secondaryText)
	// helperModel.Styles = help.Styles{
	// 	ShortDesc:      style.Copy(),
	// 	FullDesc:       style.Copy(),
	// 	ShortSeparator: style.Copy(),
	// 	FullSeparator:  style.Copy(),
	// 	ShortKey:       style.Copy(),
	// 	FullKey:        style.Copy(),
	// 	Ellipsis:       style.Copy(),
	// }

	return Model{
		Owner:  owner,
		Name:   name,
		keyMap: DefaultKeyMap(),
		cursor: cursor{
			currSectionID: 0,
			currDiscID:    0,
		},
	}
}

func getTitleWidth(viewportWidth int) int {
	return viewportWidth - usedWidth
}

func (m Model) contentWidth() int {
	return m.contentViewport.width
}

func (m Model) contentHeight() int {
	return m.contentViewport.height
}

func (m Model) previewWidth() int {
	return 50
}

func (m Model) previewHeight() int {
	return m.previewViewport.height
}

func (m *Model) cursorUp() {
	m.cursor.currDiscID = max(m.cursor.currDiscID-1, 0)
}

func (m *Model) cursorDown() {
	newCursor := min(m.cursor.currDiscID+1, len(m.data)-1)
	newCursor = max(newCursor, 0)

	m.cursor.currDiscID = newCursor
}

func (m *Model) prevSection() {
	m.cursor.currSectionID = max(m.cursor.currSectionID-1, 0)
}

func (m *Model) nextSection() {
	newCursor := min(m.cursor.currSectionID+1, len(m.sections)-1)
	newCursor = max(newCursor, 0)

	m.cursor.currSectionID = newCursor
}

func initScreen(owner, name string) tea.Cmd {
	return func() tea.Msg {
		categories, err := FetchDiscCategories(owner, name)
		if err != nil {
			log.Fatal(err)
		}

		var sections []Section

		for _, category := range categories {
			sections = append(
				sections,
				Section{Name: category.Name, CategoryID: category.ID},
			)
		}

		data, err := FetchAllDiscs(owner, name)
		if err != nil {
			log.Fatal(err)
		}

		return initMsg{
			data,
			sections,
		}
	}
}

func (m Model) Init() tea.Cmd {
	return initScreen(m.Owner, m.Name)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.help.Width = msg.Width
		verticalMargins := headerHeight + footerHeight + pagerHeight

		if !m.ready {
			m.contentViewport = contentViewport{
				width:  m.width - 50,
				height: msg.Height - verticalMargins - 1,
			}
			m.previewViewport = previewViewport{
				width:  0,
				height: msg.Height - verticalMargins + 1,
			}
			m.ready = true
		} else {
			m.contentViewport.height = msg.Height - verticalMargins - 1
			m.previewViewport.height = msg.Height - verticalMargins + 1
		}

		return m, nil

	case initMsg:
		m.sections = msg.sections
		m.data = msg.data
		return m, nil

	}

	cmds = append(cmds, m.updateSection(msg))

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateSection(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.CursorUp):
			m.cursorUp()

		case key.Matches(msg, m.keyMap.CursorDown):
			m.cursorDown()

		case key.Matches(msg, m.keyMap.PrevSection):
			m.prevSection()

		case key.Matches(msg, m.keyMap.NextSection):
			m.nextSection()

		}
	}

	return cmd
}

func (m Model) View() string {
	s := strings.Builder{}

	s.WriteString(m.tabsView())
	s.WriteString("\n")

	s.WriteString(m.sectionView())
	s.WriteString("\n")

	s.WriteString(m.helperView())

	return s.String()
}
