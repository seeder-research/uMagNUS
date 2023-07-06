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
	queueManagerContext, queueManagerKillFunc = context.WithCancel(context.Background())
	routineCount = numRoutines
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
