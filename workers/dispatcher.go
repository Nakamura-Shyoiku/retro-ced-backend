package workers

import (
	"github.com/apex/log"
)

// WorkerQueue queue of work requests
var WorkerQueue chan chan *MessageRequest

// StartDispatcher starts the work dispatcher
func StartDispatcher(nworkers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan *MessageRequest, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		log.Infof("Starting worker %d", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				log.Info("received work requeust")
				go func() {
					worker := <-WorkerQueue

					log.Info("dispatching work request")
					worker <- work
				}()
			}
		}
	}()
}
