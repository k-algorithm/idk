package question

import (
	"fmt"
	"io"

	"github.com/k-algorithm/idk/internal/tui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	qTitle, ok := listItem.(item)
	if !ok {
		return
	}

	raw = fmt.Sprintf("[Question %d] %s\n", index+1, qTitle)

	if index == m.Index() {
		line = ListSelectedListItemStyle.Render("> " + line)
	} else {
		line = ListItemStyle.Render(line)
	}

	fmt.Fprint(w, line)
}