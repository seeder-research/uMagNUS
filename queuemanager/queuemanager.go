package queuemanager

import (
	"context"
	"log"

	cl "github.com/seeder-research/uMagNUS/cl"
)

var (
	queueManagerContext  context.Context
	queueManagerKillFunc context.CancelFunc
	routineCount         = int(-1)
)

type QueueEvent struct {
	q *cl.CommandQueue
	e *cl.Event
}

func (qe *QueueEvent) GetQueue() *cl.CommandQueue {
	return qe.q
}

func (qe *QueueEvent) GetEvent() []*cl.Event {
	return []*cl.Event{qe.e}
}

func (qe *QueueEvent) SetQueue(queue *cl.CommandQueue) {
	qe.q = queue
}

func (qe *QueueEvent) SetEvent(event *cl.Event) {
	qe.e = event
}

func Init(numRoutines int, in chan *QueueEvent, queuePool chan *cl.CommandQueue) {
	// Initialize globals
	queueManagerContext, queueManagerKillFunc = context.WithCancel(context.Background())
	routineCount = numRoutines

	// Start goroutines
	for i := 0; i < numRoutines; i++ {
		go threadFunction(in, queuePool)
	}

	// Start goroutine for updating events tracked in buffers
	initEventRoutine()
}

func GetContext() context.Context {
	return queueManagerContext
}

func GetKillFunction() context.CancelFunc {
	return queueManagerKillFunc
}

func Teardown() {
	queueManagerKillFunc()
	killEventRoutine()
}

// Function that the goroutines will be running indefinitely unless the context is cancelled...
// The command queue is given to the goroutine, which waits for the queue to finish before
// returning it back to a pool of command queues
func threadFunction(in <-chan *QueueEvent, queuePool chan<- *cl.CommandQueue) {
	for {
		select {
		case item := <-in:
			err := cl.WaitForEvents(item.GetEvent())
			if err != nil {
				log.Printf("ERROR: unable to wait for event to finish...")
			}
			queuePool <- item.GetQueue()
		case <-queueManagerContext.Done():
			return
		}
	}
}
