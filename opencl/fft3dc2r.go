package opencl

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	timer "github.com/seeder-research/uMagNUS/timer"
)

// 3D single-precision real-to-complex FFT plan.
type fft3DC2RPlan struct {
	fftplan
	size [3]int
}

// 3D single-precision real-to-complex FFT plan.
func newFFT3DC2R(Nx, Ny, Nz int) fft3DC2RPlan {
	handle := cl.NewVkFFTPlan(ClCtx) // new xyz swap
	handle.VkFFTSetFFTPlanSize([]int{Nx, Ny, Nz})

	return fft3DC2RPlan{fftplan{handle}, [3]int{Nx, Ny, Nz}}
}

// Execute the FFT plan, asynchronous.
// src and dst are 3D arrays stored 1D arrays.
func (p *fft3DC2RPlan) ExecAsync(src, dst *data.Slice) error {
	if Synchronous {
		ClCmdQueue.Finish()
		timer.Start("fft")
	}
	oksrclen := p.InputLenFloats()
	if src.Len() != oksrclen {
		panic(fmt.Errorf("fft size mismatch: expecting src len %v, got %v", oksrclen, src.Len()))
	}
	okdstlen := p.OutputLenFloats()
	if dst.Len() != okdstlen {
		panic(fmt.Errorf("fft size mismatch: expecting dst len %v, got %v", okdstlen, dst.Len()))
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
			fmt.Printf("WaitForEvents failed in bwPlan.ExecAsync: %+v \n", err)
		}
	}

	err = p.handle.EnqueueBackwardTransform([]*cl.MemObject{&srcMemObj}, []*cl.MemObject{&dstMemObj})
	if Synchronous {
		ClCmdQueue.Finish()
		timer.Stop("fft")
	}
	tmpEvt, err = ClCmdQueue.EnqueueMarkerWithWaitList(nil)
	if err != nil {
		log.Printf("Failed to enqueue marker in bwPlan.ExecAsync: %+v \n", err)
	}
	dst.SetEvent(0, tmpEvt)
	src.InsertReadEvent(0, tmpEvt)
	if Debug {
		if err = cl.WaitForEvents(eventList); err != nil {
			log.Printf("WaitForEvents failed before returning bwPlan.ExecAsync: %+v \n", err)
		}
		src.RemoveReadEvent(0, tmpEvt)
	} else {
		go func(evt *c.lEvent, sl *data.Slice) {
			if err = cl.WaitForEvents(eventList); err != nil {
				log.Printf("WaitForEvents failed before returning bwPlan.ExecAsync: %+v \n", err)
			}
			src.RemoveReadEvent(0, tmpEvt)
		}(tmpEvt, src)
	}
	return err
}

// 3D size of the input array.
func (p *fft3DC2RPlan) InputSizeFloats() (Nx, Ny, Nz int) {
	return 2 * (p.size[X]/2 + 1), p.size[Y], p.size[Z]
}

// 3D size of the output array.
func (p *fft3DC2RPlan) OutputSizeFloats() (Nx, Ny, Nz int) {
	return p.size[X], p.size[Y], p.size[Z]
}

// Required length of the (1D) input array.
func (p *fft3DC2RPlan) InputLenFloats() int {
	return prod3(p.InputSizeFloats())
}

// Required length of the (1D) output array.
func (p *fft3DC2RPlan) OutputLenFloats() int {
	return prod3(p.OutputSizeFloats())
}
