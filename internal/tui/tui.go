package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/k-algorithm/idk/internal/search"
	"strings"
)

const (
	getInput = iota
	loadingQuestion
	viewQuestion
)

type gotQuestions struct {
	result []list.Item
}

type model struct {
	query        textinput.Model
	spinner      spinner.Model
	questionList list.Model

	state int
	err   error
}

func InitializeModel() tea.Model {
	query := textinput.New()
	query.Placeholder = "Question"
	query.Focus()
	query.CharLimit = 156
	query.Width = 20

	s := spinner.New()
	s.Spinner = spinner.Dot

	return model{
		query:   query,
		spinner: s,
		err:     nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.questionList.FilterState() == list.Filtering {
			break
		}
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

		if m.state == getInput {
			switch msg.Type {
			case tea.KeyEnter:
				qString := strings.TrimSpace(m.query.Value())
				m.state = loadingQuestion
				cmds = append(cmds, spinner.Tick)
				cmds = append(cmds, m.fetchQuestions(qString))
			}
		}
	case tea.WindowSizeMsg:
		switch m.state {
		case getInput:
			m.query.Width = msg.Width
		case viewQuestion:
			m.questionList.SetWidth(msg.Width)
			m.questionList.SetHeight(msg.Height)
		}

	case gotQuestions:
		m.state = viewQuestion
		m.questionList = list.New(msg.result, list.NewDefaultDelegate(), 0, 0)

	}
	switch m.state {
	case getInput:
		m.query, cmd = m.query.Update(msg)
	case loadingQuestion:
		m.spinner, cmd = m.spinner.Update(msg)
	case viewQuestion:
		m.questionList, cmd = m.questionList.Update(msg)
	}
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	// If there's an error, print it out and don't do anything else.
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	var content string
	switch m.state {
	case getInput:
		content = fmt.Sprintf(
			"\nWrite questions here...\n\n%s\n\n%s\n\n%s",
			m.query.View(),
			"(ctrl+c to quit)",
		) + "\n"
	case loadingQuestion:
		content = fmt.Sprintf("%s fetching results... please wait.", m.spinner.View())
	case viewQuestion:
		content = m.questionList.View()
	}
	return "\n" + content
}

func (m model) fetchQuestions(q string) tea.Cmd {
	return func() tea.Msg {
		searchResult := search.Google(search.GoogleParam{
			Query:    q,
			PageSize: 10,
		})
		questions := search.Questions(searchResult.QuestionIDs)
		qItem := QuestionToItem(questions)
		return gotQuestions{result: qItem}
	}
}
