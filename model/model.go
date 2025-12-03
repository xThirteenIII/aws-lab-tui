package model

import (
	"aws-iot-tui/stack"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// StateType rappresents the state type
type StateType int

const (
	StateMainMenu StateType = iota
	StateSelectJob
	StateSelectThing
	StateS3List
)

// Model is the main model of the application
type Model struct {
	// State handling
	stateHistory *stack.Stack[StateType]

	// Terminal windows size
	width  int
	height int

	// Shared data
	lastError string

	// MainMenu data
	mainMenuList list.Model

	// IoT Job Tool data
	jobInput    textinput.Model
	thingInput  textinput.Model
	suggestions suggestions

	// S3 state data
	s3List      list.Model
	s3PathStack *stack.Stack[string]
}

// Init returns an initial command that performs some I/O.
// No command is needed for now, so we return nil, which translates to no command.
func (m Model) Init() tea.Cmd {
	return nil
}

// Msg contain data from the result of a IO operation. Msgs trigger the update
// function and, henceforth, the UI.
// type Msg interface{}

// ---------------------------------------------------------------------------
// Model contains the program's state as well as its core functions.          |
//type Model interface {													  |
// Init is the first function that will be called. It returns an optional     |
// initial command. To not perform an initial command return nil.			  |
//	Init() Cmd																  |
//																			  |
// Update is called when a message is received. Use it to inspect messages    |
// and, in response, update the model and/or send a command.                  |
//	Update(Msg) (Model, Cmd)												  |
//																			  |
// View renders the program's UI, which is just a string. The view is         |
// rendered after every Update.												  |
//	View() string															  |
//}																			  |
//----------------------------------------------------------------------------

// Update method is called when "things happen".
// Its job is to look at what's happened and return an updated model in response.
// It can also return a command to execute with it.
//
// The "something happened" comes in the form of a tea.Msg, which can by any type (interface{}).
// Messages are the result of some I/O that took place, such as keypresses, timer tick, or response from a server.
//
// To figure out the type of message, use a type switch or a type assertion.

/*
// The code snippet catching msgs from bubbleTea is:
// eventLoop is the central message loop. It receives and handles the default
// Bubble Tea messages, update the model and triggers redraws.
func (p *Program) eventLoop(model Model, cmds chan Cmd) (Model, error) {
	for {
		select {
		case <-p.ctx.Done():
			return model, nil

		case err := <-p.errs:
			return model, err

		case msg := <-p.msgs:
			// Filter messages.
			if p.filter != nil {
				msg = p.filter(model, msg)
			}
			if msg == nil {
				continue
			}
*/

// Cmd is an IO operation that returns a message when it's complete. If it's
// nil it's considered a no-op. Use it for things like HTTP requests, timers,
// saving and loading from disk, and so on.
//
// Note that there's almost never a reason to use a command to send a message
// to another part of your program. That can almost always be done in the
// update function.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle global messages
	switch msg := msg.(type) {

	// i.e. setting window size and lists size
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Initialize main menu
		m.initMainMenu()

		return m, nil

		// Close app if user presses ctrl+c
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Catch error messages
	case ErrorMsg:
		m.lastError = msg.Err.Error()
		return m, nil
	}

	// Delegate the update to each state
	switch m.getCurrentState() {
	case StateMainMenu:
		return m.updateMainMenu(msg)
	case StateSelectJob:
		return m.updateSelectJob(msg)
	case StateSelectThing:
		return m.updateSelectThing(msg)
	case StateS3List:
		return m.updateS3List(msg)
	}

	return m, nil
}

// View looks at the model at its current state, and returns a string, which is the updated UI!
// Redrawing logic and stuff like that is taken care for by BubbleTea.
func (m Model) View() string {

	switch m.getCurrentState() {
	case StateMainMenu:
		return m.viewMainMenu()
	case StateSelectJob:
		return m.viewSelectJob()
	case StateSelectThing:
		return m.viewSelectThing()
	case StateS3List:
		return m.viewS3List()
	}

	return ""
}

// changeState changes current state and handles states
func (m *Model) changeState(newState StateType) {
	// Push newState state onto stack
	m.stateHistory.Push(newState)

	// Inizialize new state
	switch newState {
	case StateSelectJob:
		m.initSelectJob()
	case StateSelectThing:
		m.initSelectThing()
	case StateS3List:
		m.initS3List()
	}
}

// getCurrentState peeks at the state stack and returns its head (the current state)
func (m *Model) getCurrentState() StateType {
	return m.stateHistory.Peek()
}

// goBack goes back a state in state history
// Pops the head of the stack
func (m *Model) goBack() {
	if !m.stateHistory.IsEmpty() {
		m.stateHistory.Pop()
		m.lastError = "" // Reset error when going back
	}
}
