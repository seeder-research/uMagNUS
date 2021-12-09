package oclRAND

import (
	"github.com/seeder-research/uMagNUS/cl"

	"log"
	"unsafe"
)

type XORWOW_status_array_ptr struct {
	Ini         bool
	Seed_Arr    []uint64
	Status_buf  *cl.MemObject
	Status_size int
	GroupSize   int
	GroupCount  int
	ClCtx       *cl.Context
}

func NewXORWOWStatus() *XORWOW_status_array_ptr {
	q := new(XORWOW_status_array_ptr)
	q.Ini = false
	q.Status_size = -1
	return q
}

func (p *XORWOW_status_array_ptr) SetContext(context *cl.Context) {
	p.ClCtx = context
}

func (p *XORWOW_status_array_ptr) GetContext() *cl.Context {
	return p.ClCtx
}

func (p *XORWOW_status_array_ptr) SetStatusSize(N int) {
	p.Status_size = N
}

func (p *XORWOW_status_array_ptr) GetStatusSize() int {
	return p.Status_size
}

func (p *XORWOW_status_array_ptr) CreateStatusBuffer(context *cl.Context) {
	p.SetContext(context)
	if p.Status_size <= 0 {
		log.Fatalln("Unable to create buffer for XORWOW status array: number of PRNGs is less than 1")
	}
	var err error
	var testVar uint32
	p.Status_buf, err = p.ClCtx.CreateBufferUnsafe(cl.MemReadWrite, int(unsafe.Sizeof(testVar))*6*p.Status_size, nil)
	if err != nil {
		log.Fatalln("Unable to create buffer for XORWOW status array!")
	}
}

func (p *XORWOW_status_array_ptr) SetGroupSize(in int) {
	p.GroupSize = in
}

func (p *XORWOW_status_array_ptr) GetGroupSize() int {
	return p.GroupSize
}

func (p *XORWOW_status_array_ptr) SetGroupCount(in int) {
	p.GroupCount = in
}

func (p *XORWOW_status_array_ptr) GetGroupCount() int {
	return p.GroupCount
}

func (p *XORWOW_status_array_ptr) RecommendSize() int {
	return 8 * p.Status_size
}
