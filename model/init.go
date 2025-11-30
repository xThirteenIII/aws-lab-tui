package model

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func (m *model) initSelectJob() {
	m.jobInput = textinput.New()
	m.jobInput.Cursor.Style = cursorStyle
	m.jobInput.CharLimit = 99
	m.jobInput.Focus()
	m.jobInput.PromptStyle = focusedStyle
	m.jobInput.TextStyle = focusedStyle
	m.jobInput.ShowSuggestions = true
}

func (m *model) initSelectThing() {
	m.thingInput = textinput.New()
	m.thingInput.Cursor.Style = cursorStyle
	m.thingInput.CharLimit = 17
	m.thingInput.Focus()
	m.thingInput.PromptStyle = focusedStyle
	m.thingInput.TextStyle = focusedStyle
	m.thingInput.ShowSuggestions = true
}
