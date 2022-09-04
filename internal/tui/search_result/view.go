package search_result

import (
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (sr BubbleSearchResult) View() string {
	return docStyle.Render(sr.list.View())
}
