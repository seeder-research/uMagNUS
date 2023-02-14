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
	fillWait := make([]*cl.Event, nComp)
	for c := range ptrs {
		tmp_buf, err := ClCtx.CreateEmptyBuffer(cl.MemReadWrite, bytes)
		if err != nil {
			fmt.Printf("CreateEmptyBuffer failed: %+v \n", err)
		}
		ptrs[c] = unsafe.Pointer(tmp_buf)
		fillWait[c] = nil
		//fillWait[c], err = ClCmdQueue.EnqueueFillBuffer(tmp_buf, unsafe.Pointer(&initVal), SIZEOF_FLOAT32, 0, bytes, nil)
		//if err != nil {
		//	fmt.Printf("CreateEmptyBuffer failed: %+v \n", err)
		//}
		//if err = cl.WaitForEvents([]*cl.Event{fillWait[c]}); err != nil {
		//	fmt.Printf("Wait for EnqueueFillBuffer failed: %+v \n", err)
		//}
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
	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("MemCpyDoH failed to create command queue: %+v \n", err)
		return nil
	}
	timer.Start("memcpyDtoH")

	// execute
	var event *cl.Event
	event, err = cmdqueue.EnqueueReadBuffer((*cl.MemObject)(src), false, 0, bytes, dst, nil)
	if err != nil {
		fmt.Printf("EnqueueReadBuffer in MemCpyDtoH failed: %+v \n", err)
		return nil
	}

	// sync copy
	if Synchronous {
		err = cmdqueue.Finish()
		if err != nil {
			fmt.Printf("Wait for command to complete in MemCpyDtoH failed: %+v \n", err)
			cmdqueue.Release()
			return nil
		}
	}
	timer.Stop("memcpyDtoH")

	cmdqueue.Release()
	return []*cl.Event{event}
}

func MemCpyHtoD(dst, src unsafe.Pointer, bytes int) []*cl.Event {
	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("MemCpyDoH failed to create command queue: %+v \n", err)
		return nil
	}
	timer.Start("memcpyHtoD")

	// execute
	var event *cl.Event
	event, err = cmdqueue.EnqueueWriteBuffer((*cl.MemObject)(dst), false, 0, bytes, src, nil)
	if err != nil {
		fmt.Printf("EnqueueWriteBuffer in MemCpyHtoD failed: %+v \n", err)
		return nil
	}

	// sync copy
	if Synchronous {
		err = cmdqueue.Finish()
		if err != nil {
			fmt.Printf("Wait for command to complete in MemCpyHtoD failed: %+v \n", err)
			cmdqueue.Release()
			return nil
		}
	}

	timer.Stop("memcpyHtoD")

	cmdqueue.Release()
	return []*cl.Event{event}
}

func MemCpy(dst, src unsafe.Pointer, bytes int) []*cl.Event {
	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("MemCpyDoH failed to create command queue: %+v \n", err)
		return nil
	}
	timer.Start("memcpy")

	// execute
	var event *cl.Event
	event, err = cmdqueue.EnqueueCopyBuffer((*cl.MemObject)(src), (*cl.MemObject)(dst), 0, 0, bytes, nil)
	if err != nil {
		fmt.Printf("EnqueueCopyBuffer failed: %+v \n", err)
		return nil
	}

	// sync copy
	if Synchronous {
		err = cmdqueue.Finish()
		if err != nil {
			fmt.Printf("Wait for command to complete in MemCpy failed: %+v \n", err)
			cmdqueue.Finish()
			return nil
		}
	}
	timer.Stop("memcpy")

	cmdqueue.Release()
	returnList := make([]*cl.Event, 2)
	returnList[0], returnList[1] = event, event
	return returnList
}

// Memset sets the Slice's components to the specified values.
// To be carefully used on unified slice (need sync)
func Memset(s *data.Slice, val ...float32) {
	timer.Start("memset")

	util.Argument(len(val) == s.NComp())
	eventListFill := make([](*cl.Event), len(val))

	var wg_ sync.WaitGroup
	for c, v := range val {
		wg_.Add(1)
		if Synchronous {
			memset_func(s, c, &v, &eventListFill, wg_)
		} else {
			go memset_func(s, c, &v, &eventListFill, wg_)
		}
	}
	wg_.Wait()

	timer.Stop("memset")
}

func memset_func(s *data.Slice, comp int, v *float32, ev *[]*cl.Event, wg__ sync.WaitGroup) {
	s.ptrs[comp].Lock()
	defer s.ptrs[comp].Unlock()

	var err error
	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("MemSet failed to create command queue: %+v \n", err)
		return nil
	}

	var evt *cl.Event
	evt, err = cmdqueue.EnqueueFillBuffer((*cl.MemObject)(s.DevPtr(comp)), unsafe.Pointer(v), SIZEOF_FLOAT32, 0, s.Len()*SIZEOF_FLOAT32, [](*cl.Event){s.GetEvent(comp)})
	wg__.Done()
	if err != nil {
		fmt.Printf("MemSet failed to enqueue command: %+v \n", err)
		cmdqueue.Release()
		return nil
	}

	s.SetEvent(comp, evt)
	ev[comp] = evt

	err = cmdqueue.Finish()
	if err != nil {
		fmt.Printf("Wait for command to complete in MemCpy failed: %+v \n", err)
		cmdqueue.Release()
		return nil
	}
	cmdqueue.Release()
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
	if s.ptrs == nil {
		return
	}
	if s.ptrs[comp] == nil {
		return
	}

	var wg sync.WaitGroup()
	wg.Add(1)
	if Synchronous {
		setelem__(s, comp, index, value, wg)
	} else
		go setelem__(s, comp, index, value, wg)
		wg.Wait()
	}
}

func setelem__(s. *data.Slice, comp int, index int, value float32, wg__ sync.WaitGroup) {
	s.ptrs[comp].Lock()
	defer s.ptrs[comp].Unlock()

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("SetElem failed to create command queue: %+v \n", err)
		return nil
	}
	var event *cl.Event
	event, err := cmdqueue.EnqueueWriteBuffer((*cl.MemObject)(s.DevPtr(comp)), false, index*SIZEOF_FLOAT32, SIZEOF_FLOAT32, unsafe.Pointer(&f), [](*cl.Event){s.GetEvent(comp)})
	wg__.Done()
	if err != nil {
		fmt.Printf("EnqueueWriteBuffer failed: %+v \n", err)
		cmdqueue.Release()
		return
	}

	s.SetEvent(comp, event)
	err = cmdqueue.Finish()
	if err != nil {
		fmt.Printf("Wait for command to complete in SetElem failed: %+v \n", err)
	}

	cmdqueue.Release()
	return
}

func GetElem(s *data.Slice, comp int, index int) float32 {
	s.ptrs[comp].RLock()
	defer s.ptrs[comp].RUnlock()

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("GetElem failed to create command queue: %+v \n", err)
		return nil
	}
	var event *cl.Event
	var f float32
	event, err := cmdqueue.EnqueueReadBuffer((*cl.MemObject)(s.DevPtr(comp)), false, index*SIZEOF_FLOAT32, SIZEOF_FLOAT32, unsafe.Pointer(&f), [](*cl.Event){s.GetEvent(comp)})
	if err != nil {
		fmt.Printf("EnqueueReadBuffer failed: %+v \n", err)
		cmdqueue.Release()
		return
	}

	err = cmdqueue.Finish()
	if err != nil {
		fmt.Printf("Wait for command to complete in GetElem failed: %+v \n", err)
	}

	cmdqueue.Release()
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
