package model

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// commands.go collects all asyncronous functions.

func fetchS3FilesCmd() tea.Cmd {
	return func() tea.Msg {
		// TODO: implement s3 logic
		files := []list.Item{
			item{title: "file1.json"},
			item{title: "file2.json"},
		}
		return S3FilesMsg{Files: files}
	}
}
