package oclRAND

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/seeder-research/uMagNUS/opencl/cl"
	"github.com/seeder-research/uMagNUS/timer"
)

func (p *MTGP32dc_params_array_ptr) Init(seed uint64, events []*cl.Event) {

	SeedVal := seed << 32
	SeedVal = SeedVal >> 32
	if SeedVal == 0 {
		SeedVal = seed >> 32
	}
	event := k_mtgp32_init_seed_kernel_async(unsafe.Pointer(p.Rec_buf), unsafe.Pointer(p.Temper_buf), unsafe.Pointer(p.Flt_temper_buf), unsafe.Pointer(p.Pos_buf),
		unsafe.Pointer(p.Sh1_buf), unsafe.Pointer(p.Sh2_buf), unsafe.Pointer(p.Status_buf), uint32(SeedVal),
		&config{[]int{p.GetGroupCount() * p.GetGroupSize()}, []int{p.GetGroupSize()}}, events)

	p.Ini = true
	err := cl.WaitForEvents([]*cl.Event{event})
	if err != nil {
		fmt.Printf("WaitForEvents failed in InitRNG: %+v \n", err)
	}

}

func (p *MTGP32dc_params_array_ptr) GenerateUniform(d_data unsafe.Pointer, data_size int, events []*cl.Event) *cl.Event {

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
		timer.Start("mtgp32_uniform")
	}

	event := k_mtgp32_uniform_async(unsafe.Pointer(p.Rec_buf), unsafe.Pointer(p.Temper_buf), unsafe.Pointer(p.Flt_temper_buf), unsafe.Pointer(p.Pos_buf),
		unsafe.Pointer(p.Sh1_buf), unsafe.Pointer(p.Sh2_buf), unsafe.Pointer(p.Status_buf), d_data, tmpSize/p.GetGroupCount(),
		&config{[]int{item_num}, []int{MTGPDC_TN}}, events)

	if Synchronous { // debug
		ClCmdQueue.Finish()
		timer.Stop("mtgp32_uniform")
	}

	return event
}

func (p *MTGP32dc_params_array_ptr) GenerateNormal(d_data unsafe.Pointer, data_size int, events []*cl.Event) *cl.Event {

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
		timer.Start("mtgp32_normal")
	}

	event := k_mtgp32_normal_async(unsafe.Pointer(p.Rec_buf), unsafe.Pointer(p.Temper_buf), unsafe.Pointer(p.Flt_temper_buf), unsafe.Pointer(p.Pos_buf),
		unsafe.Pointer(p.Sh1_buf), unsafe.Pointer(p.Sh2_buf), unsafe.Pointer(p.Status_buf), d_data, tmpSize/p.GetGroupCount(),
		&config{[]int{item_num}, []int{MTGPDC_TN}}, events)

	if Synchronous { // debug
		ClCmdQueue.Finish()
		timer.Stop("mtgp32_normal")
	}

	return event
}
