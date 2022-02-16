package ui

import "github.com/charmbracelet/lipgloss"

func (m Model) renderHelper() string {
	return helperStyle.Copy().
		Width(m.width).
		Render(lipgloss.PlaceVertical(footerHeight, lipgloss.Top, m.help.View(m.keys)))
}
