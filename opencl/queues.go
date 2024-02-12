package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
)

func SyncQueues(dst, src []*cl.CommandQueue) {
	// Insert barriers in src queues, and then insert barriers with wait lists in dst queues
	var barrierWaitList []*cl.Event
	for _, q := range src {
		event, err := q.EnqueueBarrierWithWaitList(nil)
		if err != nil {
			fmt.Println("failed to enqueue barrier in src list: %+v ", err)
		}
		barrierWaitList = append(barrierWaitList, event)
	}
	for _, q := range dst {
		_, err := q.EnqueueBarrierWithWaitList(barrierWaitList)
		if err != nil {
			fmt.Println("failed to enqueue barrier in dst list: %+v ", err)
		}
	}
}

func WaitAllQueuesToFinish() error {
	var err error
	var errOut error
	errOut = nil
	for _, q := range ClCmdQueue {
		if err = q.Finish(); err != nil {
			fmt.Printf("failed to wait for queue to finish: %+v \n", err)
			errOut = err
		}
	}
	return errOut
}

func CheckoutQueue() int {
	q := <-QueueChannels
	return q
}

func CheckinQueue(idx int) {
	QueueChannels <- idx
}

func ExpandQueueList() {
	q, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("CreateCommandQueue failed: %+v \n", err)
		return
	}
	ClCmdQueue = append(ClCmdQueue, q)
	InitQueueChannels(NumQueues)
	NumQueues++
}

func InitQueueChannels(n int) {
	newChan := make(chan int, n-1)
	for i := 1; i < n; i++ {
		newChan <- i
	}
	QueueChannels = newChan
}
