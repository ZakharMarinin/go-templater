package ui

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type timeoutMsg struct{}

type statusModel struct {
	msg      string
	spinner  spinner.Model
	quitting bool
	styles   *StatusStyles
	duration time.Duration
}

type StatusStyles struct {
	statusTextStyle *lipgloss.Style
}

func (m statusModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.Tick(m.duration, func(t time.Time) tea.Msg {
			return timeoutMsg{}
		}),
	)
}

func (m statusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		
		return m, tea.Quit

	case timeoutMsg:
		m.quitting = true
		
		return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return tea.Quit()
		})

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		
		return m, cmd
	}

	return m, nil
}

func (m statusModel) View() tea.View {
	if m.quitting {
		return tea.NewView("\r\033[K")
	}

	msg := fmt.Sprintf("\r%s %s...", m.spinner.View(), m.msg)
	content := m.styles.statusTextStyle.Render(msg)
	
	return tea.NewView(content)
}

func (ui *UI) ShowStatus(msg string, duration time.Duration) error {
	var (
		statusTextStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ui.cfg.Colors.Highlight))
	)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = statusTextStyle

	m := statusModel{
		msg:      msg,
		spinner:  s,
		duration: duration,
		styles: &StatusStyles{
			statusTextStyle: &statusTextStyle,
		},
	}

	p := tea.NewProgram(m)
	_, err := p.Run()

	return err
}