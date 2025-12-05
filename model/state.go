package model

// StateType rappresents the state type
type StateType int

const (
	StateMainMenu StateType = iota
	StateSelectJob
	StateSelectThing
	StateSelectOp
	StateS3List
)
