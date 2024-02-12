package opencl

import (
	"fmt"
	"math"
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
func Sum(in *data.Slice, queue *cl.CommandQueue, events []*cl.Event) float32 {
	util.Argument(in.NComp() == 1)

	SyncQueues([]*cl.CommandQueue{ClCmdQueue[1]}, []*cl.CommandQueue{queue})
	out := reduceBuf(0, ClCmdQueue[1], events)

	if Synchronous {
		if err := queue.Finish(); err != nil {
			fmt.Printf("failed to wait for queue to finish in sum: %+v \n", err)
		}
	}

	k_reducesum_async(in.DevPtr(0), out, 0,
		in.Len(), reducecfg, ClCmdQueue[1], events)

	//	results := copybackSlice(out, ClCmdQueue[1], events)
	//	res := float32(0)
	//	for _, v := range results {
	//		res += v
	//	}
	//	return res
	results := copyback(out, ClCmdQueue[1], events)
	return results
}

// Dot product
func Dot(a, b *data.Slice, queue *cl.CommandQueue, events []*cl.Event) float32 {
	util.Argument(a.NComp() == b.NComp())
	util.Argument(a.Len() == b.Len())
	result := float32(0)
	numComp := a.NComp()
	out := make([]unsafe.Pointer, numComp)
	if err := WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("eror waiting all queues to finish in dot: %+v \n", err)
	}
	for c := 0; c < numComp; c++ {
		out[c] = reduceBuf(0, ClCmdQueue[1+c], events)
	}

	for c := 0; c < numComp; c++ {
		k_reducedot_async(a.DevPtr(c), b.DevPtr(c), out[c], 0,
			a.Len(), reducecfg, ClCmdQueue[1+c], nil) // all components add to out
	}

	for c := 0; c < a.NComp(); c++ {
		results := copyback(out[c], ClCmdQueue[1+c], nil)
		//		for _, v := range results {
		//			result += v
		//		}
		result += results
	}
	return result
}

// Maximum of absolute values of all elements.
func MaxAbs(in *data.Slice, queue *cl.CommandQueue, events []*cl.Event) float32 {
	util.Argument(in.NComp() == 1)
	SyncQueues([]*cl.CommandQueue{ClCmdQueue[1]}, []*cl.CommandQueue{queue})
	out := reduceBuf(0, ClCmdQueue[1], events)

	k_reducemaxabs_async(in.DevPtr(0), out, 0,
		in.Len(), reducecfg, ClCmdQueue[1], nil)

	//	results := copybackSlice(out, ClCmdQueue[1], nil)
	//	res := float64(results[0])
	//	for idx := 1; idx < ReduceWorkgroups; idx++ {
	//		res = math.Max(res, float64(results[idx]))
	//	}
	//	return float32(res)
	results := copyback(out, ClCmdQueue[1], nil)
	return float32(results)
}

// Maximum element-wise difference
func MaxDiff(a, b *data.Slice, queue *cl.CommandQueue, events []*cl.Event) []float32 {
	util.Argument(a.NComp() == b.NComp())
	util.Argument(a.Len() == b.Len())
	numComp := a.NComp()
	returnVal := make([]float32, numComp)
	out := make([]unsafe.Pointer, numComp)
	SyncQueues([]*cl.CommandQueue{ClCmdQueue[1], ClCmdQueue[2], ClCmdQueue[3]}, []*cl.CommandQueue{queue})
	for c := 0; c < numComp; c++ {
		out[c] = reduceBuf(0, ClCmdQueue[1+c], events)
	}

	for c := 0; c < numComp; c++ {
		k_reducemaxdiff_async(a.DevPtr(c), b.DevPtr(c), out[c], 0,
			a.Len(), reducecfg, ClCmdQueue[1+c], nil)
	}

	for c := 0; c < numComp; c++ {
		//		results := copybackSlice(out[c], ClCmdQueue[1+c], nil)
		//		tmpVal := float64(results[0])
		//		for idx := 1; idx < ReduceWorkgroups; idx++ {
		//			tmpVal = math.Max(tmpVal, float64(results[idx]))
		//		}
		//		returnVal[c] = float32(tmpVal)
		results := copyback(out[c], ClCmdQueue[1+c], nil)
		returnVal[c] = float32(results)
	}
	return returnVal
}

// Maximum of the norms of all vectors (x[i], y[i], z[i]).
//
//	max_i sqrt( x[i]*x[i] + y[i]*y[i] + z[i]*z[i] )
func MaxVecNorm(v *data.Slice, queue *cl.CommandQueue, events []*cl.Event) float64 {
	util.Argument(v.NComp() == 3)
	SyncQueues([]*cl.CommandQueue{ClCmdQueue[1]}, []*cl.CommandQueue{queue})
	out := reduceBuf(0, ClCmdQueue[1], events)

	k_reducemaxvecnorm2_async(v.DevPtr(0), v.DevPtr(1), v.DevPtr(2),
		out, 0, v.Len(), reducecfg, ClCmdQueue[1], nil)

	//	results := copybackSlice(out, ClCmdQueue[1], nil)
	//	res := float64(results[0])
	//	for idx := 1; idx < ReduceWorkgroups; idx++ {
	//		res = math.Max(res, float64(results[idx]))
	//	}
	//	return math.Sqrt(float64(res))
	results := copyback(out, ClCmdQueue[1], nil)
	return math.Sqrt(float64(results))
}

// Maximum of the norms of the difference between all vectors (x1,y1,z1) and (x2,y2,z2)
//
//	(dx, dy, dz) = (x1, y1, z1) - (x2, y2, z2)
//	max_i sqrt( dx[i]*dx[i] + dy[i]*dy[i] + dz[i]*dz[i] )
func MaxVecDiff(x, y *data.Slice, queue *cl.CommandQueue, events []*cl.Event) float64 {
	util.Argument(x.Len() == y.Len())
	util.Argument(x.NComp() == 3)
	util.Argument(y.NComp() == 3)
	SyncQueues([]*cl.CommandQueue{ClCmdQueue[1]}, []*cl.CommandQueue{queue})
	out := reduceBuf(0, ClCmdQueue[1], events)

	k_reducemaxvecdiff2_async(x.DevPtr(0), x.DevPtr(1), x.DevPtr(2),
		y.DevPtr(0), y.DevPtr(1), y.DevPtr(2),
		out, 0, x.Len(), reducecfg, ClCmdQueue[1], nil)

	//	results := copybackSlice(out, ClCmdQueue[1], nil)
	//	res := float64(results[0])
	//	for idx := 1; idx < ReduceWorkgroups; idx++ {
	//		res = math.Max(res, float64(results[idx]))
	//	}
	//	return math.Sqrt(float64(res))
	results := copyback(out, ClCmdQueue[1], nil)
	return math.Sqrt(float64(results))
}

var reduceBuffers chan (*cl.MemObject) // pool of 1-float OpenCL buffers for reduce

// return a 1-float and an N-float OPENCL reduction buffer from a pool
// initialized to initVal
func reduceBuf(initVal float32, queue *cl.CommandQueue, events []*cl.Event) unsafe.Pointer {
	if reduceBuffers == nil {
		initReduceBuf()
	}
	buf := <-reduceBuffers
	waitEvent, err := queue.EnqueueFillBuffer(buf, unsafe.Pointer(&initVal), SIZEOF_FLOAT32, 0, ReduceWorkgroups*SIZEOF_FLOAT32, events)
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
func copyback(buf unsafe.Pointer, queue *cl.CommandQueue, events []*cl.Event) float32 {
	var result float32
	if err := queue.Finish(); err != nil {
		fmt.Printf("error waiting for queue to finish in copyback: %+v \n", err)
	}
	MemCpyDtoH(unsafe.Pointer(&result), buf, SIZEOF_FLOAT32)
	// must sync since copying from d to h
	if err := D2HQueue.Finish(); err != nil {
		fmt.Printf("failed to wait for queue to finish in copyback: %+v \n", err)
	}
	reduceBuffers <- (*cl.MemObject)(buf)
	return result
}

// copy back float slice result from GPU and recycle buffer
func copybackSlice(buf unsafe.Pointer, queue *cl.CommandQueue, events []*cl.Event) []float32 {
	result := make([]float32, ReduceWorkgroups)
	if err := queue.Finish(); err != nil {
		fmt.Printf("error waiting for queue to finish in copyback: %+v \n", err)
	}
	MemCpyDtoH(unsafe.Pointer(&result[0]), buf, ReduceWorkgroups*SIZEOF_FLOAT32)
	// must sync since copying from d to h
	if err := D2HQueue.Finish(); err != nil {
		fmt.Printf("failed to wait for queue to finish in copybackSlice: %+v \n", err)
	}
	reduceBuffers <- (*cl.MemObject)(buf)
	return result
}

// initialize pool of 1-float and N-float OPENCL reduction buffers
func initReduceBuf() {
	const N = 128
	reduceBuffers = make(chan *cl.MemObject, N)
	for i := 0; i < N; i++ {
		reduceBuffers <- MemAlloc(ReduceWorkgroups * SIZEOF_FLOAT32)
	}
}

// launch configuration for reduce kernels
// 8 is typ. number of multiprocessors.
// could be improved but takes hardly ~1% of execution time
var reducecfg = &config{Grid: []int{1, 1, 1}, Block: []int{1, 1, 1}}
var ReduceWorkitems int
var ReduceWorkgroups int
