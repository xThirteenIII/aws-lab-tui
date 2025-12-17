package model

import (
	"net"
	"path"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// updateMainMenu handles MainMenu events
func (m Model) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			// Change state based on user selection
			switch m.mainMenuList.Cursor() {
			case 0: // IoT Jobs
				m.changeState(StateSelectJob)
			case 1: // Dictionary
			case 2: // Disenroll
			case 3: // Upload to S3
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.mainMenuList, cmd = m.mainMenuList.Update(msg)
	return m, cmd
}

// updateSelectJob handles Select Job events
func (m Model) updateSelectJob(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.iotTool.jobInput.SetSuggestions(m.iotTool.suggestions.jobSuggestions)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.goBack()
			return m, nil
		case "enter":
			// Salva suggestion e vai allo stato successivo
			m.iotTool.suggestions.addJobSuggestion(m.iotTool.jobInput.Value())
			m.changeState(StateSelectThing)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.iotTool.jobInput, cmd = m.iotTool.jobInput.Update(msg)
	return m, cmd
}

// updateSelectThing handles SelectThing events
func (m Model) updateSelectThing(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.iotTool.thingInput.SetSuggestions(m.iotTool.suggestions.macSuggestions)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.lastError = ""
			m.goBack()
			return m, nil
		case "enter":
			// Validate MAC address
			_, err := net.ParseMAC(m.iotTool.thingInput.Value())
			if err != nil {
				m.lastError = err.Error()
				return m, nil
			}

			// Save new suggestion in cache
			m.iotTool.suggestions.addMacSuggestion(m.iotTool.thingInput.Value())

			// And change to S3 state
			m.changeState(StateSelectOp)

			// return the updated model and a tea.Cmd that fetches S3 Files
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.iotTool.thingInput, cmd = m.iotTool.thingInput.Update(msg)
	return m, cmd
}

// updateMainMenu handles MainMenu events
func (m Model) updateSelectionOp(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch message := msg.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "q":
			return m, tea.Quit
		case "esc":
			m.goBack()
			return m, nil
		case "enter":

			// Submit user selection only if user is not filtering values.
			// If the user is filterin, "enter" is used to accept the filter value.
			if !m.operationsList.SettingFilter() {
				opMap := map[int]string{
					0: "Campaign/EPP/",
					1: "Campaign/EC3/",
					2: "Campaign/ES3/",
					3: "Campaign/DeepOTA/",
					4: "Campaign/EEPROM/",
					5: "Campaign/CEW/",
					6: "",
				}
				m.iotTool.s3PathStack.Push(opMap[m.operationsList.Cursor()])
				// And change to S3 state, this inits the s3 List
				m.changeState(StateS3List)
				return m, fetchS3FilesCmd(opMap[m.operationsList.Cursor()])
			}
		}
	}

	var cmd tea.Cmd
	m.operationsList, cmd = m.operationsList.Update(msg)
	return m, cmd
}

// updateS3List handles S3 events
func (m Model) updateS3List(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch message := msg.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "esc":
			m.goBack()
			return m, nil
		case "b":
			m.iotTool.s3PathStack.Pop()
			// Fetch files into that folder
			return m, fetchS3FilesCmd(m.iotTool.s3PathStack.Peek())
		case "enter":
			// TODO: select file
			// Check if selected item is a json file
			// Selected item is the last element in the Stack
			if !m.iotTool.s3List.SettingFilter() {
				if strings.HasSuffix(m.iotTool.s3List.SelectedItem().(item).title, ".json") {
					m.changeState(StateSendJob)
					// send job here
					return m, nil
				} else {
					// Push completePath+selected folder into the path stack
					lastEl := m.iotTool.s3PathStack.Peek()
					// WARNING: if paginator takes more than 1 page, Cursor() value is
					// always 0 < Cursor.value < len(page)
					// use SelectedItem() to return the item selected
					m.iotTool.s3PathStack.Push(lastEl + m.iotTool.s3List.SelectedItem().(item).title + "/")
					// Fetch files into that folder
					return m, fetchS3FilesCmd(m.iotTool.s3PathStack.Peek())
				}
			}
		}

		// Asyncrously catch S3 files
	case S3FilesMsg:

		// Reset List of Items everytime something is fetched from S3
		// This is necessary to reset cursor, pagination and everything related to it
		m.initS3List()

		// Are there no files?
		if len(message.Files) == 0 {
			// Then show it
			m.iotTool.s3List.SetItems([]list.Item{item{title: "no item fetched"}})
		} else {

			// Show just base path in the list, for cleaner UI
			basePaths := []string{}
			for _, file := range message.Files {
				basePaths = append(basePaths, path.Base(file.(item).title))
			}
			m.iotTool.s3List.SetItems(message.Files)
		}
		h, v := docStyle.GetFrameSize()
		m.iotTool.s3List.SetSize(m.width-h, m.height-v)
	}

	var cmd tea.Cmd
	m.iotTool.s3List, cmd = m.iotTool.s3List.Update(msg)
	return m, cmd
}

func (m *Model) updateSendJob(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch message := msg.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "ctrl+c":
			return m, tea.Quit
			// cancel job
		case "esc":
			m.goBack()
			return m, nil
		case "enter":
			sendIoTJob(m)

		}
	case IoTJobMsg:

		// Add current job to the job stack.
		// More accurately, job description input object
		m.iotTool.jobStack.Push(*message.JobExeInput)

	}

	return m, nil
}
