package tasker

type TaskFilter struct {
	Name   string
	Type   string
	Tags   []string
	Status string
}

// GetTasks returns all tasks that match the filter
func (tm *TaskManager) GetTasks(filter TaskFilter) []*Task {
	// Get all
	if filter.Name == "*" || filter.Type == "*" {
		tasks := make([]*Task, len(tm.tasks))
		i := 0
		for _, task := range tm.tasks {
			tasks[i] = task
			i++
		}
		return tasks
	}

	tasks := make([]*Task, 0)
	for _, task := range tm.tasks {
		if filter.Name != "" && filter.Name != task.Name {
			continue
		}
		if filter.Type != "" && filter.Type != task.Type {
			continue
		}
		if len(filter.Tags) > 0 {
			for _, tag := range filter.Tags {
				if !stringInSlice(tag, task.Tags) {
					continue
				}
			}
		}
		if filter.Status != "" && filter.Status != task.Status {
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks
}
