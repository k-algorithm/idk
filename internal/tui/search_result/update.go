package search_result

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/k-algorithm/idk/internal/search"
	"github.com/k-algorithm/idk/internal/tui/common"
)

func (sr BubbleSearchResult) Update(msg tea.Msg) (BubbleSearchResult, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		horizontal, vertical := common.ListStyle.GetFrameSize()
		paginatorHeight := lipgloss.Height(sr.list.Paginator.View())

		sr.list.SetSize(msg.Width-horizontal, msg.Height-vertical-paginatorHeight)
		sr.viewport = viewport.New(msg.Width, msg.Height)
		// sr.viewport.SetContent(sr.detailView())
	}

	var cmd tea.Cmd

	searchResult := search.Google(search.GoogleParam{
		Query:    sr.query,
		PageSize: 10,
	})
	questions := search.Questions(searchResult.QuestionIDs)

	for i, question := range questions {
		sr.list.InsertItem(i, Item{
			Title:    question.Title,
			Question: question,
		})
	}

	return sr, cmd
}
