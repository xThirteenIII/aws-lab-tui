// package model defines the Model struct, which defines application states.
// It implements three methods which operate on the model structure.
//
// Init: a function that returns an initial command for the application to run
// Update: a function that updates incoming events and updates the model accordingly
// View: a function that renders UI based on the data in the model

package model

import (
	"aws-iot-tui/stack"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {

	// stateStack is a Stack structure that holds states history
	// It allows to go back and forth in the menu.
	stateStack   *stack.Stack[state]
	currentState state

	choices []string // items on the tool list
	cursor  int      // which tool item the cursor is pointing at
}

// initialModel defines the initial state of the application.
// Defining a function to return the initial model, but could use a variable elsewhere, too.
func InitialModel() model {
	// initStack creates a new empty stack
	initStack := stack.NewStack[state]()

	// push mainMenu as default first state onto the Stack
	initStack.Push(mainMenu)
	return model{

		// WARNING: to set the currentState is crucial, otherwise program will create unlimited goRoutines!
		// runtime: goroutine stack exceeds 1000000000-byte limit
		// runtime: sp=0xc020360390 stack=[0xc020360000, 0xc040360000] fatal error: stack overflow

		// TODO:  Investigate why :D
		currentState: mainMenu,
		stateStack:   initStack,

		// List all the IoT Tools usable with the app.
		choices: []string{"Send IoT Jobs", "Dictionary", "Disenroll Inverter", "Upload .json to AWS S3"},
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

	// In which state are we currently in?
	switch m.currentState {

	// In the main menu?
	case mainMenu:
		// Then call updateMainMenu and pass the current tea.Msg to it.
		return m.updateMainMenu(msg)
	case selectIoTJob:
		// Then call updateSelectJob and pass the current tea.Msg to it.
		return m.updateSelectJob(msg)
	}

	return m, nil
}

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

// View looks at the model at its current state, and returns a string, which is the updated UI!
// Redrawing logic and stuff like that is taken care for by BubbleTea.
func (m model) View() string {

	switch m.currentState {
	case mainMenu:
		return m.viewMainMenu()
	case selectIoTJob:
		return m.viewSelectJob()
	}
	return ""
}

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
