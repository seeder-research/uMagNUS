package queuemanager

import (
	"context"
	"log"

	cl "github.com/seeder-research/uMagNUS/cl"
)

var (
	queueManagerContext	context.Context
	queueManagerKillFunc	context.CancelFunc
	routineCount		= int(-1)
)

func Init(numRoutines int, queueInput, queuePool chan *cl.CommandQueue) {
	// Initialize globals
	queueManagerContext, queueManagerKillFunc = context.WithCancel(context.Background())
	routineCount = numRoutines

	// Start goroutines
	for i := 0; i < numRoutines; i++ {
		go threadFunction(queueInput, queuePool)
	}
}

func GetContext() context.Context {
	return queueManagerContext
}

func GetKillFunction() context.CancelFunc {
	return queueManagerKillFunc
}

func Teardown() {
	queueManagerKillFunc()
}

// Function that the goroutines will be running indefinitely unless the context is cancelled...
// The command queue is given to the goroutine, which waits for the queue to finish before
// returning it back to a pool of command queues
func threadFunction(queueIn <-chan *cl.CommandQueue, queuePool chan<- *cl.CommandQueue) {
	for {
		select {
		case cmdQueue := <- queueIn:
			err := cmdQueue.Finish()
			if err != nil {
				log.Printf("ERROR: unable to wait for command queue to finish...")
			}
			queuePool <-cmdQueue
		case <-queueManagerContext.Done():
			return
		}
	}
}
