package tasker

import "fmt"

var ()

type ErrTaskInvalid struct {
	Task   Task
	Reason string
}

func (e ErrTaskInvalid) Error() string {
	return "task is invalid: " + e.Reason
}

type ErrTaskNotFound struct {
	Filter TaskFilter
}

func (e ErrTaskNotFound) Error() string {
	return fmt.Sprintf("task not found: %#v", e.Filter)
}
