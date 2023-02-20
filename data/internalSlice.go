package data

import (
	"sync"
	"unsafe"
)

// The pointer to the actual memory resides in internalSlice.
// For storing a single component in a data.Slice
type internalSlice struct {
	Ptr     unsafe.Pointer
	Size    [3]int
	MemType int8
	Wg      sync.WaitGroup
	sync.RWMutex
}

func newInternalSlice() *internalSlice {
	iS := new(internalSlice)
	iS.SetPtr(nil)
	iS.SetSize([3]int{0, 0, 0})
	iS.SetMemType(int8(0))
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
