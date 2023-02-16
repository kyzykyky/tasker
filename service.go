package tasker

type TaskManagerConf struct {
	EnablePanic  bool
	EnableLogger bool
	ErrorHandler func(err error)
}

// Task manager context
type taskManager struct {
	tasks map[string]*Task
}

func NewTaskManager(conf TaskManagerConf) *taskManager {
	return &taskManager{
		tasks: make(map[string]*Task),
	}
}
