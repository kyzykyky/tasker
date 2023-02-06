package tasker_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/kyzykyky/tasker"
)

func TestTasker(t *testing.T) {
	// Create a new task manager
	tm := tasker.TaskManager{}
	// Add a task
	tm.AddTasks(tasker.Task{
		Name: "test1",
		Type: "test",
		Func: func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					log.Println("Task done")
					return
				default:
					log.Println("Task 1 running")
					time.Sleep(1 * time.Second)
				}
			}
		},
	}, tasker.Task{
		Name: "test2",
		Type: "test",
		Tags: []string{"parallel"},
		Func: func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					log.Println("Task 2 done")
					return
				default:
					log.Println("Task 2 running", time.Now().UnixNano())
					time.Sleep(1 * time.Second)
				}
			}
		},
		Paralleled: true,
		Workers:    2,
	})

	// Start all tasks
	errs := tm.StartTasks(tasker.TaskFilter{
		Name: "*",
	})
	if len(errs) > 0 {
		t.Error("Expected no errors, got", errs)
	}

	// Wait for 3 seconds
	time.Sleep(3 * time.Second)
	// Stop parallel task
	errs = tm.StopTasks(tasker.TaskFilter{Tags: []string{"parallel"}})
	if len(errs) > 0 {
		t.Error("Expected no errors, got", errs)
	}

	// Wait for 2 seconds after stopping
	time.Sleep(2 * time.Second)

	stopped := tm.GetTasks(tasker.TaskFilter{Status: tasker.StatusStopped})
	if len(stopped) != 2 {
		t.Error("Expected 2 stopped task, got", len(stopped))
	}

	tm.AddTasks(tasker.Task{
		Name: "test3",
		Type: "test",
		Func: func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					log.Println("Task 3 done")
					return
				default:
					log.Println("Task 3 running")
					time.Sleep(1 * time.Second)
				}
			}
		},
	})

	// Start added task
	tm.StartTasks(tasker.TaskFilter{Name: "test3"})

	// Wait for 3 seconds
	time.Sleep(3 * time.Second)

	time.Sleep(2 * time.Second)
	running := tm.GetTasks(tasker.TaskFilter{Status: tasker.StatusStarted})
	if len(running) != 1 {
		t.Error("Expected 1 running tasks, got", len(running))
	}
}
