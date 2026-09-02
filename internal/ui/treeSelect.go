package ui

import (
	"fmt"
	"github.com/ZakharMarinin/go-templater/internal/config"
	"github.com/ZakharMarinin/go-templater/internal/domain/entity"
	"github.com/ZakharMarinin/go-templater/pkg/response"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type treeItem struct {
	node  *entity.Node
	depth int
}

func (t treeItem) FilterValue() string { return t.node.Name }

func flattenNodes(nodes []*entity.Node, depth int, parentNode *entity.Node, parent map[*entity.Node]*entity.Node, items *[]treeItem) {
	for _, n := range nodes {
		parent[n] = parentNode
		*items = append(*items, treeItem{node: n, depth: depth})

		if n.IsDir {
			flattenNodes(n.Children, depth+1, n, parent, items)
		}
	}
}

func seedDefaultExclusions(items []treeItem) map[*entity.Node]bool {
	excluded := make(map[*entity.Node]bool)

	for _, it := range items {
		if !it.node.IsDir && strings.HasPrefix(it.node.Name, ".") {
			excluded[it.node] = true
		}
	}

	return excluded
}

func isExcluded(node *entity.Node, parent map[*entity.Node]*entity.Node, explicit map[*entity.Node]bool) bool {
	for cur := node; cur != nil; cur = parent[cur] {
		if explicit[cur] {
			return true
		}
	}

	return false
}

func hasContentCandidate(node *entity.Node) bool {
	return node.IsDir || node.Content != ""
}

type treeDelegate struct {
	parent   map[*entity.Node]*entity.Node
	explicit map[*entity.Node]bool
	colors   config.Colors
}

func (d treeDelegate) Height() int  { return 1 }
func (d treeDelegate) Spacing() int { return 0 }

func (d treeDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d treeDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	it, ok := listItem.(treeItem)
	if !ok {
		return
	}

	eligible := hasContentCandidate(it.node)

	checkbox := "[ ]"
	if eligible && !isExcluded(it.node, d.parent, d.explicit) {
		checkbox = "[x]"
	}

	label := it.node.Name
	switch {
	case it.node.IsDir:
		label += "/"
	case !eligible:
		label += " (no content)"
	}

	line := fmt.Sprintf("%s%s %s", strings.Repeat("  ", it.depth), checkbox, label)

	style := lipgloss.NewStyle()
	if !eligible {
		style = style.Foreground(lipgloss.Color("240"))
	}

	if index == m.Index() {
		style = style.Foreground(lipgloss.Color(d.colors.Highlight))
	}

	fmt.Fprint(w, style.Render(line)) //nolint: errcheck
}

type treeModel struct {
	list     list.Model
	parent   map[*entity.Node]*entity.Node
	explicit map[*entity.Node]bool
	finished bool
}

func (ui *UI) SelectContentExclusions(nodes []*entity.Node) error {
	if len(nodes) == 0 {
		return nil
	}

	var items []treeItem

	parent := make(map[*entity.Node]*entity.Node)
	flattenNodes(nodes, 0, nil, parent, &items)

	explicit := seedDefaultExclusions(items)

	listItems := make([]list.Item, 0, len(items))
	for _, it := range items {
		listItems = append(listItems, it)
	}

	delegate := treeDelegate{parent: parent, explicit: explicit, colors: ui.cfg.Colors}

	l := list.New(listItems, delegate, 0, 0)
	l.Title = "Choose what to exclude from content copying (space: toggle, enter: confirm)"
	l.SetShowStatusBar(false)
	l.SetHeight(20)
	l.SetWidth(70)

	m := treeModel{list: l, parent: parent, explicit: explicit}

	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	res := finalModel.(treeModel)
	if !res.finished {
		return response.ErrCanceled
	}

	applyExclusions(items, res.parent, res.explicit)

	return nil
}

func applyExclusions(items []treeItem, parent map[*entity.Node]*entity.Node, explicit map[*entity.Node]bool) {
	for _, it := range items {
		if !it.node.IsDir && isExcluded(it.node, parent, explicit) {
			it.node.Content = ""
		}
	}
}

func (m treeModel) Init() tea.Cmd {
	return nil
}

func (m treeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case cancel, "q", esc:
			return m, tea.Quit
		case enter:
			m.finished = true

			return m, tea.Quit
		case "space":
			if it, ok := m.list.SelectedItem().(treeItem); ok && hasContentCandidate(it.node) {
				m.explicit[it.node] = !m.explicit[it.node]
			}

			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m treeModel) View() tea.View {
	v := tea.NewView(lipgloss.NewStyle().Margin(1, 2).Render(m.list.View()))
	v.AltScreen = true

	return v
}
