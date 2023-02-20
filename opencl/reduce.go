package opencl

import (
	"fmt"
	"math"
	"sync"
	"unsafe"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	util "github.com/seeder-research/uMagNUS/util"
)

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
func Sum(in *data.Slice) float32 {
	util.Argument(in.NComp() == 1)
	out := reduceBuf(0)
	// check input slice for event to synchronize (if any)
	in.RLock(0)
	defer in.RUnlock(0)
	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("sum failed to create command queue: %+v \n", err)
		return -1.0
	}
	defer cmdqueue.Release()

	event := k_reducesum_async(in.DevPtr(0), out, 0,
		in.Len(), reducecfg, cmdqueue, nil)
	// Must synchronize since out is copied from device back to host
	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in sum: %+v \n", err)
	}
	//	results := copybackSlice(out)
	//	res := float32(0)
	//	for _, v := range results {
	//		res += v
	//	}
	//	return res
	results := copyback(out)
	return results
}

// Dot product
func Dot(a, b *data.Slice) float32 {
	util.Argument(a.NComp() == b.NComp())
	util.Argument(a.Len() == b.Len())
	result := float32(0)
	numComp := a.NComp()
	out := make([]unsafe.Pointer, numComp)
	for c := 0; c < numComp; c++ {
		out[c] = reduceBuf(0)
	}
	hostResult := make([]float32, numComp)
	var wg sync.WaitGroup
	// async over components
	for c := 0; c < numComp; c++ {
		wg.Add(1)
		if Synchronous {
			dot__(a, b, out[c], &hostResult[c], c, &wg)
		} else {
			go dot__(a, b, out[c], &hostResult[c], c, &wg)
		}
	}
	// Must synchronize since result is copied from device back to host
	wg.Wait()
	for _, oVal := range hostResult {
		result += oVal
	}
	return result
}

func dot__(a, b *data.Slice, outBufferPtr unsafe.Pointer, res *float32, idx int, wg_ *sync.WaitGroup) {
	defer wg_.Done()

	a.RLock(idx)
	b.RLock(idx)
	defer a.RUnlock(idx)
	defer b.RUnlock(idx)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("adddotproduct failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	event := k_reducedot_async(a.DevPtr(idx), b.DevPtr(idx), outBufferPtr, 0,
		a.Len(), reducecfg, cmdqueue, nil) // all components add to out

	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed at index %d in adddotproduct: %+v \n", idx, err)
	}
	//			results := copybackSlice(outBufferPtr)
	//			tmp := float32(0)
	//			for _, v := range results {
	//				tmp += v
	//			}
	//			*res = tmp
	results := copyback(outBufferPtr)
	*res = results
}

// Maximum of absolute values of all elements.
func MaxAbs(in *data.Slice) float32 {
	util.Argument(in.NComp() == 1)
	out := reduceBuf(0)

	in.RLock(0)
	defer in.RUnlock(0)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("adddotproduct failed to create command queue: %+v \n", err)
		return -1.0
	}
	defer cmdqueue.Release()

	event := k_reducemaxabs_async(in.DevPtr(0), out, 0,
		in.Len(), reducecfg, cmdqueue, nil)

	// Must synchronize since out is copied from device back to host
	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in maxabs: %+v \n", err)
	}
	//	results := copybackSlice(out)
	//	res := float64(results[0])
	//	for idx := 1; idx < ReduceWorkgroups; idx++ {
	//		res = math.Max(res, float64(results[idx]))
	//	}
	//	return float32(res)
	results := copyback(out)
	return float32(results)
}

// Maximum element-wise difference
func MaxDiff(a, b *data.Slice) []float32 {
	util.Argument(a.NComp() == b.NComp())
	util.Argument(a.Len() == b.Len())
	numComp := a.NComp()
	returnVal := make([]float32, numComp)
	out := make([]unsafe.Pointer, numComp)
	for c := 0; c < numComp; c++ {
		out[c] = reduceBuf(0)
	}
	var wg sync.WaitGroup
	for c := 0; c < numComp; c++ {
		wg.Add(1)
		if Synchronous {
			maxdiff__(a, b, out[c], &returnVal[c], c, &wg)
		} else {
			go maxdiff__(a, b, out[c], &returnVal[c], c, &wg)
		}
	}
	// Must synchronize since returnVal is copied from device back to host
	wg.Wait()
	return returnVal
}

func maxdiff__(a, b *data.Slice, outBufferPtr unsafe.Pointer, res *float32, idx int, wg_ *sync.WaitGroup){
	defer wg_.Done()

	a.RLock(idx)
	b.RLock(idx)
	defer a.RUnlock(idx)
	defer b.RUnlock(idx)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("maxabs failed to create command queue: %+v \n", err)
		return
	}
	defer cmdqueue.Release()

	event := k_reducemaxdiff_async(a.DevPtr(idx), b.DevPtr(idx), outBufferPtr, 0,
		a.Len(), reducecfg, cmdqueue, nil)

	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in maxabs: %+v \n", err)
	}

	//			results := copybackSlice(outBufferPtr)
	//			tmp := float64(results[0])
	//			for idx := 1; idx < ReduceWorkgroups; idx++ {
	//				tmp = math.Max(tmp, float64(results[idx]))
	//			}
	//			*res = float32(tmp)
	results := copyback(outBufferPtr)
	*res = float32(results)
}

// Maximum of the norms of all vectors (x[i], y[i], z[i]).
//
//	max_i sqrt( x[i]*x[i] + y[i]*y[i] + z[i]*z[i] )
func MaxVecNorm(v *data.Slice) float64 {
	util.Argument(v.NComp() == 3)
	out := reduceBuf(0)

	v.RLock(X)
	v.RLock(Y)
	v.RLock(Z)
	defer v.RLock(X)
	defer v.RLock(Y)
	defer v.RLock(Z)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("maxvecnorm failed to create command queue: %+v \n", err)
		return -1.0
	}
	defer cmdqueue.Release()

	event := k_reducemaxvecnorm2_async(v.DevPtr(0), v.DevPtr(1), v.DevPtr(2),
		out, 0, v.Len(), reducecfg, cmdqueue, nil)

	// Must synchronize since out is copied from device back to host
	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in maxvecnorm: %+v \n", err)
	}
	//	results := copybackSlice(out)
	//	res := float64(results[0])
	//	for idx := 1; idx < ReduceWorkgroups; idx++ {
	//		res = math.Max(res, float64(results[idx]))
	//	}
	//	return math.Sqrt(float64(res))
	results := copyback(out)
	return math.Sqrt(float64(results))
}

// Maximum of the norms of the difference between all vectors (x1,y1,z1) and (x2,y2,z2)
//
//	(dx, dy, dz) = (x1, y1, z1) - (x2, y2, z2)
//	max_i sqrt( dx[i]*dx[i] + dy[i]*dy[i] + dz[i]*dz[i] )
func MaxVecDiff(x, y *data.Slice) float64 {
	util.Argument(x.Len() == y.Len())
	util.Argument(x.NComp() == 3)
	util.Argument(y.NComp() == 3)
	out := reduceBuf(0)

	x.RLock(X)
	x.RLock(Y)
	x.RLock(Z)
	y.RLock(X)
	y.RLock(Y)
	y.RLock(Z)
	defer x.RUnlock(X)
	defer x.RUnlock(Y)
	defer x.RUnlock(Z)
	defer y.RUnlock(X)
	defer y.RUnlock(Y)
	defer y.RUnlock(Z)

	// Create the command queue to execute the command
	cmdqueue, err := ClCtx.CreateCommandQueue(ClDevice, 0)
	if err != nil {
		fmt.Printf("maxvecdiff failed to create command queue: %+v \n", err)
		return -1.0
	}
	defer cmdqueue.Release()

	event := k_reducemaxvecdiff2_async(x.DevPtr(0), x.DevPtr(1), x.DevPtr(2),
		y.DevPtr(0), y.DevPtr(1), y.DevPtr(2),
		out, 0, x.Len(), reducecfg, cmdqueue, nil)

	// Must synchronize since out is copied from device back to host
	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in maxvecdiff: %+v \n", err)
	}
	//	results := copybackSlice(out)
	//	res := float64(results[0])
	//	for idx := 1; idx < ReduceWorkgroups; idx++ {
	//		res = math.Max(res, float64(results[idx]))
	//	}
	//	return math.Sqrt(float64(res))
	results := copyback(out)
	return math.Sqrt(float64(results))
}

var reduceBuffers chan (*cl.MemObject) // pool of 1-float OpenCL buffers for reduce

// return a 1-float and an N-float OPENCL reduction buffer from a pool
// initialized to initVal
func reduceBuf(initVal float32) unsafe.Pointer {
	if reduceBuffers == nil {
		initReduceBuf()
	}
	buf := <-reduceBuffers
	//	waitEvent, err := ClCmdQueue.EnqueueFillBuffer(buf, unsafe.Pointer(&initVal), SIZEOF_FLOAT32, 0, ReduceWorkgroups*SIZEOF_FLOAT32, nil)
	waitEvent, err := ClCmdQueue.EnqueueFillBuffer(buf, unsafe.Pointer(&initVal), SIZEOF_FLOAT32, 0, SIZEOF_FLOAT32, nil)
	if err != nil {
		fmt.Printf("reduceBuf failed: %+v \n", err)
		return nil
	}
	err = cl.WaitForEvents([]*cl.Event{waitEvent})
	if err != nil {
		fmt.Printf("First WaitForEvents in reduceBuf failed: %+v \n", err)
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
