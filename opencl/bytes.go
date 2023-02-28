package opencl

// This file provides GPU byte slices, used to store regions.

import (
	"log"
	"sync"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	util "github.com/seeder-research/uMagNUS/util"
)

// 3D byte slice, used for region lookup.
type Bytes struct {
	Ptr   unsafe.Pointer
	Len   int
	sync.RWMutex
}

// Construct new byte slice with given length,
// initialised to zeros.
func NewBytes(Len int) *Bytes {
	ptr, err1 := ClCtx.CreateEmptyBuffer(cl.MemReadWrite, Len)
	if err1 != nil {
		panic(err1)
	}

	outByte := new(Bytes)
	outByte.Ptr = unsafe.Pointer(ptr)
	outByte.Len = Len

	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		byte_zero__(outByte, &wg)
	} else {
		go func() {
			byte_zero__(outByte, &wg)
		}()
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
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	log.Printf("bytes.zero failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	//var event *cl.Event
	event, err := cmdqueue.EnqueueFillBuffer((*cl.MemObject)(dst.Ptr), unsafe.Pointer(&zeroPattern), 1, 0, dst.Len, nil)
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
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		MemCpyHtoD(dst.Ptr, unsafe.Pointer(&src[0]), dst.Len, &wg)
		dst.Unlock()
	} else {
		go func() {
			MemCpyHtoD(dst.Ptr, unsafe.Pointer(&src[0]), dst.Len, &wg)
			dst.Unlock()
		}()
	}
	wg.Wait()
}

// Copy on device: dst = src.
func (dst *Bytes) Copy(src *Bytes) {
	util.Argument(dst.Len == src.Len)
	dst.Lock()
	src.RLock()
	var wg sync.WaitGroup
	wg.Add(1)
	if Synchronous {
		MemCpy(dst.Ptr, src.Ptr, dst.Len, &wg)
		dst.Unlock()
		src.RUnlock()
	} else {
		go func() {
			MemCpy(dst.Ptr, src.Ptr, dst.Len, &wg)
			dst.Unlock()
			src.RUnlock()
		}()
	}
	wg.Wait()
}

// Copy to host: dst = src.
func (src *Bytes) Download(dst []byte) {
	util.Argument(src.Len == len(dst))
	src.RLock()
	var wg sync.WaitGroup
	wg.Add(1)
	MemCpyDtoH(unsafe.Pointer(&dst[0]), src.Ptr, src.Len, &wg)
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
		go func() {
			bytes_set__(dst, index, value, &wg)
		}()
	}
	wg.Wait()
}

func bytes_set__(dst *Bytes, index int, value byte, wg_ *sync.WaitGroup) {
	dst.Lock()
	defer dst.Unlock()

	src := value

	// Create the command queue to execute the command
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	log.Panicf("MemCpyDoH failed to create command queue: %+v \n", err)
	//	return
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

	//var event *cl.Event
	event, err := cmdqueue.EnqueueWriteBuffer((*cl.MemObject)(dst.Ptr), false, index, 1, unsafe.Pointer(&src), nil)
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
	//cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	//if err != nil {
	//	log.Panicf("MemCpyDoH failed to create command queue: %+v \n", err)
	//	return byte(0)
	//}
	//defer cmdqueue.Release()
	cmdqueue := checkoutQueue()
	defer checkinQueue(cmdqueue)

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
}
