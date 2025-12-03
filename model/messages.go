package model

import "github.com/charmbracelet/bubbles/list"

// ErrorMsg represents an error
type ErrorMsg struct {
	Err error
}

// S3FilesMsg holds S3 files loaded
type S3FilesMsg struct {
	Files []list.Item
}
