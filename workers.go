package ellie

import (
	"fmt"
	"reflect"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Worker represents a background worker which picks up tasks and communicates
// its progress on its set channels
type Worker struct {
	id       int
	hash     uuid.UUID
	tasks    TaskQueue
	workers  chan *Worker
	complete chan *Task
}

// SpawnWorkers creates the number of workers in the config and starts them as
// goroutines listening for jobs to pick up.
func SpawnWorkers() {
	for i := 0; i < AppConfig.NumWorkers; i++ {
		worker := &Worker{
			id:       i,
			hash:     uuid.NewV4(),
			tasks:    AppConfig.ScheduledTasks,
			workers:  AppConfig.WorkerPool,
			complete: AppConfig.FinishedTasks,
		}

		go worker.Start()
		LogWorkerStarted(worker)
	}
}

// Process takes a task and does the work on it.
func (w *Worker) Process(t *Task) {
	LogTaskStarted(w, t)

	fn := reflect.ValueOf(t.fn)
	fnType := fn.Type()
	if fnType.Kind() != reflect.Func && fnType.NumIn() != len(t.args) {
		panic("Expected a function")
	}

	var args []reflect.Value
	for _, arg := range t.args {
		args = append(args, reflect.ValueOf(arg))
	}

	res := fn.Call(args)
	for _, val := range res {
		fmt.Println("Response:", val.Interface())
	}

	if t.repeat {
		EnqueueEvery(t.interval, t.fn, t.args)
	}

	w.complete <- t
	LogTaskFinished(w, t)
}

// Sleep pauses the worker before its next run
func (w *Worker) Sleep() {
	LogWorkerSleeping(w)
	time.Sleep(time.Duration(AppConfig.WorkInterval) * time.Second)
	w.workers <- w
}

// Start begins a selected worker's scanning loop waiting for tasks to come
// in. When a task comes in, we first check if it is scheduled to be dequeued.
// If so, we don't run it and remove it. If it is ready to be run, it
// processes it. If it isn't ready to be run, it reschedules the task to check
// again. If the worker doesn't find anything within 100 milliseconds, it
// sends the worker into sleep mode for the set interval.
func (w *Worker) Start() {
	for {
		select {
		case <-AppConfig.NewTasks:
			if w.tasks.Len() > 0 {
				task := w.tasks.Pop()
				if ok := AppConfig.CancelledTasks.Get(task.hash); ok {
					LogTaskDequeued(task)
					AppConfig.CancelledTasks.Remove(task.hash)
				} else if time.Since(task.nextRun) > 0 {
					w.Process(task)
				} else {
					w.tasks.Push(task)
				}
			}
		case <-time.After(100 * time.Millisecond):
			w.Sleep()
		}
	}
}
