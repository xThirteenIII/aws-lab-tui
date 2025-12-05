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
	b.WriteString(m.jobInput.View())
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("esc → go back | enter → submit"))
	return docStyle.Render(b.String())
}

func (m Model) viewSelectThing() string {
	var b strings.Builder
	b.WriteString("Type the MAC address of the Thing\n\n")
	b.WriteString(m.thingInput.View())
	b.WriteString("\n\n\n")

	if m.lastError != "" {
		b.WriteString(errStyle.Render(m.lastError))
		b.WriteString("\n\n")
	}

	b.WriteString(helpStyle.Render("esc → go back | enter → submit"))
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
	m.s3List.FullHelp()
	return docStyle.Render(m.s3List.View() + b.String())
}

func (m Model) viewError() string {
	return errStyle.Render(m.lastError)
}
