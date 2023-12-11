package opencl

import (
	"log"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
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
	handle := cl.NewVkFFTPlan(ClCtx)
	handle.VkFFTSetFFTPlanSize([]int{Nx, Ny, Nz})

	return fft3DR2CPlan{fftplan{handle}, [3]int{Nx, Ny, Nz}}
}

// Execute the FFT plan, asynchronous.
// src and dst are 3D arrays stored 1D arrays.
func (p *fft3DR2CPlan) ExecAsync(src, dst *data.Slice, queue *cl.CommandQueue, events []*cl.Event) error {

	if Synchronous {
		if err := queue.Finish(); err != nil {
			log.Panicf("failed to wait for queue to finish in beginning of fft3dr2c.execasync: %+v \n", err)
		}
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

	err := p.handle.EnqueueForwardTransform([]*cl.MemObject{&srcMemObj}, []*cl.MemObject{&dstMemObj})

	if err != nil {
		log.Printf("Failed to enqueue fwFFT: %+v \n", err)
	}

	if Synchronous {
		if err := queue.Finish(); err != nil {
			log.Panicf("failed to wait for queue to finish at end of fft3dr2c.execasync: %+v \n", err)
		}
		timer.Stop("fft")
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

// Return command queue associated with the plan
func (p *fft3DR2CPlan) GetCommandQueue() *cl.CommandQueue {
	return p.handle.GetCommandQueue()
}
