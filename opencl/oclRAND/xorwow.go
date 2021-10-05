package oclRAND

import (
	"github.com/seeder-research/uMagNUS/opencl/cl"
	"log"
)

func NewXORWOWStatus() *XORWOW_status_array_ptr {
	q := new(XORWOW_status_array_ptr)
	q.Ini = false
	q.Status_size = -1
	return q
}

func (p *XORWOW_status_array_ptr) SetSize(N int) {
	p.Status_buf, err = context.CreateBufferUnsafe(cl.MemReadWrite, int(unsafe.Sizeof(uint32))*6*N, nil)
	if err != nil {
		log.Fatalln("Unable to create buffer for XORWOW status array!")
	}
	p.Status_size = N
}

func (p *XORWOW_status_array_ptr) CreateStatusBuf() {
	if p.Status_size <= 0 {
		log.Fatalln("Unable to create buffer for XORWOW status array: number of PRNGs is less than 1")
	}
	p.Status_buf, err = context.CreateBufferUnsafe(cl.MemReadWrite, int(unsafe.Sizeof(uint32))*6*p.Status_size, nil)
	if err != nil {
		log.Fatalln("Unable to create buffer for XORWOW status array!")
	}
}
