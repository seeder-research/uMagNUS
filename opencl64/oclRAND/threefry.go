package oclRAND

import (
	"github.com/seeder-research/uMagNUS/cl"

	"log"
	"unsafe"
)

type THREEFRY_status_array_ptr struct {
	Ini            bool
	Seed_Arr       []uint64
	Status_key     *cl.MemObject
	Status_counter *cl.MemObject
	Status_result  *cl.MemObject
	Status_tracker *cl.MemObject
	Status_size    int
	GroupSize      int
	GroupCount     int
	ClCtx          *cl.Context
}

func NewTHREEFRYStatus() *THREEFRY_status_array_ptr {
	q := new(THREEFRY_status_array_ptr)
	q.Ini = false
	q.Status_size = -1
	return q
}

func (p *THREEFRY_status_array_ptr) SetContext(context *cl.Context) {
	p.ClCtx = context
}

func (p *THREEFRY_status_array_ptr) GetContext() *cl.Context {
	return p.ClCtx
}

func (p *THREEFRY_status_array_ptr) SetStatusSize(N int) {
	p.Status_size = N
}

func (p *THREEFRY_status_array_ptr) GetStatusSize() int {
	return p.Status_size
}

func (p *THREEFRY_status_array_ptr) CreateStatusBuffer(context *cl.Context) {
	p.SetContext(context)
	if p.Status_size <= 0 {
		log.Fatalln("Unable to create buffer for THREEFRY status array: number of PRNGs is less than 1")
	}
	var err error
	var testVar uint32
	p.Status_key, err = p.ClCtx.CreateBufferUnsafe(cl.MemReadWrite, int(unsafe.Sizeof(testVar))*SIZEOF_FLOAT64*p.Status_size, nil)
	if err != nil {
		log.Fatalln("Unable to create buffer for THREEFRY status key array!")
	}
	p.Status_counter, err = p.ClCtx.CreateBufferUnsafe(cl.MemReadWrite, int(unsafe.Sizeof(testVar))*SIZEOF_FLOAT64*p.Status_size, nil)
	if err != nil {
		log.Fatalln("Unable to create buffer for THREEFRY status counter array!")
	}
	p.Status_result, err = p.ClCtx.CreateBufferUnsafe(cl.MemReadWrite, int(unsafe.Sizeof(testVar))*SIZEOF_FLOAT64*p.Status_size, nil)
	if err != nil {
		log.Fatalln("Unable to create buffer for THREEFRY status result array!")
	}
	p.Status_tracker, err = p.ClCtx.CreateBufferUnsafe(cl.MemReadWrite, int(unsafe.Sizeof(testVar))*p.Status_size, nil)
	if err != nil {
		log.Fatalln("Unable to create buffer for THREEFRY status tracker array!")
	}
}

func (p *THREEFRY_status_array_ptr) SetGroupSize(in int) {
	p.GroupSize = in
}

func (p *THREEFRY_status_array_ptr) GetGroupSize() int {
	return p.GroupSize
}

func (p *THREEFRY_status_array_ptr) SetGroupCount(in int) {
	p.GroupCount = in
}

func (p *THREEFRY_status_array_ptr) GetGroupCount() int {
	return p.GroupCount
}

func (p *THREEFRY_status_array_ptr) RecommendSize() int {
	return 16 * p.Status_size
}
