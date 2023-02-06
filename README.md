# Tasker

**Tasker** - Go package for managing background tasks.

---

## Examples

### Create a new task manager and add tasks

```go
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
```

### Start all tasks

```go
errs := tm.StartTasks(tasker.TaskFilter{
  Name: "*",
 })
```

### Stop specific task

```go
errs := tm.StopTasks(tasker.TaskFilter{Tags: []string{"parallel"}})
```
