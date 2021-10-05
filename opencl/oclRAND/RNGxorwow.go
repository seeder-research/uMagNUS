package oclRAND

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/seeder-research/uMagNUS/opencl/cl"
	"math/rand"
)

func (p *XORWOW_status_array_ptr) Init(seed uint32, events []*cl.Event) {
	rand.Seed((int64)(seed))
	totalCount := p.GetStatusSize()
	seed_arr := make([]uint32, totalCount)
	for idx := 0; idx < totalCount; idx++ {
		tmpNum := rand.Uint32()
		for tmpNum == 0 {
			tmpNum = rand.Uint32()
		}
		seed_arr[idx] = tmpNum
	}
	context := p.GetContext()
	seed_buf, err := context.CreateBufferUnsafe(cl.MemReadWrite, int(unsafe.Sizeof(seed))*p.Status_size, nil)
	if err != nil {
		log.Fatalln("Unable to create buffer for XORWOW seed array!")
	}
	var seed_event *cl.Event
	seed_event, err = ClCmdQueue.EnqueueWriteBuffer(seed_buf, true, 0, totalCount, unsafe.Pointer(&seed_arr[0]), nil)
	if err != nil {
		log.Fatalln("Unable to write seed buffer to device: ", err)
	}
	err = cl.WaitForEvents(events)
	if err != nil {
		fmt.Printf("First WaitForEvents failed in InitRNG: %+v \n", err)
	}
	event := k_xorwow_seed_async(unsafe.Pointer(p.Status_buf), unsafe.Pointer(seed_buf), &config{[]int{p.GetGroupCount() * p.GetGroupSize()}, []int{p.GetGroupSize()}}, []*cl.Event{seed_event})

	p.Ini = true
	err = cl.WaitForEvents([]*cl.Event{event})
	if err != nil {
		fmt.Printf("Second WaitForEvents failed in InitRNG: %+v \n", err)
	}

}
