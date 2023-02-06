package tasker

import "log"

func (task *Task) run() {
	if task.Paralleled && task.Workers > 1 && task.Func != nil {
		// If paralleled is true and workers is greater than 1, then run the task
		task.Status = StatusStarted
		for i := uint(0); i < task.Workers; i++ {
			go task.Func(task.Ctx)
		}
	} else if task.Paralleled && task.Workers == 1 && task.Func != nil {
		// If paralleled is true but workers is 1, then run the task as a normal
		task.Paralleled = false
		task.Status = StatusStarted
		task.Func(task.Ctx)
	} else if task.Func != nil {
		// If function is not nil, then run the task
		task.Status = StatusStarted
		task.Func(task.Ctx)
	} else {
		// If function is nil, then set the status to error
		task.Status = StatusError
		task.err = "Task function is nil"
	}
}

// Run with recover
func (task *Task) runSafe() {
	defer task.recoverTask()
	task.run()
}

// Recover from panic of goroutine
func (task *Task) recoverTask() {
	if p := recover(); p != nil {
		for _, cleanup := range task.Cleanups {
			cleanup()
		}
		task.Status = StatusError
		task.err = p
		log.Println("Recovered from panic:", p)
		// TODO: Implement auto restart
		// if task.AutoRestart {
		// }
	}
}

// Stop the task with cleanups if exist
func (task *Task) stop() error {
	if task.Status == StatusStopped {
		return nil
	} else if task.Status == StatusError {
		return ErrTaskInvalid{Task: *task, Reason: "Task is in error state"}
	} else if task.Status == StatusPaused {
		return ErrTaskInvalid{Task: *task, Reason: "Task is in paused state"}
	} else if task.Status == StatusStarted {

		for _, cleanup := range task.Cleanups {
			cleanup()
		}
		task.Status = StatusStopped
		task.Cancel()
		return nil
	}
	return ErrTaskInvalid{Task: *task, Reason: "Task is in unknown state"}
}
