package model

import (
	"aws-iot-tui/stack"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

// InitialModel creates initial Model for the app
// Necessary fields for the app to start are initialized
func InitialModel() Model {

	initStack := stack.NewStack[StateType]()
	initS3Path := stack.NewStack[string]()
	initStack.Push(StateMainMenu)

	m := Model{
		stateHistory: initStack,
		s3PathStack:  initS3Path,
	}

	// Load env variables
	err := godotenv.Load()
	if err != nil {
		m.lastError = err.Error()
	}

	// Push root s3 path
	m.s3PathStack.Push("")

	// Load suggestions from cache
	m.suggestions.cacheFile = "cache.bin"
	m.suggestions.loadFromCache()

	return m
}

// initMainMenu initializes main menu data
func (m *Model) initMainMenu() {
	items := []list.Item{
		item{title: "IoT Jobs", desc: "Send an IoT Job to a Thing"},
		item{title: "Dictionary", desc: "Download commands from HeidiDB database"},
		item{title: "Disenroll Inverter", desc: "Disenroll an Haier Inverter"},
		item{title: "Upload to S3", desc: "Upload firmware files to S3"},
	}

	delegate := list.NewDefaultDelegate()
	m.mainMenuList = list.New(items, delegate, 0, 0)
	m.mainMenuList.Title = "Select the IoT Tool you want to use"
	// Set lists size
	// TODO: this has to be somewhere else, otherwise when all tools are developed,
	// this will be a long list (no pun intended)
	// WARNING: m.width and m.height MUST BE INITIALIZED before calling this
	h, v := docStyle.GetFrameSize()
	m.mainMenuList.SetSize(m.width-h, m.height-v)
}

// initSelectJob initializes Select Job state data
func (m *Model) initSelectJob() {
	m.jobInput = textinput.New()
	m.jobInput.Cursor.Style = cursorStyle
	m.jobInput.PromptStyle = focusedStyle
	m.jobInput.TextStyle = focusedStyle
	m.jobInput.ShowSuggestions = true
	m.jobInput.Focus()
	m.jobInput.CharLimit = 99
	m.jobInput.SetValue(time.Now().Format("20060201"))
}

// initSelectThing initializes Select Thing state data
func (m *Model) initSelectThing() {
	m.thingInput = textinput.New()
	m.thingInput.Cursor.Style = cursorStyle
	m.thingInput.PromptStyle = focusedStyle
	m.thingInput.TextStyle = focusedStyle
	m.thingInput.ShowSuggestions = true
	m.thingInput.Focus()
	m.thingInput.CharLimit = 17
}

// initS3List inizializes S3 state data
func (m *Model) initS3List() tea.Cmd {
	items := []list.Item{}

	delegate := list.NewDefaultDelegate()
	newSpinner := spinner.New()
	newSpinner.Style = spinnerStyle
	newSpinner.Spinner = spinner.Dot
	m.s3List = list.New(items, delegate, 0, 0)
	m.s3List.Title = "Select S3 Document"
	// TODO: set item help

	// TODO: figure out why spinner is not shown where the list is, but top right of the screen
	h, v := docStyle.GetFrameSize()
	newSpinner.Style.Align(lipgloss.Position(0))
	m.s3List.SetSize(m.width-h, m.height-v)
	m.s3List.SetSpinner(newSpinner.Spinner)

	// To return a tea.Cmd might be useless
	return m.s3List.StartSpinner()
}

// initMainMenu initializes main menu data
func (m *Model) initSelectOp() {
	items := []list.Item{
		item{title: "EPP", desc: "Send an OTA to an ESP32 EPP microcontroller"},
		item{title: "EC3", desc: "Send an OTA to an ESP32 EC3 microcontroller"},
		item{title: "ES3", desc: "Send an OTA to an ESP32 ES3 microcontroller"},
		item{title: "Deep OTA", desc: "Send a Deep OTA"},
		item{title: "OTA EEPROM", desc: "Send an OTA EEPROM"},
		item{title: "CEW OTA", desc: "Send an CEW OTA"},
		item{title: "Other", desc: "Navigate S3 Bucket"},
	}

	delegate := list.NewDefaultDelegate()
	m.operationsList = list.New(items, delegate, 0, 0)
	m.operationsList.Title = "Select the Operation you want to perform"
	// Set lists size
	// TODO: this has to be somewhere else, otherwise when all tools are developed,
	// this will be a long list (no pun intended)
	// WARNING: m.width and m.height MUST BE INITIALIZED before calling this
	h, v := docStyle.GetFrameSize()
	m.operationsList.SetSize(m.width-h, m.height-v)
}
