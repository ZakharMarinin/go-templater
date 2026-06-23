package ui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type confirmModel struct {
	msg       string
	confirmed bool
	quitting  bool
	styles *ConfirmStyles
}

type ConfirmStyles struct {
	confirmBorderStyle *lipgloss.Style
	confirmTextStyle *lipgloss.Style
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "Y":
			m.confirmed = true
			m.quitting = true
			
			return m, tea.Quit
		case "n", esc, cancel:
			m.quitting = true
			
			return m, tea.Quit
		}
	}

	return m, nil
}

func (ui *UI) ConfirmOverwrite(fileName string) (bool, error) {
	var (
		confirmBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(ui.cfg.Colors.Highlight)).
		Padding(1)
		confirmTextStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ui.cfg.Colors.Black)).
		Bold(false)
	)

	m := confirmModel{
		msg: fmt.Sprintf("Warning: File '%s' already exists.\nOverwrite it? (Y/n)", fileName),
		styles: &ConfirmStyles{
			confirmBorderStyle: &confirmBorderStyle,
			confirmTextStyle: &confirmTextStyle,
		},
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return false, err
	}

	res := finalModel.(confirmModel)
	
	return res.confirmed, nil
}

func (m confirmModel) View() tea.View {
	content := m.styles.confirmTextStyle.Render(m.msg)

	if m.quitting {
		return tea.NewView(m.styles.confirmTextStyle.Render())
	}
	
	return tea.NewView(m.styles.confirmBorderStyle.Render(content))
}