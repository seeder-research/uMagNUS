package opencl

import (
	"fmt"
	"sync"
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

func newSlice(nComp int, size [3]int, memType int8) *data.Slice {
	length := prod(size)
	bytes := length * SIZEOF_FLOAT32
	ptrs := make([]unsafe.Pointer, nComp)
	for c := range ptrs {
		tmp_buf, err := ClCtx.CreateEmptyBuffer(cl.MemReadWrite, bytes)
		if err != nil {
			fmt.Printf("CreateEmptyBuffer failed: %+v \n", err)
		}
		ptrs[c] = unsafe.Pointer(tmp_buf)
	}

	dataPtr := data.SliceFromPtrs(size, memType, ptrs)
	return dataPtr
}

// wrappers for data.EnableGPU arguments

func memFree(ptr unsafe.Pointer) {
	if ptr != nil {
		buf := (*cl.MemObject)(ptr)
		buf.Release()
	}
}

func MemCpyDtoH(dst, src unsafe.Pointer, bytes int, wg *sync.WaitGroup) {
	// execute
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)
	//event, err := ClD2HQueue.EnqueueReadBuffer((*cl.MemObject)(src), false, 0, bytes, dst, nil)
	event, err := cmdqueue.EnqueueReadBuffer((*cl.MemObject)(src), false, 0, bytes, dst, nil)
	if err != nil {
		fmt.Printf("EnqueueReadBuffer in MemCpyDtoH failed: %+v \n", err)
	}
	wg.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in MemCpyDtoH: %+v \n", err)
	}
}

func MemCpyHtoD(dst, src unsafe.Pointer, bytes int, wg *sync.WaitGroup) {
	// execute
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)
	//event, err := ClH2DQueue.EnqueueWriteBuffer((*cl.MemObject)(dst), false, 0, bytes, src, nil)
	event, err := cmdqueue.EnqueueWriteBuffer((*cl.MemObject)(dst), false, 0, bytes, src, nil)
	if err != nil {
		fmt.Printf("EnqueueWriteBuffer in MemCpyHtoD failed: %+v \n", err)
	}
	wg.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in MemCpyHtoD: %+v \n", err)
	}
}

func MemCpy(dst, src unsafe.Pointer, bytes int, wg *sync.WaitGroup) {
	// execute
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)
	//event, err := ClCmdQueue.EnqueueCopyBuffer((*cl.MemObject)(src), (*cl.MemObject)(dst), 0, 0, bytes, nil)
	event, err := cmdqueue.EnqueueCopyBuffer((*cl.MemObject)(src), (*cl.MemObject)(dst), 0, 0, bytes, nil)
	if err != nil {
		fmt.Printf("EnqueueCopyBuffer failed: %+v \n", err)
	}
	wg.Done()

	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents failed in MemCpy: %+v \n", err)
	}
}

// Memset sets the Slice's components to the specified values.
// To be carefully used on unified slice (need sync)
func Memset(s *data.Slice, val ...float32) {
	timer.Start("memset")

	util.Argument(len(val) == s.NComp())
	//eventListFill := make([](*cl.Event), len(val))

	var wg_ sync.WaitGroup
	for c, v := range val {
		wg_.Add(1)
		if Synchronous {
			memset_func(s, c, &v, &wg_)
		} else {
			idx := c
			v0 := v
			go func() {
				memset_func(s, idx, &v0, &wg_)
			}()
		}
	}
	wg_.Wait()

	timer.Stop("memset")
}

func memset_func(s *data.Slice, comp int, v *float32, wg__ *sync.WaitGroup) {
	s.Lock(comp)
	defer s.Unlock(comp)

	//var err error
	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//defer cmdqueue.Release()
	//if err != nil {
	//	fmt.Printf("MemSet failed to create command queue: %+v \n", err)
	//	return
	//}
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	//var event *cl.Event
	event, err := cmdqueue.EnqueueFillBuffer((*cl.MemObject)(s.DevPtr(comp)), unsafe.Pointer(v), SIZEOF_FLOAT32, 0, s.Len()*SIZEOF_FLOAT32, nil)
	wg__.Done()
	if err != nil {
		fmt.Printf("MemSet failed to enqueue command: %+v \n", err)
		cmdqueue.Release()
		return
	}

	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("Wait for command to complete in MemCpy failed: %+v \n", err)
		return
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
	s.Lock(comp)
	defer s.Unlock(comp)

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		setelem__(s, comp, index, value, &wg)
	} else {
		go func() {
			setelem__(s, comp, index, value, &wg)
		}()
	}
	wg.Wait()
}

func setelem__(s *data.Slice, comp, index int, value float32, wg__ *sync.WaitGroup) {
	s.Lock(comp)
	defer s.Unlock(comp)

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//defer cmdqueue.Release()
	//if err != nil {
	//	fmt.Printf("SetElem failed to create command queue: %+v \n", err)
	//	return
	//}
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)
	//var event *cl.Event
	f := value
	event, err := cmdqueue.EnqueueWriteBuffer((*cl.MemObject)(s.DevPtr(comp)), false, index*SIZEOF_FLOAT32, SIZEOF_FLOAT32, unsafe.Pointer(&f), nil)
	wg__.Done()
	if err != nil {
		fmt.Printf("EnqueueWriteBuffer failed: %+v \n", err)
		return
	}

	if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("Wait for command to complete in SetElem failed: %+v \n", err)
	}

}

func GetElem(s *data.Slice, comp int, index int) float32 {
	s.RLock(comp)
	defer s.RUnlock(comp)

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//defer cmdqueue.Release()
	//if err != nil {
	//	fmt.Printf("GetElem failed to create command queue: %+v \n", err)
	//	return -1.0
	//}
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)
	//var event *cl.Event
	var f float32
	event, err := cmdqueue.EnqueueReadBuffer((*cl.MemObject)(s.DevPtr(comp)), false, index*SIZEOF_FLOAT32, SIZEOF_FLOAT32, unsafe.Pointer(&f), nil)
	if err != nil {
		fmt.Printf("EnqueueReadBuffer failed: %+v \n", err)
		return -1.0
	}

	if err = cl.WaitForEvents([]*cl.Event{event});err != nil {
		fmt.Printf("WaitForEvents in GetElem failed: %+v \n", err)
	}

	return f
}

func GetCell(s *data.Slice, comp, ix, iy, iz int) float32 {
	return GetElem(s, comp, s.Index(ix, iy, iz))
}
