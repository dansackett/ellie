package ellie

import (
	"container/list"
	"errors"
	"math/rand"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

// TaskDequeue is a threadsafe container for tasks to be dequeued
type TaskDequeue struct {
	Tasks map[uuid.UUID]bool
	Lock  *sync.Mutex
}

// NewDequeue returns a new instance of a TaskQueue
func NewDequeue() TaskDequeue {
	return TaskDequeue{
		Tasks: make(map[uuid.UUID]bool),
		Lock:  &sync.Mutex{},
	}
}

// Push adds a new task into the front of the TaskQueue
func (q *TaskDequeue) Push(hash uuid.UUID) (uuid.UUID, error) {
	q.Lock.Lock()
	defer q.Lock.Unlock()

	if _, ok := q.Tasks[hash]; !ok {
		q.Tasks[hash] = true
		return hash, nil
	}

	// @TODO use proper error
	return hash, errors.New("Task is already scheduled to be dequeued")
}

// Get checks if a key exists in our dequeued task list
func (q *TaskDequeue) Get(hash uuid.UUID) bool {
	q.Lock.Lock()
	defer q.Lock.Unlock()

	if _, ok := q.Tasks[hash]; ok {
		return true
	}
	return false
}

// Remove deletes the dequeued entry once we are done with it
func (q *TaskDequeue) Remove(hash uuid.UUID) {
	q.Lock.Lock()
	defer q.Lock.Unlock()

	delete(q.Tasks, hash)
}

// TaskQueue is a threadsafe container for tasks to be processed
type TaskQueue struct {
	Tasks *list.List
	Lock  *sync.Mutex
}

// NewQueue returns a new instance of a TaskQueue
func NewQueue() TaskQueue {
	return TaskQueue{
		Tasks: list.New(),
		Lock:  &sync.Mutex{},
	}
}

// Push adds a new task into the front of the TaskQueue
func (q *TaskQueue) Len() int {
	q.Lock.Lock()
	defer q.Lock.Unlock()

	return q.Tasks.Len()
}

// Push adds a new task into the front of the TaskQueue
func (q *TaskQueue) Push(t *Task) {
	q.Lock.Lock()
	defer q.Lock.Unlock()

	q.Tasks.PushFront(t)

	NewTasksQueue <- true
}

// Pop grabs the last task from the TaskQueue
func (q *TaskQueue) Pop() *Task {
	q.Lock.Lock()
	defer q.Lock.Unlock()

	task := q.Tasks.Remove(q.Tasks.Back())
	return task.(*Task)
}

// Task represents a task to run. It can be scheduled to run later or right away.
type Task struct {
	id       int
	hash     uuid.UUID
	nextRun  time.Time
	interval time.Duration
	repeat   bool
	fn       interface{}
	args     []interface{}
}

// newTask is an internal function to create a basic new task object.
func newTask(fn interface{}, args ...interface{}) *Task {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	task := &Task{
		id:     random.Intn(10000),
		hash:   uuid.NewV4(),
		repeat: false,
		fn:     fn,
		args:   args,
	}

	return task
}

// enqueue is an internal function used to asynchronously push a task onto the
// queue and log the state to the terminal.
func enqueue(task *Task) {
	Tasks.Push(task)
	LogTaskScheduled(task)
}

// Enqueue schedules a task to run as soon as the next worker is available.
func Dequeue(hash uuid.UUID) {
	go func() {
		if _, err := TasksDequeue.Push(hash); err != nil {
			// @TODO handle this properly
			panic(err)
		}
	}()
}

// Enqueue schedules a task to run as soon as the next worker is available.
func Enqueue(fn interface{}, args ...interface{}) uuid.UUID {
	task := newTask(fn, args...)
	task.nextRun = time.Now()
	go enqueue(task)
	return task.hash
}

// EnqueueIn schedules a task to run a certain amount of time from the current
// time. This allows us to schedule tasks to run in intervals.
func EnqueueIn(period time.Duration, fn interface{}, args ...interface{}) uuid.UUID {
	task := newTask(fn, args...)
	task.nextRun = time.Now().Add(period)
	go enqueue(task)
	return task.hash
}

// EnqueueAt schedules a task to run at a certain time in the future.
func EnqueueAt(period time.Time, fn interface{}, args ...interface{}) uuid.UUID {
	task := newTask(fn, args...)
	task.nextRun = period
	go enqueue(task)
	return task.hash
}

// EnqueueEvery schedules a task to run and reschedule itself on a regular
// interval. It works like EnqueueIn but repeats
func EnqueueEvery(period time.Duration, fn interface{}, args ...interface{}) uuid.UUID {
	task := newTask(fn, args...)
	task.nextRun = time.Now().Add(period)
	task.interval = period
	task.repeat = true
	go enqueue(task)
	return task.hash
}
