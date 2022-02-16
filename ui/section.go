package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) syncMainViewport() {
	m.mainViewport.model.Width = m.width - getSidebarWidth()
	discussions := m.renderDiscussionsList()
	m.mainViewport.model.SetContent(discussions)
}

func (m Model) makeRenderDiscussionCmd() tea.Cmd {
	return func() tea.Msg {
		return discussionsRenderedMsg{
			content: m.renderDiscussionsList(),
		}
	}
}

func (m *Model) renderTableHeader() string {
	upvotesCell := titleCellStyle.Copy().Width(upvotesCellWidth).Render("Upvotes")
	answeredCell := titleCellStyle.Copy().Width(answeredCellWidth).Render("Answered")
	titleCell := titleCellStyle.Copy().Width(getTitleWidth(m.mainViewport.model.Width)).MaxWidth(getTitleWidth(m.mainViewport.model.Width)).Render("Title")
	authorCell := titleCellStyle.Copy().Width(authorCellWidth).Render("Author")
	answeredByCell := titleCellStyle.Copy().Width(answeredByCellWidth).Render("Answered By")

	return headerStyle.
		PaddingLeft(mainContentPadding).
		PaddingRight(mainContentPadding).
		Width(m.mainViewport.model.Width).
		MaxWidth(m.mainViewport.model.Width).
		Render(lipgloss.JoinHorizontal(
			lipgloss.Left,
			upvotesCell,
			answeredCell,
			titleCell,
			authorCell,
			answeredByCell,
		))
}

func (m *Model) renderDiscussionsList() string {
	var renderedDiscussions []string

	for dID, d := range *m.data {
		isSelected := m.cursor.currDiscId == dID
		renderedDiscussions = append(renderedDiscussions, d.render(isSelected, m.mainViewport.model.Width))
	}

	return lipgloss.NewStyle().Render(lipgloss.JoinVertical(lipgloss.Left, renderedDiscussions...))
}

func (m *Model) renderCurrentSection() string {
	return lipgloss.NewStyle().
		PaddingLeft(mainContentPadding).
		PaddingRight(mainContentPadding).
		MaxWidth(m.mainViewport.model.Width).
		Render(m.renderMainViewport())
}

func (m *Model) renderMainViewport() string {
	pagerContent := ""

	numDiscs := len(*m.data)
	if numDiscs > 0 {
		pagerContent = fmt.Sprintf("Discussions %v/%v", m.cursor.currDiscId+1, numDiscs)
	}

	return lipgloss.JoinVertical(lipgloss.Top, m.mainViewport.model.View(), pagerStyle.Copy().Render(pagerContent))

}
