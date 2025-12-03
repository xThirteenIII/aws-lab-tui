package main

import (
	"aws-iot-tui/refactor"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	program := tea.NewProgram(refactor.InitialModel(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
