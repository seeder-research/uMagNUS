package engine

import (
	"fmt"

	cl "github.com/seeder-research/uMagNUS/cl"
	data "github.com/seeder-research/uMagNUS/data"
	opencl "github.com/seeder-research/uMagNUS/opencl"
	util "github.com/seeder-research/uMagNUS/util"
)

func sInRegion(q Quantity, r int) ScalarField {
	return AsScalarField(inRegion(q, r))
}

func vInRegion(q Quantity, r int) VectorField {
	return AsVectorField(inRegion(q, r))
}

func sOneRegion(q Quantity, r int) *sOneReg {
	util.Argument(q.NComp() == 1)
	return &sOneReg{oneReg{q, r}}
}

func vOneRegion(q Quantity, r int) *vOneReg {
	util.Argument(q.NComp() == 3)
	return &vOneReg{oneReg{q, r}}
}

type sOneReg struct{ oneReg }

func (q *sOneReg) Average() float64 { return q.average()[0] }

type vOneReg struct{ oneReg }

func (q *vOneReg) Average() data.Vector { return unslice(q.average()) }

// represents a new quantity equal to q in the given region, 0 outside.
type oneReg struct {
	parent Quantity
	region int
}

func inRegion(q Quantity, region int) Quantity {
	return &oneReg{q, region}
}

func (q *oneReg) NComp() int             { return q.parent.NComp() }
func (q *oneReg) Name() string           { return fmt.Sprint(NameOf(q.parent), ".region", q.region) }
func (q *oneReg) Unit() string           { return UnitOf(q.parent) }
func (q *oneReg) Mesh() *data.Mesh       { return MeshOf(q.parent) }
func (q *oneReg) EvalTo(dst *data.Slice) { EvalTo(q, dst) }

// returns a new slice equal to q in the given region, 0 outside.
func (q *oneReg) Slice() (*data.Slice, bool) {
	// sync in the beginning
	if err := opencl.WaitAllQueuesToFinish(); err != nil {
		fmt.Printf("error waiting for queue to finish in onereg.slice: %+v \n", err)
	}
	// checkout queues
	seqQueue := opencl.ClCmdQueue[0]
	q1idx, q2idx, q3idx := opencl.CheckoutQueue(), opencl.CheckoutQueue(), opencl.CheckoutQueue()
	defer opencl.CheckinQueue(q1idx)
	defer opencl.CheckinQueue(q2idx)
	defer opencl.CheckinQueue(q3idx)
	queues := []*cl.CommandQueue{opencl.ClCmdQueue[q1idx], opencl.ClCmdQueue[q2idx], opencl.ClCmdQueue[q3idx]}

	src := ValueOf(q.parent)
	defer opencl.Recycle(src)
	out := opencl.Buffer(q.NComp(), q.Mesh().Size())
	opencl.RegionSelect(out, src, regions.Gpu(), byte(q.region), queues, nil)
	// sync before returning
	opencl.SyncQueues([]*cl.CommandQueue{seqQueue}, queues)
	if err := seqQueue.Finish(); err != nil {
		fmt.Printf("error waiting for queue to finish after onereg.slice: %+v \n", err)
	}
	return out, true
}

func (q *oneReg) average() []float64 {
	slice, r := q.Slice()
	if r {
		defer opencl.Recycle(slice)
	}
	avg := sAverageUniverse(slice)
	sDiv(avg, regions.volume(q.region))
	return avg
}

func (q *oneReg) Average() []float64 { return q.average() }

// slice division
func sDiv(v []float64, x float64) {
	for i := range v {
		v[i] /= x
	}
}
