package logging

type DbAction string

const (
	ActionInsert DbAction = "insert"
	ActionUpdate DbAction = "update"
)
