package opencl

import (
	"fmt"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	qm "github.com/seeder-research/uMagNUS/queuemanager"
	util "github.com/seeder-research/uMagNUS/util"
)

// Stores the necessary state to perform FFT-accelerated convolution
// with magnetostatic kernel (or other kernel of same symmetry).
type DemagConvolution struct {
	inputSize        [3]int            // 3D size of the input/output data
	realKernSize     [3]int            // Size of kernel and logical FFT size.
	fftKernLogicSize [3]int            // logic size FFTed kernel, real parts only, we store less
	fftRBuf          [3]*data.Slice    // FFT input buf; 2D: Z shares storage with X.
	fftCBuf          [3]*data.Slice    // FFT output buf; 2D: Z shares storage with X.
	kern             [3][3]*data.Slice // FFT kernel on device
	fwPlan           [3]fft3DR2CPlan   // Forward FFT (1 component)
	bwPlan           [3]fft3DC2RPlan   // Backward FFT (1 component)
}

// Initializes a convolution to evaluate the demag field for the given mesh geometry.
// Sanity-checked if test == true (slow-ish for large meshes).
func NewDemag(inputSize, PBC [3]int, kernel [3][3]*data.Slice, test bool) *DemagConvolution {
	c := new(DemagConvolution)
	c.inputSize = inputSize
	c.realKernSize = kernel[X][X].Size()
	c.init(kernel)
	if test {
		testConvolution(c, PBC, kernel)
	}
	return c
}

// Calculate the demag field of m * vol * Bsat, store result in B.
//
//	m:    magnetization normalized to unit length
//	vol:  unitless mask used to scale m's length, may be nil
//	Bsat: saturation magnetization in Tesla
//	B:    resulting demag field, in Tesla
func (c *DemagConvolution) Exec(B, m, vol *data.Slice, Msat MSlice) {
	util.Argument(B.Size() == c.inputSize && m.Size() == c.inputSize)
	if c.is2D() {
		c.exec2D(B, m, vol, Msat)
	} else {
		c.exec3D(B, m, vol, Msat)
	}
}

func (c *DemagConvolution) exec3D(outp, inp, vol *data.Slice, Msat MSlice) {
	for i := 0; i < 3; i++ { // FW FFT
		c.fwFFT(i, inp, vol, Msat)
	}

	// Create event marker in fwFFT queues for synchronization
	q1, q2, q3 := c.fwPlan[0].GetCommandQueue(), c.fwPlan[1].GetCommandQueue(), c.fwPlan[2].GetCommandQueue()
	var ev1, ev2, ev3 *cl.Event
	var err error
	ev1, err = q1.EnqueueMarkerWithWaitList([]*cl.Event{})
	if err != nil {
		fmt.Printf("EnqueueMarkerWithWaitList failed on q1 in exec3D: %+v \n", err)
	}
	ev2, err = q2.EnqueueMarkerWithWaitList([]*cl.Event{})
	if err != nil {
		fmt.Printf("EnqueueMarkerWithWaitList failed on q2 in exec3D: %+v \n", err)
	}
	ev3, err = q3.EnqueueMarkerWithWaitList([]*cl.Event{})
	if err != nil {
		fmt.Printf("EnqueueMarkerWithWaitList failed on q3 in exec3D: %+v \n", err)
	}

	// Checkout new queue to launch kernMulRSymm3D_async kernel
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, nil)

	// Launch kernMulRSymm3D_async kernel with wait list on
	// the synchronization events in all fwPlan queues
	kernMulRSymm3D_async(c.fftCBuf,
		c.kern[X][X], c.kern[Y][Y], c.kern[Z][Z],
		c.kern[Y][Z], c.kern[X][Z], c.kern[X][Y],
		c.fftKernLogicSize[X], c.fftKernLogicSize[Y], c.fftKernLogicSize[Z],
		[]*cl.Event{ev1, ev2, ev3},
		tmpQueue)

	// Get event marker for synchronizing bwPlan queues
	var ev *cl.Event
	ev, err = tmpQueue.EnqueueMarkerWithWaitList([]*cl.Event{})
	event := []*cl.Event{ev}

	// Check in queue after kernel launch
	qwg := qm.NewQueueWaitGroup(tmpQueue, nil)
	ReturnQueuePool <- qwg

	// Insert synchronization event into all bwPlan queues
	q1, q2, q3 = c.bwPlan[0].GetCommandQueue(), c.bwPlan[1].GetCommandQueue(), c.bwPlan[2].GetCommandQueue()
	_, err = q1.EnqueueMarkerWithWaitList(event)
	if err != nil {
		fmt.Printf("EnqueueMarkerWithWaitList failed on q1 in exec3D: %+v \n", err)
	}
	_, err = q2.EnqueueMarkerWithWaitList(event)
	if err != nil {
		fmt.Printf("EnqueueMarkerWithWaitList failed on q2 in exec3D: %+v \n", err)
	}
	_, err = q3.EnqueueMarkerWithWaitList(event)
	if err != nil {
		fmt.Printf("EnqueueMarkerWithWaitList failed on q3 in exec3D: %+v \n", err)
	}

	for i := 0; i < 3; i++ { // BW FFT
		c.bwFFT(i, outp)
	}
}

func (c *DemagConvolution) exec2D(outp, inp, vol *data.Slice, Msat MSlice) {
	// Convolution is separated into
	// a 1D convolution for z and a 2D convolution for xy.
	// So only 2 FFT buffers are needed at the same time.
	Nx, Ny := c.fftKernLogicSize[X], c.fftKernLogicSize[Y]

	q1, q2, q3 := c.fwPlan[X].GetCommandQueue(), c.fwPlan[Y].GetCommandQueue(), c.fwPlan[Z].GetCommandQueue()
	var ev1, ev2 *cl.Event
	var err error

	// Z
	c.fwFFT(Z, inp, vol, Msat)

	// Create event in fwFFT (Z) queue for synchronization
	ev1, err = q3.EnqueueMarkerWithWaitList([]*cl.Event{})
	if err != nil {
		fmt.Printf("EnqueueMarkerWithWaitList failed on q3 in exec2D (fwFFT, Z): %+v \n", err)
	}

	// Checkout new queue to launch kernMulRSymm2Dz_async kernel
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, nil)

	// Launch kernel with wait list to sync with fwFFT (Z)
	kernMulRSymm2Dz_async(c.fftCBuf[Z], c.kern[Z][Z], Nx, Ny, []*cl.Event{ev1}, tmpQueue)

	// Get event to synchronize bwFFT (Z) queue
	ev1, err = tmpQueue.EnqueueMarkerWithWaitList([]*cl.Event{})
	if err != nil {
		fmt.Printf("EnqueueMarkerWithWaitList failed on q3 in exec2D (MulRSymm2Dz): %+v \n", err)
	}

	// Check in queue post kernel execution
	qwg1 := qm.NewQueueWaitGroup(tmpQueue, nil)
	ReturnQueuePool <- qwg1

	// Insert synchronization event into bwPlan (Z) queue
	q3 = c.bwPlan[Z].GetCommandQueue()
	_, err = tmpQueue.EnqueueMarkerWithWaitList([]*cl.Event{})
	if err != nil {
		fmt.Printf("EnqueueMarkerWithWaitList failed on q3 in exec2D (bwFFT, Z): %+v \n", err)
	}

	c.bwFFT(Z, outp)

	// XY
	c.fwFFT(X, inp, vol, Msat)
	c.fwFFT(Y, inp, vol, Msat)

	// Create event in fwFFT queues for synchronization
	ev1, err = q1.EnqueueMarkerWithWaitList([]*cl.Event{})
	if err != nil {
		fmt.Printf("EnqueueMarkerWithWaitList failed on q1 in exec2D (fwFFT, X): %+v \n", err)
	}
	ev2, err = q2.EnqueueMarkerWithWaitList([]*cl.Event{})
	if err != nil {
		fmt.Printf("EnqueueMarkerWithWaitList failed on q2 in exec2D (fwFFTm Y): %+v \n", err)
	}

	// Checkout new queue to launch kernMulRSymm2Dz_async kernel
	tmpQueue = qm.CheckoutQueue(CmdQueuePool, nil)

	// Launch kernel with wait list to sync with fwFFT (Z)
	kernMulRSymm2Dxy_async(c.fftCBuf[X], c.fftCBuf[Y],
		c.kern[X][X], c.kern[Y][Y], c.kern[X][Y], Nx, Ny, []*cl.Event{ev1, ev2}, tmpQueue)

	// Get event to synchronize bwFFT (Z) queue
	ev1, err = tmpQueue.EnqueueMarkerWithWaitList([]*cl.Event{})
	if err != nil {
		fmt.Printf("EnqueueMarkerWithWaitList failed on q3 in exec2D (MulRSymm2Dxy): %+v \n", err)
	}

	// Check in queue post kernel execution
	qwg2 := qm.NewQueueWaitGroup(tmpQueue, nil)
	ReturnQueuePool <- qwg2

	// Insert synchronization event into bwPlan (Z) queue
	q1, q2 = c.bwPlan[X].GetCommandQueue(), c.bwPlan[Y].GetCommandQueue()

	// Create event in fwFFT queues for synchronization
	ev1, err = q1.EnqueueMarkerWithWaitList([]*cl.Event{})
	if err != nil {
		fmt.Printf("EnqueueMarkerWithWaitList failed on q1 in exec2D (bwFFT, X): %+v \n", err)
	}
	ev2, err = q2.EnqueueMarkerWithWaitList([]*cl.Event{})
	if err != nil {
		fmt.Printf("EnqueueMarkerWithWaitList failed on q2 in exec2D (bwFFT, Y): %+v \n", err)
	}

	c.bwFFT(X, outp)
	c.bwFFT(Y, outp)
}

func (c *DemagConvolution) is2D() bool {
	return c.inputSize[Z] == 1
}

// zero 1-component slice
func zero1_async(dst *data.Slice, q *cl.CommandQueue, ewl []*cl.Event) {
	val := float32(0.0)
	if dst == nil {
		panic("ERROR (zero1_async): dst pointer cannot be nil")
	}

	// Launch kernel
	event, err := q.EnqueueFillBuffer((*cl.MemObject)(dst.DevPtr(0)), unsafe.Pointer(&val), SIZEOF_FLOAT32, 0, dst.Len()*SIZEOF_FLOAT32, ewl)

	if err != nil {
		fmt.Printf("EnqueueFillBuffer failed: %+v \n", err)
	}

	if Debug {
		if err = cl.WaitForEvents([](*cl.Event){event}); err != nil {
			fmt.Printf("WaitForEvents failed in zero1_async: %+v \n", err)
		}
	}
}

// forward FFT component i
func (c *DemagConvolution) fwFFT(i int, inp, vol *data.Slice, Msat MSlice) {
	// Get queue from fwPlan
	tmpQueue := c.fwPlan[i].GetCommandQueue()
	zero1_async(c.fftRBuf[i], tmpQueue, []*cl.Event{})
	in := inp.Comp(i)
	copyPadMul(c.fftRBuf[i], in, vol, c.realKernSize, c.inputSize, Msat, tmpQueue, []*cl.Event{})
	if err := c.fwPlan[i].ExecAsync(c.fftRBuf[i], c.fftCBuf[i]); err != nil {
		fmt.Printf("Error enqueuing forward fft: %+v \n", err)
	}
}

// backward FFT component i
func (c *DemagConvolution) bwFFT(i int, outp *data.Slice) {
	// Get queue from bwPlan
	tmpQueue := c.bwPlan[i].GetCommandQueue()
	if err := c.bwPlan[i].ExecAsync(c.fftCBuf[i], c.fftRBuf[i]); err != nil {
		fmt.Printf("Error enqueuing backward fft: %+v", err)
	}
	out := outp.Comp(i)
	copyUnPad(out, c.fftRBuf[i], c.inputSize, c.realKernSize, tmpQueue, []*cl.Event{})
}

func (c *DemagConvolution) init(realKern [3][3]*data.Slice) {
	// init device buffers
	// 2D re-uses fftBuf[X] as fftBuf[Z], 3D needs all 3 fftBufs.
	nc := fftR2COutputSizeFloats(c.realKernSize)
	c.fftCBuf[X] = NewSlice(1, nc)
	c.fftCBuf[Y] = NewSlice(1, nc)
	if c.is2D() {
		c.fftCBuf[Z] = c.fftCBuf[X]
	} else {
		c.fftCBuf[Z] = NewSlice(1, nc)
	}

	c.fftRBuf[X] = NewSlice(1, c.realKernSize)
	c.fftRBuf[Y] = NewSlice(1, c.realKernSize)
	if c.is2D() {
		c.fftRBuf[Z] = c.fftRBuf[X]
	} else {
		c.fftRBuf[Z] = NewSlice(1, c.realKernSize)
	}

	// init FFT plans
	c.fwPlan[0] = newFFT3DR2C(c.realKernSize[X], c.realKernSize[Y], c.realKernSize[Z])
	c.fwPlan[1] = newFFT3DR2C(c.realKernSize[X], c.realKernSize[Y], c.realKernSize[Z])
	c.fwPlan[2] = newFFT3DR2C(c.realKernSize[X], c.realKernSize[Y], c.realKernSize[Z])
	c.bwPlan[0] = newFFT3DC2R(c.realKernSize[X], c.realKernSize[Y], c.realKernSize[Z])
	c.bwPlan[1] = newFFT3DC2R(c.realKernSize[X], c.realKernSize[Y], c.realKernSize[Z])
	c.bwPlan[2] = newFFT3DC2R(c.realKernSize[X], c.realKernSize[Y], c.realKernSize[Z])

	// init FFT kernel

	// logic size of FFT(kernel): store real parts only
	c.fftKernLogicSize = fftR2COutputSizeFloats(c.realKernSize)
	util.Assert(c.fftKernLogicSize[X]%2 == 0)
	c.fftKernLogicSize[X] /= 2

	// physical size of FFT(kernel): store only non-redundant part exploiting Y, Z mirror symmetry
	// X mirror symmetry already exploited: FFT(kernel) is purely real.
	physKSize := [3]int{c.fftKernLogicSize[X], c.fftKernLogicSize[Y]/2 + 1, c.fftKernLogicSize[Z]/2 + 1}

	output := c.fftCBuf[0]
	input := c.fftRBuf[0]
	fftKern := data.NewSlice(1, physKSize)
	kfull := data.NewSlice(1, output.Size()) // not yet exploiting symmetry
	kfulls := kfull.Scalars()
	kCSize := physKSize
	kCSize[X] *= 2                     // size of kernel after removing Y,Z redundant parts, but still complex
	kCmplx := data.NewSlice(1, kCSize) // not yet exploiting X symmetry
	kc := kCmplx.Scalars()

	for i := 0; i < 3; i++ {
		for j := i; j < 3; j++ { // upper triangular part
			if realKern[i][j] != nil { // ignore 0's
				// FW FFT
				data.Copy(input, realKern[i][j])
				err := c.fwPlan[j].ExecAsync(input, output)
				if err != nil {
					fmt.Printf("error enqueuing forward fft in init: %+v \n ", err)
				}
				data.Copy(kfull, output)

				// extract non-redundant part (Y,Z symmetry)
				for iz := 0; iz < kCSize[Z]; iz++ {
					for iy := 0; iy < kCSize[Y]; iy++ {
						for ix := 0; ix < kCSize[X]; ix++ {
							kc[iz][iy][ix] = kfulls[iz][iy][ix]
						}
					}
				}

				// extract real parts (X symmetry)
				scaleRealParts(fftKern, kCmplx, 1/float32(c.fwPlan[j].InputLen()))
				c.kern[i][j] = GPUCopy(fftKern)
			}
		}
	}
}

func (c *DemagConvolution) Free() {
	if c == nil {
		return
	}
	c.inputSize = [3]int{}
	c.realKernSize = [3]int{}
	for i := 0; i < 3; i++ {
		c.fftCBuf[i].Free()
		c.fftRBuf[i].Free()
		c.fftCBuf[i] = nil
		c.fftRBuf[i] = nil

		for j := 0; j < 3; j++ {
			c.kern[i][j].Free()
			c.kern[i][j] = nil
		}
		c.fwPlan[i].Free()
		c.bwPlan[i].Free()
	}
}
