package ellie

const (
	DEFAULT_WORKER_COUNT  = 5
	DEFAULT_WORK_INTERVAL = 5 // in seconds
)

// appConfig is the configuration object to use within the actual module.
var AppConfig *Config

// Config contains the base configuration for the work queue.
type Config struct {
	// NumWorkers specifies the maximum number of active workers to run at any
	// given time.
	NumWorkers int
	// WorkInterval is the time it takes for a worker to sleep before it
	// checks the task queue for more work to do.
	WorkInterval int
	// ScheduledTasks is the default queue used to decide what is available
	// for the workers to consume.
	ScheduledTasks TaskQueue
	// CancelledTasks is a queue which is checked before a task is executed to
	// see if the task has been cancelled.
	CancelledTasks TaskDequeue
	// NewTasks is a signal channel to express that a new task has been pushed
	// to the ScheduledTasks queue.
	NewTasks chan bool
	// WorkerPool in a channel to wait for a worker when a job comes in and
	// we send workers back into it when they are done.
	WorkerPool chan *Worker
	// FinishedTasks is a channel which cleans up after a task has finished.
	FinishedTasks chan *Task
}

// Configure sets up the base application confiuration options.
func Configure(numWorkers, workInterval int) *Config {
	config := &Config{
		NumWorkers:     numWorkers,
		WorkInterval:   workInterval,
		ScheduledTasks: NewQueue(),
		CancelledTasks: NewDequeue(),
		NewTasks:       make(chan bool, 100),
		FinishedTasks:  make(chan *Task, 100),
		WorkerPool:     make(chan *Worker, numWorkers),
	}
	AppConfig = config
	return config
}

// DefaultConfig uses the defaults to configure the application.
func DefaultConfig() *Config {
	return Configure(DEFAULT_WORKER_COUNT, DEFAULT_WORK_INTERVAL)
}
