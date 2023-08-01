package opencl

import (
	"context"
	"fmt"
	"math"
	"sync"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	qm "github.com/seeder-research/uMagNUS/queuemanager"
	util "github.com/seeder-research/uMagNUS/util"
)

// Use a small bunch of goroutines to handle copybacks

// Data type to pass to goroutines
type reduceEventAndPointer struct {
	event []*cl.Event
	cmdQ  *cl.CommandQueue
	ptr   unsafe.Pointer
	idx   int
}

type reduceOutput struct {
	res float32
	idx int
}

// Variables for managing the goroutines
var (
	reduceManagerContext  context.Context
	reduceManagerKillFunc context.CancelFunc
	reduceThreadCount     = int(-1)
	reduceItem            chan reduceEventAndPointer
	reduceOutput          chan float32
	reduceWaitGroup       sync.WaitGroup
)

// function executed by goroutines
func threadedCopyBack(in <-chan reduceEventAndPointer, out chan<- reduceOutput, bufPool chan<- (*cl.MemObject), wg_ *sync.WaitGroup) {
	tmpQueue, err := ClCxt.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("Failed to create command queue in goroutine: %+v \n", err)
	}
	for {
		select {
		case item <- in:
			if err = cl.WaitForEvents(in.event); err != nil {
				fmt.Printf("Failed to WaitForEvents: %+v \n", err)
			}
			var val float32
			_, err = tmpQueue.EnqueueReadBuffer((*cl.MemObject)(item.ptr), false, 0, SIZEOF_FLOAT32, (unsafe.Pointer)(&val), nil)
			if err != nil {
				fmt.Printf("Failed to enqueuereadbuffer: %+v \n", err)
			}
			if err = tmpQueue.Finish(); err != nil {
				fmt.Printf("Failed to wait for enqueuereadbuffer to finish: %+v \n", err)
			}
			out <- reduceOutput{val, item.idx}
			bufPool <- (*cl.MemObject)(item.ptr)
			wg_.Done()
		case <-reduceManagerContext.Done():
			err = tmpQueue.Finish()
			if err != nil {
				fmt.Printf("Failed to wait for queue to finish in goroutine: %+v \n", err)
			}
			tmpQueue.Release()
			return
		}
	}
}

// Initialize the goroutines for reduction output
func reduceInit(n int) {
	// Initialize globals
	reduceManagerContext, reduceManagerKillFunc = context.WithCancel(context.Background())
	reduceThreadCount = n
	reduceItem = make(reduceEventAndPointer, reduceThreadCount)
	reduceOutput = make(float32, reduceThreadCount)

	for j := 0; j < reduceThreads; j++ {
		go threadedCopyBack(reduceIten, reduceOutput, reduceBuffers, &reduceWaitGroup)
	}
}

// Teardown of goroutines for reduction output
func reduceTeardown() {
	if reduceThreadCount > 0 {
		reduceManagerKillFunc()
	}
	reduceThreadCount = -1
}

/*
Need to update the reduction sum and dot algorithms to balance the distribution of inputs to the work-groups
Less of a problem for max and min because they are direct comparisons

for sum and dot, the magnitude of intermediate values depend on how many input values were summed to obtain
them. Thus, the entire input data should be distributed in a binary tree for the summation, and the work-
groups should "synchronize" at a fixed level of the tree. Small input sizes will need fewer work-groups
that can efficiently calculate for trees having short depths. Larger input sizes will need larger number of
possibly small work-groups.
*/

// Sum of all elements.
func Sum(in *data.Slice, q *cl.CommandQueue, ewl []*cl.Event) float32 {
	util.Argument(in.NComp() == 1)
	out := reduceBuf(0)

	// Ensure no other reduction kernel is running
	reduceWaitGroup.Wait()

	// Launch kernel
	event := k_reducesum_async(in.DevPtr(0), out, 0,
		in.Len(), reducecfg, ewl, q)

	// Copy back to host in goroutine
	reduceWaitGroup.Add(1)
	reduceItem <- reduceEventAndPointer{[]*cl.Event{event}, q, out, 0}

	// Ensure all reduction kernel has completed
	reduceWaitGroup.Wait()
	tmp := <-reduceOutput

	return tmp.res
}

// Dot product
func Dot(a, b *data.Slice, q []*cl.CommandQueue, ewl []*cl.Event) float32 {
	util.Argument(a.NComp() == b.NComp())
	util.Argument(a.Len() == b.Len())
	util.Argument(a.NComp() == len(q))
	result := float32(0)
	numComp := a.NComp()
	out := make([]unsafe.Pointer, numComp)
	for c := 0; c < numComp; c++ {
		out[c] = reduceBuf(0)
	}
	hostResult := make([]float32, numComp)

	// Ensure no other reduction kernel is running
	reduceWaitGroup.Wait()
	for c := 0; c < numComp; c++ {
		// Launch kernel
		event := k_reducedot_async(a.DevPtr(c), b.DevPtr(c), out[c], 0,
			a.Len(), reducecfg, ewl, q[c]) // all components add to out

		// Copy back to host in goroutine
		reduceWaitGroup.Add(1)
		reduceItem <- reduceEventAndPointer{[]*cl.Event{event}, q[c], out[c], c}
	}

	// Ensure all reduction kernel has completed
	reduceWaitGroup.Wait()
	for c := 0; c < numComp; c++ {
		tmp := <-reduceOutput
		result += tmp.res
	}

	return result
}

// Maximum of absolute values of all elements.
func MaxAbs(in *data.Slice, q *cl.CommandQueue, ewl []*cl.Event) float32 {
	util.Argument(in.NComp() == 1)
	out := reduceBuf(0)

	// Ensure no other reduction kernel is running
	reduceWaitGroup.Wait()

	// Launch kernel
	event = k_reducemaxabs_async(in.DevPtr(0), out, 0,
		in.Len(), reducecfg, ewl, q)

	// Copy back to host in goroutine
	reduceWaitGroup.Add(1)
	reduceItem <- reduceEventAndPointer{[]*cl.Event{event}, q, out, 0}

	// Ensure all reduction kernel has completed
	reduceWaitGroup.Wait()
	tmp := <-reduceOutput

	return float32(tmp.res)
}

// Maximum element-wise difference
func MaxDiff(a, b *data.Slice, q *cl.CommandQueue, ewl []*cl.Event) []float32 {
	util.Argument(a.NComp() == b.NComp())
	util.Argument(a.Len() == b.Len())
	numComp := a.NComp()
	returnVal := make([]float32, numComp)
	out := make([]unsafe.Pointer, numComp)
	for c := 0; c < numComp; c++ {
		out[c] = reduceBuf(0)
	}

	// Ensure no other reduction kernel is running
	reduceWaitGroup.Wait()

	for c := 0; c < numComp; c++ {
		// Launch kernel
		event := k_reducemaxdiff_async(a.DevPtr(c), b.DevPtr(c), out[c], 0,
			a.Len(), reducecfg, ewl, q)

		// Copy back to host in goroutine
		reduceWaitGroup.Add(1)
		reduceItem <- reduceEventAndPointer{[]*cl.Event{event}, q[c], out[c], c}
	}

	// Must synchronize since returnVal is copied from device back to host
	reduceWaitGroup.Wait()
	for c := 0; c < numComp; c++ {
		tmp := <-reduceOutput
		returnVal[tmp.idx] = tmp.res
	}

	return returnVal
}

// Maximum of the norms of all vectors (x[i], y[i], z[i]).
//
//	max_i sqrt( x[i]*x[i] + y[i]*y[i] + z[i]*z[i] )
func MaxVecNorm(v *data.Slice, q *cl.CommandQueue, ewl []*cl.Event) float64 {
	util.Argument(v.NComp() == 3)
	out := reduceBuf(0)

	// Ensure no other reduction kernel is running
	reduceWaitGroup.Wait()

	// Launch kernel
	event := k_reducemaxvecnorm2_async(v.DevPtr(0), v.DevPtr(1), v.DevPtr(2),
		out, 0, v.Len(), reducecfg, ewl, q)

	// Copy back to host in goroutine
	reduceWaitGroup.Add(1)
	reduceItem <- reduceEventAndPointer{[]*cl.Event{event}, q, out, 0}

	// Ensure all reduction kernel has completed
	reduceWaitGroup.Wait()
	tmp := <-reduceOutput

	return math.Sqrt(float64(tmp.res))
}

// Maximum of the norms of the difference between all vectors (x1,y1,z1) and (x2,y2,z2)
//
//	(dx, dy, dz) = (x1, y1, z1) - (x2, y2, z2)
//	max_i sqrt( dx[i]*dx[i] + dy[i]*dy[i] + dz[i]*dz[i] )
func MaxVecDiff(x, y *data.Slice, q *cl.CommandQueue, ewl []*cl.Event) float64 {
	util.Argument(x.Len() == y.Len())
	util.Argument(x.NComp() == 3)
	util.Argument(y.NComp() == 3)
	out := reduceBuf(0)

	// Ensure no other reduction kernel is running
	reduceWaitGroup.Wait()

	// Launch kernel
	event := k_reducemaxvecdiff2_async(x.DevPtr(0), x.DevPtr(1), x.DevPtr(2),
		y.DevPtr(0), y.DevPtr(1), y.DevPtr(2),
		out, 0, x.Len(), reducecfg, ewl, q)

	// Copy back to host in goroutine
	reduceWaitGroup.Add(1)
	reduceItem <- reduceEventAndPointer{[]*cl.Event{event}, q, out, 0}

	// Ensure all reduction kernel has completed
	reduceWaitGroup.Wait()
	tmp := <-reduceOutput

	return math.Sqrt(float64(tmp.res))
}

var reduceBuffers chan (*cl.MemObject) // pool of 1-float OpenCL buffers for reduce

// return a 1-float and an N-float OPENCL reduction buffer from a pool
// initialized to initVal
func reduceBuf(initVal float32) unsafe.Pointer {
	if reduceBuffers == nil {
		initReduceBuf()
	}
	buf := <-reduceBuffers

	// Checkout queue
	tmpQueue := qm.CheckoutQueue(CmdQueuePool, nil)

	// Launch kernel
	//	waitEvent, err := ClCmdQueue.EnqueueFillBuffer(buf, unsafe.Pointer(&initVal), SIZEOF_FLOAT32, 0, ReduceWorkgroups*SIZEOF_FLOAT32, nil)
	waitEvent, err := tmpQueue.EnqueueFillBuffer(buf, unsafe.Pointer(&initVal), SIZEOF_FLOAT32, 0, SIZEOF_FLOAT32, nil)
	if err != nil {
		fmt.Printf("reduceBuf failed: %+v \n", err)
		return nil
	}

	// always synchronize reduceBuf()
	err = cl.WaitForEvents([]*cl.Event{waitEvent})

	// Check in queue post execution
	qwg := qm.NewQueueWaitGroup(tmpQueue, nil)
	ReturnQueuePool <- qwg

	if err != nil {
		fmt.Printf("First WaitForEvents in reduceBuf failed: %+v \n", err)
		reduceBuffers <- buf
		return nil
	}

	return (unsafe.Pointer)(buf)
}

// copy back single float result from GPU and recycle buffer
func copyback(buf unsafe.Pointer) float32 {
	var result float32
	MemCpyDtoH(unsafe.Pointer(&result), buf, SIZEOF_FLOAT32)
	reduceBuffers <- (*cl.MemObject)(buf)
	return result
}

// copy back float slice result from GPU and recycle buffer
func copybackSlice(buf unsafe.Pointer) []float32 {
	result := make([]float32, ReduceWorkgroups)
	MemCpyDtoH(unsafe.Pointer(&result[0]), buf, ReduceWorkgroups*SIZEOF_FLOAT32)
	reduceBuffers <- (*cl.MemObject)(buf)
	return result
}

// initialize pool of 1-float and N-float OPENCL reduction buffers
func initReduceBuf() {
	const N = 128
	reduceBuffers = make(chan *cl.MemObject, N)
	for i := 0; i < N; i++ {
		//		reduceBuffers <- MemAlloc(ReduceWorkgroups * SIZEOF_FLOAT32)
		reduceBuffers <- MemAlloc(SIZEOF_FLOAT32)
	}
}

// launch configuration for reduce kernels
// 8 is typ. number of multiprocessors.
// could be improved but takes hardly ~1% of execution time
var reducecfg = &config{Grid: []int{1, 1, 1}, Block: []int{1, 1, 1}}
var ReduceWorkitems int
var ReduceWorkgroups int
