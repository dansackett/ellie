# Ellie

A distributed task queue written in Go (like Celery).

This is my first real Go program so I'm using it to learn. Hopefully it turns into something cool.

It requires ZMQ >= 4.2.

## Getting Started

Ellie uses ZeroMQ as a broker handling network sockets. To work with this
project, you will need ZeroMQ with a version >= 4.2. We can get this on Ubuntu
here:

```
$ add-apt-repository ppa:chris-lea/zeromq
$ apt-get update
$ apt-get install zeromq-bin libzmq-dev libzmq0
```

Once we have that, we can grab the project for your own project using `go
get`:

```
$ go get github.com/dansackett/ellie
```

You're good to go!

## Working With This Project

Flesh out how the API works basically.
