package model

import "strings"

func (m Model) viewMainMenu() string {
	return docStyle.Render(m.mainMenuList.View())
}

func (m Model) viewSelectJob() string {
	// A Builder is used to efficiently build a string using [Builder.Write] methods.
	// It minimizes memory copying. The zero value is ready to use.
	// Do not copy a non-zero Builder.
	var b strings.Builder
	b.WriteString("Type the IoT Job name\n\n")
	b.WriteString(m.iotTool.jobInput.View())
	b.WriteString("\n\n")
	b.WriteString(helpCmdStyle.Render("esc ") +
		helpStyle.Render("back • ") +
		helpCmdStyle.Render("enter ") +
		helpStyle.Render("submit"))
	return docStyle.Render(b.String())
}

func (m Model) viewSelectThing() string {
	var b strings.Builder
	b.WriteString("Type the MAC address of the Thing\n\n")
	b.WriteString(m.iotTool.thingInput.View())
	b.WriteString("\n\n\n")

	if m.lastError != "" {
		b.WriteString(errStyle.Render(m.lastError))
		b.WriteString("\n\n")
	}

	b.WriteString(helpCmdStyle.Render("esc ") +
		helpStyle.Render("back • ") +
		helpCmdStyle.Render("enter ") +
		helpStyle.Render("submit"))
	return docStyle.Render(b.String())
}

func (m Model) viewOpList() string {
	var b strings.Builder
	if m.lastError != "" {
		b.WriteString(errStyle.Render(m.lastError))
		b.WriteString("\n\n")
	}
	return docStyle.Render(m.operationsList.View())
}

func (m Model) viewS3List() string {
	var b strings.Builder
	if m.lastError != "" {
		b.WriteString(errStyle.Render(m.lastError))
		b.WriteString("\n\n")
	}
	return docStyle.Render(m.iotTool.s3List.View() + b.String())
}

func (m Model) viewError() string {
	return errStyle.Render(m.lastError)
}

func (m Model) viewSendJob() string {
	return docStyle.Render("sending job")
}
