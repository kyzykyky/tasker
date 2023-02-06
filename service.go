package tasker

import (
	"context"
)

// Task manager context
type TaskManager struct {
	tasks map[string]*Task
}

func (tm *TaskManager) AddTasks(tasks ...Task) []error {
	errs := make([]error, 0)

	if tm.tasks == nil {
		tm.tasks = make(map[string]*Task, len(tasks))
	} else if len(tm.tasks) > 0 {
		for _, task := range tasks {
			if _, ok := tm.tasks[task.Name]; ok {
				errs = append(errs, ErrTaskInvalid{Task: task, Reason: "Task already exists"})
			}
		}
	}
	for _, task := range tasks {
		t := &Task{
			Name:       task.Name,
			Type:       task.Type,
			Tags:       task.Tags,
			Func:       task.Func,
			Paralleled: task.Paralleled,
			Workers:    task.Workers,
		}
		tm.tasks[t.Name] = t
	}
	return errs
}

func (tm *TaskManager) StartTasks(filter TaskFilter) []error {
	errs := make([]error, 0)
	tasks := tm.GetTasks(filter)
	if len(tasks) == 0 {
		errs = append(errs, ErrTaskNotFound{Filter: filter})
		return errs
	}

	for i := range tasks {
		if tasks[i].Func == nil {
			errs = append(errs, ErrTaskInvalid{Task: *tasks[i], Reason: "Missing function, context or stop function"})
		}
		tasks[i].Ctx, tasks[i].Cancel = context.WithCancel(context.Background())
		go tasks[i].runSafe()
	}
	return errs
}

func (tm *TaskManager) StopTasks(filter TaskFilter) []error {
	errs := make([]error, 0)
	tasks := tm.GetTasks(filter)
	for i := range tasks {
		if tasks[i].Cancel != nil {
			err := tasks[i].stop()
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

// TODO: Implement restart
// func (tm *TaskManager) RestartTasks(filter TaskFilter) {
// 	tasks := tm.GetTasks(filter)
// 	for i := range tasks {
// 		if tasks[i].Func == nil || tasks[i].Ctx == nil || tasks[i].Cancel == nil {
// 			continue
// 		}
// 		if tasks[i].Cancel != nil {
// 			tasks[i].Stop()
// 		}
// 		tasks[i].Ctx, tasks[i].Cancel = context.WithCancel(context.Background())
// 		tasks[i].restarts++
// 		go tasks[i].RunSafe()
// 	}
// }
