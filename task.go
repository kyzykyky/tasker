package tasker

import (
	"context"
)

// type TaskFuncFactory func(svcCtx interface{}) TaskFunc
type TaskFunc func(ctx context.Context)

type Task struct {
	Name string
	Type string
	Tags []string
	Func TaskFunc

	AutoRestart  bool
	MaxRestarts  uint
	restarts     uint
	Status       string
	err          interface{}
	ErrorHandler func(err interface{})

	Paralleled bool
	Workers    uint

	// Cleanups are called when the task panics or before the task is stopped
	Cleanups []func()
	// Add your own context here if required
	Ctx    context.Context
	Cancel context.CancelFunc
}

var (
	StatusAdded   = ""
	StatusStarted = "started"
	StatusStopped = "stopped"
	StatusError   = "error"
	StatusPaused  = "paused"
)
