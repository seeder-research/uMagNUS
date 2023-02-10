package opencl64

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
	timer "github.com/seeder-research/uMagNUS/timer"
	util "github.com/seeder-research/uMagNUS/util"
)

// Make a GPU Slice with nComp components each of size length.
func NewSlice(nComp int, size [3]int) *data.Slice {
	return newSlice(nComp, size, data.GPUMemory)
}

func newSlice(nComp int, size [3]int, memType int8) *data.Slice {
	length := prod(size)
	bytes := length * SIZEOF_FLOAT64
	ptrs := make([]unsafe.Pointer, nComp)
	initVal := float64(0.0)
	fillWait := make([]*cl.Event, nComp)
	for c := range ptrs {
		tmp_buf, err := ClCtx.CreateEmptyBuffer(cl.MemReadWrite, bytes)
		if err != nil {
			fmt.Printf("CreateEmptyBuffer failed: %+v \n", err)
		}
		ptrs[c] = unsafe.Pointer(tmp_buf)
		fillWait[c], err = ClCmdQueue.EnqueueFillBuffer(tmp_buf, unsafe.Pointer(&initVal), SIZEOF_FLOAT64, 0, bytes, nil)
		if err != nil {
			fmt.Printf("CreateEmptyBuffer failed: %+v \n", err)
		}
		if err = cl.WaitForEvents([]*cl.Event{fillWait[c]}); err != nil {
			fmt.Printf("Wait for EnqueueFillBuffer failed: %+v \n", err)
		}
	}

	dataPtr := data.SliceFromPtrs(size, memType, ptrs)
	dataPtr.SetEvents(fillWait)
	return dataPtr
}

// wrappers for data.EnableGPU arguments

func memFree(ptr unsafe.Pointer) {
	if ptr != nil {
		buf := (*cl.MemObject)(ptr)
		buf.Release()
	}
}

func MemCpyDtoH(dst, src unsafe.Pointer, bytes int) []*cl.Event {
	timer.Start("memcpyDtoH")

	// execute
	event, err := ClCmdQueue.EnqueueReadBuffer((*cl.MemObject)(src), false, 0, bytes, dst, nil)
	if err != nil {
		fmt.Printf("EnqueueReadBuffer failed: %+v \n", err)
		return nil
	}

	// sync copy
	if Synchronous {
		err = cl.WaitForEvents([]*cl.Event{event})
		if err != nil {
			fmt.Printf("Second WaitForEvents in MemCpyDtoH failed: %+v \n", err)
			return nil
		}
	}
	timer.Stop("memcpyDtoH")

	return []*cl.Event{event}
}

func MemCpyHtoD(dst, src unsafe.Pointer, bytes int) []*cl.Event {
	timer.Start("memcpyHtoD")

	// execute
	event, err := ClCmdQueue.EnqueueWriteBuffer((*cl.MemObject)(dst), false, 0, bytes, src, nil)
	if err != nil {
		fmt.Printf("EnqueueWriteBuffer failed: %+v \n", err)
		return nil
	}

	// sync copy
	if Synchronous {
		err = cl.WaitForEvents([]*cl.Event{event})
		if err != nil {
			fmt.Printf("Second WaitForEvents in MemCpyHtoD failed: %+v \n", err)
			return nil
		}
	}
	timer.Stop("memcpyHtoD")

	return []*cl.Event{event}
}

func MemCpy(dst, src unsafe.Pointer, bytes int) []*cl.Event {
	timer.Start("memcpy")

	// execute
	event, err := ClCmdQueue.EnqueueCopyBuffer((*cl.MemObject)(src), (*cl.MemObject)(dst), 0, 0, bytes, nil)
	if err != nil {
		fmt.Printf("EnqueueCopyBuffer failed: %+v \n", err)
		return nil
	}

	// sync copy
	if Synchronous {
		err = cl.WaitForEvents([]*cl.Event{event})
		if err != nil {
			fmt.Printf("First WaitForEvents in MemCpy failed: %+v \n", err)
			return nil
		}
	}
	timer.Stop("memcpy")

	returnList := make([]*cl.Event, 2)
	returnList[0], returnList[1] = event, event
	return returnList
}

// Memset sets the Slice's components to the specified values.
// To be carefully used on unified slice (need sync)
func Memset(s *data.Slice, val ...float64) {
	eventList := make([](*cl.Event), s.NComp())
	err := cl.WaitForEvents(nil)

	if Synchronous { // debug
		for c := range eventList {
			eventList[c] = s.GetEvent(c)
		}
		eventBar, errBar := ClCmdQueue.EnqueueBarrierWithWaitList(eventList)
		errBar = cl.WaitForEvents([](*cl.Event){eventBar})
		if errBar != nil {
			fmt.Printf("First WaitForEvents in MemSet failed: %+v \n", err)
		}
		timer.Start("memset")
	}
	util.Argument(len(val) == s.NComp())
	eventListFill := make([](*cl.Event), len(val))
	for c, v := range val {
		eventListFill[c], err = ClCmdQueue.EnqueueFillBuffer((*cl.MemObject)(s.DevPtr(c)), unsafe.Pointer(&v), SIZEOF_FLOAT64, 0, s.Len()*SIZEOF_FLOAT64, [](*cl.Event){s.GetEvent(c)})
		s.SetEvent(c, eventListFill[c])
		if err != nil {
			fmt.Printf("EnqueueFillBuffer failed: %+v \n", err)
		}
	}
	if Synchronous { //debug
		eventBar, errBar := ClCmdQueue.EnqueueBarrierWithWaitList(eventListFill)
		errBar = cl.WaitForEvents([](*cl.Event){eventBar})
		if errBar != nil {
			fmt.Printf("Second WaitForEvents in MemSet failed: %+v \n", err)
		}
		timer.Stop("memset")
	}
}

// Set all elements of all components to zero.
func Zero(s *data.Slice) {
	Memset(s, make([]float64, s.NComp())...)
}

func SetCell(s *data.Slice, comp int, ix, iy, iz int, value float64) {
	SetElem(s, comp, s.Index(ix, iy, iz), value)
}

func SetElem(s *data.Slice, comp int, index int, value float64) {
	f := value
	event, err := ClCmdQueue.EnqueueWriteBuffer((*cl.MemObject)(s.DevPtr(comp)), false, index*SIZEOF_FLOAT64, SIZEOF_FLOAT64, unsafe.Pointer(&f), [](*cl.Event){s.GetEvent(comp)})
	if err != nil {
		fmt.Printf("EnqueueWriteBuffer failed: %+v \n", err)
		return
	}
	s.SetEvent(comp, event)
	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents in SetElem failed: %+v \n", err)
	}
}

func GetElem(s *data.Slice, comp int, index int) float64 {
	var f float64
	event, err := ClCmdQueue.EnqueueReadBuffer((*cl.MemObject)(s.DevPtr(comp)), false, index*SIZEOF_FLOAT64, SIZEOF_FLOAT64, unsafe.Pointer(&f), [](*cl.Event){s.GetEvent(comp)})
	if err != nil {
		fmt.Printf("EnqueueReadBuffer failed: %+v \n", err)
	}
	s.SetEvent(comp, event)
	err = cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents in GetElem failed: %+v \n", err)
	}
	return f
}

func GetCell(s *data.Slice, comp, ix, iy, iz int) float64 {
	return GetElem(s, comp, s.Index(ix, iy, iz))
}

func updateSlicesWithRdEvent(s []*data.Slice, e *cl.Event) {
	for _, ds := range s {
		if ds != nil {
			for idx := 0; idx < ds.NComp(); idx++ {
				ds.InsertReadEvent(idx, e)
			}
		}
	}
}

func removeRdEventFromSlices(s []*data.Slice, e *cl.Event) {
	for _, ds := range s {
		if ds != nil {
			for idx := 0; idx < ds.NComp(); idx++ {
				ds.RemoveReadEvent(idx, e)
			}
		}
	}
}
