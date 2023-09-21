package opencl

// TODO: update call to reducesum and reducedot to
// launch kernels with number of workgroups and workitems
// optimized to size of input buffer

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

var reduceInterBuffers chan (*cl.MemObject)

/*
Need to update the reduction sum and dot algorithms to balance the distribution of inputs to the work-groups
Less of a problem for max and min because they are direct comparisons

for sum and dot, the magnitude of intermediate values depend on how many input values were summed to obtain
them. Thus, the entire input data should be distributed in a binary tree for the summation, and the work-
groups should "synchronize" at a fixed level of the tree. Small input sizes will need fewer work-groups
that can efficiently calculate for trees having short depths. Larger input sizes will need larger number of
possibly small work-groups.
*/
var sumImpl func(*data.Slice, unsafe.Pointer, float32, *cl.CommandQueue, []*cl.Event) *cl.Event
var dotImpl func(*data.Slice, *data.Slice, unsafe.Pointer, float32, *cl.CommandQueue, []*cl.Event) *cl.Event

type reduceDotSumStage interface {
	reduceDotKernel(unsafe.Pointer, unsafe.Pointer, unsafe.Pointer, float32, *cl.CommandQueue, []*cl.Event) *cl.Event
	reduceSumKernel(unsafe.Pointer, unsafe.Pointer, float32, *cl.CommandQueue, []*cl.Event) *cl.Event
	SetFac(int)
	SetGrpN(int)
	SetNElem(int)
	SetConfig(config)
}

type reduceFunc struct {
	fac   int
	grp_n int
	nElem int
	cfg   config
}

func (r *reduceFunc) SetFac(i int) {
	r.fac = i
}

func (r *reduceFunc) SetGrpN(i int) {
	r.grp_n = i
}

func (r *reduceFunc) SetNElem(i int) {
	r.nElem = i
}

func (r *reduceFunc) SetConfig(c config) {
	r.cfg = c
}

// Use reduceSumStage function to build various reducesum functions.
// If number of elements is small, call reduceSumStage with:
//
//	outPtr = nil
//	gridCfg with one workgroup
//	fac = 0
//	group_n = in.Len()
//
// If number of elements is large side of small, call reduceSumStage with:
//
//	outPtr = nil
//	gridCfg with two workgroups
//	fac = 0
//	group_n = ceil(in.Len() / 2)
//
// If number of elements is large side of small, call reduceSumStage with:
//
//	outPtr = nil
//	gridCfg with two workgroups
//	fac = 0
//	group_n = ceil(in.Len() / 2)
func (r *reduceFunc) reduceSumKernel(inPtr, outPtr unsafe.Pointer, initVal float32, q *cl.CommandQueue, ewl []*cl.Event) *cl.Event {

	// Launch kernel
	return k_reducesum_async(inPtr, outPtr, initVal, r.fac, r.grp_n,
		r.nElem, r.cfg, ewl, q)
}

func (r *reduceFunc) reduceDotKernel(inPtr0, inPtr1, outPtr unsafe.Pointer, initVal float32, q *cl.CommandQueue, ewl []*cl.Event) *cl.Event {

	// Launch kernel
	return k_reducedot_async(inPtr0, inPtr1, outPtr, initVal, r.fac, r.grp_n,
		r.nElem, r.cfg, ewl, q)
}

type reduceDotSumPlans struct {
	Plans  []reduceDotSumStage
	inLen  int
	bufSz  int
	nStage int
}

var globalReduceDotSumPlan reduceDotSumPlans

func (r *reduceDotSumPlans) ZeroBuffer(q *cl.CommandQueue) *cl.Event {
	zero := uint8(0)
	ev, err := q.EnqueueFillBuffer(r.buf, (unsafe.Pointer)(&zero), 1, 0, r.bufSz)
	if err != nil {
		fmt.Println("failed to zero buffer in reduceSumPlans: %+v ", err)
		return nil
	}
	return ev
}

func CreateTwoStageDotSum(bufSize int) *reduceDotSumPlans {
	plans := make([]reduceDotSumStage, 2)
	return &reduceDotSumPlans{Plans: plans, bufSz: bufSize, inLen: -1, nStage: 2}
}

func CreateOneStageDotSum() *reduceDotSumPlans {
	plans := make([]reduceDotSumStage, 1)
	return &reduceDotSumPlans{Plans: plans, bufSz: -1, inLen: -1, nStage: 1}
}

func (p *reduceDotSumPlans) SetLength(i int) {
	p.inLen = i
}

func (p *reduceDotSumPlans) GetLength() int {
	return p.inLen
}

func (p *reduceDotSumPlans) GetStageCount() int {
	return p.nStage
}

func newReducePlan(i int) reduceDotSumPlans {
	defaultBlockSize := 64
	var outPlan reduceDotSumPlans
	if i < 65536 {
		// Configuration for reducesum and reducedot
		outPlan = CreateOneStageDotSum()
		outPlan.Plans[0].SetFac(0)
		outPlan.Plans[0].SetGrpN(i)
		outPlan.Plans[0].SetNElem(i)
		outPlan.Plans[0].SetConfig(config{Grid: {defaultBlockSize, 1, 1}, Block: {defaultBlockSize, 1, 1}})

		initReduceInterBuffers(-1)
		sumImpl = sumOneStageImpl
		dotImpl = dotOneStageImpl
	} else if i < 524288 {
		// Configuration for reducesum and reducedot
		outPlan = CreateOneStageDotSum()
		outPlan.Plans[0].SetFac(0)
		outPlan.Plans[0].SetGrpN(int(math.Ceil((float64)(i) * float64(0.5))))
		outPlan.Plans[0].SetNElem(i)
		outPlan.Plans[0].SetConfig(config{Grid: {2 * defaultBlockSize, 1, 1}, Block: {defaultBlockSize, 1, 1}})

		initReduceInterBuffers(-1)
		sumImpl = sumOneStageImpl
		dotImpl = dotOneStageImpl
	} else if i < 2097152 {
		// Input buffer is large enough that we need two stages.

		// First stage
		// 32x65536: to handle with 128 per workitems, which means each workgroup
		// manages up to 8192 items
		groupCount := int(math.Ceil((float64)(i) / (float64)(8192)))
		grpN := int(math.Ceil((float64)(i) / (float64)(groupCount)))
		outPlan = CreateTwoStageDotSum(groupCount * SIZEOF_FLOAT32)
		if groupCount > 128 {
			outPlan.Plans[0].SetFac(2)
		} else {
			outPlan.Plans[0].SetFac(1)
		}
		outPlan.Plans[0].SetGrpN(grpN)
		outPlan.Plans[0].SetNElem(i)
		outPlan.Plans[0].SetConfig(config{Grid: {groupCount * defaultBlockSize, 1, 1}, Block: {defaultBlockSize, 1, 1}})

		// Second stage
		outPlan.Plans[1].SetFac(0)
		outPlan.Plans[1].SetGrpN(groupCount)
		outPlan.Plans[1].SetNElem(groupCount)
		outPlan.Plans[1].SetConfig(config{Grid: {defaultBlockSize, 1, 1}, Block: {defaultBlockSize, 1, 1}})

		initReduceInterBuffers(groupCount)
		sumImpl = sumTwoStageImpl
		dotImpl = dotTwoStageImpl
	} else if i < 33554432 {
		// Input buffer is large enough that we need two stages.

		// First stage
		// 512x65536: to handle with 2048 per workitems, which means each workgroup
		// manages up to 131072 items
		groupCount := int(math.Ceil((float64)(i) / (float64)(131072)))
		grpN := int(math.Ceil((float64)(i) / (float64)(groupCount)))
		outPlan = CreateTwoStageDotSum(groupCount * SIZEOF_FLOAT32)
		if groupCount > 128 {
			outPlan.Plans[0].SetFac(2)
		} else {
			outPlan.Plans[0].SetFac(1)
		}
		outPlan.Plans[0].SetGrpN(grpN)
		outPlan.Plans[0].SetNElem(i)
		outPlan.Plans[0].SetConfig(config{Grid: {groupCount * defaultBlockSize, 1, 1}, Block: {defaultBlockSize, 1, 1}})

		// Second stage
		outPlan.Plans[1].SetFac(0)
		outPlan.Plans[1].SetGrpN(groupCount)
		outPlan.Plans[1].SetNElem(groupCount)
		outPlan.Plans[1].SetConfig(config{Grid: {defaultBlockSize, 1, 1}, Block: {defaultBlockSize, 1, 1}})

		initReduceInterBuffers(groupCount)
		sumImpl = sumTwoStageImpl
		dotImpl = dotTwoStageImpl
	} else {
		// Input buffer is large enough that we need two stages.
		// The extreme large case

		// Check how many items can be managed by 512 groups.
		groupCount := 512
		grpN := int(math.Ceil((float64)(i) / (float64)(512)))
		if grpN > 4194304 {
			groupCount = int(math.Ceil((float64)(i) / (float64)(4194304)))
			grpN = int(math.Ceil((float64)(i) / (float64)(groupCount)))
		}
		outPlan = CreateTwoStageDotSum(groupCount * SIZEOF_FLOAT32)
		outPlan.Plans[0].SetGrpN(grpN)
		outPlan.Plans[0].SetFac(1)
		outPlan.Plans[0].SetNElem(i)
		outPlan.Plans[0].SetConfig(config{Grid: {groupCount * defaultBlockSize, 1, 1}, Block: {defaultBlockSize, 1, 1}})

		// Second stage
		outPlan.Plans[1].SetFac(0)
		outPlan.Plans[1].SetGrpN(groupCount)
		outPlan.Plans[1].SetNElem(groupCount)
		outPlan.Plans[1].SetConfig(config{Grid: {defaultBlockSize, 1, 1}, Block: {defaultBlockSize, 1, 1}})

		initReduceInterBuffers(groupCount)
		sumImpl = sumTwoStageImpl
		dotImpl = dotTwoStageImpl
	}

	outPlan.SetLength(i)

	return outPlan
}

// TODO: Need to call this everything the mesh is updated
func UpdateReduceSumLen(i int) {
	globalReducePlan = newReducePlan(i)
}

func (p *reduceDotSumPlans) sumOneStageImpl_(in *data.Slice, outPtr unsafe.Pointer, initVal float32, q *cl.CommandQueue, ewl []*cl.Event) *cl.Event {
	util.Argument(in.Len() == p.GetLength())
	if len(ewl) != 0 {
		_, err := q.EnqueueBarrierWithWaitList(ewl)
		if err != nil {
			fmt.Println("ERROR failed to enqueue barrier in beginning of Sum_impl: %+v ", err)
			return nil
		}
	}
	return p.Plan[0].reduceSumKernel(in.DevPtr(0), outPtr, initVal, q, nil)
}

func (p *reduceDotSumPlans) dotOneStageImpl_(in0, in1 *data.Slice, outPtr unsafe.Pointer, initVal float32, q *cl.CommandQueue, ewl []*cl.Event) *cl.Event {
	util.Argument(in0.Len() == p.GetLength())
	util.Argument(in1.Len() == p.GetLength())
	if len(ewl) != 0 {
		_, err := q.EnqueueBarrierWithWaitList(ewl)
		if err != nil {
			fmt.Println("ERROR failed to enqueue barrier in beginning of Sum_impl: %+v ", err)
			return nil
		}
	}
	return p.Plan[0].reduceDotKernel(in0.DevPtr(0), in1.DevPtr(0), outPtr, initVal, q, nil)
}

func (p *reduceDotSumPlans) sumTwoStageImpl_(in *data.Slice, outPtr unsafe.Pointer, initVal float32, q *cl.CommandQueue, ewl []*cl.Event) *cl.Event {
	util.Argument(in.Len() == p.GetLength())
	util.Argument(p.GetStageCount() == 2)
	if len(ewl) != 0 {
		_, err := q.EnqueueBarrierWithWaitList(ewl)
		if err != nil {
			fmt.Println("ERROR failed to enqueue barrier in beginning of Sum_impl: %+v ", err)
			return nil
		}
	}
	buf_ := <-reduceInterBuffers // Needs to be checked in after reduction ends
	ev := p.Plans[0].reduceSumKernel(in.DevPtr(0), buf_, initVal, q, nil)
	ev = p.Plans[1].reduceSumKernel(buf_, outPtr, float32(0.0), q, nil)

	// Copy back to host in goroutine
	reduceWaitGroup.Add(1)
	reduceItem <- reduceEventAndPointer{event: []*cl.Event{event}, cmdQ: q, ptr: out, idx: 0, bufr: buf_}

	return ev
}

func (p *reduceDotSumPlans) dotTwoStageImpl_(in0, in1 *data.Slice, outPtr unsafe.Pointer, initVal float32, q *cl.CommandQueue, ewl []*cl.Event) *cl.Event {
	util.Argument(in0.Len() == p.GetLength())
	util.Argument(in1.Len() == p.GetLength())
	util.Argument(p.GetStageCount() == 2)
	if len(ewl) != 0 {
		_, err := q.EnqueueBarrierWithWaitList(ewl)
		if err != nil {
			fmt.Println("ERROR failed to enqueue barrier in beginning of Sum_impl: %+v ", err)
			return nil
		}
	}
	buf_ := <-reduceInterBuffers // Needs to be checked in after reduction ends
	ev := p.Plans[0].reduceDotKernel(in0.DevPtr(0), in1.DevPtr(0), (unsafe.Pointer)(buf_), initVal, q, nil)
	ev = p.Plans[1].reduceSumKernel((unsafe.Pointer)(buf_), outPtr, float32(0.0), q, nil)

	// Copy back to host in goroutine
	reduceWaitGroup.Add(1)
	reduceItem <- reduceEventAndPointer{event: []*cl.Event{event}, cmdQ: q, ptr: out, idx: 0, bufr: buf_}

	return ev
}

// When mesh is updated, need to setup the correct function as the sumImpl function
func sumOneStageImpl(in *data.Slice, outPtr unsafe.Pointer, initVal float32, q *cl.CommandQueue, ewl []*cl.Event) *cl.Event {
	return globalReducePlan.sumOneStageImpl_(in, outPtr, initVal, q, ewl)
}

func sumTwoStageImpl(in *data.Slice, outPtr unsafe.Pointer, initVal float32, q *cl.CommandQueue, ewl []*cl.Event) *cl.Event {
	return globalReducePlan.sumTwoStageImpl_(in, outPtr, initVal, q, ewl)
}

func dotOneStageImpl(in0, in1 *data.Slice, outPtr unsafe.Pointer, initVal float32, q *cl.CommandQueue, ewl []*cl.Event) *cl.Event {
	return globalReducePlan.dotOneStageImpl_(in0, in1, outPtr, initVal, q, ewl)
}

func dotTwoStageImpl(in0, in1 *data.Slice, outPtr unsafe.Pointer, initVal float32, q *cl.CommandQueue, ewl []*cl.Event) *cl.Event {
	return globalReducePlan.dotTwoStageImpl_(in0, in1, outPtr, initVal, q, ewl)
}

// Sum of all elements.
func Sum(in *data.Slice, q *cl.CommandQueue, ewl []*cl.Event) float32 {
	util.Argument(in.NComp() == 1)
	out := reduceBuf(0)

	// Ensure no other reduction kernel is running
	reduceWaitGroup.Wait()

	// Launch kernel
	event := sumImpl(in.DevPtr(0), out, float32(0), q, ewl)

	// Ensure all reduction kernel has completed
	reduceWaitGroup.Wait()
	tmp := <-reduceRes

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
	//hostResult := make([]float32, numComp)

	// Ensure no other reduction kernel is running
	reduceWaitGroup.Wait()
	for c := 0; c < numComp; c++ {
		// Launch kernel
		event := dotImpl(a.DevPtr(c), b.DevPtr(c), out[c], float32(0.0), q[c], ewl) // all components add to out
	}

	// Ensure all reduction kernel has completed
	reduceWaitGroup.Wait()
	for c := 0; c < numComp; c++ {
		tmp := <-reduceRes
		result += tmp.res
	}

	return result
}

// initialize pool of N-float OPENCL reduction buffers
func initReduceInterBuffer(num int) {
	const N = 32
	if num < 0 {
		// Clear intermediate buffers
		for i := 0; i < N; i++ {
			buf_ := <-reduceInterBuffers
			buf_.Release()
			reduceInterBuffers <- nil
		}
	} else {
		reduceInterBuffer = make(chan *cl.MemObject, N)
		for i := 0; i < N; i++ {
			buf_ := <-reduceInterBuffers
			buf_.Release()
			reduceInterBuffers <- MemAlloc(num * SIZEOF_FLOAT32)
		}
	}
}
