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

func newSlice(nComp int, size [3]int, memType int8) *data.Slice {
	length := prod(size)
	bytes := length * SIZEOF_FLOAT32
	ptrs := make([]unsafe.Pointer, nComp)
	initVal := float32(0.0)
	for c := range ptrs {
		tmp_buf, err := ClCtx.CreateEmptyBuffer(cl.MemReadWrite, bytes)
		if err != nil {
			fmt.Printf("CreateEmptyBuffer failed: %+v \n", err)
		}
		ptrs[c] = unsafe.Pointer(tmp_buf)
		fillWait, err := H2DQueue.EnqueueFillBuffer(tmp_buf, unsafe.Pointer(&initVal), SIZEOF_FLOAT32, 0, bytes, nil)
		if err != nil {
			fmt.Printf("CreateEmptyBuffer failed: %+v \n", err)
		}
		if Synchronous {
			if err = cl.WaitForEvents([]*cl.Event{fillWait}); err != nil {
				fmt.Printf("Wait for EnqueueFillBuffer failed: %+v \n", err)
			}
		}
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

func MemCpyDtoH(dst, src unsafe.Pointer, bytes int) {
	var err error
	var event *cl.Event

	// debug
	if Synchronous {
		// sync previous kernels
		SyncQueues([]*cl.CommandQueue{D2HQueue}, append(ClCmdQueue, H2DQueue))
		if err = D2HQueue.Finish(); err != nil {
			fmt.Printf("failed to wait for queue to finish: %+v \n", err)
		}
		timer.Start("memcpyDtoH")
	}

	// execute
	if event, err = D2HQueue.EnqueueReadBuffer((*cl.MemObject)(src), false, 0, bytes, dst, nil); err != nil {
		fmt.Printf("EnqueueReadBuffer failed: %+v \n", err)
	}

	// debug
	if Synchronous {
		// sync copy
		if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("Second WaitForEvents in memcpyDtoH failed: %+v \n", err)
		}
		timer.Stop("memcpyDtoH")
	}
}

func MemCpyHtoD(dst, src unsafe.Pointer, bytes int) {
	var err error
	var event *cl.Event

	// debug
	if Synchronous {
		// sync previous kernels
		SyncQueues([]*cl.CommandQueue{H2DQueue}, append(ClCmdQueue, D2HQueue))
		if err = H2DQueue.Finish(); err != nil {
			fmt.Printf("failed to wait for queue to finish: %+v \n", err)
		}
		timer.Start("memcpyHtoD")
	}

	// execute
	event, err = H2DQueue.EnqueueWriteBuffer((*cl.MemObject)(dst), false, 0, bytes, src, nil)
	if err != nil {
		fmt.Printf("EnqueueWriteBuffer failed: %+v \n", err)
	}

	if Synchronous {
		// sync copy
		if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("WaitForEvents in memcpyHtoD failed: %+v \n", err)
		}
		timer.Stop("memcpyHtoD")
	}
}

func MemCpy(dst, src unsafe.Pointer, bytes int) {
	var err error

	// debug
	if Synchronous {
		// sync kernels
		var queues []*cl.CommandQueue = ClCmdQueue[1:(NumQueues - 1)]
		SyncQueues([]*cl.CommandQueue{ClCmdQueue[0]}, append(append(queues, D2HQueue), H2DQueue))
		if err = ClCmdQueue[0].Finish(); err != nil {
			fmt.Printf("failed to wait for queue to finish: %+v \n", err)
		}
		timer.Start("memcpy")
	}

	// execute
	event, err := ClCmdQueue[0].EnqueueCopyBuffer((*cl.MemObject)(src), (*cl.MemObject)(dst), 0, 0, bytes, nil)
	if err != nil {
		fmt.Printf("EnqueueCopyBuffer failed: %+v \n", err)
	}

	if Synchronous {
		// sync copy
		if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
			fmt.Printf("First WaitForEvents in memcpy failed: %+v \n", err)
		}
		timer.Stop("memcpy")
	}
}

// Memset sets the Slice's components to the specified values.
// To be carefully used on unified slice (need sync)
func Memset(s *data.Slice, val ...float32) {
	var err error

	// debug
	if Synchronous {
		SyncQueues([]*cl.CommandQueue{H2DQueue}, append(ClCmdQueue, D2HQueue))
		if err = H2DQueue.Finish(); err != nil {
			fmt.Printf("failed to wait for queue to finish in beginning of memset: %+v \n", err)
		}
		timer.Start("memset")
	}

	util.Argument(len(val) == s.NComp())

	for c, v := range val {
		event, err := H2DQueue.EnqueueFillBuffer((*cl.MemObject)(s.DevPtr(c)), unsafe.Pointer(&v), SIZEOF_FLOAT32, 0, s.Len()*SIZEOF_FLOAT32, nil)
		if err != nil {
			fmt.Printf("EnqueueFillBuffer failed: %+v \n", err)
		}
		if Synchronous { // debug
			if err = cl.WaitForEvents([]*cl.Event{event}); err != nil {
				fmt.Printf("Second WaitForEvents in memset failed: %+v \n", err)
			}
		}
	}

	// debug
	if Synchronous {
		timer.Stop("memset")
	}

	SyncQueues([]*cl.CommandQueue{ClCmdQueue[0]}, []*cl.CommandQueue{H2DQueue})
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
	event, err := H2DQueue.EnqueueWriteBuffer((*cl.MemObject)(s.DevPtr(comp)), false, index*SIZEOF_FLOAT32, SIZEOF_FLOAT32, unsafe.Pointer(&f), nil)
	if err != nil {
		fmt.Printf("EnqueueWriteBuffer failed: %+v \n", err)
		return
	}
	SyncQueues([]*cl.CommandQueue{ClCmdQueue[0]}, []*cl.CommandQueue{H2DQueue})
	if Synchronous {
		if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents in SetElem failed: %+v \n", err)
		}
	}
}

func GetElem(s *data.Slice, comp int, index int) float32 {
	var f float32
	// sync previous kernels
	SyncQueues([]*cl.CommandQueue{D2HQueue}, append(ClCmdQueue, H2DQueue))
	event, err := D2HQueue.EnqueueReadBuffer((*cl.MemObject)(s.DevPtr(comp)), false, index*SIZEOF_FLOAT32, SIZEOF_FLOAT32, unsafe.Pointer(&f), nil)
	if err != nil {
		fmt.Printf("EnqueueReadBuffer failed: %+v \n", err)
	}
	// Must sync
	if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
		fmt.Printf("WaitForEvents in GetElem failed: %+v \n", err)
	}
	return f
}

func GetCell(s *data.Slice, comp, ix, iy, iz int) float32 {
	return GetElem(s, comp, s.Index(ix, iy, iz))
}
