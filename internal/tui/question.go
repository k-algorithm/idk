package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/k-algorithm/idk/internal/search"
)

type item string

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	qTitle, ok := listItem.(item)
	if !ok {
		return
	}

	raw := fmt.Sprintf("[Question %d] %s\n", index+1, qTitle)

	if index == m.Index() {
		line := ListSelectedListItemStyle.Render("> " + raw)
	} else {
		line := ListItemStyle.Render(raw)
	}

	fmt.Fprint(w, line)
}

func InitialModel(qArr []search.Question) QuestionModel {
	qItems := questionToItem(qArr)
	l := list.New(qItems, itemDelegate{}, 0, 0)
	l.Title = "Questions"
	l.SetShowHelp(false)

	return QuestionModel{list: l}
}

func questionToItem(qArr []search.Question) []list.Item {
	items := make([]list.Item, len(qArr))
	for i, q := range qArr {
		items[i] = item(q.Title)
	}
	return items
}

type QuestionModel struct {
	list     list.Model
	viewport viewport.Model
}

func (qm QuestionModel) Init() tea.Cmd {
	return nil
}

func (qm QuestionModel) Update(msg tea.Msg) (QuestionModel, tea.Cmd) {
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

func (qm QuestionModel) View() string {
	return "\n" + qm.list.View()
}
