# Ellie

A distributed task queue written in Go.

## Installation

Grab the project for your own project using `go get`:

```
$ go get github.com/OrlandoGolang/ellie
```

## Examples

```
package main

import (
	"time"

	"github.com/OrlandoGolang/ellie"
)

func Sum(x, y int) int {
    return x + y
}

func main() {
	// Configure the application to run 10 workers with 5 seconds of sleep between each run.
	ellie.Configure(10, 5)

	// Enqueue a task to run now
	ellie.Enqueue(Sum, 3, 4)

	// Enqueue a task to run in 30 seconds
	ellie.EnqueueIn(30*time.Second, Sum, 3, 4)

	// Enqueue a task to run in 2 minutes
	ellie.EnqueueAt(time.Now().Add(2*time.Minute), Sum, 3, 4)

	// Enqueue a task to run every minute and a half
	ellie.EnqueueEvery((1*time.Minute)+(30*time.Second), Sum, 3, 4)

	// Enqueue a task to run that we intend to cancel
	cancelHash := ellie.EnqueueIn(5*time.Minute, Sum, 3, 4)

	// Dequeue a task from running
	ellie.Dequeue(cancelHash)

	// Start the workers and watch for new tasks
	ellie.RunServer()
}
```
