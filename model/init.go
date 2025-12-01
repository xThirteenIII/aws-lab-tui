package model

import (
	"time"

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
	m.jobInput.SetValue(time.Now().Format("20060201"))
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

func (m *model) initSelectS3File() {
	// Set view size
	h, v := docStyle.GetFrameSize()
	m.s3FilesList.SetSize(m.width-h, m.height-v)
}
