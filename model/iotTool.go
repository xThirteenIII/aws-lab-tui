package model

import (
	"aws-iot-tui/stack"

	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

type IoTToolData struct {
	jobInput   textinput.Model
	thingInput textinput.Model

	JobName     string
	ThingName   string
	DocumentKey string // Full S3 path to the document to send

	JobExpiresInSeconds int64

	suggestions suggestions

	s3List      list.Model
	s3PathStack *stack.Stack[string]

	jobStack *stack.Stack[iot.DescribeJobExecutionInput]

	IsJobRunning bool
	JobStatus    string
}

func (iot *IoTToolData) Reset() {
	iot.JobName = ""
	iot.ThingName = ""
	iot.DocumentKey = ""
	iot.IsJobRunning = false
	iot.JobStatus = ""
}

func (iot *IoTToolData) IsValid() bool {
	return iot.JobName != "" &&
		iot.ThingName != "" &&
		iot.DocumentKey != ""
}
