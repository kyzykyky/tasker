package tasker

import "context"

func (tm *taskManager) AddTasks(tasks ...Task) []error {
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
			Name:        task.Name,
			Type:        task.Type,
			Tags:        task.Tags,
			Func:        task.Func,
			Paralleled:  task.Paralleled,
			Workers:     task.Workers,
			AutoRestart: task.AutoRestart,
			MaxRestarts: task.MaxRestarts,
		}
		tm.tasks[t.Name] = t
	}
	return errs
}

func (tm *taskManager) StartTasks(filter TaskFilter) []error {
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

func (tm *taskManager) StopTasks(filter TaskFilter) []error {
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

// TODO: Implement restart -- Test
func (tm *taskManager) RestartTasks(filter TaskFilter) {
	tasks := tm.GetTasks(filter)
	for i := range tasks {
		if tasks[i].Func == nil || tasks[i].Ctx == nil || tasks[i].Cancel == nil {
			continue
		}
		if tasks[i].Cancel != nil {
			tasks[i].stop()
		}
		tasks[i].Ctx, tasks[i].Cancel = context.WithCancel(context.Background())
		tasks[i].restarts++
		go tasks[i].runSafe()
	}
}
