package oclRAND

import (
	"github.com/seeder-research/uMagNUS/cl"
	"log"
	"unsafe"
)

type config struct {
	Grid, Block []int
}

var (
	Synchronous bool
	ClCmdQueue  *cl.CommandQueue
	KernList    = map[string]*cl.Kernel{} // Store pointers to all compiled kernels
)

func Init(q *cl.CommandQueue, synch bool, kList map[string]*cl.Kernel) {
	Synchronous = synch
	ClCmdQueue = q
	KernList = kList
}

func LaunchKernel(kernname string, gridDim, workDim []int, queue *cl.CommandQueue, events []*cl.Event) *cl.Event {
	if KernList[kernname] == nil {
		log.Panic("Kernel " + kernname + " does not exist!")
		return nil
	}
	KernEvent, err := queue.EnqueueNDRangeKernel(KernList[kernname], nil, gridDim, workDim, events)
	if err != nil {
		log.Fatal(err)
		return nil
	} else {
		return KernEvent
	}
}

func SetKernelArgWrapper(kernname string, index int, arg interface{}) {
	if KernList[kernname] == nil {
		log.Panic("Kernel " + kernname + " does not exist!")
	}
	switch val := arg.(type) {
	default:
		if err := KernList[kernname].SetArg(index, val); err != nil {
			log.Fatal(err)
		}
	case unsafe.Pointer:
		memBufHandle, flag := arg.(unsafe.Pointer)
		if memBufHandle == unsafe.Pointer(uintptr(0)) {
			if err := KernList[kernname].SetArgUnsafe(index, 8, memBufHandle); err != nil {
				log.Fatal(err)
			}
		} else {
			if flag {
				if err := KernList[kernname].SetArg(index, (*cl.MemObject)(memBufHandle)); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal("Unable to change argument type to *cl.MemObject")
			}
		}
	case int:
		if err := KernList[kernname].SetArg(index, (int32)(val)); err != nil {
			log.Fatal(err)
		}
	}
}
