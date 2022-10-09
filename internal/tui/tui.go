package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/k-algorithm/idk/internal/search"
	"github.com/k-algorithm/idk/internal/tui/question"
	"google.golang.org/appengine/log"
)

type state int

const (
	showIndexView = iota
	showQuestionView
)

type model struct {
	state         state
	query         textinput.Model
	questions     []search.Question
	searchResult  string
	bq            question.BubbleQuestion
	width, height int
	err           error
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			// Note(marv): temporary code to test question view
			if m.state == showIndexView && m.searchResult != "" {
				m.state = showQuestionView
				m.bq = question.InitializeModel(m.questions[0])
				m.bq.Update(msg)
				return m, cmd
			}

			questionString := ""
			searchResult := search.Google(search.GoogleParam{
				Query:    m.query.Value(),
				PageSize: 10,
			})
			m.questions = search.Questions(searchResult.QuestionIDs)
			for i, question := range m.questions {
				questionString += fmt.Sprintf("[Question %d] %s\n", i+1, question.Title)
			}
			m.searchResult = questionString
			return m, nil
		}
	}
	m.query, cmd = m.query.Update(msg)
	return m, cmd
}

func (m model) View() string {
	// If there's an error, print it out and don't do anything else.
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}
	switch m.state {
	case showQuestionView:
		return m.bq.View()
	case showIndexView:
		return fmt.Sprintf(
			"\nWrite questions here...\n\n%s\n\n%s\n\n%s",
			m.query.View(),
			m.searchResult,
			"(ctrl+c to quit)",
		) + "\n"
	default:
		log.Errorf(nil, "Unknown state: %d", m.state)
		return ""
	}
}

func InitializeModel() tea.Model {
	query := textinput.New()
	query.Placeholder = "Question"
	query.Focus()
	query.CharLimit = 156
	query.Width = 20

	return model{
		query: query,
		err:   nil,
	}
}
