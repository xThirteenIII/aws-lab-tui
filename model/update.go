package model

import (
	"net"

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
	m.jobInput.SetSuggestions(m.suggestions.jobSuggestions)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.goBack()
			return m, nil
		case "enter":
			// Salva suggestion e vai allo stato successivo
			m.suggestions.addJobSuggestion(m.jobInput.Value())
			m.changeState(StateSelectThing)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.jobInput, cmd = m.jobInput.Update(msg)
	return m, cmd
}

// updateSelectThing handles SelectThing events
func (m Model) updateSelectThing(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.thingInput.SetSuggestions(m.suggestions.macSuggestions)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.lastError = ""
			m.goBack()
			return m, nil
		case "enter":
			// Validate MAC address
			_, err := net.ParseMAC(m.thingInput.Value())
			if err != nil {
				m.lastError = err.Error()
				return m, nil
			}

			// Save new suggestion in cache
			m.suggestions.addMacSuggestion(m.thingInput.Value())

			// And change to S3 state
			m.changeState(StateSelectOp)

			// return the updated model and a tea.Cmd that fetches S3 Files
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.thingInput, cmd = m.thingInput.Update(msg)
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
		case "enter":
			// TODO: select file
			return m, nil
		}

		// Asyncrously catch S3 files
	case S3FilesMsg:

		// is s3List existing?
		m.s3List.StopSpinner()
		if len(message.Files) == 0 {
			m.s3List.SetItems([]list.Item{item{title: "no item fetched"}})
		} else {
			m.s3List.SetItems(message.Files)
		}
		h, v := docStyle.GetFrameSize()
		m.s3List.SetSize(m.width-h, m.height-v)
	}

	var cmd tea.Cmd
	m.s3List, cmd = m.s3List.Update(msg)
	return m, cmd
}

// updateMainMenu handles MainMenu events
func (m Model) updateSelectionOp(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "esc":
			m.goBack()
			return m, nil
		case "enter":
			opMap := map[int]string{
				0: "Campaign/EPP/",
				1: "Campaign/EC3/",
				2: "Campaign/ES3/",
				3: "Campaign/DeepOTA/",
				4: "Campaign/EEPROM/",
				5: "Campaign/CEW/",
				6: "/",
			}
			m.s3PathStack.Push(opMap[m.operationsList.Cursor()])
			// And change to S3 state
			m.changeState(StateS3List)
			return m, fetchS3FilesCmd(opMap[m.operationsList.Cursor()])
		}
	}

	var cmd tea.Cmd
	m.operationsList, cmd = m.operationsList.Update(msg)
	return m, cmd
}
