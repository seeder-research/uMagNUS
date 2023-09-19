package opencl

// Generation of Magnetic Force Microscopy images.

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	mag "github.com/seeder-research/uMagNUS/mag"
	qm "github.com/seeder-research/uMagNUS/queuemanager"
)

// Stores the necessary state to perform FFT-accelerated convolution
type MFMConvolution struct {
	size        [3]int          // 3D size of the input/output data
	kernSize    [3]int          // Size of kernel and logical FFT size.
	fftKernSize [3]int          //
	fftRBuf     [3]*data.Slice  // FFT input buf for FFT, shares storage with fftCBuf.
	fftCBuf     [3]*data.Slice  // FFT output buf, shares storage with fftRBuf
	gpuFFTKern  [3]*data.Slice  // FFT kernel on device
	fwPlan      [3]fft3DR2CPlan // Forward FFT (1 component)
	bwPlan      [3]fft3DC2RPlan // Backward FFT (1 component)
	kern        [3]*data.Slice  // Real-space kernel (host)
	mesh        *data.Mesh
}

func (c *MFMConvolution) Free() {
	if c == nil {
		return
	}
	c.size = [3]int{}
	c.kernSize = [3]int{}

	for j := 0; j < 3; j++ {
		c.gpuFFTKern[j].Free()
		c.gpuFFTKern[j] = nil
		c.kern[j] = nil
		c.fwPlan[j].Free()
		c.bwPlan[j].Free()
		c.fftCBuf[j].Free() // shared with fftRbuf
		c.fftCBuf[j] = nil
		c.fftRBuf[j] = nil
	}
}

func (c *MFMConvolution) init() {
	// init FFT plans
	padded := c.kernSize
	c.fwPlan[X] = newFFT3DR2C(padded[X], padded[Y], padded[Z])
	c.fwPlan[Y] = newFFT3DR2C(padded[X], padded[Y], padded[Z])
	c.fwPlan[Z] = newFFT3DR2C(padded[X], padded[Y], padded[Z])
	c.bwPlan[X] = newFFT3DC2R(padded[X], padded[Y], padded[Z])
	c.bwPlan[Y] = newFFT3DC2R(padded[X], padded[Y], padded[Z])
	c.bwPlan[Z] = newFFT3DC2R(padded[X], padded[Y], padded[Z])

	// init device buffers
	nc := fftR2COutputSizeFloats(c.kernSize)
	c.fftCBuf[X] = NewSlice(1, nc)
	c.fftCBuf[Y] = NewSlice(1, nc)
	c.fftCBuf[Z] = NewSlice(1, nc)
	c.fftRBuf[X] = NewSlice(1, c.kernSize)
	c.fftRBuf[Y] = NewSlice(1, c.kernSize)
	c.fftRBuf[Z] = NewSlice(1, c.kernSize)

	c.gpuFFTKern[X] = NewSlice(1, nc)
	c.gpuFFTKern[Y] = NewSlice(1, nc)
	c.gpuFFTKern[Z] = NewSlice(1, nc)

	c.initFFTKern3D()
}

func (c *MFMConvolution) initFFTKern3D() {
	c.fftKernSize = fftR2COutputSizeFloats(c.kernSize)

	for i := 0; i < 3; i++ {
		// Attempt to use a single command queue for each loop
		// since the commands are sequential

		// Use queue for the fwFFT
		tmpQueue := c.fwPlan[i].GetCommandQueue()

		// Launch the zero1 kernel
		zero1_async(c.fftRBuf[i], tmpQueue, []*cl.Event{})
		// Launch the copy kernel, which runs in a separate queue
		data.Copy(c.fftRBuf[i], c.kern[i])

		// Get the event for the copy and insert into the
		// main execution queue
		ewl := c.fftRBuf[i].GetAllEvents(1)
		_, err := tmpQueue.EnqueueMarkerWithWaitList(ewl)
		if err != nil {
			fmt.Printf("error enqueuing sync marker in initfftkern3d: %+v \n", err)
		}

		// Launch fwFFT kernel in the main queue
		err = c.fwPlan[i].ExecAsync(c.fftRBuf[i], c.fftCBuf[i])
		if err != nil {
			fmt.Printf("error enqueuing forward fft in initfftkern3d: %+v \n", err)
		}
		scale := 2 / float32(c.fwPlan[i].InputLen()) // ??

		// Checkout new queue for zero1 of gpuFFTKern and launch
		// Possible for it to run concurrently with fwFFT
		zq := qm.CheckoutQueue(CmdQueuePool, nil)
		zero1_async(c.gpuFFTKern[i], zq, []*cl.Event{})

		// Get marker for synchronizing madd2
		var ev1 *cl.Event
		ev1, err = zq.EnqueueMarkerWithWaitList([]*cl.Event{})
		if err != nil {
			fmt.Printf("Failed to enqueue marker in initFFTKern3d: %+v \n", err)
		}

		// Checkin queue post execution
		qwg := qm.NewQueueWaitGroup(zq, nil)
		ReturnQueuePool <- qwg

		// Launch this Madd2 in the fwFFT queue, which will produce fftCBuf
		// But the kernel will need to sync to the event for zero1 of gpuFFTKernel
		Madd2(c.gpuFFTKern[i], c.gpuFFTKern[i], c.fftCBuf[i], 0, scale, []*cl.CommandQueue{tmpQueue}, []*cl.Event{ev1})
	}
}

// store MFM image in output, based on magnetization in inp.
func (c *MFMConvolution) Exec(outp, inp, vol *data.Slice, Msat MSlice) {
	for i := 0; i < 3; i++ {
		// Each loop will use one main queue (fwFFT queue)
		// New queues will be used for any kernel that
		// can run concurrently with previously queued
		// kernels

		// Use queue for the fwFFT
		tmpQueue := c.fwPlan[i].GetCommandQueue()

		// Launch zero1 kernel
		zero1_async(c.fftRBuf[i], tmpQueue, []*cl.Event{})

		// Launch copyPadMul kernel
		copyPadMul(c.fftRBuf[i], inp.Comp(i), vol, c.kernSize, c.size, Msat, tmpQueue, []*cl.Event{})

		var err error
		if err = c.fwPlan[i].ExecAsync(c.fftRBuf[i], c.fftCBuf[i]); err != nil {
			fmt.Printf("error enqueuing forward fft in mfmconv exec: %+v \n", err)
		}

		Nx, Ny := c.fftKernSize[X]/2, c.fftKernSize[Y] //   ??
		kernMulC_async(c.fftCBuf[i], c.gpuFFTKern[i], Nx, Ny, []*cl.Event{}, tmpQueue)

		// Now we need to insert marker in queue for bwFFT to sync
		var ev1 *cl.Event
		ev1, err = tmpQueue.EnqueueMarkerWithWaitList([]*cl.Event{})
		tmpQueue = c.bwPlan[i].GetCommandQueue()
		_, err = tmpQueue.EnqueueMarkerWithWaitList([]*cl.Event{ev1})

		// Launch bwFFT
		if err = c.bwPlan[i].ExecAsync(c.fftCBuf[i], c.fftRBuf[i]); err != nil {
			fmt.Printf("error enqueuing backward fft in mfmconv exec: %+v \n", err)
		}

		// Launch copyUnPad in bwFFT queue
		copyUnPad(outp.Comp(i), c.fftRBuf[i], c.size, c.kernSize, tmpQueue, []*cl.Event{})
	}
}

func (c *MFMConvolution) Reinit(lift, tipsize float64, cachedir string) {
	c.kern = mag.MFMKernel(c.mesh, lift, tipsize, cachedir)
	c.initFFTKern3D()
}

// Initializes a convolution to evaluate the demag field for the given mesh geometry.
func NewMFM(mesh *data.Mesh, lift, tipsize float64, cachedir string) *MFMConvolution {
	k := mag.MFMKernel(mesh, lift, tipsize, cachedir)
	size := mesh.Size()
	c := new(MFMConvolution)
	c.size = size
	c.kern = k
	c.kernSize = k[X].Size()
	c.init()
	c.mesh = mesh
	return c
}
