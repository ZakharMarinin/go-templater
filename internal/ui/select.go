package ui

import (
	"fmt"
	"go-templater/internal/domain/entity"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type item struct {
	title    string
	desc     string
	template entity.Template
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type selectModel struct {
	list     list.Model
	selected item
	finished bool
}

func (ui *UI) Select(templates []*entity.Template) (*entity.Template, error) {
	items := make([]list.Item, 0, len(templates))

	for _, t := range templates {
		items = append(items, item{
			title:    t.Name,
			desc:     t.Description,
			template: *t,
		})
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a Template"
	l.SetHeight(10)
	l.SetWidth(40)

	m := selectModel{list: l}
	
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	res := finalModel.(selectModel)
	if res.finished {
		return &res.selected.template, nil
	}

	return nil, fmt.Errorf("no template selected")
}

func (m selectModel) Init() tea.Cmd {
	return nil
}

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case cancel, "q":
			return m, tea.Quit
		case enter:
			i := m.list.SelectedItem()
			m.selected = i.(item)
			m.finished = true

			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m selectModel) View() tea.View {
	v := tea.NewView(lipgloss.NewStyle().Margin(1, 2).Render(m.list.View()))
	v.AltScreen = true

	return v
}
