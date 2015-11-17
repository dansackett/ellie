package ellie

var (
	numWorkers    = 5
	workInterval  = 5
	Tasks         = NewQueue()
	TasksDequeue  = NewDequeue()
	NewTasksQueue = make(chan bool, 100)
	WorkerQueue   = make(chan *Worker, numWorkers)
	FinishedQueue = make(chan *Task, 100)
)
