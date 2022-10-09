package question

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/k-algorithm/idk/internal/search"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

type BubbleQuestion struct {
	ready         bool // TODO(marv): implement loading screen
	viewport      viewport.Model
	question      search.Question
	answer        []search.Answer
	codeSnippsets []string
}

func InitializeModel(question search.Question) BubbleQuestion {
	bq := BubbleQuestion{}
	bq.answer = search.Answers(fmt.Sprint(question.Id))
	return bq
}

func (bq BubbleQuestion) Init() tea.Cmd {
	return nil
}

func (bq BubbleQuestion) Update(msg tea.Msg) (BubbleQuestion, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(bq.headerView())
		footerHeight := lipgloss.Height(bq.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !bq.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			bq.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			bq.viewport.YPosition = headerHeight
			bq.viewport.HighPerformanceRendering = false
			bq.viewport.SetContent(search.AnswersToString(bq.answer))
			bq.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			bq.viewport.YPosition = headerHeight + 1
		} else {
			bq.viewport.Width = msg.Width
			bq.viewport.Height = msg.Height - verticalMarginHeight
		}
	case tea.KeyMsg:
		switch msg.Type {
		default:
			switch msg.String() {
			case "cmd+c":
				clipboard.WriteAll(bq.codeSnippsets[0])
			}
		}
	}
	bq.viewport, cmd = bq.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return bq, tea.Batch(cmds...)
}

func (bq BubbleQuestion) View() string {
	return fmt.Sprintf("%s\n%s\n%s", bq.headerView(), search.AnswersToString(bq.answer), bq.footerView())
}

func (bq BubbleQuestion) headerView() string {
	title := titleStyle.Render(
		fmt.Sprintf("Question: %s ", bq.question.Title),
	)
	line := strings.Repeat("─", max(0, bq.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (bq BubbleQuestion) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", bq.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, bq.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
