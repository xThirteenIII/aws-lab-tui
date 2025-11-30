// here are all the update functions for each different state
package model

import (
	"net"

	tea "github.com/charmbracelet/bubbletea"
)

// updateMainMenu updates the model when the user is in the main menu titlescreen.
func (m model) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {

	// Main Menu Title
	m.mainMenuList.Title = "Select the IoT Tool you want to use"

	switch message := msg.(type) {

	// Is the message a keyPress?
	case tea.KeyMsg:

		// Cool, what key was pressed?
		switch message.String() {

		// close the program
		case "ctrl+c", "q":
			return m, tea.Quit

		// go to selectIoTJob state when pressing enter
		case "enter":
			m.currentState = selectIoTJob
			m.stateStack.Push(selectIoTJob)
			m.initSelectJob()
		}
	// WindowSizeMsg is used to report the terminal size. It's sent to Update once
	// initially and then on every terminal resize. Note that Windows does not
	// have support for reporting when resizes occur as it does not support the
	// SIGWINCH signal.
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.mainMenuList.SetSize(message.Width-h, message.Height-v)
	}

	var cmd tea.Cmd
	m.mainMenuList, cmd = m.mainMenuList.Update(msg)
	return m, cmd
}

// updateSelectJob updates the model when the user is in the selectIoTJob state.
func (m model) updateSelectJob(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.jobInput.SetSuggestions(m.suggestions.jobSuggestions)
	switch message := msg.(type) {

	// Was a key pressed?
	case tea.KeyMsg:

		// Cool, which one?
		switch message.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		// Go back in the state history when pressing esc
		case "esc":
			m.stateStack.Pop()
			m.currentState = m.stateStack.Peek()
		case "enter":
			m.suggestions.addJobSuggestion(m.jobInput.Value())

			m.currentState = selectThing
			m.stateStack.Push(selectThing)
		}
	}

	// call updateInputs to update input typing
	cmd := m.updateInputs(msg)
	return m, cmd
}

// updateSelectJob updates the model when the user is in the selectIoTJob state.
func (m model) updateSelectThing(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch message := msg.(type) {

	// Was a key pressed?
	case tea.KeyMsg:

		// Cool, which one?
		switch message.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		// Go back in the state history when pressing esc
		case "esc":
			// Avoid getting the same error when we go back in history
			m.err = ""
			m.stateStack.Pop()
			m.currentState = m.stateStack.Peek()
		case "enter":
			// We do not use input.Validate() for now
			// Because don't know what to do with Model.Err field
			_, err := net.ParseMAC(m.thingInput.Value())
			if err != nil {
				m.currentState = selectThing
				m.initSelectThing()
				m.err = "Please type a valid mac address"
				break
			}
			m.suggestions.addMacSuggestion(m.thingInput.Value())
			m.thingInput.SetSuggestions(m.suggestions.macSuggestions)
		}
	}

	// call updateInputs to update input typing
	cmd := m.updateThingInputs(msg)
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	// Update is the Bubble Tea update loop.
	// TODO: useless to update both inputs, need to be refactored
	m.jobInput, cmd = m.jobInput.Update(msg)
	return cmd
}

func (m *model) updateThingInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	// Update is the Bubble Tea update loop.
	// TODO: useless to update both inputs, need to be refactored
	m.thingInput, cmd = m.thingInput.Update(msg)
	return cmd
}
