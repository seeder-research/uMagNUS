package queuemanager

import (
	"context"
	"log"
	"sync"

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

type QueueWaitGroup struct {
	q   *cl.CommandQueue
	wg_ *sync.WaitGroup
}

func NewQueueWaitGroup(cq *cl.CommandQueue, wg *sync.WaitGroup) *QueueWaitGroup {
	tmp := new(QueueWaitGroup)
	tmp.q = cq
	tmp.wg_ = wg
	return tmp
}

func (qw *QueueWaitGroup) SetCommandQueue(cq *cl.CommandQueue) {
	qw.q = cq
}

func (qw *QueueWaitGroup) SetWaitGroup(wg *sync.WaitGroup) {
	qw.wg_ = wg
}

func (qw *QueueWaitGroup) GetCommandQueue() *cl.CommandQueue {
	return qw.q
}

func (qw *QueueWaitGroup) GetWaitGroup() *sync.WaitGroup {
	return qw.wg_
}

func Init(numRoutines int, in chan QueueWaitGroup, queuePool chan *cl.CommandQueue) {
	// Initialize globals
	queueManagerContext, queueManagerKillFunc = context.WithCancel(context.Background())
	routineCount = numRoutines

	// Start goroutines
	for i := 0; i < numRoutines; i++ {
		// go threadFunction(in, queuePool)
		go signalWorkgroupOnQueueFinish(in, queuePool, &ThreadSignals)
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
func checkinQueueOnEventFunction(in <-chan *QueueEvent, queuePool chan<- *cl.CommandQueue) {
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

func CheckoutQueue(queuePool <-chan *cl.CommandQueue, wg_ *sync.WaitGroup) *cl.CommandQueue {
	wg_.Add(1)
	return <-queuePool
}

func signalWorkgroupOnQueueFinish(in <-chan *QueueWaitGroup, queuePool chan<- *cl.CommandQueue, wg_ *sync.WaitGroup) {
	for {
		select {
		case item := <-in:
			q := item.GetQueue()
			err := q.Finish()
			if err != nil {
				log.Printf("ERROR: unable to wait for command queue to finish...")
			}
			queuePool <- q
			if item.GetWaitGroup() != nil {
				item.GetWaitGroup().Done()
			}
		case <-queueManagerContext.Done():
			return
		}
	}
}

func AddToThreadSignals(in int) {
	ThreadSignals.Add(in)
}

func WaitThreadSignals() {
	ThreadSignals.Wait()
}
