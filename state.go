package ellie

import (
	"log"
	"time"

	"github.com/fatih/color"
)

var (
	// Worker Channels
	WorkerStarted  = make(chan *Worker)
	WorkerSleeping = make(chan *Worker)

	// Task Channels
	TaskScheduled = make(chan *Task)
	TaskDequeued  = make(chan *Task)
	TaskStarted   = make(chan map[*Worker]*Task)
	TaskFinished  = make(chan map[*Worker]*Task)
)

// StateMonitor provides a sane way to listen for state changes in the
// application. New state is passed via channels outputting logs from anywhere
// in the application.
func StateMonitor() {
	for {
		select {
		case worker := <-WorkerStarted:
			color.Set(color.Bold, color.FgBlue)
			log.Println("Started Worker", worker.id)
			color.Unset()
		case worker := <-WorkerSleeping:
			color.Set(color.Faint)
			log.Println("Worker", worker.id, "sleeping for", AppConfig.WorkInterval, "seconds")
			color.Unset()
		case task := <-TaskScheduled:
			color.Set(color.Bold, color.FgYellow)
			log.Println("Task", task.id, "scheduled to run at", task.nextRun.Format(time.UnixDate))
			color.Unset()
		case task := <-TaskDequeued:
			color.Set(color.Bold, color.FgRed)
			log.Println("Task", task.id, "dequeued")
			color.Unset()
		case data := <-TaskStarted:
			color.Set(color.Bold)
			for worker, task := range data {
				log.Println("Worker", worker.id, "picked up task", task.id)
			}
			color.Unset()
		case data := <-TaskFinished:
			color.Set(color.Bold, color.FgGreen)
			for worker, task := range data {
				log.Println("Worker", worker.id, "finished task", task.id)
			}
			color.Unset()
		}
	}
}

// LogWorkerStarted sends a signal to the WorkerStarted channel triggering the
// output text.
func LogWorkerStarted(w *Worker) {
	WorkerStarted <- w
}

// LogWorkerSleeping sends a signal to the WorkerSleeping channel triggering the
// output text.
func LogWorkerSleeping(w *Worker) {
	WorkerSleeping <- w
}

// LogTaskScheduled sends a signal to the TaskScheduled channel triggering the
// output text.
func LogTaskScheduled(t *Task) {
	TaskScheduled <- t
}

// LogTaskDequeued sends a signal to the TaskDequeued channel triggering the
// output text.
func LogTaskDequeued(t *Task) {
	TaskDequeued <- t
}

// LogTaskStarted sends a signal to the TaskStarted channel triggering the
// output text.
func LogTaskStarted(w *Worker, t *Task) {
	TaskStarted <- map[*Worker]*Task{w: t}
}

// LogTaskFinished sends a signal to the TaskFinished channel triggering the
// output text.
func LogTaskFinished(w *Worker, t *Task) {
	TaskFinished <- map[*Worker]*Task{w: t}
}
