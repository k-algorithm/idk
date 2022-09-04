package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/k-algorithm/idk/internal/tui/search_result"
)

type state int

const (
	showSearchView = iota
	showQuestionView
)

type model struct {
	state         state
	query         textinput.Model
	searchResult  search_result.BubbleSearchResult
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
			m.searchResult.UpdateQuery(m.query.Value())
			m.searchResult, cmd = m.searchResult.Update(msg)
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

	return fmt.Sprintf(
		"\nWrite questions here...\n\n%s\n\n%s\n\n%s",
		m.query.View(),
		"(ctrl+c to quit)",
		m.searchResult.View(),
	) + "\n"
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
