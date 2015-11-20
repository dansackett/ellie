package main

import (
	"time"

	ellie "./ellie"
)

func main() {
	// Enqueue a task to run now
	ellie.Enqueue(func(x, y int) int { return x * y }, 3, 4)

	// Enqueue a task to run in 30 seconds
	ellie.EnqueueIn(30*time.Second, func(x, y int) int { return x + y }, 3, 4)

	// Enqueue a task to run in 2 minutes
	ellie.EnqueueAt(time.Now().Add(2*time.Minute), func(x, y int) int { return x - y }, 3, 4)

	// Enqueue a task to run every minute and a half
	ellie.EnqueueEvery((1*time.Minute)+(30*time.Second), func(x, y int) int { return x % y }, 3, 4)

	// Enqueue a task to run that we intend to cancel
	cancelHash := ellie.EnqueueIn(5*time.Minute, func(x, y int) int { return x + y }, 3, 4)

	// Dequeue a task from running
	ellie.Dequeue(cancelHash)

	// Start the workers and watch for new tasks
	ellie.RunServer()
}
