package oclRAND

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/seeder-research/uMagNUS/cl"
	"github.com/seeder-research/uMagNUS/timer"
	"math/rand"
)

func (p *MTGP64dc_params_array_ptr) Init(seed uint64, events []*cl.Event) {

	// Generate random seed array to seed the PRNG
	rand.Seed((int64)(seed))
	totalCount := p.GetGroupCount()
	seed_arr := make([]uint64, totalCount)
	for idx := 0; idx < totalCount; idx++ {
		tmpNum := rand.Uint32()
		for tmpNum == 0 {
			tmpNum = rand.Uint64()
		}
		seed_arr[idx] = tmpNum
	}

	// Copy random seed array to GPU
	context := p.GetContext()
	seed_buf, err := context.CreateBufferUnsafe(cl.MemReadWrite, int(unsafe.Sizeof(seed_arr[0]))*totalCount, nil)
	defer seed_buf.Release()
	if err != nil {
		log.Fatalln("Unable to create buffer for mtgp64 seed array!")
	}
	var seed_event *cl.Event
	seed_event, err = ClCmdQueue.EnqueueWriteBuffer(seed_buf, false, 0, int(unsafe.Sizeof(seed_arr[0]))*totalCount, unsafe.Pointer(&seed_arr[0]), nil)
	if err != nil {
		log.Fatalln("Unable to write seed buffer to device: ", err)
	}
	if events != nil {
		err = cl.WaitForEvents(events)
		if err != nil {
			fmt.Printf("First WaitForEvents failed in InitRNG: %+v \n", err)
		}
	}

	event := k_mtgp64_init_seed_kernel_async(unsafe.Pointer(p.Rec_buf), unsafe.Pointer(p.Temper_buf), unsafe.Pointer(p.Flt_temper_buf), unsafe.Pointer(p.Pos_buf),
		unsafe.Pointer(p.Sh1_buf), unsafe.Pointer(p.Sh2_buf), unsafe.Pointer(p.Status_buf), unsafe.Pointer(seed_buf),
		&config{[]int{p.GetGroupCount() * p.GetGroupSize()}, []int{p.GetGroupSize()}}, []*cl.Event{seed_event})

	p.Ini = true
	err = cl.WaitForEvents([]*cl.Event{event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in InitRNG: %+v \n", err)
	}

}

func (p *MTGP64dc_params_array_ptr) GenerateUniform(d_data unsafe.Pointer, data_size int, events []*cl.Event) *cl.Event {

	if p.Ini == false {
		log.Fatalln("Generator has not been initialized!")
	}

	item_num := MTGPDC_TN * p.GetGroupCount()
	min_size := MTGPDC_LS * p.GetGroupCount()
	tmpSize := data_size
	if data_size%min_size != 0 {
		tmpSize = (data_size/min_size + 1) * min_size
	}

	if Synchronous { // debug
		ClCmdQueue.Finish()
		timer.Start("mtgp64_uniform")
	}

	event := k_mtgp64_uniform_async(unsafe.Pointer(p.Rec_buf), unsafe.Pointer(p.Temper_buf), unsafe.Pointer(p.Flt_temper_buf), unsafe.Pointer(p.Pos_buf),
		unsafe.Pointer(p.Sh1_buf), unsafe.Pointer(p.Sh2_buf), unsafe.Pointer(p.Status_buf), d_data, tmpSize/p.GetGroupCount(),
		&config{[]int{item_num}, []int{MTGPDC_TN}}, events)

	if Synchronous { // debug
		ClCmdQueue.Finish()
		timer.Stop("mtgp64_uniform")
	}

	return event
}

func (p *MTGP64dc_params_array_ptr) GenerateNormal(d_data unsafe.Pointer, data_size int, events []*cl.Event) *cl.Event {

	if p.Ini == false {
		log.Fatalln("Generator has not been initialized!")
	}

	item_num := MTGPDC_TN * p.GetGroupCount()
	min_size := MTGPDC_LS * p.GetGroupCount()
	tmpSize := data_size
	if data_size%min_size != 0 {
		tmpSize = (data_size/min_size + 1) * min_size
	}

	if Synchronous { // debug
		ClCmdQueue.Finish()
		timer.Start("mtgp64_normal")
	}

	event := k_mtgp64_normal_async(unsafe.Pointer(p.Rec_buf), unsafe.Pointer(p.Temper_buf), unsafe.Pointer(p.Flt_temper_buf), unsafe.Pointer(p.Pos_buf),
		unsafe.Pointer(p.Sh1_buf), unsafe.Pointer(p.Sh2_buf), unsafe.Pointer(p.Status_buf), d_data, tmpSize/p.GetGroupCount(),
		&config{[]int{item_num}, []int{MTGPDC_TN}}, events)

	if Synchronous { // debug
		ClCmdQueue.Finish()
		timer.Stop("mtgp64_normal")
	}

	return event
}
