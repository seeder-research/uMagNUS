package opencl

// This file provides GPU byte slices, used to store regions.

import (
	"log"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	util "github.com/seeder-research/uMagNUS/util"
)

// 3D byte slice, used for region lookup.
type Bytes struct {
	Ptr unsafe.Pointer
	Len int
	Evt *cl.Event
}

// Construct new byte slice with given length,
// initialised to zeros.
func NewBytes(Len int) *Bytes {
	ptr, err := ClCtx.CreateEmptyBuffer(cl.MemReadWrite, Len)
	if err != nil {
		panic(err)
	}
	zeroPattern := uint8(0)
	var event *cl.Event
	event, err = ClCmdQueue.EnqueueFillBuffer(ptr, unsafe.Pointer(&zeroPattern), 1, 0, Len, nil)
	if err != nil {
		panic(err)
	}
	if Debug {
		if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
			log.Panic("WaitForEvents failed in NewBytes:", err)
		}
	}
	return &Bytes{unsafe.Pointer(ptr), Len, event}
}

// Upload src (host) to dst (gpu).
func (dst *Bytes) Upload(src []byte) {
	util.Argument(dst.Len == len(src))
	dstEvt := dst.GetEvent()
	if dstEvt != nil {
		if err := cl.WaitForEvents([](*cl.Event){dstEvt}); err != nil {
			log.Panic("WaitForEvents failed in Upload:", err)
		}
	}
	MemCpyHtoD(dst.Ptr, unsafe.Pointer(&src[0]), dst.Len)
}

// Copy on device: dst = src.
func (dst *Bytes) Copy(src *Bytes) {
	util.Argument(dst.Len == src.Len)
	eventWaitList := []*cl.Event{}
	tmpEvt := dst.GetEvent()
	if tmpEvt != nil {
		eventWaitList = append(eventWaitList, tmpEvt)
	}
	tmpEvt = src.GetEvent()
	if tmpEvt != nil {
		eventWaitList = append(eventWaitList, tmpEvt)
	}
	if len(eventWaitList) > 0 {
		if err := cl.WaitForEvents(eventWaitList); err != nil {
			log.Panic("WaitForEvents failed in Copy:", err)
		}
	}
	MemCpy(dst.Ptr, src.Ptr, dst.Len)
}

// Copy to host: dst = src.
func (src *Bytes) Download(dst []byte) {
	util.Argument(src.Len == len(dst))
	srcEvt := src.GetEvent()
	if srcEvt != nil {
		if err := cl.WaitForEvents([](*cl.Event){srcEvt}); err != nil {
			log.Panic("WaitForEvents failed in Download:", err)
		}
	}
	MemCpyDtoH(unsafe.Pointer(&dst[0]), src.Ptr, src.Len)
}

// Set one element to value.
// data.Index can be used to find the index for x,y,z.
func (dst *Bytes) Set(index int, value byte) {
	if index < 0 || index >= dst.Len {
		log.Panic("Bytes.Set: index out of range:", index)
	}
	src := value
	dstEvt := dst.GetEvent()
	var event *cl.Event
	var err error
	if dstEvt != nil {
		event, err = ClCmdQueue.EnqueueWriteBuffer((*cl.MemObject)(dst.Ptr), false, index, 1, unsafe.Pointer(&src), []*cl.Event{dstEvt})
	} else {
		event, err = ClCmdQueue.EnqueueWriteBuffer((*cl.MemObject)(dst.Ptr), false, index, 1, unsafe.Pointer(&src), nil)
	}
	if err != nil {
		panic(err)
	}
	dst.SetEvent(event)
	if Debug {
		if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
			log.Panic("WaitForEvents failed in Bytes.Set():", err)
		}
	}
}

// Get one element.
// data.Index can be used to find the index for x,y,z.
func (src *Bytes) Get(index int) byte {
	if index < 0 || index >= src.Len {
		log.Panic("Bytes.Set: index out of range:", index)
	}
	dst := make([]byte, 1)
	srcEvent := src.GetEvent()
	var event *cl.Event
	var err error
	if srcEvent != nil {
		event, err = ClCmdQueue.EnqueueReadBufferByte((*cl.MemObject)(src.Ptr), false, index, dst, []*cl.Event{srcEvent})
	} else {
		event, err = ClCmdQueue.EnqueueReadBufferByte((*cl.MemObject)(src.Ptr), false, index, dst, nil)
	}
	if err != nil {
		panic(err)
	}
	src.SetEvent(event)
	// Must synchronize
	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		log.Panic("WaitForEvents failed in Bytes.Set():", err)
	}
	return dst[0]
}

// Frees the GPU memory and disables the slice.
func (b *Bytes) Free() {
	if b.Ptr != nil {
		tmpObj := (*cl.MemObject)(b.Ptr)
		tmpObj.Release()
	}
	b.Ptr = nil
	b.Len = 0
	b.Evt = nil
}

// Set the event to synchonize the buffer of bytes
func (b *Bytes) SetEvent(e *cl.Event) {
	b.Evt = e
}

// Get the event to synchonize the buffer of bytes
func (b *Bytes) GetEvent() *cl.Event {
	return b.Evt
}
