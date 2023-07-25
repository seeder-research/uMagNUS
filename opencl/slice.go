package opencl

import (
	"fmt"
	"sync"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	qm "github.com/seeder-research/uMagNUS/queuemanager"
	timer "github.com/seeder-research/uMagNUS/timer"
	util "github.com/seeder-research/uMagNUS/util"
)

// Make a GPU Slice with nComp components each of size length.
func NewSlice(nComp int, size [3]int) *data.Slice {
	return newSlice(nComp, size, data.GPUMemory)
}

func newSlice(nComp int, size [3]int, memType int8) *data.Slice {
	length := prod(size)
	bytes := length * SIZEOF_FLOAT32
	ptrs := make([]unsafe.Pointer, nComp)
	initVal := float32(0.0)
	fillWait := make([]*cl.Event, nComp)
	var newSliceSyncWaitGroup sync.WaitGroup

	// Synchronize command queues so that all kernel launches occur in-sequence
	if Synchronous {
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}

	// Zero buffers in slice
	for c := range ptrs {
		tmp_buf, err := ClCtx.CreateEmptyBuffer(cl.MemReadWrite, bytes)
		if err != nil {
			fmt.Printf("CreateEmptyBuffer failed: %+v \n", err)
		}
		ptrs[c] = unsafe.Pointer(tmp_buf)

		// Checkout command queue from pool and launch kernel
		tmpQueue := qm.CheckoutQueue(CmdQueuePool, &newSliceSyncWaitGroup)
		fillWait[c], err = tmpQueue.EnqueueFillBuffer(tmp_buf, unsafe.Pointer(&initVal), SIZEOF_FLOAT32, 0, bytes, nil)
		if err != nil {
			fmt.Printf("CreateEmptyBuffer failed: %+v \n", err)
		}

		// Check command queue back into main pool
		qwg := qm.NewQueueWaitGroup(tmpQueue, &newSliceSyncWaitGroup)
		ReturnQueuePool <- qwg

		// Synchronize command queues so that all kernel launches occur in-sequence
		if Synchronous {
			newSliceSyncWaitGroup.Done()
		}
	}

	dataPtr := data.SliceFromPtrs(size, memType, ptrs)
	if Synchronous == false {
		newSliceSyncWaitGroup.Done()
	}
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
	// Synchronize command queues so that all kernel launches occur in-sequence
	if Synchronous {
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}
	var memCpyDtoHSyncWaitGroup sync.WaitGroup
	timer.Start("memcpyDtoH")

	// Checkout command queue from pool and execute
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &memCpyDtoHSyncWaitGroup)
	eventList[0], err = tmpQueue.EnqueueReadBuffer((*cl.MemObject)(src), false, 0, bytes, dst, nil)
	if err != nil {
		fmt.Printf("EnqueueReadBuffer failed: %+v \n", err)
		return nil
	}

	// sync copy
	qwg := qm.NewQueueWaitGroup(tmpQueue, &memCpyDtoHSyncWaitGroup)
	ReturnQueuePool <- qwg
	err = cl.WaitForEvents(eventList)
	timer.Stop("memcpyDtoH")
	if err != nil {
		fmt.Printf("Second WaitForEvents in MemCpyDtoH failed: %+v \n", err)
		return nil
	}

	return eventList
}

func MemCpyHtoD(dst, src unsafe.Pointer, bytes int) []*cl.Event {
	// Synchronize command queues so that all kernel launches occur in-sequence
	if Synchronous {
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}
	var memCpyHtoDSyncWaitGroup sync.WaitGroup
	timer.Start("memcpyHtoD")

	// Checkout command queue from pool and execute
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &memCpyHtoDSyncWaitGroup)
	eventList[0], err = tmpQueue.EnqueueWriteBuffer((*cl.MemObject)(dst), false, 0, bytes, src, nil)
	if err != nil {
		fmt.Printf("EnqueueWriteBuffer failed: %+v \n", err)
		return nil
	}

	// sync copy
	qwg := qm.NewQueueWaitGroup(tmpQueue, &memCpyHtoDSyncWaitGroup)
	ReturnQueuePool <- qwg
	err = cl.WaitForEvents(eventList)
	timer.Stop("memcpyHtoD")
	if err != nil {
		fmt.Printf("Second WaitForEvents in MemCpyHtoD failed: %+v \n", err)
		return nil
	}

	return eventList
}

func MemCpy(dst, src unsafe.Pointer, bytes int) []*cl.Event {
	// Synchronize command queues so that all kernel launches occur in-sequence
	if Synchronous {
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}
	var memCpySyncWaitGroup sync.WaitGroup
	timer.Start("memcpy")

	// Checkout command queue from pool and execute
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &memCpySyncWaitGroup)
	eventList[0], err = tmpQueue.EnqueueCopyBuffer((*cl.MemObject)(src), (*cl.MemObject)(dst), 0, 0, bytes, nil)
	if err != nil {
		fmt.Printf("EnqueueCopyBuffer failed: %+v \n", err)
		return nil
	}

	// sync copy
	qwg := qm.NewQueueWaitGroup(tmpQueue, &memCpyHtoDSyncWaitGroup)
	ReturnQueuePool <- qwg
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
	if Synchronous {
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}
	var memSetSyncWaitGroup sync.WaitGroup

	timer.Start("memset")
	util.Argument(len(val) == s.NComp())
	eventListFill := make([](*cl.Event), len(val))
	for c, v := range val {
		// Checkout command queue from pool and execute
		tmpQueue := qm.CheckoutQueue(CmdQueuePool, &memSetSyncWaitGroup)
		eventListFill[c], err = tmpQueue.EnqueueFillBuffer((*cl.MemObject)(s.DevPtr(c)), unsafe.Pointer(&v), SIZEOF_FLOAT32, 0, s.Len()*SIZEOF_FLOAT32, [](*cl.Event){s.GetEvent(c)})
		s.SetEvent(c, eventListFill[c])
		if err != nil {
			fmt.Printf("EnqueueFillBuffer failed: %+v \n", err)
		}
		qwg := qm.NewQueueWaitGroup(tmpQueue, &memCpyHtoDSyncWaitGroup)
		ReturnQueuePool <- qwg

		// Synchronize command queues so that all kernel launches occur in-sequence
		if Synchronous { //debug
			memSetSyncWaitGroup.Wait()
		}
	}

	// Synchronize command queues
	memSetSyncWaitGroup.Wait()
	timer.Stop("memset")
}

// Set all elements of all components to zero.
func Zero(s *data.Slice) {
	Memset(s, make([]float32, s.NComp())...)
}

func SetCell(s *data.Slice, comp int, ix, iy, iz int, value float32) {
	SetElem(s, comp, s.Index(ix, iy, iz), value)
}

func SetElem(s *data.Slice, comp int, index int, value float32) {
	if Synchronous { // debug
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}
	var setElemSyncWaitGroup sync.WaitGroup
	f := value

	// Checkout command queue from pool and execute
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &setElemSyncWaitGroup)
	event, err := tmpQueue.EnqueueWriteBuffer((*cl.MemObject)(s.DevPtr(comp)), false, index*SIZEOF_FLOAT32, SIZEOF_FLOAT32, unsafe.Pointer(&f), [](*cl.Event){s.GetEvent(comp)})
	if err != nil {
		fmt.Printf("EnqueueWriteBuffer failed: %+v \n", err)
		return
	}
	s.SetEvent(comp, event)

	// Checkin command queue post execution
	qwg := qm.NewQueueWaitGroup(tmpQueue, &memSetSyncWaitGroup)
	ReturnQueuePool <- qwg
	setElemSyncWaitGroup.Wait()
}

func GetElem(s *data.Slice, comp int, index int) float32 {
	if Synchronous { // debug
		for len(CmdQueuePool) < QueuePoolSz {
		}
	}
	var getElemSyncWaitGroup sync.WaitGroup
	var f float32

	// Checkout command queue from pool and execute
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, &getElemSyncWaitGroup)
	event, err := tmpQueue.EnqueueReadBuffer((*cl.MemObject)(s.DevPtr(comp)), false, index*SIZEOF_FLOAT32, SIZEOF_FLOAT32, unsafe.Pointer(&f), [](*cl.Event){s.GetEvent(comp)})
	if err != nil {
		fmt.Printf("EnqueueReadBuffer failed: %+v \n", err)
	}
	s.SetEvent(comp, event)

	// Checkin command queue post execution
	qwg := qm.NewQueueWaitGroup(tmpQueue, &getElemSyncWaitGroup)
	ReturnQueuePool <- qwg

	err = cl.WaitForEvents([](*cl.Event){event})
	if err != nil {
		fmt.Printf("WaitForEvents in GetElem failed: %+v \n", err)
	}
	return f
}

func GetCell(s *data.Slice, comp, ix, iy, iz int) float32 {
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
