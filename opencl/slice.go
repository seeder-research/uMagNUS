package opencl

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	timer "github.com/seeder-research/uMagNUS/timer"
	util "github.com/seeder-research/uMagNUS/util"
)

// Make a GPU Slice with nComp components each of size length.
func NewSlice(nComp int, size [3]int) *data.Slice {
	return newSlice(nComp, size, data.GPUMemory)
}

// Make a GPU Slice with nComp components each of size length.
//func NewUnifiedSlice(nComp int, m *data.Mesh) *data.Slice {
//	return newSlice(nComp, m, cu.MemAllocHost, data.UnifiedMemory)
//}

func newSlice(nComp int, size [3]int, memType int8) *data.Slice {
	length := prod(size)
	bytes := length * SIZEOF_FLOAT32
	ptrs := make([]unsafe.Pointer, nComp)
	initVal := float32(0.0)
	fillWait := make([]*cl.Event, nComp)
	for c := range ptrs {
		tmp_buf, err := ClCtx.CreateEmptyBuffer(cl.MemReadWrite, bytes)
		if err != nil {
			fmt.Printf("CreateEmptyBuffer failed: %+v \n", err)
		}
		ptrs[c] = unsafe.Pointer(tmp_buf)
		fillWait[c], err = ClCmdQueue.EnqueueFillBuffer(tmp_buf, unsafe.Pointer(&initVal), SIZEOF_FLOAT32, 0, bytes, nil)
		if err != nil {
			fmt.Printf("CreateEmptyBuffer failed: %+v \n", err)
		}
		err = cl.WaitForEvents([]*cl.Event{fillWait[c]})
		if err != nil {
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
	// sync previous kernels
	eventList := make([](*cl.Event), 1)
	waitList, err := ClCmdQueue.EnqueueBarrierWithWaitList(nil)
	if err != nil {
		fmt.Printf("EnqueueBarrierWithWaitList failed: %+v \n", err)
		return nil
	}
	eventList[0] = waitList
	err = cl.WaitForEvents(eventList)
	if err != nil {
		fmt.Printf("First WaitForEvents in MemCpyDtoH failed: %+v \n", err)
		return nil
	}
	timer.Start("memcpyDtoH")

	// execute
	eventList[0], err = ClCmdQueue.EnqueueReadBuffer((*cl.MemObject)(src), false, 0, bytes, dst, nil)
	if err != nil {
		fmt.Printf("EnqueueReadBuffer failed: %+v \n", err)
		return nil
	}

	// sync copy
	err = cl.WaitForEvents(eventList)
	timer.Stop("memcpyDtoH")
	if err != nil {
		fmt.Printf("Second WaitForEvents in MemCpyDtoH failed: %+v \n", err)
		return nil
	}

	return eventList
}

func MemCpyHtoD(dst, src unsafe.Pointer, bytes int) []*cl.Event {
	// sync previous kernels
	eventList := make([](*cl.Event), 1)
	waitList, err := ClCmdQueue.EnqueueBarrierWithWaitList(nil)
	if err != nil {
		fmt.Printf("EnqueueBarrierWithWaitList failed: %+v \n", err)
		return nil
	}
	eventList[0] = waitList
	err = cl.WaitForEvents(eventList)
	if err != nil {
		fmt.Printf("First WaitForEvents in MemCpyHtoD failed: %+v \n", err)
		return nil
	}
	timer.Start("memcpyHtoD")

	// execute
	eventList[0], err = ClCmdQueue.EnqueueWriteBuffer((*cl.MemObject)(dst), false, 0, bytes, src, nil)
	if err != nil {
		fmt.Printf("EnqueueWriteBuffer failed: %+v \n", err)
		return nil
	}

	// sync copy
	err = cl.WaitForEvents(eventList)
	timer.Stop("memcpyHtoD")
	if err != nil {
		fmt.Printf("Second WaitForEvents in MemCpyHtoD failed: %+v \n", err)
		return nil
	}

	return eventList
}

func MemCpy(dst, src unsafe.Pointer, bytes int) []*cl.Event {
	// sync kernels
	eventList := make([](*cl.Event), 1)
	waitList, err := ClCmdQueue.EnqueueBarrierWithWaitList(nil)
	if err != nil {
		fmt.Printf("EnqueueBarrierWithWaitList failed: %+v \n", err)
		return nil
	}
	eventList[0] = waitList
	err = cl.WaitForEvents(eventList)
	if err != nil {
		fmt.Printf("First WaitForEvents in MemCpy failed: %+v \n", err)
		return nil
	}
	timer.Start("memcpy")

	// execute
	eventList[0], err = ClCmdQueue.EnqueueCopyBuffer((*cl.MemObject)(src), (*cl.MemObject)(dst), 0, 0, bytes, nil)
	if err != nil {
		fmt.Printf("EnqueueCopyBuffer failed: %+v \n", err)
		return nil
	}

	// sync copy
	err = cl.WaitForEvents(eventList)
	timer.Stop("memcpy")
	if err != nil {
		fmt.Printf("First WaitForEvents in MemCpy failed: %+v \n", err)
		return nil
	}

	returnList := make([]*cl.Event, 2)
	returnList[0], returnList[1] = eventList[0], eventList[0]
	return returnList
}

// Memset sets the Slice's components to the specified values.
// To be carefully used on unified slice (need sync)
func Memset(s *data.Slice, val ...float32) {
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
		eventListFill[c], err = ClCmdQueue.EnqueueFillBuffer((*cl.MemObject)(s.DevPtr(c)), unsafe.Pointer(&v), SIZEOF_FLOAT32, 0, s.Len()*SIZEOF_FLOAT32, [](*cl.Event){s.GetEvent(c)})
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
	Memset(s, make([]float32, s.NComp())...)
}

func SetCell(s *data.Slice, comp int, ix, iy, iz int, value float32) {
	SetElem(s, comp, s.Index(ix, iy, iz), value)
}

func SetElem(s *data.Slice, comp int, index int, value float32) {
	f := value
	event, err := ClCmdQueue.EnqueueWriteBuffer((*cl.MemObject)(s.DevPtr(comp)), false, index*SIZEOF_FLOAT32, SIZEOF_FLOAT32, unsafe.Pointer(&f), [](*cl.Event){s.GetEvent(comp)})
	if err != nil {
		fmt.Printf("EnqueueWriteBuffer failed: %+v \n", err)
		return
	}
	s.SetEvent(comp, event)
}

func GetElem(s *data.Slice, comp int, index int) float32 {
	var f float32
	event, err := ClCmdQueue.EnqueueReadBuffer((*cl.MemObject)(s.DevPtr(comp)), false, index*SIZEOF_FLOAT32, SIZEOF_FLOAT32, unsafe.Pointer(&f), [](*cl.Event){s.GetEvent(comp)})
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

func GetCell(s *data.Slice, comp, ix, iy, iz int) float32 {
	return GetElem(s, comp, s.Index(ix, iy, iz))
}
