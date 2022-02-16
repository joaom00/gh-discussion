package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	width           int
	ready           bool
	mainViewport    MainViewport
	sidebarViewport viewport.Model
	cursor          cursor
	data            *[]Discussion
	keys            keyMap
	help            help.Model
}

type MainViewport struct {
	model viewport.Model
}

type cursor struct {
	currDiscId int
}

type discussionsRenderedMsg struct {
	content string
}

func NewModel(data *[]Discussion) Model {
	helperModel := help.NewModel()
	style := lipgloss.NewStyle().Foreground(secondaryText)
	helperModel.Styles = help.Styles{
		ShortDesc:      style.Copy(),
		FullDesc:       style.Copy(),
		ShortSeparator: style.Copy(),
		FullSeparator:  style.Copy(),
		ShortKey:       style.Copy(),
		FullKey:        style.Copy(),
		Ellipsis:       style.Copy(),
	}
	return Model{
		help: helperModel,
		data: data,
		keys: keys,
		cursor: cursor{
			currDiscId: 0,
		},
	}
}

func getTitleWidth(viewportWidth int) int {
	return viewportWidth - usedWidth
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Down):
			m.nextDisc()
			m.syncMainViewport()
			m.syncSidebarViewport()
			return m, nil
		case key.Matches(msg, m.keys.Up):
			m.prevDisc()
			m.syncMainViewport()
			m.syncSidebarViewport()
			return m, nil
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.help.Width = msg.Width
		verticalMargins := headerHeight + footerHeight + pagerHeight

		if !m.ready {
			m.mainViewport = MainViewport{
				model: viewport.Model{
					Width:  m.width - getSidebarWidth(),
					Height: msg.Height - verticalMargins - 1,
				},
			}
			m.sidebarViewport = viewport.Model{
				Width:  0,
				Height: msg.Height - verticalMargins + 1,
			}
			m.ready = true
		} else {
			m.mainViewport.model.Height = msg.Height - verticalMargins - 1
			m.sidebarViewport.Height = msg.Height - verticalMargins + 1
			m.syncMainViewport()
			m.syncSidebarViewport()
		}

		return m, m.makeRenderDiscussionCmd()

	case discussionsRenderedMsg:
		m.mainViewport.model.SetContent(msg.content)
		m.syncSidebarViewport()
		return m, nil
	}

	m.sidebarViewport, cmd = m.sidebarViewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	paddedContentStyle := lipgloss.NewStyle().Padding(0, mainContentPadding)

	s := strings.Builder{}
	table := paddedContentStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top, m.renderTableHeader(), m.renderCurrentSection()),
	)
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, table, m.renderSidebar()))
	s.WriteString("\n")
	s.WriteString(m.renderHelper())

	return s.String()
}
