package opencl

// INTERNAL
// Base implementation for all FFT plans.

import (
	cl "github.com/seeder-research/uMagNUS/cl"
)

// Base implementation for all FFT plans.
type fftplan struct {
	//	handle *cl.ClFFTPlan
	handle *cl.VkfftPlan
}

func prod3(x, y, z int) int {
	return x * y * z
}

// Releases all resources associated with the FFT plan.
func (p *fftplan) Free() {
	if p.handle != nil {
		p.handle.Destroy()
		p.handle = nil
	}
}
