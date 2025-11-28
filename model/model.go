// package model defines the Model struct, which defines application states.
// It implements three methods which operate on the model structure.
//
// Init: a function that returns an initial command for the application to run
// Update: a function that updates incoming events and updates the model accordingly
// View: a function that renders UI based on the data in the model

package model

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string         // items on the tool list
	cursor   int              // which tool item the cursor is pointing at
	selected map[int]struct{} // which tool items are selected
}

// initialModel defines the initial state of the application.
// Defining a function to return the initial model, but could use a variable elsewhere, too.
func initialModel() model {
	return model{

		// List all the IoT Tools usable with the app.
		choices: []string{"Send IoT Jobs", "Dictionary", "Disenroll Inverter", "Upload .json to AWS S3"},

		// A map which indicates which choices are selected.
		// For now, only one choice is possible, for later states (e.g. Downloading commands from Dictionary), we can select multiples.
		// The keys refer to the indexes of the "choices" slice, above.
		selected: make(map[int]struct{}),
	}
}

// Init returns an initial command that performs some I/O.
// No command is needed for now, so we return nil, which translates to no command.
func (m model) Init() tea.Cmd {

	// nil means "no I/O for now, please"
	return nil
}

// Update method is called when "things happen".
// Its job is to look at what's happened and return an updated model in response.
// It can also return a command to execute with it.
// For now, when user presses keys to navigate through Tool List, it will update the model cursor.
//
// The "something happened" comes in the form of a tea.Msg, which can by any type (interface{}).
// Messages are the result of some I/O that took place, such as keypresses, timer tick, or response from a server.
//
// To figure out the type of message, use a type switch or a type assertion.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	// tea.KeyMsg are automatically sent to the Update function when keys are pressed.
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up.
		case "up", "k":
			if m.cursor > 0 {

				// Cursor goes back in the slice when moving up, not forward!
				// [0]
				// [1]
				// [2]
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down.
		case "down", "j":
			if m.cursor < len(m.choices)-1 {

				m.cursor++
			}

		// The "enter" key toggles
		// the selected state for the item that the cursor is pointing at.
		case "enter":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		}

	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// No tea.Cmd returned.
	return m, nil
}

// View looks at the model at its current state, and returns a string, which is the updated UI!
// Redrawing logic and stuff like that is taken care for by BubbleTea.
func (m model) View() string {
	// The header
	ui := "Select the IoT Tool you want to use\n\n"

	// Iterate over choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not slected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected
		}

		// render the row
		ui += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	ui += "\nPress q to quit.\n"

	// send the UI for rendering
	return ui
}
