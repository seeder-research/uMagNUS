package opencl

// This file provides GPU byte slices, used to store regions.

import (
	"log"
	"sync"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	qm "github.com/seeder-research/uMagNUS/queuemanager"
	util "github.com/seeder-research/uMagNUS/util"
)

// 3D byte slice, used for region lookup.
type Bytes struct {
	Ptr   unsafe.Pointer
	Len   int
	Evt   *cl.Event
	RdEvt *data.SliceEventMap
}

// Construct new byte slice with given length,
// initialised to zeros.
func NewBytes(Len int) *Bytes {
	if Synchronous { // debug
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}

	ptr, err := ClCtx.CreateEmptyBuffer(cl.MemReadWrite, Len)
	if err != nil {
		panic(err)
	}
	zeroPattern := uint8(0)

	// Checkout command queue from pool and launch kernel
	var newBytesSyncWaitGroup sync.WaitGroup
	var event *cl.Event
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &newBytesSyncWaitGroup)
	event, err = tmpQueue.EnqueueFillBuffer(ptr, unsafe.Pointer(&zeroPattern), 1, 0, Len, nil)
	if err != nil {
		panic(err)
	}

	// Check in queue post execution
	qwg := qm.NewQueueWaitGroup(tmpQueue, &newBytesSyncWaitGroup)
	ReturnQueuePool <- qwg

	if Synchronous || Debug {
		if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
			log.Panic("WaitForEvents failed in NewBytes:", err)
		}
	}
	emptyMap := new(data.SliceEventMap)
	emptyMap.Init()
	return &Bytes{unsafe.Pointer(ptr), Len, event, emptyMap}
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
	// synchronization should be done by code calling
	// this function

	util.Argument(dst.Len == src.Len)
	MemCpy(dst.Ptr, src.Ptr, dst.Len)
}

// Copy to host: dst = src.
func (src *Bytes) Download(dst []byte) {
	// synchronization should be done by code calling
	// this function

	util.Argument(src.Len == len(dst))
	MemCpyDtoH(unsafe.Pointer(&dst[0]), src.Ptr, src.Len)
}

// Set one element to value.
// data.Index can be used to find the index for x,y,z.
func (dst *Bytes) Set(index int, value byte) {
	// synchronization should be done by code calling
	// this function

	if index < 0 || index >= dst.Len {
		log.Panic("Bytes.Set: index out of range:", index)
	}
	src := value

	// Checkout command queue from pool and launch kernel
	var setByteSyncWaitGroup sync.WaitGroup
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &setByteSyncWaitGroup)
	event, err := tmpQueue.EnqueueWriteBuffer((*cl.MemObject)(dst.Ptr), false, index, 1, unsafe.Pointer(&src), nil)
	if err != nil {
		panic(err)
	}

	// Check in queue post execution
	qwg := qm.NewQueueWaitGroup(tmpQueue, &setByteSyncWaitGroup)
	ReturnQueuePool <- qwg

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
	// synchronization should be done by code calling
	// this function

	if index < 0 || index >= src.Len {
		log.Panic("Bytes.Set: index out of range:", index)
	}
	dst := make([]byte, 1)

	// Checkout command queue from pool and launch kernel
	var getByteSyncWaitGroup sync.WaitGroup
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &getByteSyncWaitGroup)
	event, err := tmpQueue.EnqueueReadBufferByte((*cl.MemObject)(src.Ptr), false, index, dst, nil)
	if err != nil {
		panic(err)
	}

	// Check in queue post execution
	qwg := qm.NewQueueWaitGroup(tmpQueue, &getByteSyncWaitGroup)
	ReturnQueuePool <- qwg

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
