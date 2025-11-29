package model

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func (m *model) initSelectJob() {
	m.input = textinput.New()
	m.input.Cursor.Style = cursorStyle
	m.input.CharLimit = 99
	m.input.Focus()
	m.input.PromptStyle = focusedStyle
	m.input.TextStyle = focusedStyle
}
