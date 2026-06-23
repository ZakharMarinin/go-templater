package ui

import (
	"fmt"
	"go-templater/internal/domain/entity"
	"go-templater/pkg/response"
	"strings"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type textModel struct {
	focusIndex  int
	inputs      []textinput.Model
	cursorMode  cursor.Mode
	quitting    bool
	interupting bool
	noName      bool
	styles      *Styles
}

type Styles struct {
	focusedStyle *lipgloss.Style
	blurredStyle *lipgloss.Style
	focusedButton string
	blurredButton string
}

func (ui *UI) Input() (*entity.Variables, error) {
	var (
		focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ui.cfg.Colors.Highlight))
		blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ui.cfg.Colors.Black))

		focusedButton = focusedStyle.Render("[ Submit ]")
		blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	)
	
	m := textModel{
		focusIndex: 0,
		inputs:     make([]textinput.Model, 2),
		cursorMode: cursor.CursorBlink,
		styles: &Styles{
			focusedStyle: &focusedStyle,
			blurredStyle: &blurredStyle,
			focusedButton: focusedButton,
			blurredButton: blurredButton,
		},
	}

	for i := range m.inputs {
		t := textinput.New()
		t.CharLimit = 32

		s := t.Styles()
		s.Cursor.Color = lipgloss.Color(ui.cfg.Colors.Highlight)
		s.Focused.Prompt = focusedStyle
		s.Focused.Text = focusedStyle
		s.Blurred.Prompt = blurredStyle
		s.Blurred.Text = blurredStyle
		t.SetStyles(s)

		switch i {
		case 0:
			t.Placeholder = "Name (required)"
			t.SetWidth(32)
			t.Focus()
		case 1:
			t.Placeholder = "Description (optional)"
			t.SetWidth(64)
			t.CharLimit = 64
		}
		m.inputs[i] = t
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	res := finalModel.(textModel)

	if res.interupting {
		return nil, response.ErrCanceled
	}

	if res.inputs[0].Value() == "" {
		return nil, fmt.Errorf("name is required")
	}

	return &entity.Variables{
		Name:        res.inputs[0].Value(),
		Description: res.inputs[1].Value(),
	}, nil
}

func (m textModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m textModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case cancel, esc:
			m.interupting = true

			return m, tea.Quit

		case tab, shift+tab, enter, up, down:
			s := msg.String()

			if s == enter && m.focusIndex == len(m.inputs) {
				if m.inputs[0].Value() == "" {
					m.noName = true

					return m, nil
				}
				m.quitting = true

				return m, tea.Quit
			}

			if s == up || s == shift+tab {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()

					continue
				}
				m.inputs[i].Blur()
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *textModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m textModel) View() tea.View {
	var b strings.Builder
	var c *tea.Cursor

	for i, in := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
		if m.cursorMode != cursor.CursorHide && in.Focused() {
			c = in.Cursor()
			if c != nil {
				c.Y += i
			}
		}
	}

	button := &m.styles.blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &m.styles.focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n", *button)

	if m.noName {
		b.WriteString(m.styles.focusedStyle.Render("\nName is required!"))
	}

	v := tea.NewView(b.String())
	v.Cursor = c

	return v
}
