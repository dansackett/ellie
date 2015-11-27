package ellie

// RunServer starts a blocking loop allowing the goroutines to communicate
// without the program closing. We spawn the workers here and also fire off
// the StateMonitor to listen for state changes while processing.
func RunServer() {
	StateMonitor()
	SpawnWorkers()

	for {
		select {
		case <-FinishedQueue:
		case <-WorkerQueue:
		}
	}
}
