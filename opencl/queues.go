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
