// package model defines the Model struct, which defines application states.
// It implements three methods which operate on the model structure.
//
// Init: a function that returns an initial command for the application to run
// Update: a function that updates incoming events and updates the model accordingly
// View: a function that renders UI based on the data in the model

package model

import (
	"aws-iot-tui/stack"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Order struct from biggest Bytes to lowest, to minimize padding and optimize memory.
type model struct {
	// list represents a list of items to show.
	mainMenuList list.Model // A lot of Bytes
	jobInput     textinput.Model
	thingInput   textinput.Model

	suggestions suggestions
	// stateStack is a Stack structure that holds states history
	// It allows to go back and forth in the menu.
	stateStack *stack.Stack[state] // states history, pointer: 8B

	currentState state // int: 8B

	cursor int // which tool item the cursor is pointing at, int: 8B
}

// initialModel defines the initial state of the application.
// Defining a function to return the initial model, but could use a variable elsewhere, too.
func InitialModel() model {

	// initStack creates a new empty stack
	initStack := stack.NewStack[state]()

	// push mainMenu as default first state onto the Stack
	initStack.Push(mainMenu)

	initModel := model{

		// WARNING: to set the currentState is crucial, otherwise program will create unlimited goRoutines!
		// runtime: goroutine stack exceeds 1000000000-byte limit
		// runtime: sp=0xc020360390 stack=[0xc020360000, 0xc040360000] fatal error: stack overflow

		// TODO:  Investigate why :D
		currentState: mainMenu,
		stateStack:   initStack,

		mainMenuList: list.New(
			[]list.Item{
				item{title: "Iot Jobs", desc: "Send an IoT Job to a Thing"},
				item{title: "Dictionary", desc: "Download commands from HeidiDB database, in js format"},
				item{title: "Disenroll Inverter", desc: "Disenroll an Haier Inverter from database and dynamoDB"},
				item{title: "Upload to S3", desc: "Select a new firmware version and upload its .json files to S3"},
			},
			list.NewDefaultDelegate(), 0, 0),
	}
	return initModel
}

// Init returns an initial command that performs some I/O.
// No command is needed for now, so we return nil, which translates to no command.
// TODO: make textinput.Blink work as soon as we're in the selectJobState
func (m model) Init() tea.Cmd {
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
	case selectThing:
		return m.updateSelectThing(msg)
	}

	return m, nil
}

// View looks at the model at its current state, and returns a string, which is the updated UI!
// Redrawing logic and stuff like that is taken care for by BubbleTea.
func (m model) View() string {

	// Update view based on current state of the app
	switch m.currentState {
	case mainMenu:
		return m.viewMainMenu()
	case selectIoTJob:
		return m.viewSelectJob()
	case selectThing:
		return m.viewSelectThing()
	}
	return ""
}
