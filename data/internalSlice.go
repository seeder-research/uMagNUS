package data

import (
	"sync"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
)

type SliceEventMap struct {
	ReadEvents map[*cl.Event]int8
	sync.RWMutex
}

func (sem *SliceEventMap) Init() {
	sem.ReadEvents = make(map[*cl.Event]int8)
}

// The pointer to the actual memory resides in internalSlice.
// For storing a single component in a data.Slice
type internalSlice struct {
	Ptr     unsafe.Pointer
	Size    [3]int
	MemType int8
	Event   *cl.Event
	RdEvent *SliceEventMap
	Wg      sync.WaitGroup
	sync.RWMutex
}

func newInternalSlice() *internalSlice {
	iS := new(internalSlice)
	iS.SetPtr(nil)
	iS.SetSize([3]int{0, 0, 0})
	iS.SetMemType(int8(0))
	iS.SetEvent(nil)
	iS.RdEvent = new(SliceEventMap)
	iS.RdEvent.Init()
	return iS
}

func (iS *internalSlice) Free() {
	if iS.Ptr == nil {
		return
	}
	memFree(iS.Ptr)
}

// Gets the pointer of the slice
func (iS *internalSlice) GetPtr() unsafe.Pointer {
	return iS.Ptr
}

// Sets the pointer of the slice
func (iS *internalSlice) SetPtr(ptr unsafe.Pointer) {
	iS.Ptr = ptr
}

// Gets the size of the slice
func (iS *internalSlice) GetSize() [3]int {
	return iS.Size
}

// Sets the size of the slice
func (iS *internalSlice) SetSize(size [3]int) {
	iS.Size = size
}

// Gets the MemType of the slice
func (iS *internalSlice) GetMemType() int8 {
	return iS.MemType
}

// Sets the MemType of the slice
func (iS *internalSlice) SetMemType(memType int8) {
	iS.MemType = memType
}

// Functions for manipulating associated Events
// Gets the event of the slice
func (iS *internalSlice) GetEvent() *cl.Event {
	return iS.Event
}

// Sets the event of the slice
func (iS *internalSlice) SetEvent(ev *cl.Event) {
	iS.Event = ev
}

// Functions for manipulating associated ReadEvents
func (iS *internalSlice) ClearAllEvents() {
	iS.SetEvent(nil)
	iS.ClearReadEvents()
}

// Emtpies ReadEvents of the slice
func (iS *internalSlice) ClearReadEvents() {
	iS.RdEvent.Lock()
	iS.RdEvent.Init()
	iS.RdEvent.Unlock()
}

// Inserts a cl.Event into ReadEvents of the slice
func (iS *internalSlice) InsertReadEvent(event *cl.Event) {
	iS.RdEvent.Lock()
	if _, ok := iS.RdEvent.ReadEvents[event]; ok == false {
		iS.RdEvent.ReadEvents[event] = 1
	}
	iS.RdEvent.Unlock()
}

// Removes a cl.Event from ReadEvents of the slice
func (iS *internalSlice) RemoveReadEvent(event *cl.Event) {
	iS.RdEvent.Lock()
	if _, ok := iS.RdEvent.ReadEvents[event]; ok {
		delete(iS.RdEvent.ReadEvents, event)
	}
	iS.RdEvent.Unlock()
}

// Gets the ReadEvents of the slice as a slice
func (iS *internalSlice) GetReadEvents() []*cl.Event {
	iS.RdEvent.RLock()
	evList := []*cl.Event{}
	for k, _ := range iS.RdEvent.ReadEvents {
		if k != nil {
			evList = append(evList, k)
		}
	}
	iS.RdEvent.RUnlock()
	return evList
}

// Sets the ReadEvents of the slice
func (iS *internalSlice) SetReadEvents(eventList []*cl.Event) {
	iS.RdEvent.Lock()
	iS.RdEvent.Init()
	for _, e := range eventList {
		iS.RdEvent.ReadEvents[e] = 1
	}
	iS.RdEvent.Unlock()
}

// Returns all events of the slice (for syncing kernels writing to the slice)
func (iS *internalSlice) GetAllEvents() []*cl.Event {
	evList := iS.GetReadEvents()
	evList = append(evList, iS.GetEvent())
	return evList
}

// Functions associated with sync.WaitGroup
// Increment the WaitGroup
func (iS *internalSlice) Add(i int) {
	iS.Wg.Add(i)
}

// Decrement the WaitGroup
func (iS *internalSlice) Done() {
	iS.Wg.Done()
}

// Wait on the WaitGroup
func (iS *internalSlice) Wait() {
	iS.Wg.Wait()
}
