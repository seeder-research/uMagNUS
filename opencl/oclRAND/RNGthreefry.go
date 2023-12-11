package oclRAND

import (
	"fmt"
	"log"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	timer "github.com/seeder-research/uMagNUS/timer"
	"math/rand"
)

func (p *THREEFRY_status_array_ptr) Init(seed uint64, queue *cl.CommandQueue, events []*cl.Event) {
	// Generate random seed array to seed the PRNG
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

	// Copy random seed array to GPU
	context := p.GetContext()
	seed_buf, err := context.CreateBufferUnsafe(cl.MemReadWrite, int(unsafe.Sizeof(seed_arr[0]))*totalCount, nil)
	defer seed_buf.Release()
	if err != nil {
		log.Fatalln("Unable to create buffer for THREEFRY seed array!")
	}
	var seed_event *cl.Event
	seed_event, err = queue.EnqueueWriteBuffer(seed_buf, false, 0, int(unsafe.Sizeof(seed_arr[0]))*totalCount, unsafe.Pointer(&seed_arr[0]), nil)
	if err != nil {
		log.Fatalln("Unable to write seed buffer to device: ", err)
	}
	if events != nil {
		err = cl.WaitForEvents(events)
		if err != nil {
			fmt.Printf("First WaitForEvents failed in InitRNG: %+v \n", err)
		}
	}

	// Seed the RNG
	event := k_threefry_seed_async(unsafe.Pointer(p.Status_key), unsafe.Pointer(p.Status_counter),
		unsafe.Pointer(p.Status_result), unsafe.Pointer(p.Status_tracker), unsafe.Pointer(seed_buf),
		&config{[]int{totalCount}, []int{p.GetGroupSize()}}, queue, []*cl.Event{seed_event})

	p.Ini = true
	err = cl.WaitForEvents([]*cl.Event{event})
	if err != nil {
		fmt.Printf("Second WaitForEvents failed in InitRNG: %+v \n", err)
	}

}

func (p *THREEFRY_status_array_ptr) GenerateUniform(d_data unsafe.Pointer, data_size int, queue *cl.CommandQueue, events []*cl.Event) *cl.Event {

	if p.Ini == false {
		log.Fatalln("Generator has not been initialized!")
	}

	if Synchronous { // debug
		if err := queue.Finish(); err != nil {
			fmt.Printf("failed to wait for queue to finish in beginning of generateuniform: %+v \n", err)
		}
		timer.Start("threefry_uniform")
	}

	event := k_threefry_uniform_async(unsafe.Pointer(p.Status_key), unsafe.Pointer(p.Status_counter),
		unsafe.Pointer(p.Status_result), unsafe.Pointer(p.Status_tracker), d_data, data_size,
		&config{[]int{p.GetStatusSize()}, []int{p.GetGroupSize()}}, queue, events)

	if Synchronous { // debug
		if err := queue.Finish(); err != nil {
			fmt.Printf("failed to wait for queue to finish at end of generateuniform: %+v \n", err)
		}
		timer.Stop("threefry_uniform")
	}

	return event
}

func (p *THREEFRY_status_array_ptr) GenerateNormal(d_data unsafe.Pointer, data_size int, queue *cl.CommandQueue, events []*cl.Event) *cl.Event {

	if p.Ini == false {
		log.Fatalln("Generator has not been initialized!")
	}

	if Synchronous { // debug
		if err := queue.Finish(); err != nil {
			fmt.Printf("failed to wait for queue to finish in beginning of generateuniform: %+v \n", err)
		}
		timer.Start("threefry_normal")
	}

	event := k_threefry_normal_async(unsafe.Pointer(p.Status_key), unsafe.Pointer(p.Status_counter),
		unsafe.Pointer(p.Status_result), unsafe.Pointer(p.Status_tracker), d_data, data_size,
		&config{[]int{p.GetStatusSize()}, []int{p.GetGroupSize()}}, queue, events)

	if Synchronous { // debug
		if err := queue.Finish(); err != nil {
			fmt.Printf("failed to wait for queue to finish in beginning of generateuniform: %+v \n", err)
		}
		timer.Stop("threefry_normal")
	}

	return event
}
