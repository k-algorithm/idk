package question

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/k-algorithm/idk/internal/search"
)

type BubbleQuestion struct {
	viewport viewport.Model
	ready    bool
	question search.Question
	answer   []search.Answer
}

func InitializeModel(question search.Question) BubbleQuestion {
	answer = search.Answers(question.Id)
	return BubbleQuestion{}
}

func (bq BubbleQuestion) Init() tea.Cmd {
	return nil
}

func (bq BubbleQuestion) Update(msg tea.Msg) (BubbleQuestion, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return bq, tea.Batch()
		}
	}
}

func (bq BubbleQuestion) View() string {
	return ""
}
