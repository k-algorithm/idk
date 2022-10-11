package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/k-algorithm/idk/internal/search"
	"io"
)

type item string

func (i item) FilterValue() string { return string(i) }

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                               { return 1 }
func (d ItemDelegate) Spacing() int                              { return 0 }
func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	qTitle, ok := listItem.(item)
	if !ok {
		return
	}

	raw := fmt.Sprintf("[Question %d] %s\n", index+1, qTitle)
	fn := ListItemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return ListSelectedListItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(raw))
}

func QuestionToItem(qArr []search.Question) []list.Item {
	items := make([]list.Item, len(qArr))
	for i, q := range qArr {
		items[i] = item(q.Title)
	}
	return items
}

type QuestionModel struct {
	list list.Model
}

func InitialModel(qArr []search.Question) QuestionModel {
	qItems := QuestionToItem(qArr)
	l := list.New(qItems, ItemDelegate{}, 0, 0)
	l.Title = "Questions"
	l.SetShowHelp(true)

	return QuestionModel{list: l}
}

func (qm QuestionModel) Init() tea.Cmd {
	return nil
}

func (qm QuestionModel) Update(msg tea.Msg) (QuestionModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		qm.list.SetWidth(msg.Width)
		return qm, nil
	}

	var cmd tea.Cmd
	qm.list, cmd = qm.list.Update(msg)

	return qm, cmd
}

func (qm QuestionModel) View() string {
	return "\n" + qm.list.View()
}
