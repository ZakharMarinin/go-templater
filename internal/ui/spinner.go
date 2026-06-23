package ui

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type TaskResultMsg struct {
	Err error
}

type SpinnerModel struct {
	spinner spinner.Model
	loading bool
	err     error
	title   string
	task    tea.Cmd
	styles  *SpinnerStyle
}

type SpinnerStyle struct {
	errorStyle    *lipgloss.Style
	completeStyle *lipgloss.Style
	textStyle     *lipgloss.Style
}

func (ui *UI) NewSpinner(title string, task func() error) error {
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ui.cfg.Colors.Error))
	completeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ui.cfg.Colors.Complete))
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ui.cfg.Colors.Highlight))

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = textStyle

	taskCmd := func() tea.Msg {
		return TaskResultMsg{Err: task()}
	}

	spinner := SpinnerModel{
		spinner: s,
		loading: true,
		title:   title,
		task:    taskCmd,
		styles: &SpinnerStyle{
			errorStyle:    &errorStyle,
			completeStyle: &completeStyle,
			textStyle:     &textStyle,
		},
	}

	p := tea.NewProgram(spinner)
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	m, ok := finalModel.(SpinnerModel)
	if !ok {
		return m.Error()
	}

	if m.Error() != nil {
		fmt.Println(m.styles.errorStyle.Render(fmt.Sprintf("❌ Error: %v", m.Error())))

		return m.Error()
	}

	fmt.Println(m.styles.completeStyle.Render(fmt.Sprintf("✓ %s complete!", title)))

	return nil
}

func (m SpinnerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.task)
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case cancel, "q":
			return m, tea.Quit
		}

	case TaskResultMsg:
		m.loading = false
		m.err = msg.Err

		if m.err != nil {
			return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return tea.Quit()
			})
		}

		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m SpinnerModel) View() tea.View {
	if m.loading {
		return tea.NewView(fmt.Sprintf("%s %s...\n", m.spinner.View(), m.title))
	}

	return tea.NewView("")
}

func (m SpinnerModel) Error() error {
	return m.err
}
