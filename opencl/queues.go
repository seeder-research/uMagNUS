package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
)

var cmdQueueIdx chan (int)               // token indices for grabbing OpenCL command queue to work use
var cmdQueueArr []*cl.CommandQueue       // Array of OpenCL command queues
var cmdQueueMap map[*cl.CommandQueue]int // Map of index associated with each command queue

func initCmdQueues(context *cl.Context, device *cl.Device) error {
	const cmdQueueCnt = 8 // number of concurrently executing OpenCL command queues available
	var err error

	cmdQueueIdx = make(chan int, cmdQueueCnt)
	cmdQueueArr = make([]*cl.CommandQueue, cmdQueueCnt)
	cmdQueueMap = make(map[*cl.CommandQueue]int)
	for i := 0; i < cmdQueueCnt; i++ {
		cmdQueueIdx <- i
		cmdQueueArr[i], err = context.CreateCommandQueue(device, 0)
		if err != nil {
			fmt.Printf("CreateCommandQueue[%+v] failed: %+v \n", i, err)
			for j := 0; j < i; j++ {
				cmdQueueArr[j].Release()
			}
			return err
		}
		cmdQueueMap[cmdQueueArr[i]] = i
	}
	return err
}

func freeCmdQueues() {
	for _, queue := range cmdQueueArr {
		delete(cmdQueueMap, queue)
		queue.Release()
	}
}

func checkoutQueue() *cl.CommandQueue {
	return cmdQueueArr[<-cmdQueueIdx]
}

func checkinQueue(q *cl.CommandQueue) {
	cmdQueueIdx <-cmdQueueMap[q]
}
