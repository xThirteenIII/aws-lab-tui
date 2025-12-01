package model

type state int

const (
	mainMenu state = iota + 1
	selectIoTJob
	selectThing
	selectS3File
)
