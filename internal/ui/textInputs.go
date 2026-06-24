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
	validationFailed      bool
	styles      *Styles
	fields      []*entity.FieldConfig
}

type Styles struct {
	focusedStyle  *lipgloss.Style
	blurredStyle  *lipgloss.Style
	focusedButton string
	blurredButton string
}

func (ui *UI) DynamicInput(title string, fields []*entity.FieldConfig) (map[string]string, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields provided")
	}

	var (
		focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ui.cfg.Colors.Highlight))
		blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ui.cfg.Colors.Black))

		focusedButton = focusedStyle.Render("[ Submit ]")
		blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	)

	m := textModel{
		focusIndex: 0,
		inputs:     make([]textinput.Model, len(fields)),
		cursorMode: cursor.CursorBlink,
		styles: &Styles{
			focusedStyle:  &focusedStyle,
			blurredStyle:  &blurredStyle,
			focusedButton: focusedButton,
			blurredButton: blurredButton,
		},
	}

	m.fields = fields

	for i, f := range fields {
		t := textinput.New()

		if f.CharLimit > 0 {
			t.CharLimit = f.CharLimit
		} else {
			t.CharLimit = 32
		}

		if f.Width > 0 {
			t.SetWidth(f.Width)
		} else {
			t.SetWidth(32)
		}

		s := t.Styles()
		s.Cursor.Color = lipgloss.Color(ui.cfg.Colors.Highlight)
		s.Focused.Prompt = focusedStyle
		s.Focused.Text = focusedStyle
		s.Blurred.Prompt = blurredStyle
		s.Blurred.Text = blurredStyle
		t.SetStyles(s)

		t.Placeholder = f.Placeholder

		if i == 0 {
			t.Focus()
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

	resultMap := make(map[string]string)
	for i, f := range fields {
		val := res.inputs[i].Value()

		if f.Required && val == "" {
			return nil, fmt.Errorf("field '%s' is required", f.Placeholder)
		}
		resultMap[f.Key] = val
	}

	return resultMap, nil
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

		case tab, shift + tab, enter, up, down:
			s := msg.String()

			if s == enter && m.focusIndex == len(m.inputs) {
			    for i, f := range m.fields {
			        if f.Required && m.inputs[i].Value() == "" {
			            m.validationFailed = true
															
			            return m, nil
			        }
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

	if m.validationFailed {
		b.WriteString(m.styles.focusedStyle.Render("\nPlease fill in all required fields!"))
	}

	v := tea.NewView(b.String())
	v.Cursor = c

	return v
}
