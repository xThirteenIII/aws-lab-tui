package model

import (
	"aws-iot-tui/stack"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
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
	}
	m.iotTool.s3PathStack = initS3Path

	// Load env variables
	err := godotenv.Load()
	if err != nil {
		m.lastError = err.Error()
	}

	// Push root s3 path
	m.iotTool.s3PathStack.Push("")

	// Load suggestions from cache
	m.iotTool.suggestions.cacheFile = "cache.bin"
	m.iotTool.suggestions.loadFromCache()

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
	m.iotTool.jobInput = textinput.New()
	m.iotTool.jobInput.Cursor.Style = cursorStyle
	m.iotTool.jobInput.PromptStyle = focusedStyle
	m.iotTool.jobInput.TextStyle = focusedStyle
	m.iotTool.jobInput.ShowSuggestions = true
	m.iotTool.jobInput.Focus()
	m.iotTool.jobInput.CharLimit = 99
	m.iotTool.jobInput.SetValue(time.Now().Format("20060201"))
}

// initSelectThing initializes Select Thing state data
func (m *Model) initSelectThing() {
	m.iotTool.thingInput = textinput.New()
	m.iotTool.thingInput.Cursor.Style = cursorStyle
	m.iotTool.thingInput.PromptStyle = focusedStyle
	m.iotTool.thingInput.TextStyle = focusedStyle
	m.iotTool.thingInput.ShowSuggestions = true
	m.iotTool.thingInput.Focus()
	m.iotTool.thingInput.CharLimit = 17
}

// initS3List inizializes S3 state data
func (m *Model) initS3List() {
	items := []list.Item{item{title: "Fetching items..."}}

	delegate := list.NewDefaultDelegate()
	m.iotTool.s3List = list.New(items, delegate, 0, 0)
	m.iotTool.s3List.Title = "Select S3 Document"
	// TODO: set item help

	// TODO: figure out why spinner is not shown where the list is, but top right of the screen
	h, v := docStyle.GetFrameSize()
	m.iotTool.s3List.SetSize(m.width-h, m.height-v)

}

// initMainMenu initializes main menu data
func (m *Model) initSelectOp() {
	// reset stack path
	m.iotTool.s3PathStack = stack.NewStack[string]()
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

func (m *Model) initSendJob() {
	// WARNING: stack initialization must be done once only
	m.iotTool.jobStack = stack.NewStack[iot.DescribeJobExecutionInput]()
}
