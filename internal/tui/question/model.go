package question

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/k-algorithm/idk/internal/search"
)

type QuestionModel struct {
	list list.Model
	viewport viewport.Model
}

type item string

type itemDelegate struct{}

func InitialModel(qArr []search.Question) QuestionModel {
	qItems := questionToItem(qArr)
	l := list.New(qItems, itemDelegate{}, 0, 0)
	l.Title = "Questions"
	l.SetShowHelp(false)

	return QuestionModel{list: l}
}

func questionToItem(qArr []search.Question) []list.Item {
	items := make([]list.Item, len(q))
	for i, q := range qArr {
		items[i] = item(q.Title)
	}
	return items
}