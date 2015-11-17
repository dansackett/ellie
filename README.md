# Ellie

A distributed task queue written in Go.

This is my first real Go program so I'm using it to learn. Hopefully it turns into something cool.

## Getting Started

Grab the project for your own project using `go get`:

```
$ go get github.com/dansackett/ellie
```

You're good to go!

# TODO

- [X] Use a queue to store tasks so task pulling can be done in order
- [X] Add functions to enqueue a task at different times and durations
- [X] Add function to dequeue a task from running
- [ ] Add ability to set config options
- [ ] Add proper error handling
- [ ] Use base config as place to hold communicate channels and other vars
- [ ] Add ability to set a function for the work to be done
- [ ] Add ability to expect any input for args
- [ ] Add Redis as a backend to store tasks
