// Here are all the view functions for each different state
package model

func (m model) viewMainMenu() string {

	return docStyle.Render(m.mainMenuList.View())
}

func (m model) viewSelectJob() string {

	// The header
	selectJobMenu := "Type the name you want to give to the IoT Job\n\n"

	selectJobMenu += "\nPress q to quit. Esc to back.\n"

	// send the UI for rendering
	return selectJobMenu
}
