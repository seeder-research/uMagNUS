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
	out, intermed := reduceBuf(0)
	// check input slice for event to synchronize (if any)
	var intEvent *cl.Event
	var event *cl.Event
	syncEvent := in.GetEvent(0)
	if reduceintcfg == nil {
		if syncEvent == nil {
			event = k_reducesum_async(in.DevPtr(0), out, 0,
				in.Len(), reducecfg, nil)
		} else {
			event = k_reducesum_async(in.DevPtr(0), out, 0,
				in.Len(), reducecfg, []*cl.Event{syncEvent})
		}
	} else {
		if syncEvent == nil {
			intEvent = k_reducesum_async(in.DevPtr(0), intermed, 0,
				in.Len(), reduceintcfg, nil)
		} else {
			intEvent = k_reducesum_async(in.DevPtr(0), intermed, 0,
				in.Len(), reduceintcfg, [](*cl.Event){syncEvent})
		}
		event = k_reducesum_async(intermed, out, 0, ClMaxWGSize, reducecfg, [](*cl.Event){intEvent})
	}
	// Must synchronize since out is copied from device back to host
	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in sum: %+v \n", err)
	}
	reduceIntBuffers <- (*cl.MemObject)(intermed)
	return copyback(out)
}

// Dot product
func Dot(a, b *data.Slice) float32 {
	util.Argument(a.NComp() == b.NComp())
	util.Argument(a.Len() == b.Len())
	result := float32(0)
	numComp := a.NComp()
	out := make([]unsafe.Pointer, numComp)
	intermed := make([]unsafe.Pointer, numComp)
	for c := 0; c < numComp; c++ {
		out[c], intermed[c] = reduceBuf(0)
	}
	eventSync := make([]*cl.Event, numComp)
	hostResult := make([]float32, numComp)
	var wg sync.WaitGroup
	// async over components
	for c := 0; c < numComp; c++ {
		eventIntList := []*cl.Event{}
		tmpEvt := a.GetEvent(c)
		if tmpEvt != nil {
			eventIntList = append(eventIntList, tmpEvt)
		}
		tmpEvt = b.GetEvent(c)
		if tmpEvt != nil {
			eventIntList = append(eventIntList, tmpEvt)
		}
		if reduceintcfg == nil {
			if len(eventIntList) > 0 {
				eventSync[c] = k_reducedot_async(a.DevPtr(c), b.DevPtr(c), out[c], 0,
					a.Len(), reducecfg, eventIntList) // all components add to out
			} else {
				eventSync[c] = k_reducedot_async(a.DevPtr(c), b.DevPtr(c), out[c], 0,
					a.Len(), reducecfg, nil) // all components add to out
			}
		} else {
			var barInt *cl.Event
			if len(eventIntList) > 0 {
				barInt = k_reducedot_async(a.DevPtr(c), b.DevPtr(c), intermed[c], 0,
					a.Len(), reduceintcfg, eventIntList) // all components add to intermed
			} else {
				barInt = k_reducedot_async(a.DevPtr(c), b.DevPtr(c), intermed[c], 0,
					a.Len(), reduceintcfg, nil) // all components add to intermed
			}
			eventSync[c] = k_reducesum_async(intermed[c], out[c], 0,
				ClMaxWGSize, reducecfg, []*cl.Event{barInt}) // all components add to out
		}
		wg.Add(1)
		go func(idx int, eventList []*cl.Event, bufferChannel chan<- *cl.MemObject, inBufferPtr, outBufferPtr unsafe.Pointer, res *float32) {
			defer wg.Done()
			if err := cl.WaitForEvents(eventList); err != nil {
				fmt.Printf("WaitForEvents failed at index %d in dot: %+v \n", idx, err)
			}
			bufferChannel <- (*cl.MemObject)(inBufferPtr)
			*res = copyback(outBufferPtr)
		}(c, []*cl.Event{eventSync[c]}, reduceIntBuffers, intermed[c], out[c], &hostResult[c])
	}
	// Must synchronize since result is copied from device back to host
	wg.Wait()
	for _, oVal := range hostResult {
		result += oVal
	}
	return result
}

// Maximum of absolute values of all elements.
func MaxAbs(in *data.Slice) float32 {
	util.Argument(in.NComp() == 1)
	out, intermed := reduceBuf(0)
	// check input slice for event to synchronize (if any)
	var intEvent *cl.Event
	var event *cl.Event
	syncEvent := in.GetEvent(0)
	if reduceintcfg == nil {
		if syncEvent == nil {
			event = k_reducemaxabs_async(in.DevPtr(0), out, 0,
				in.Len(), reducecfg, nil)
		} else {
			event = k_reducemaxabs_async(in.DevPtr(0), out, 0,
				in.Len(), reducecfg, [](*cl.Event){syncEvent})
		}
	} else {
		if syncEvent == nil {
			intEvent = k_reducemaxabs_async(in.DevPtr(0), intermed, 0,
				in.Len(), reduceintcfg, nil)
		} else {
			intEvent = k_reducemaxabs_async(in.DevPtr(0), intermed, 0,
				in.Len(), reduceintcfg, [](*cl.Event){syncEvent})
		}
		event = k_reducemaxabs_async(intermed, out, 0, ClMaxWGSize, reducecfg, [](*cl.Event){intEvent})
	}
	// Must synchronize since out is copied from device back to host
	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in maxabs: %+v \n", err)
	}
	reduceIntBuffers <- (*cl.MemObject)(intermed)
	return copyback(out)
}

// Maximum element-wise difference
func MaxDiff(a, b *data.Slice) []float32 {
	util.Argument(a.NComp() == b.NComp())
	util.Argument(a.Len() == b.Len())
	numComp := a.NComp()
	returnVal := make([]float32, numComp)
	out := make([]unsafe.Pointer, numComp)
	intermed := make([]unsafe.Pointer, numComp)
	for c := 0; c < numComp; c++ {
		out[c], intermed[c] = reduceBuf(0)
	}
	eventSync := make([]*cl.Event, numComp)
	var wg sync.WaitGroup
	for c := 0; c < numComp; c++ {
		eventIntList := []*cl.Event{}
		tmpEvent := a.GetEvent(c)
		if tmpEvent != nil {
			eventIntList = append(eventIntList, tmpEvent)
		}
		tmpEvent = b.GetEvent(c)
		if tmpEvent != nil {
			eventIntList = append(eventIntList, tmpEvent)
		}
		if reduceintcfg == nil {
			if len(eventIntList) > 0 {
				eventSync[c] = k_reducemaxdiff_async(a.DevPtr(c), b.DevPtr(c), out[c], 0,
					a.Len(), reducecfg, eventIntList)
			} else {
				eventSync[c] = k_reducemaxdiff_async(a.DevPtr(c), b.DevPtr(c), out[c], 0,
					a.Len(), reducecfg, nil)
			}
		} else {
			var intEvent *cl.Event
			if len(eventIntList) > 0 {
				intEvent = k_reducemaxdiff_async(a.DevPtr(c), b.DevPtr(c), intermed[c], 0,
					a.Len(), reduceintcfg, eventIntList)
			} else {
				intEvent = k_reducemaxdiff_async(a.DevPtr(c), b.DevPtr(c), intermed[c], 0,
					a.Len(), reduceintcfg, nil)
			}
			eventSync[c] = k_reducemaxabs_async(intermed[c], out[c], 0,
				ClMaxWGSize, reducecfg, [](*cl.Event){intEvent})
		}
		wg.Add(1)
		go func(eventList []*cl.Event, bufferChannel chan<- *cl.MemObject,
			inBufferPtr, outBufferPtr unsafe.Pointer, res *float32) {
			defer wg.Done()
			if err := cl.WaitForEvents(eventList); err != nil {
				fmt.Printf("WaitForEvents failed in maxabs: %+v \n", err)
			}
			bufferChannel <- (*cl.MemObject)(inBufferPtr)
			*res = copyback(outBufferPtr)
		}([]*cl.Event{eventSync[c]}, reduceIntBuffers, intermed[c], out[c], &returnVal[c])
	}
	// Must synchronize since returnVal is copied from device back to host
	wg.Wait()
	return returnVal
}

// Maximum of the norms of all vectors (x[i], y[i], z[i]).
// 	max_i sqrt( x[i]*x[i] + y[i]*y[i] + z[i]*z[i] )
func MaxVecNorm(v *data.Slice) float64 {
	util.Argument(v.NComp() == 3)
	out, intermed := reduceBuf(0)
	// check input slice for events to synchronize (if any)
	var intEvent *cl.Event
	var event *cl.Event
	syncEvent := []*cl.Event{}
	for c := 0; c < v.NComp(); c++ {
		tmpEvent := v.GetEvent(c)
		if tmpEvent != nil {
			syncEvent = append(syncEvent, tmpEvent)
		}
	}
	if reduceintcfg == nil {
		if len(syncEvent) > 0 {
			event = k_reducemaxvecnorm2_async(v.DevPtr(0), v.DevPtr(1), v.DevPtr(2),
				out, 0, v.Len(), reducecfg, syncEvent)
		} else {
			event = k_reducemaxvecnorm2_async(v.DevPtr(0), v.DevPtr(1), v.DevPtr(2),
				out, 0, v.Len(), reducecfg, nil)
		}
	} else {
		if len(syncEvent) > 0 {
			intEvent = k_reducemaxvecnorm2_async(v.DevPtr(0), v.DevPtr(1), v.DevPtr(2),
				intermed, 0, v.Len(), reduceintcfg, syncEvent)
		} else {
			intEvent = k_reducemaxvecnorm2_async(v.DevPtr(0), v.DevPtr(1), v.DevPtr(2),
				intermed, 0, v.Len(), reduceintcfg, nil)
		}
		event = k_reducemaxabs_async(intermed, out, 0, ClMaxWGSize, reducecfg, [](*cl.Event){intEvent})
	}
	// Must synchronize since out is copied from device back to host
	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in maxvecnorm: %+v \n", err)
	}
	reduceIntBuffers <- (*cl.MemObject)(intermed)
	return math.Sqrt(float64(copyback(out)))
}

// Maximum of the norms of the difference between all vectors (x1,y1,z1) and (x2,y2,z2)
// 	(dx, dy, dz) = (x1, y1, z1) - (x2, y2, z2)
// 	max_i sqrt( dx[i]*dx[i] + dy[i]*dy[i] + dz[i]*dz[i] )
func MaxVecDiff(x, y *data.Slice) float64 {
	util.Argument(x.Len() == y.Len())
	util.Argument(x.NComp() == 3)
	util.Argument(y.NComp() == 3)
	out, intermed := reduceBuf(0)
	// check input slice for event to synchronize (if any)
	var intEvent *cl.Event
	var event *cl.Event
	syncEvent := []*cl.Event{}
	for c := 0; c < 3; c++ {
		tmpEvent := x.GetEvent(c)
		if tmpEvent != nil {
			syncEvent = append(syncEvent, tmpEvent)
		}
		tmpEvent = y.GetEvent(c)
		if tmpEvent != nil {
			syncEvent = append(syncEvent, tmpEvent)
		}
	}
	if reduceintcfg == nil {
		if len(syncEvent) > 0 {
			event = k_reducemaxvecdiff2_async(x.DevPtr(0), x.DevPtr(1), x.DevPtr(2),
				y.DevPtr(0), y.DevPtr(1), y.DevPtr(2),
				out, 0, x.Len(), reducecfg, syncEvent)
		} else {
			event = k_reducemaxvecdiff2_async(x.DevPtr(0), x.DevPtr(1), x.DevPtr(2),
				y.DevPtr(0), y.DevPtr(1), y.DevPtr(2),
				out, 0, x.Len(), reducecfg, nil)
		}
	} else {
		if len(syncEvent) > 0 {
			intEvent = k_reducemaxvecdiff2_async(x.DevPtr(0), x.DevPtr(1), x.DevPtr(2),
				y.DevPtr(0), y.DevPtr(1), y.DevPtr(2),
				intermed, 0, x.Len(), reduceintcfg, syncEvent)
		} else {
			intEvent = k_reducemaxvecdiff2_async(x.DevPtr(0), x.DevPtr(1), x.DevPtr(2),
				y.DevPtr(0), y.DevPtr(1), y.DevPtr(2),
				intermed, 0, x.Len(), reduceintcfg, nil)
		}
		event = k_reducemaxabs_async(intermed, out, 0, ClMaxWGSize, reducecfg, [](*cl.Event){intEvent})
	}

	// Must synchronize since out is copied from device back to host
	if err := cl.WaitForEvents([]*cl.Event{event}); err != nil {
		fmt.Printf("WaitForEvents failed in maxvecdiff: %+v \n", err)
	}
	reduceIntBuffers <- (*cl.MemObject)(intermed)
	return math.Sqrt(float64(copyback(out)))
}

var reduceBuffers chan (*cl.MemObject)    // pool of 1-float OpenCL buffers for reduce
var reduceIntBuffers chan (*cl.MemObject) // pool of 1-float OpenCL buffers for reduce

// return a 1-float and an N-float OPENCL reduction buffer from a pool
// initialized to initVal
func reduceBuf(initVal float32) (unsafe.Pointer, unsafe.Pointer) {
	if reduceBuffers == nil {
		initReduceBuf()
	}
	buf := <-reduceBuffers
	interBuf := <-reduceIntBuffers
	waitEvent, err := ClCmdQueue.EnqueueFillBuffer(buf, unsafe.Pointer(&initVal), SIZEOF_FLOAT32, 0, SIZEOF_FLOAT32, nil)
	if err != nil {
		fmt.Printf("reduceBuf failed: %+v \n", err)
		return nil, nil
	}
	err = cl.WaitForEvents([]*cl.Event{waitEvent})
	if err != nil {
		fmt.Printf("First WaitForEvents in reduceBuf failed: %+v \n", err)
		return nil, nil
	}
	waitEvent, err = ClCmdQueue.EnqueueFillBuffer(interBuf, unsafe.Pointer(&initVal), SIZEOF_FLOAT32, 0, ClMaxWGSize*SIZEOF_FLOAT32, nil)
	if err != nil {
		fmt.Printf("reduceBuf failed: %+v \n", err)
		return nil, nil
	}
	err = cl.WaitForEvents([]*cl.Event{waitEvent})
	if err != nil {
		fmt.Printf("Second WaitForEvents in reduceBuf failed: %+v \n", err)
		return nil, nil
	}
	return (unsafe.Pointer)(buf), (unsafe.Pointer)(interBuf)
}

// copy back single float result from GPU and recycle buffer
func copyback(buf unsafe.Pointer) float32 {
	var result float32
	MemCpyDtoH(unsafe.Pointer(&result), buf, SIZEOF_FLOAT32)
	reduceBuffers <- (*cl.MemObject)(buf)
	return result
}

// initialize pool of 1-float and N-float OPENCL reduction buffers
func initReduceBuf() {
	const N = 128
	reduceBuffers = make(chan *cl.MemObject, N)
	reduceIntBuffers = make(chan *cl.MemObject, N)
	for i := 0; i < N; i++ {
		reduceBuffers <- MemAlloc(SIZEOF_FLOAT32)
		reduceIntBuffers <- MemAlloc(ClMaxWGSize * SIZEOF_FLOAT32)
	}
}

// launch configuration for reduce kernels
// 8 is typ. number of multiprocessors.
// could be improved but takes hardly ~1% of execution time
var reducecfg = &config{Grid: []int{1, 1, 1}, Block: []int{1, 1, 1}}
var reduceintcfg = &config{Grid: []int{8, 1, 1}, Block: []int{1, 1, 1}}

func UpdateReduceConfigs(c []int) {
	numItems := c[0] * c[1] * c[2] // total number of items to sum

	// Set up reduce config for final reduce step
	reducecfg = &config{Grid: []int{ClMaxWGSize, 1, 1}, Block: []int{ClMaxWGSize, 1, 1}}

	// Find reduce config for intermediate reduce step
	if numItems <= ClMaxWGSize {
		reduceintcfg = nil
	} else {
		if numItems >= ClMaxWGSize*ClMaxWGNum {
			reduceintcfg = &config{Grid: []int{ClMaxWGSize * ClMaxWGNum, 1, 1}, Block: []int{ClMaxWGSize, 1, 1}}
		} else {
			for ii0 := ClMaxWGSize; ii0 < numItems; ii0 += ClMaxWGSize {
				reduceintcfg = &config{Grid: []int{ii0, 1, 1}, Block: []int{ClMaxWGSize, 1, 1}}
			}
		}
	}
}
