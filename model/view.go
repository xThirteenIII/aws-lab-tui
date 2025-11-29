// Here are all the view functions for each different state
package model

import "fmt"

func (m model) viewMainMenu() string {

	// The header
	mainMenu := "Select the IoT Tool you want to use\n\n"

	// Iterate over choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// render the row
		mainMenu += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	// The footer
	mainMenu += "\nPress q to quit.\n"

	// send the UI for rendering
	return mainMenu
}

func (m model) viewSelectJob() string {

	// The header
	selectJobMenu := "Type the name you want to give to the IoT Job\n\n"

	selectJobMenu += "\nPress q to quit. Esc to back.\n"

	// send the UI for rendering
	return selectJobMenu
}
