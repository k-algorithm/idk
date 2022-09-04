package search_result

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/k-algorithm/idk/internal/search"
)

type Item struct {
	Title    string
	Question search.Question
}

func (i Item) FilterValue() string { return i.Title }

type BubbleSearchResult struct {
	query    string
	list     list.Model
	viewport viewport.Model
}

func (sr BubbleSearchResult) UpdateQuery(query string) {
	sr.query = query
}

func InitialModel() BubbleSearchResult {
	return BubbleSearchResult{}
}
