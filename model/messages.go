package model

import (
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/charmbracelet/bubbles/list"
)

// ErrorMsg represents an error
type ErrorMsg struct {
	Err error
}

// S3FilesMsg holds S3 files loaded
type S3FilesMsg struct {
	Files []list.Item
}

type IoTJobMsg struct {
	JobExeInput *iot.DescribeJobExecutionInput
}
