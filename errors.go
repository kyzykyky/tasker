package tasker

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
	return "task not found"
}
