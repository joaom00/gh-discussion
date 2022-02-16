package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Discussion struct {
	Title       string
	URL         string
	UpvoteCount int
	Body        string
	Answer      struct {
		Author struct {
			Login string
		}
		IsAnswer bool
	}
	Author struct {
		Login string
	}
}

func (d Discussion) renderUpvotes(isSelected bool) string {
	return makeCellStyle(
		isSelected,
	).Width(upvotesCellWidth).
		MaxWidth(upvotesCellWidth).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("↑ %d", d.UpvoteCount))
}

func (d Discussion) renderAnswerStatus(isSelected bool) string {
	answerCellStyle := makeCellStyle(
		isSelected,
	).Width(answeredCellWidth).
		MaxWidth(answeredCellWidth).
		Align(lipgloss.Center)

	if d.Answer.IsAnswer {
		return answerCellStyle.Copy().Foreground(successText).Render("")
	}

	return answerCellStyle.Copy().Foreground(faintText).Render("")
}

func (d Discussion) renderTitle(viewportWidth int, isSelected bool) string {
	totalWidth := getTitleWidth(viewportWidth)

	title := makeCellStyle(isSelected).Render(truncateString(d.Title, totalWidth))

	return makeCellStyle(isSelected).
		Width(totalWidth).
		MaxWidth(totalWidth).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, title))
}

func (d Discussion) renderAuthor(isSelected bool) string {
	return makeCellStyle(
		isSelected,
	).Width(authorCellWidth).
		Render(truncateString(d.Author.Login, authorCellWidth-cellPadding))
}

func (d Discussion) renderAnsweredBy(isSelected bool) string {
	if d.Answer.IsAnswer {
		return makeCellStyle(
			isSelected,
		).Width(answeredByCellWidth).
			Render(truncateString(d.Answer.Author.Login, answeredByCellWidth-cellPadding))
	}

	return makeCellStyle(isSelected).Width(answeredByCellWidth).Render("--")
}

func (d Discussion) render(isSelected bool, viewportWidth int) string {
	upvotesCell := d.renderUpvotes(isSelected)
	answeredCell := d.renderAnswerStatus(isSelected)
	titleCell := d.renderTitle(viewportWidth, isSelected)
	authorCell := d.renderAuthor(isSelected)
	answeredByCell := d.renderAnsweredBy(isSelected)

	rowStyle := discussionStyle.Copy()

	return rowStyle.
		Width(viewportWidth).
		MaxWidth(viewportWidth).
		MaxHeight(discRowHeight).
		Render(lipgloss.JoinHorizontal(lipgloss.Left,
			upvotesCell,
			answeredCell,
			titleCell,
			authorCell,
			answeredByCell,
		))
}
