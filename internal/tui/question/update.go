package question

import (
	"github.com/k-algorithm/idk/internal/tui"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (qm QuestionModel) Update(msg tea.Msg) (BubbleGroup, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		horizontal, vertical := ListStyle.GetFrameSize()
		paginatorHeight := lipgloss.Height(qm.list.Paginator.View())

		qm.list.SetSize(msg.Width-horizontal, msg.Height-vertical-paginatorHeight)
	}

	var cmd tea.Cmd
	qm.list, cmd = qm.list.Update(msg)

	return qm, cmd
}