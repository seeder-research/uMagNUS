package opencl64

import (
	"log"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
	timer "github.com/seeder-research/uMagNUS/timer"
	util "github.com/seeder-research/uMagNUS/util"
)

// 3D single-precision real-to-complex FFT plan.
type fft3DR2CPlan struct {
	fftplan
	size [3]int
}

// 3D single-precision real-to-complex FFT plan.
func newFFT3DR2C(Nx, Ny, Nz int) fft3DR2CPlan {
	handle := cl.NewVkFFTPlanDouble(ClCtx)
	handle.VkFFTSetFFTPlanSize([]int{Nx, Ny, Nz})

	return fft3DR2CPlan{fftplan{handle}, [3]int{Nx, Ny, Nz}}
}

// Execute the FFT plan, asynchronous.
// src and dst are 3D arrays stored 1D arrays.
func (p *fft3DR2CPlan) ExecAsync(src, dst *data.Slice) error {
	if Synchronous {
		ClCmdQueue.Finish()
		timer.Start("fft")
	}
	util.Argument(src.NComp() == 1 && dst.NComp() == 1)
	oksrclen := p.InputLen()
	if src.Len() != oksrclen {
		log.Panicf("fft size mismatch: expecting src len %v, got %v", oksrclen, src.Len())
	}
	okdstlen := p.OutputLen()
	if dst.Len() != okdstlen {
		log.Panicf("fft size mismatch: expecting dst len %v, got %v", okdstlen, dst.Len())
	}
	tmpPtr := src.DevPtr(0)
	srcMemObj := *(*cl.MemObject)(tmpPtr)
	tmpPtr = dst.DevPtr(0)
	dstMemObj := *(*cl.MemObject)(tmpPtr)

	// Synchronize in the beginning
	var err error
	eventList := []*cl.Event{}
	tmpEvt := src.GetEvent(0)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
	}
	tmpEvtL := dst.GetAllEvents(0)
	if len(tmpEvtL) > 0 {
		eventList = append(eventList, tmpEvtL...)
	}
	if len(eventList) != 0 {
		if err = cl.WaitForEvents(eventList); err != nil {
			log.Printf("WaitForEvents failed in fwPlan.ExecAsync: %+v \n", err)
		}
	}

	err = p.handle.EnqueueForwardTransform([]*cl.MemObject{&srcMemObj}, []*cl.MemObject{&dstMemObj})
	if Synchronous {
		ClCmdQueue.Finish()
		timer.Stop("fft")
	}
	tmpEvt, err = ClCmdQueue.EnqueueMarkerWithWaitList(nil)
	if err != nil {
		log.Printf("Failed to enqueue marker in fwPlan.ExecAsync: %+v \n", err)
	}
	dst.SetEvent(0, tmpEvt)
	src.InsertReadEvent(0, tmpEvt)
	if Debug {
		if err0 := cl.WaitForEvents([]*cl.Event{tmpEvt}); err0 != nil {
			log.Printf("WaitForEvents failed before returning fwPlan.ExecAsync: %+v \n", err0)
		}
		src.RemoveReadEvent(0, tmpEvt)
	} else {
		go func(evt *cl.Event, sl *data.Slice) {
			if err1 := cl.WaitForEvents([]*cl.Event{evt}); err1 != nil {
				log.Printf("WaitForEvents failed before returning fwPlan.ExecAsync: %+v \n", err1)
			}
			sl.RemoveReadEvent(0, evt)
		}(tmpEvt, src)
	}
	return err
}

// 3D size of the input array.
func (p *fft3DR2CPlan) InputSizeFloats() (Nx, Ny, Nz int) {
	return p.size[X], p.size[Y], p.size[Z]
}

// 3D size of the output array.
func (p *fft3DR2CPlan) OutputSizeFloats() (Nx, Ny, Nz int) {
	return 2 * (p.size[X]/2 + 1), p.size[Y], p.size[Z]
}

// Required length of the (1D) input array.
func (p *fft3DR2CPlan) InputLen() int {
	return prod3(p.InputSizeFloats())
}

// Required length of the (1D) output array.
func (p *fft3DR2CPlan) OutputLen() int {
	return prod3(p.OutputSizeFloats())
}
