// here are all the update functions for each different state
package model

import tea "github.com/charmbracelet/bubbletea"

// updateMainMenu updates the model when the user is in the main menu titlescreen.
func (m model) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch message := msg.(type) {

	// Is the message a keyPress?
	case tea.KeyMsg:

		// Cool, what key was pressed?
		switch message.String() {

		// close the program
		case "ctrl+c", "q":
			return m, tea.Quit
		// navigate up
		// cursor is decremented even if counterintuitive!
		//
		// choices[0]
		// choices[1]
		// choices[2]
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		// navigate down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.currentState = selectIoTJob
			m.stateStack.Push(selectIoTJob)
		}

	}

	return m, nil
}

// updateSelectJob updates the model when the user is in the selectIoTJob state.
func (m model) updateSelectJob(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch message := msg.(type) {

	// Was a key pressed?
	case tea.KeyMsg:

		// Cool, which one?
		switch message.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			m.stateStack.Pop()
			m.currentState = m.stateStack.Peek()
		}
	}
	return m, nil
}
