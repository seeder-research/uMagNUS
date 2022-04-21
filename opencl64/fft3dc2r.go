package opencl64

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data64"
	timer "github.com/seeder-research/uMagNUS/timer"
)

// 3D single-precision real-to-complex FFT plan.
type fft3DC2RPlan struct {
	fftplan
	size [3]int
}

// 3D single-precision real-to-complex FFT plan.
func newFFT3DC2R(Nx, Ny, Nz int) fft3DC2RPlan {
	handle := cl.NewVkFFTPlanDouble(ClCtx) // new xyz swap
	handle.VkFFTSetFFTPlanSize([]int{Nx, Ny, Nz})

	return fft3DC2RPlan{fftplan{handle}, [3]int{Nx, Ny, Nz}}
}

// Execute the FFT plan, asynchronous.
// src and dst are 3D arrays stored 1D arrays.
func (p *fft3DC2RPlan) ExecAsync(src, dst *data.Slice) error {
	ClCmdQueue.Finish()
	if Synchronous {
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
	tmpEvt = dst.GetEvent(0)
	if tmpEvt != nil {
		eventList = append(eventList, tmpEvt)
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
