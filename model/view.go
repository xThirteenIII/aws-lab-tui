// Here are all the view functions for each different state
package model

import (
	"strings"
)

func (m model) viewMainMenu() string {

	return docStyle.Render(m.mainMenuList.View())
}

func (m model) viewS3List() string {

	return docStyle.Render(m.s3FilesList.View())
}

func (m model) viewSelectJob() string {

	// A Builder is used to efficiently build a string using [Builder.Write] methods.
	// It minimizes memory copying. The zero value is ready to use.
	// Do not copy a non-zero Builder.
	var b strings.Builder
	b.WriteString("Type the IoT Job name\n\n")
	b.WriteString(m.jobInput.View())
	b.WriteString("\n\n\n")
	b.WriteString(m.err)
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("esc: go back | "))
	b.WriteString(helpStyle.Render("enter: submit name"))

	// send the UI for rendering
	return docStyle.Render(b.String())
}

func (m model) viewSelectThing() string {

	// A Builder is used to efficiently build a string using [Builder.Write] methods.
	// It minimizes memory copying. The zero value is ready to use.
	// Do not copy a non-zero Builder.
	var b strings.Builder
	b.WriteString("Type the mac address of the Thing\n\n")
	b.WriteString(m.thingInput.View())
	b.WriteString("\n\n\n")
	b.WriteString(errStyle.Render(m.err))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("esc: go back | "))
	b.WriteString(helpStyle.Render("enter: submit mac"))

	// send the UI for rendering
	return docStyle.Render(b.String())
}
