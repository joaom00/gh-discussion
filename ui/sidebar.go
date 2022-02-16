package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Sidebar struct {
	model Model
	disc  Discussion
}

func (m *Model) syncSidebarViewport() {
	m.sidebarViewport.Width = getSidebarWidth()
	m.setSidebarViewportContent()
}

func (m *Model) renderSidebar() string {
	height := m.sidebarViewport.Height + pagerHeight
	style := sidebarStyle.Copy().
		Height(height).
		MaxHeight(height).
		Width(getSidebarWidth()).
		MaxWidth(getSidebarWidth())

	disc := &(*m.data)[m.cursor.currDiscId]

	if disc == nil {
		return style.Copy().
			Align(lipgloss.Center).
			Render(lipgloss.PlaceVertical(height, lipgloss.Center, "Select a Discussion..."))
	}

	return style.Copy().Render(lipgloss.JoinVertical(
		lipgloss.Left,
		m.sidebarViewport.View(),
		pagerStyle.Copy().Render(fmt.Sprintf("%d%%", int(m.sidebarViewport.ScrollPercent()*100))),
	))
}

func (m *Model) setSidebarViewportContent() {
	disc := &(*m.data)[m.cursor.currDiscId]
	if disc == nil {
		return
	}

	sidebar := Sidebar{
		model: *m,
		disc:  *disc,
	}

	s := strings.Builder{}
	s.WriteString("\n\n")
	s.WriteString(sidebar.renderTitle())
	s.WriteString("\n\n")
	s.WriteString(sidebar.renderBody())

	m.sidebarViewport.SetContent(s.String())
}

func (s *Sidebar) renderTitle() string {
	return mainTextStyle.Copy().Width(getSidebarWidth() - 6).
		Render(s.disc.Title)
}

func (s Sidebar) renderBody() string {
	body, _ := glamour.Render(s.disc.Body, "dark")
	return mainTextStyle.Copy().Width(getSidebarWidth() - 6).Render(body)
}

func getSidebarWidth() int {
	return 50
}
