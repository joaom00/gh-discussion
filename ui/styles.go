package ui

import "github.com/charmbracelet/lipgloss"

var (
	headerHeight       = 6
	footerHeight       = 2
	discRowHeight      = 2
	singleRuneWidth    = 4
	mainContentPadding = 1
	pagerHeight        = 2
	cellPadding        = cellStyle.GetPaddingLeft() + cellStyle.GetPaddingRight()

	authorCellWidth     = 15
	answeredByCellWidth = 15
	answeredCellWidth   = lipgloss.Width(cellStyle.Render("Answered"))
	upvotesCellWidth    = lipgloss.Width(cellStyle.Render("Upvotes"))
	updatedAtCellWidth  = lipgloss.Width(cellStyle.Render(" Updated"))
	usedWidth           = authorCellWidth + answeredByCellWidth + answeredCellWidth + upvotesCellWidth + updatedAtCellWidth

	indigo             = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#383B5B"}
	subtleIndigo       = lipgloss.AdaptiveColor{Light: "#5A57B5", Dark: "#242347"}
	selectedBackground = lipgloss.AdaptiveColor{Light: subtleIndigo.Light, Dark: subtleIndigo.Dark}
	border             = lipgloss.AdaptiveColor{Light: indigo.Light, Dark: indigo.Dark}
	secondaryBorder    = lipgloss.AdaptiveColor{Light: indigo.Light, Dark: "#39386B"}
	faintBorder        = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#28283B"}
	mainText           = lipgloss.AdaptiveColor{Light: "#242347", Dark: "#E2E1ED"}
	secondaryText      = lipgloss.AdaptiveColor{Light: indigo.Light, Dark: indigo.Dark}
	faintText          = lipgloss.AdaptiveColor{Light: indigo.Light, Dark: "#3E4057"}
	warningText        = lipgloss.AdaptiveColor{Light: "#F23D5C", Dark: "#F23D5C"}
	successText        = lipgloss.AdaptiveColor{Light: "#3DF294", Dark: "#3DF294"}

	cellStyle = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).MaxHeight(1)

	titleCellStyle = cellStyle.Copy().Background(selectedBackground)

	selectedCellStyle = cellStyle.Copy().Background(selectedBackground)

	sidebarStyle = lipgloss.NewStyle().
			Padding(0, 2).
			BorderLeft(true).
			BorderStyle(lipgloss.Border{
			Top:         "",
			Bottom:      "",
			Left:        "│",
			Right:       "",
			TopLeft:     "",
			TopRight:    "",
			BottomRight: "",
			BottomLeft:  "",
		}).
		BorderForeground(border)

	headerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(secondaryBorder).
			BorderBottom(true)

	discussionStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(faintBorder).
			BorderBottom(true)

	pagerStyle = lipgloss.NewStyle().
			MarginTop(1).
			Bold(true).
			Foreground(faintText)

	mainTextStyle = lipgloss.NewStyle().
			Foreground(mainText).
			Bold(true)

	helperStyle = lipgloss.NewStyle().
			Height(footerHeight).
			BorderTop(true).BorderStyle(lipgloss.NormalBorder()).BorderForeground()
)

func makeCellStyle(isSelected bool) lipgloss.Style {
	if isSelected {
		return selectedCellStyle.Copy()
	}

	return cellStyle.Copy()
}
