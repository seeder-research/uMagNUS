package opencl

// This file provides GPU byte slices, used to store regions.

import (
	"log"
	"sync"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

// 3D byte slice, used for region lookup.
type Bytes struct {
	Ptr   unsafe.Pointer
	Len   int
	Evt   *cl.Event
	RdEvt *data.SliceEventMap
	sync.RWMutex
}

// Construct new byte slice with given length,
// initialised to zeros.
func NewBytes(Len int) *Bytes {
	ptr, err1 := ClCtx.CreateEmptyBuffer(cl.MemReadWrite, Len)
	if err1 != nil {
		panic(err1)
	}

	emptyMap := new(data.SliceEventMap)
	emptyMap.Init()

	outByte := new(Bytes)
	outByte.Ptr = unsafe.Pointer(ptr)
	outByte.Len = Len
	outByte.RdEvt = emptyMap

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		byte_zero__(outByte, &wg)
	} else {
		go byte_zero__(outByte, &wg)
	}
	wg.Wait()
	return outByte
}

func byte_zero__(b *Bytes, wg_ *sync.WaitGroup) {
	b.Zero(wg_)
}

func (dst *Bytes) Zero(wg_ *sync.WaitGroup) {
	dst.Lock()
	defer dst.Unlock()

	zeroPattern := uint8(0)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		log.Printf("bytes.zero failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	var event *cl.Event
	event, err = cmdqueue.EnqueueFillBuffer((*cl.MemObject)(dst.Ptr), unsafe.Pointer(&zeroPattern), 1, 0, dst.Len, nil)
	wg_.Done()
	if err != nil {
		panic(err)
	}

	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		log.Panicf("WaitForEvents failed in bytes.zero:", err)
	}
}

// Upload src (host) to dst (gpu).
func (dst *Bytes) Upload(src []byte) {
	util.Argument(dst.Len == len(src))
	dst.Lock()
	ev := MemCpyHtoD(dst.Ptr, unsafe.Pointer(&src[0]), dst.Len)
	if Synchronous {
		if err := cl.WaitForEvents(ev); err != nil {
			panic(err)
		}
		dst.Unlock()
	} else {
		go func() {
			if err := cl.WaitForEvents(ev); err != nil {
				panic(err)
			}
			dst.Unlock()
		}()
	}
}

// Copy on device: dst = src.
func (dst *Bytes) Copy(src *Bytes) {
	util.Argument(dst.Len == src.Len)
	dst.Lock()
	src.RLock()
	ev := MemCpy(dst.Ptr, src.Ptr, dst.Len)
	if Synchronous {
		if err := cl.WaitForEvents(ev); err != nil {
			panic(err)
		}
		dst.Unlock()
		src.RUnlock()
	} else {
		go func() {
			if err := cl.WaitForEvents(ev); err != nil {
				panic(err)
			}
			dst.Unlock()
			src.RUnlock()
		}()
	}
}

// Copy to host: dst = src.
func (src *Bytes) Download(dst []byte) {
	util.Argument(src.Len == len(dst))
	src.RLock()
	ev := MemCpyDtoH(unsafe.Pointer(&dst[0]), src.Ptr, src.Len)
	if err := cl.WaitForEvents(ev); err != nil {
		panic(err)
	}
	src.RUnlock()
}

// Set one element to value.
// data.Index can be used to find the index for x,y,z.
func (dst *Bytes) Set(index int, value byte) {
	if index < 0 || index >= dst.Len {
		log.Panic("Bytes.Set: index out of range:", index)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		bytes_set__(dst, index, value, &wg)
	} else {
		go bytes_set__(dst, index, value, &wg)
	}
	wg.Wait()
}

func bytes_set__(dst *Bytes, index int, value byte, wg_ *sync.WaitGroup) {
	dst.Lock()
	defer dst.Unlock()

	src := value

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		log.Panicf("MemCpyDoH failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	var event *cl.Event
	event, err = cmdqueue.EnqueueWriteBuffer((*cl.MemObject)(dst.Ptr), false, index, 1, unsafe.Pointer(&src), nil)
	wg_.Done()
	if err != nil {
		panic(err)
	}
	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		log.Panicf("WaitForEvents failed in Bytes.Set():", err)
	}
}

// Get one element.
// data.Index can be used to find the index for x,y,z.
func (src *Bytes) Get(index int) byte {
	if index < 0 || index >= src.Len {
		log.Panic("Bytes.Set: index out of range:", index)
	}
	dst := make([]byte, 1)
	src.RLock()
	defer src.RUnlock()

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		log.Panicf("MemCpyDoH failed to create command queue: %+v \n", err)
		return byte(0)
	}
	defer cmdqueue.Release()

	event, err := cmdqueue.EnqueueReadBufferByte((*cl.MemObject)(src.Ptr), false, index, dst, nil)
	if err != nil {
		panic(err)
	}
	// Must synchronize
	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		log.Panic("WaitForEvents failed in Bytes.Get():", err)
	}
	return dst[0]
}

// Frees the GPU memory and disables the slice.
func (b *Bytes) Free() {
	b.Lock()
	b.Unlock()
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

// Sets the rdEvent of the slice
func (b *Bytes) SetReadEvents(eventList []*cl.Event) {
	b.RdEvt.Lock()
	for _, e := range eventList {
		if _, ok := b.RdEvt.ReadEvents[e]; ok == false {
			b.RdEvt.ReadEvents[e] = 1
		}
	}
	b.RdEvt.Unlock()
}

// Insert a cl.Event to rdEvent of the slice
func (b *Bytes) InsertReadEvent(event *cl.Event) {
	b.RdEvt.Lock()
	if _, ok := b.RdEvt.ReadEvents[event]; ok == false {
		b.RdEvt.ReadEvents[event] = 1
	}
	b.RdEvt.Unlock()
}

// Remove a cl.Event from rdEvent of the slice
func (b *Bytes) RemoveReadEvent(event *cl.Event) {
	b.RdEvt.Lock()
	if _, ok := b.RdEvt.ReadEvents[event]; ok == false {
		delete(b.RdEvt.ReadEvents, event)
	}
	b.RdEvt.Unlock()
}

// Returns rdEvent of the slice as a slice
func (b *Bytes) GetReadEvents() []*cl.Event {
	b.RdEvt.RLock()
	evList := []*cl.Event{}
	for k, _ := range b.RdEvt.ReadEvents {
		if k != nil {
			evList = append(evList, k)
		}
	}
	b.RdEvt.RUnlock()
	return evList
}

// Returns all events of the slice (for syncing kernels writing to the slice)
func (b *Bytes) GetAllEvents() []*cl.Event {
	eventList := b.GetReadEvents()
	if b.Evt != nil {
		eventList = append(eventList, b.Evt)
	}
	return eventList
}
