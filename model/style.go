package model

import "github.com/charmbracelet/lipgloss"

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// item implements list.Item struct, implementing FilterValue() method
type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
